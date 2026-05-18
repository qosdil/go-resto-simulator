# Go Resto Simulator

A concurrent restaurant reservation simulator written in Go. This project simulates a restaurant environment where customers arrive randomly, attempt to make reservations, and dine at available tables.

## Features

- **Concurrent Simulation**: Uses Go goroutines and synchronization primitives to handle multiple concurrent customers
- **Table Management**: Manages a pool of restaurant tables with reservation tracking
- **Random Arrivals**: Simulates realistic customer arrival patterns with randomized intervals
- **Dynamic Dining Duration**: Each customer has a random dining duration between configured limits

## Running

To run the simulator:

```bash
go run main.go
```

The simulator will start and continuously log customer arrivals, table reservations, and any reservation failures.

To check for possible race conditions, you can run the simulator with the race detector:

```bash
go run -race main.go
```

## Configuration

You can modify the following constants in `main.go` to customize the simulation:

- `custArrivalMinSecs`: Minimum seconds between customer arrivals (default: 0)
- `custArrivalMaxSecs`: Maximum seconds between customer arrivals (default: 20)
- `custMinDiningDuration`: Minimum dining duration in seconds (default: 60)
- `custMaxDiningDuration`: Maximum dining duration in seconds (default: 90)
- `numTables`: Number of tables in the restaurant (default: 10)

## How It Works

1. The simulation initializes a fixed number of tables
2. Customers arrive at random intervals defined by the arrival time constants
3. Each arriving customer attempts to reserve an available table
4. If a table is available, the reservation succeeds and is logged
5. If all tables are reserved, the customer leaves without dining
6. The simulation continues indefinitely, logging all events

## Requirements

- Go 1.26 or later

## Concurrency Model

The simulator uses Go's concurrency primitives:
- **goroutines**: For handling concurrent customer arrivals
- **sync.RWMutex**: For thread-safe table access
- **sync/atomic**: For atomic operations on table availability and customer numbering
