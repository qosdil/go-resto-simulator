package customer

import (
	"errors"
	"math/rand/v2"
	"sync/atomic"
	"time"
)

var (
	lastNumber        atomic.Uint32
	minDiningDuration atomic.Uint32
	maxDiningDuration atomic.Uint32
)

type Customer struct {
	Number uint8
}

// New creates a new customer with a unique number.
func New() (*Customer, error) {
	if minDiningDuration.Load() == 0 || maxDiningDuration.Load() == 0 {
		return nil, errors.New("min, max of dining duration not set properly")
	}

	lastNumber.Add(1)
	return &Customer{Number: uint8(lastNumber.Load())}, nil
}

// Dine simulates the dining activity of a customer by sleeping for a random duration between the minimum and maximum dining durations.
func (c *Customer) Dine() int {
	var min, max = int(minDiningDuration.Load()), int(maxDiningDuration.Load())
	sec := min + rand.IntN(max-min)
	time.Sleep(time.Duration(sec) * time.Second)
	return sec
}

func SetMinMaxDiningDurations(min, max uint8) {
	minDiningDuration.Store(uint32(min))
	maxDiningDuration.Store(uint32(max))
}
