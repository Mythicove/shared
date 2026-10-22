// Package queue 定义任务队列抽象。
// 低档实现为内联同步执行，高档实现为 asynq（Redis）。
package queue

import (
	"context"
	"encoding/json"
)

// Task 是一个异步任务。
type Task struct {
	Type    string
	Payload json.RawMessage
}

// Handler 处理一类任务。
type Handler func(ctx context.Context, t Task) error

// Queue 是任务投递与注册的抽象。
type Queue interface {
	// Enqueue 投递任务。
	Enqueue(ctx context.Context, t Task) error
	// Register 为某类任务注册处理器。
	Register(taskType string, h Handler)
}
