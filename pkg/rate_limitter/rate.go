package ratelimitter

import (
	"fmt"
	"time"
)

type Nexter interface {
	Next()
}

type RateLimit struct {
	exitch   chan struct{}
	ipch     chan string
	errorch  chan error
	bucket   map[string]int
	duration time.Duration
	capacity int
}

func New(capacity int, duration time.Duration) *RateLimit {
	r := &RateLimit{
		bucket:   make(map[string]int),
		exitch:   make(chan struct{}),
		errorch:  make(chan error),
		ipch:     make(chan string),
		capacity: capacity,
		duration: duration,
	}

	go r.refill()
	return r
}

func (r *RateLimit) Start(ip string, nexter Nexter) error {
	r.ipch <- ip
	if err := <-r.errorch; err != nil {
		return err
	}

	nexter.Next()
	return nil
}

func (r *RateLimit) Stop() {
	close(r.exitch)
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
		case <-r.exitch:
			ticker.Stop()
			return
		}
	}
}
