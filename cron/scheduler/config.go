package scheduler

import (
	"errors"
	"time"
)

// SchedulerConfig 定义了调度器的配置选项。
type SchedulerConfig struct {
	TaskTimeout            time.Duration // 任务超时时间
	TaskIterations         int           // 任务循环次数
	LogFile                string        // 日志输出文件名
	DBDriver               string        // 数据库驱动
	DBConnString           string        // 数据库连接字符串
	APIEndpoint            string        // API 接口地址
	BackoffInitialInterval time.Duration // 退避策略的初始等待时间
	BackoffMultiplier      float64       // 退避策略的乘数
	MaxBackoffTime         time.Duration // 退避策略的最大等待时间
}

// Validate 验证配置选项是否合法。
func (c *SchedulerConfig) Validate() error {
	if c == nil {
		return errors.New("nil configuration")
	} else if c.APIEndpoint == "" {
		return errors.New("API endpoint must not be empty")
	} else if c.BackoffInitialInterval <= 0 {
		return errors.New("backoff initial interval must be positive")
	} else if c.BackoffMultiplier <= 1 {
		return errors.New("backoff multiplier must be greater than 1")
	} else if c.MaxBackoffTime <= c.BackoffInitialInterval {
		return errors.New("max backoff time must be greater than backoff initial interval")
	}
	return nil
}
