// Package ratelimit 定义限流抽象。
// 低档实现为进程内计数，高档实现为 Redis（令牌桶/固定窗口由实现决定）。
package ratelimit

import (
	"context"
	"time"
)

// Limiter 判定一次请求是否放行。
type Limiter interface {
	// Allow 判断 key 在 window 内的请求数是否超过 limit。
	Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error)
}
