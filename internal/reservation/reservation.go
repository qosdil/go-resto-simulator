package reservation

import (
	"errors"
	"go-resto-simulator/internal/customer"
	"log"
	"sync"
	"sync/atomic"
)

var (
	areTablesSet         atomic.Bool
	ErrAllTablesReserved = errors.New("all tables are reserved")
	numAvailableTables   atomic.Uint32
	tables               *[]table
)

type table struct {
	IsAvailable     bool
	Number          uint8
	NumReservations uint8
}

type Reservation struct {
	Customer *customer.Customer
	Table    *table
}

func NumAvailableTables() uint32 {
	return numAvailableTables.Load()
}

// Reserve attempts to reserve a table for the given customer. If all tables are reserved, it returns an error.
// Otherwise, it reserves a table and starts a goroutine to simulate the dining activity of the customer.
func Reserve(cust *customer.Customer, lock *sync.RWMutex) (*Reservation, error) {
	if numAvailableTables.Load() == 0 {
		return nil, ErrAllTablesReserved
	}

	var reservationTable *table

	// Find the first available table and reserve it for the customer. Protect the operation with a write lock to
	// ensure concurrency safety.
	lock.Lock()
	derefTables := *tables
	for i, t := range derefTables {
		if t.IsAvailable == true {
			derefTables[i].IsAvailable = false
			derefTables[i].NumReservations++
			reservationTable = &derefTables[i]
			break
		}
	}
	tables = &derefTables
	lock.Unlock()

	numAvailableTables.Store(numAvailableTables.Load() - 1)
	r := &Reservation{Customer: cust, Table: reservationTable}

	// Simulate the dining activity of the customer in a separate goroutine.
	go func(r *Reservation, numAvailableTables *atomic.Uint32, lock *sync.RWMutex) {
		sec := cust.Dine()

		// After dining, the customer leaves and the table becomes available again.
		numAvailableTables.Add(1)
		lock.Lock()
		r.Table.IsAvailable = true
		lock.Unlock()

		// Seconds as minutes for simulation
		log.Printf("customer %d is done after dining for %d minutes, table %d is now available", cust.Number, sec,
			reservationTable.Number)
	}(r, &numAvailableTables, lock)

	return r, nil
}

// SetTables initializes the tables for the restaurant. It creates a specified number of tables and marks them as available.
func SetTables(numTables uint32) {
	if areTablesSet.Load() {
		return
	}

	i := 1
	newTables := []table{}
	for range numTables {
		newTables = append(newTables, table{IsAvailable: true, Number: uint8(i)})
		i++
	}

	tables = &newTables
	numAvailableTables.Add(numTables)
	areTablesSet.Store(true)
	log.Printf("resto is ready to serve with %d tables", numTables)
}

func Tables() *[]table {
	return tables
}
