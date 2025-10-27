package vust

import (
	"github.com/sirupsen/logrus"
)

type ErrorListener interface {
	OnRead(ctx StepContext, log *logrus.Logger, err error)
	OnProcess(ctx StepContext, log *logrus.Logger, item any, err error)
	OnWrite(ctx StepContext, log *logrus.Logger, batch *Batch, err error)
}

type DefaultErrorListener struct{}

func NewDefaultErrorListener() ErrorListener {
	return &DefaultErrorListener{}
}

func (l *DefaultErrorListener) OnRead(ctx StepContext, log *logrus.Logger, err error) {
	log.Error(err)
}

func (l *DefaultErrorListener) OnProcess(ctx StepContext, log *logrus.Logger, item any, err error) {
	log.Errorf("Failed to process item\nItem\t: %+v\nError\t: %s\n", item, err.Error())
}

func (l *DefaultErrorListener) OnWrite(ctx StepContext, log *logrus.Logger, batch *Batch, err error) {
	log.Errorf("Failed to write batch\nBatchID\t: %d\nError\t: %s\n", batch.ID, err.Error())
}
