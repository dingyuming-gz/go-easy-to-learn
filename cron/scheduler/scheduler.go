package scheduler

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/example/api"
	"github.com/example/database"
	"github.com/example/logger"
	"github.com/go-co-op/gocron"
)

// Scheduler 定义了调度器结构体。
type Scheduler struct {
	ctx     context.Context
	config  *SchedulerConfig
	db      *database.Database
	api     *api.API
	logger  *logger.Logger
	counter map[string]int64
	mu      sync.Mutex
}

// NewScheduler 创建一个新的调度器实例。
func NewScheduler(ctx context.Context, config *SchedulerConfig) (*Scheduler, error) {
	err := config.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	logFile, err := os.OpenFile(config.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	logger := logger.NewLogger(logFile)

	db := database.NewDatabase(config.DBDriver, config.DBConnString)

	api := api.NewAPI(config.APIEndpoint)

	return &Scheduler{
		ctx:     ctx,
		config:  config,
		db:      db,
		api:     api,
		logger:  logger,
		counter: make(map[string]int64),
	}, nil
}

// Start 启动调度器并开始执行任务。
func (s *Scheduler) Start() error {
	cron := gocron.NewScheduler(time.UTC)
	cron.TagsUnique()

	task := cron.Every(s.config.TaskTimeout).Seconds().Tag("task")

	for i := 1; i <= s.config.TaskIterations; i++ {
		n := i
		task.Do(func() {
			s.runTask(n)
		})
	}

	cron.StartBlocking()

	return nil
}

// runTask 执行任务操作。
func (s *Scheduler) runTask(iteration int) {
	backoff := s.config.BackoffInitialInterval
	for i := 0; i < iteration; i++ {
		if i != 0 {
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * s.config.BackoffMultiplier)
			if backoff > s.config.MaxBackoffTime {
				backoff = s.config.MaxBackoffTime
			}
		}

		err := s.db.Connect()
		if err != nil {
			s.logger.Logf("failed to connect to database: %v", err)
			continue
		}

		body, err := s.api.Call(s.ctx)
		if err != nil {
			s.logger.Logf("failed to call API: %v", err)
			continue
		}

		err = s.saveData(body)
		if err != nil {
			s.logger.Logf("failed to save data to database: %v", err)
			continue
		}

		s.mu.Lock()
		s.counter[string(body)]++
		s.mu.Unlock()

		s.logger.Logf("saved data #%d: %s", i+1, body)
	}

	err := s.db.Close()
	if err != nil {
		s.logger.Logf("failed to close database connection: %v", err)
	}
}

// saveData 将数据保存到数据库中。
func (s *Scheduler) saveData(data []byte) error {
	return s.db.Save(data)
}

// GetCounters 获取 API 调用计数器的副本。
func (s *Scheduler) GetCounters() map[string]int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	counters := make(map[string]int64)
	for k, v := range s.counter {
		counters[k] = v
	}
	return counters
}

var (
	ErrDatabaseNotConnected = errors.New("database not connected")
	ErrInvalidStatusCode    = errors.New("invalid status code")
)
