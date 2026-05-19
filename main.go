package main

import (
	"context"
	"fmt"
	"go-resto-simulator/internal/customer"
	"go-resto-simulator/internal/reservation"
	"log"
	"math/rand/v2"
	"slices"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// Change these constants to customize the simulation.
	custArrivalMinSecs    = 0
	custArrivalMaxSecs    = 20
	custMinDiningDuration = 60
	custMaxDiningDuration = 90
	numTables             = 10

	// Minutes to simulate hours of resto's "open" status, e.g. 11 hours from 10.00a to 09.00p.
	servingDurationMins = 11
)

var (
	isOpen atomic.Bool // whether resto is open or closed
	lock   sync.RWMutex
)

func simulateActivities(ctx context.Context, startTime time.Time, custArrivalMinSecs, custArrivalMaxSecs int,
	lock *sync.RWMutex) {
start:
	for {
		select {
		case <-ctx.Done():
			logText := "we have served for %d hours, closing soon after all dining customers are done, new customers"
			logText += " can come back tomorrow"
			log.Printf(logText, servingDurationMins)

			// Set the resto's status to closed, so that no new customers will come in.
			// The resto will be fully closed after all dining customers are done.
			isOpen.Store(false)

			break start
		default:
			sec := custArrivalMinSecs + rand.IntN(custArrivalMaxSecs-custArrivalMinSecs)
			time.Sleep(time.Duration(sec) * time.Second)
			simulateActivity(startTime, lock)
		}
	}

	log.Println("simulateActivities() is done")
}

func simulateActivity(startTime time.Time, lock *sync.RWMutex) {
	c, err := customer.New()
	if err != nil {
		panic(err)
	}

	log.Printf("customer %d arrived on minute %d", c.Number, int(time.Since(startTime).Seconds()))
	r, err := reservation.Reserve(c, lock)
	if err == reservation.ErrAllTablesReserved {
		log.Printf("customer %d left because they failed to reserve: %v", c.Number, err)
		return
	}

	if err != nil {
		log.Printf("unknown error: %v", err)
		return
	}

	log.Printf("customer %d has just reserved table %d", c.Number, r.Table.Number)
}

func main() {
	reservation.SetTables(numTables)
	customer.SetMinMaxDiningDurations(custMinDiningDuration, custMaxDiningDuration)

	// Open the resto, start serving customers.
	isOpen.Store(true)

	ctx, cancelFunc := context.WithTimeout(context.Background(), servingDurationMins*time.Minute)
	defer cancelFunc()

	// Simulate the arrival of customers and their dining activities in a separate goroutine.
	go simulateActivities(ctx, time.Now(), custArrivalMinSecs, custArrivalMaxSecs, &lock)

	// Periodically log the status of tables and customers dining every 3 seconds.
	for {
		var availableTableNumbers []uint8
		var logs []string

		lock.Lock()
		tables := *reservation.Tables()

		// Iterate through the tables to find available ones and log their reservation counts.
		for _, t := range tables {
			if t.IsAvailable {
				availableTableNumbers = append(availableTableNumbers, t.Number)
			}

			logs = append(logs, fmt.Sprintf("reservation count of table %d: %d", t.Number, t.NumReservations))
		}
		lock.Unlock()

		for _, logText := range logs {
			log.Println(logText)
		}

		// Log the number of customers currently dining and the number of available tables.
		numAvailableTables := reservation.NumAvailableTables()
		numCustomersDining := numTables - numAvailableTables
		log.Printf("# of customers dining: %v", numCustomersDining)
		log.Printf("# of available tables: %v", numAvailableTables)

		// Sort the available tables by their number and print them.
		slices.Sort(availableTableNumbers)
		log.Printf("available table numbers: %v", availableTableNumbers)

		if isOpen.Load() == false && numCustomersDining == 0 {
			log.Println("resto is closed now, bye!")
			break
		}

		time.Sleep(3 * time.Second)
	}
}
