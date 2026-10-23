package ratelimit

import (
	"context"
	"testing"
	"time"
)

func TestMemoryAllow(t *testing.T) {
	m := NewMemory()
	ctx := context.Background()
	for i := 0; i < 3; i++ {
		ok, err := m.Allow(ctx, "ip1", 3, time.Minute)
		if err != nil || !ok {
			t.Fatalf("request %d: ok=%v err=%v", i, ok, err)
		}
	}
	if ok, _ := m.Allow(ctx, "ip1", 3, time.Minute); ok {
		t.Fatal("expected limit exceeded")
	}
	if ok, _ := m.Allow(ctx, "ip2", 3, time.Minute); !ok {
		t.Fatal("other key should be allowed")
	}
}

func TestMemoryWindowReset(t *testing.T) {
	m := NewMemory()
	ctx := context.Background()
	if ok, _ := m.Allow(ctx, "k", 1, 10*time.Millisecond); !ok {
		t.Fatal("first should pass")
	}
	if ok, _ := m.Allow(ctx, "k", 1, 10*time.Millisecond); ok {
		t.Fatal("second should be limited")
	}
	time.Sleep(15 * time.Millisecond)
	if ok, _ := m.Allow(ctx, "k", 1, 10*time.Millisecond); !ok {
		t.Fatal("after window should pass")
	}
}

func TestMemorySweep(t *testing.T) {
	m := NewMemory()
	ctx := context.Background()
	_, _ = m.Allow(ctx, "old", 5, time.Nanosecond)
	m.sweepAt = time.Time{}
	time.Sleep(time.Millisecond)
	_, _ = m.Allow(ctx, "new", 5, time.Minute)
	m.mu.Lock()
	n := len(m.buckets)
	m.mu.Unlock()
	if n != 1 {
		t.Fatalf("expected expired bucket swept, got %d buckets", n)
	}
}
