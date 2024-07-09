package ratelimitter

import (
	"errors"
	"fmt"
	"time"

	"github.com/cryptoPickle/rate_limitter/types"
	"github.com/sirupsen/logrus"
)

type Bucketer interface {
	Set(string, int) error
	Get(string) (int, error)
	Has(string) bool
	DecrementAll() error
}
type RateLimit struct {
	ipch     chan string
	ratech   chan error
	bucket   Bucketer
	duration time.Duration
	capacity int
}

func New(capacity int, duration time.Duration, bucket Bucketer) *RateLimit {
	r := &RateLimit{
		bucket:   bucket,
		ratech:   make(chan error),
		ipch:     make(chan string),
		capacity: capacity,
		duration: duration,
	}

	go r.refill()
	return r
}

func (r *RateLimit) Start(nexter types.Nexter) error {
	val, ok := nexter.Get("ip")

	if !ok {
		return errors.New("can't get ip")
	}

	if ip, ok := val.(string); ok {
		r.ipch <- ip
		if err := <-r.ratech; err != nil {
			return err
		}

		nexter.Next()
		return nil
	}
	return errors.New("ip adress should be string")
}

func (r *RateLimit) refill() {
	ticker := time.NewTicker(r.duration)

	for {
		select {
		case <-ticker.C:
			if err := r.bucket.DecrementAll(); err != nil {
				logrus.Errorf("can not decrement tokens... error: %s\n", err)
			}
		case ip := <-r.ipch:
			reqCount, err := r.getRequestCount(ip)
			if err != nil {
				logrus.Error(err)
			}
			if err := r.updateRate(ip, reqCount); err != nil {
				logrus.Error(err)
			}
		}
	}
}

func (r *RateLimit) updateRate(ip string, reqCount int) error {
	if reqCount >= r.capacity {
		r.ratech <- fmt.Errorf("rate limit exceed")
	} else {
		if err := r.bucket.Set(ip, reqCount+1); err != nil {
			return fmt.Errorf("can not update request count err: %s ", err)
		}
		r.ratech <- nil
	}
	return nil
}

func (r *RateLimit) getRequestCount(ip string) (int, error) {
	if ok := r.bucket.Has(ip); !ok {
		if err := r.bucket.Set(ip, 0); err != nil {
			return 0, fmt.Errorf("can not initialise the ip count err: %s", err)
		}
	}
	reqCount, err := r.bucket.Get(ip)
	if err != nil {
		return 0, fmt.Errorf("can not get the current request count for %v, error: %s", ip, err)
	}

	return reqCount, nil
}
