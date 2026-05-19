# Go Resto Simulator

A concurrent restaurant reservation simulator written in Go. This project simulates a restaurant environment where customers arrive randomly, attempt to make reservations, and dine at available tables.

## Features

- **Concurrent Simulation**: Uses Go goroutines and synchronization primitives to handle multiple concurrent customers
- **Table Management**: Manages a pool of restaurant tables with reservation tracking
- **Random Arrivals**: Simulates realistic customer arrival patterns with randomized intervals
- **Dynamic Dining Duration**: Each customer has a random dining duration between configured limits
- **Operating Hours**: Simulates realistic restaurant operating hours with graceful shutdown
- **Context-based Control**: Uses Go context for coordinated shutdown across goroutines
- **Status Monitoring**: Periodically logs restaurant status including available tables and dining customers

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
- `servingDurationMins`: How many minutes the restaurant operates (default: 11)

## Simulated Time

The simulation runs in real time. Durations configured in the code correspond directly to actual time elapsed:
- Customer arrival intervals are in seconds
- Dining durations are in seconds  
- Restaurant operating hours are in minutes

For example, the default `servingDurationMins = 11` means the restaurant operates for 11 actual minutes before gracefully closing.

## How It Works

1. The simulation initializes a fixed number of tables and sets the restaurant status to "open"
2. Customers arrive at random intervals defined by the arrival time constants
3. Each arriving customer attempts to reserve an available table
4. If a table is available, the reservation succeeds, the customer dines, and the event is logged
5. If all tables are reserved, the customer leaves without dining
6. The restaurant operates for the configured `servingDurationMins` duration
7. When operating hours expire, the restaurant stops accepting new customers and gracefully closes after all current diners finish
8. The simulation logs periodic status updates including available tables and current reservations

## Requirements

- Go 1.26 or later

## Concurrency Model

The simulator uses Go's concurrency primitives:
- **goroutines**: For handling concurrent customer arrivals and dining activities
- **sync.RWMutex**: For concurrency-safe table access and modifications
- **sync/atomic**: For atomic operations on table availability, customer numbering, and restaurant status
- **context.Context**: For coordinated graceful shutdown of the simulation when operating hours expire
