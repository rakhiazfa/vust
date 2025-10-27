package vust

import (
	"time"

	"github.com/sirupsen/logrus"
)

var (
	startTimeKey = "startTime"
)

type JobListener interface {
	Before(ctx JobContext, log *logrus.Logger)
	After(ctx JobContext, log *logrus.Logger)
}

type DefaultJobListener struct{}

func NewDefaultJobListener() JobListener {
	return &DefaultJobListener{}
}

func (l *DefaultJobListener) Before(ctx JobContext, log *logrus.Logger) {
	startTime := time.Now()
	ctx.Set(startTimeKey, startTime)

	log.Info("Job started")
}

func (l *DefaultJobListener) After(ctx JobContext, log *logrus.Logger) {
	startTime := ctx.GetTime(startTimeKey)

	log.Infof("Job finished in %s", time.Since(startTime))
}
