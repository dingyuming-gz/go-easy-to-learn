package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/example/scheduler/scheduler"
)

var (
	apiEndpoint            = flag.String("api-endpoint", "", "API endpoint")
	dbDriver               = flag.String("db-driver", "", "database driver name")
	dbConnString           = flag.String("db-conn-string", "", "database connection string")
	logFile                = flag.String("log-file", "", "log file name")
	taskTimeout            = flag.Duration("task-timeout", 0, "task timeout")
	taskIterations         = flag.Int("task-iterations", 0, "task iterations")
	backoffInitialInterval = flag.Duration("backoff-initial-interval", 0, "initial interval of backoff strategy")
	backoffMultiplier      = flag.Float64("backoff-multiplier", 0, "multiplier of backoff strategy")
	maxBackoffTime         = flag.Duration("max-backoff-time", 0, "maximum interval of backoff strategy")
)

func main() {
	flag.Parse()

	config := &scheduler.SchedulerConfig{
		TaskTimeout:            *taskTimeout,
		TaskIterations:         *taskIterations,
		LogFile:                *logFile,
		DBDriver:               *dbDriver,
		DBConnString:           *dbConnString,
		APIEndpoint:            *apiEndpoint,
		BackoffInitialInterval: *backoffInitialInterval,
		BackoffMultiplier:      *backoffMultiplier,
		MaxBackoffTime:         *maxBackoffTime,
	}

	ctx := context.Background()

	s, err := scheduler.NewScheduler(ctx, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create scheduler: %v\n", err)
		os.Exit(1)
	}

	err = s.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start scheduler: %v\n", err)
		os.Exit(1)
	}

	counters := s.GetCounters()
	for k, v := range counters {
		fmt.Printf("%s: %d\n", k, v)
	}
}
