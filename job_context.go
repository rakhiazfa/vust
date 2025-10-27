package vust

import (
	"context"
	"maps"
	"sync"
	"time"
)

type JobContext interface {
	context.Context

	Keys() []string
	Get(key string) (any, bool)
	GetString(key string, def ...string) string
	GetInt(key string, def ...int) int
	GetFloat(key string, def ...float64) float64
	GetBool(key string, def ...bool) bool
	GetTime(key string, def ...time.Time) time.Time
	GetDuration(key string, def ...time.Duration) time.Duration
	Set(key string, value any)
	Delete(key string)

	WithCancel() (JobContext, context.CancelFunc)
	WithTimeout(d time.Duration) (JobContext, context.CancelFunc)
	WithDeadline(t time.Time) (JobContext, context.CancelFunc)
}

type jobContext struct {
	ctx    context.Context
	values map[string]any
	mu     sync.RWMutex
}

func NewJobContext(ctx context.Context) JobContext {
	jc := &jobContext{
		ctx:    context.Background(),
		values: make(map[string]any),
	}

	if ctx != nil {
		jc.ctx = ctx
	}

	return jc
}

func (c *jobContext) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *jobContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *jobContext) Err() error {
	return c.ctx.Err()
}

func (c *jobContext) Value(key any) any {
	if s, ok := key.(string); ok {
		if v, ok := c.Get(s); ok {
			return v
		}
	}

	return c.ctx.Value(key)
}

func (c *jobContext) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.values))
	for k := range c.values {
		keys = append(keys, k)
	}

	return keys
}

func (c *jobContext) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.values[key]

	return v, ok
}

func (c *jobContext) GetString(key string, def ...string) string {
	if v, ok := c.Get(key); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}

	if len(def) > 0 {
		return def[0]
	}

	return ""
}

func (c *jobContext) GetInt(key string, def ...int) int {
	if v, ok := c.Get(key); ok {
		if i, ok := v.(int); ok {
			return i
		}
	}

	if len(def) > 0 {
		return def[0]
	}

	return 0
}

func (c *jobContext) GetFloat(key string, def ...float64) float64 {
	if v, ok := c.Get(key); ok {
		switch val := v.(type) {
		case float64:
			return val
		case float32:
			return float64(val)
		case int:
			return float64(val)
		case int64:
			return float64(val)
		}
	}

	if len(def) > 0 {
		return def[0]
	}

	return 0
}

func (c *jobContext) GetBool(key string, def ...bool) bool {
	if v, ok := c.Get(key); ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}

	if len(def) > 0 {
		return def[0]
	}

	return false
}

func (c *jobContext) GetTime(key string, def ...time.Time) time.Time {
	if v, ok := c.Get(key); ok {
		if t, ok := v.(time.Time); ok {
			return t
		}
	}

	if len(def) > 0 {
		return def[0]
	}

	return time.Time{}
}

func (c *jobContext) GetDuration(key string, def ...time.Duration) time.Duration {
	if v, ok := c.Get(key); ok {
		if d, ok := v.(time.Duration); ok {
			return d
		}
	}

	if len(def) > 0 {
		return def[0]
	}

	return 0
}

func (c *jobContext) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.values[key] = value
}

func (c *jobContext) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.values, key)
}

func (c *jobContext) WithCancel() (JobContext, context.CancelFunc) {
	ctx, cancel := context.WithCancel(c.ctx)
	return c.wrap(ctx), cancel
}

func (c *jobContext) WithTimeout(d time.Duration) (JobContext, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(c.ctx, d)
	return c.wrap(ctx), cancel
}

func (c *jobContext) WithDeadline(t time.Time) (JobContext, context.CancelFunc) {
	ctx, cancel := context.WithDeadline(c.ctx, t)
	return c.wrap(ctx), cancel
}

func (c *jobContext) wrap(ctx context.Context) JobContext {
	jc := &jobContext{
		ctx:    ctx,
		values: make(map[string]any),
	}

	c.mu.RLock()
	maps.Copy(jc.values, c.values)
	c.mu.RUnlock()

	return jc
}
