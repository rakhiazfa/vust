package vust

import (
	"github.com/sirupsen/logrus"
)

type ExampleItemWriter struct{}

func NewExampleItemWriter() ItemWriter {
	return &ExampleItemWriter{}
}

func (w *ExampleItemWriter) Write(ctx StepContext, log *logrus.Logger, batch *Batch) error {
	return nil
}
