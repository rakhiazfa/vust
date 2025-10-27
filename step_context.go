package vust

import (
	"context"
	"maps"
	"sync"
	"time"
)

type StepContext interface {
	context.Context

	JobContext() JobContext
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

	WithCancel() (StepContext, context.CancelFunc)
	WithTimeout(d time.Duration) (StepContext, context.CancelFunc)
	WithDeadline(t time.Time) (StepContext, context.CancelFunc)
}

type stepContext struct {
	ctx        context.Context
	jobContext JobContext
	values     map[string]any
	mu         sync.RWMutex
}

func NewStepContext(jobContext JobContext) StepContext {
	sc := &stepContext{
		ctx:        context.Background(),
		jobContext: NewJobContext(context.Background()),
		values:     make(map[string]any),
	}

	if jobContext != nil {
		sc.jobContext = jobContext
	}

	return sc
}

func (c *stepContext) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *stepContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *stepContext) Err() error {
	return c.ctx.Err()
}

func (c *stepContext) Value(key any) any {
	if s, ok := key.(string); ok {
		if v, ok := c.Get(s); ok {
			return v
		}
	}

	return c.ctx.Value(key)
}

func (c *stepContext) JobContext() JobContext {
	return c.jobContext
}

func (c *stepContext) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.values))
	for k := range c.values {
		keys = append(keys, k)
	}

	return keys
}

func (c *stepContext) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.values[key]

	return v, ok
}

func (c *stepContext) GetString(key string, def ...string) string {
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

func (c *stepContext) GetInt(key string, def ...int) int {
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

func (c *stepContext) GetFloat(key string, def ...float64) float64 {
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

func (c *stepContext) GetBool(key string, def ...bool) bool {
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

func (c *stepContext) GetTime(key string, def ...time.Time) time.Time {
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

func (c *stepContext) GetDuration(key string, def ...time.Duration) time.Duration {
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

func (c *stepContext) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.values[key] = value
}

func (c *stepContext) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.values, key)
}

func (c *stepContext) WithCancel() (StepContext, context.CancelFunc) {
	ctx, cancel := context.WithCancel(c.ctx)
	return c.wrap(ctx), cancel
}

func (c *stepContext) WithTimeout(d time.Duration) (StepContext, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(c.ctx, d)
	return c.wrap(ctx), cancel
}

func (c *stepContext) WithDeadline(t time.Time) (StepContext, context.CancelFunc) {
	ctx, cancel := context.WithDeadline(c.ctx, t)
	return c.wrap(ctx), cancel
}

func (c *stepContext) wrap(ctx context.Context) StepContext {
	sc := &stepContext{
		ctx:    ctx,
		values: make(map[string]any),
	}

	c.mu.RLock()
	maps.Copy(sc.values, c.values)
	c.mu.RUnlock()

	return sc
}
