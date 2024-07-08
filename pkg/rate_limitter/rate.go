package ratelimitter

import (
	"errors"
	"fmt"
	"time"

	"github.com/cryptoPickle/rate_limitter/types"
)

type RateLimit struct {
	ipch     chan string
	errorch  chan error
	bucket   map[string]int
	duration time.Duration
	capacity int
}

func New(capacity int, duration time.Duration) *RateLimit {
	r := &RateLimit{
		bucket:   make(map[string]int),
		errorch:  make(chan error),
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
		if err := <-r.errorch; err != nil {
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
			for key, val := range r.bucket {
				if val > 0 {
					r.bucket[key]--
				}
			}
		case val := <-r.ipch:
			if _, exists := r.bucket[val]; !exists {
				r.bucket[val] = 0
			}
			if r.bucket[val] >= r.capacity {
				r.errorch <- fmt.Errorf("rate limit exceed")
			} else {
				r.bucket[val]++
				r.errorch <- nil
			}
		}
	}
}
