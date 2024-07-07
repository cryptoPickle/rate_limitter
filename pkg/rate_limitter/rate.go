package ratelimitter

import "time"

type Nexter interface {
	Next()
	AbortWithStatusJSON(int, any)
}

type RateLimit struct {
	exitch chan struct{}
	bucket map[string]int
	Nexter
	duration time.Duration
	capacity int
}

func New(capacity int, duration time.Duration, next Nexter) *RateLimit {
	r := &RateLimit{
		bucket:   make(map[string]int),
		capacity: capacity,
		duration: duration,
	}
	go r.refill()
	return r
}

func (r *RateLimit) Start() {
}

func (r *RateLimit) Stop() {
	close(r.exitch)
}

func (r *RateLimit) refill(ip string) {
	ticker := time.NewTicker(r.duration)

	for {
		select {
		case <-ticker.C:
			if r.bucket[ip] > r.capacity && r.bucket[ip] != 0 {
				r.bucket[ip] = r.bucket[ip] - 1
			}
		case <-r.exitch:
			return
		}
	}
}
