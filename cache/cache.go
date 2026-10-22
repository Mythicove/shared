// Package cache 定义缓存抽象。
// 低档实现为进程内内存，高档实现为 Redis。
package cache

import (
	"context"
	"time"
)

// Cache 是键值缓存的抽象。
type Cache interface {
	// Get 返回键对应的值；不存在时 ok 为 false。
	Get(ctx context.Context, key string) (val []byte, ok bool, err error)
	// Set 写入带 TTL 的键值。
	Set(ctx context.Context, key string, val []byte, ttl time.Duration) error
	// Delete 删除键。
	Delete(ctx context.Context, key string) error
}
