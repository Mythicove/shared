package ratelimit

import (
	"context"
	"sync"
	"time"
)

// Memory 是进程内固定窗口限流实现（低档）。
// 过期条目在 Allow 调用时惰性清扫，避免无界增长。
type Memory struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	sweepAt time.Time
}

type bucket struct {
	count int
	reset time.Time
}

var _ Limiter = (*Memory)(nil)

// sweepInterval 是两次清扫之间的最小间隔。
const sweepInterval = time.Minute

// NewMemory 构造进程内限流器。
func NewMemory() *Memory {
	return &Memory{buckets: map[string]*bucket{}}
}

// Allow 判断 key 在 window 内的请求数是否超过 limit。
func (m *Memory) Allow(_ context.Context, key string, limit int, window time.Duration) (bool, error) {
	now := time.Now()
	m.mu.Lock()
	defer m.mu.Unlock()
	if now.After(m.sweepAt) {
		m.sweepAt = now.Add(sweepInterval)
		for k, b := range m.buckets {
			if now.After(b.reset) {
				delete(m.buckets, k)
			}
		}
	}
	b, ok := m.buckets[key]
	if !ok || now.After(b.reset) {
		b = &bucket{reset: now.Add(window)}
		m.buckets[key] = b
	}
	b.count++
	return b.count <= limit, nil
}
