package main

import (
	"fmt"
	"go-resto-simulator/internal/customer"
	"go-resto-simulator/internal/reservation"
	"log"
	"math/rand/v2"
	"slices"
	"sync"
	"time"
)

const (
	custArrivalMinSecs    = 0
	custArrivalMaxSecs    = 20
	custMinDiningDuration = 60
	custMaxDiningDuration = 90
	numTables             = 10
)

var lock sync.RWMutex

func simulateActivities(startTime time.Time, custArrivalMinSecs, custArrivalMaxSecs int, lock *sync.RWMutex) {
	for {
		sec := custArrivalMinSecs + rand.IntN(custArrivalMaxSecs-custArrivalMinSecs)
		time.Sleep(time.Duration(sec) * time.Second)
		c, err := customer.New()
		if err != nil {
			panic(err)
		}

		log.Printf("customer %d arrived on minute %d", c.Number, int(time.Since(startTime).Seconds()))
		r, err := reservation.Reserve(c, lock)
		if err == reservation.ErrAllTablesReserved {
			log.Printf("customer %d left because they failed to reserve: %v", c.Number, err)
			continue
		}

		if err != nil {
			log.Printf("unknown error: %v", err)
			return
		}

		log.Printf("customer %d has just reserved table %d", c.Number, r.Table.Number)
	}
}

func main() {
	reservation.SetTables(numTables)
	customer.SetMinMaxDiningDurations(custMinDiningDuration, custMaxDiningDuration)

	// Simulate the arrival of customers and their dining activities in a separate goroutine.
	go simulateActivities(time.Now(), custArrivalMinSecs, custArrivalMaxSecs, &lock)

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

			var logText string
			switch t.NumReservations {
			case 0:
				logText = fmt.Sprintf("table %d has never been reserved", t.Number)
			case 1:
				logText = fmt.Sprintf("table %d has been reserved once", t.Number)
			default:
				logText = fmt.Sprintf("table %d has been reserved for %d times", t.Number, t.NumReservations)
			}

			logs = append(logs, logText)
		}
		lock.Unlock()

		for _, logText := range logs {
			log.Println(logText)
		}

		// Log the number of customers currently dining and the number of available tables.
		numAvailableTables := reservation.NumAvailableTables()
		log.Printf("# of customers dining: %v", numTables-numAvailableTables)
		log.Printf("# of available tables: %v", numAvailableTables)

		// Sort the available tables by their number and print them.
		slices.Sort(availableTableNumbers)
		log.Printf("available table numbers: %v", availableTableNumbers)

		time.Sleep(3 * time.Second)
	}
}
