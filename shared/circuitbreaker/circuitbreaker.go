package circuitbreaker

import (
	"time"

	"github.com/sony/gobreaker"
)

type CircuitBreaker struct {
	*gobreaker.CircuitBreaker
}

func NewCircuitBreaker(name string) *CircuitBreaker {
	cbSettings := gobreaker.Settings{
		Name:        name,
		MaxRequests: 5,
		Interval:    30 * time.Second,
		Timeout:     10 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 3
		},
	}
	return &CircuitBreaker{
		CircuitBreaker: gobreaker.NewCircuitBreaker(cbSettings),
	}
}
