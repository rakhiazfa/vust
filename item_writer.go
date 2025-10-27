package vust

import (
	"github.com/sirupsen/logrus"
)

type ItemWriter interface {
	Write(ctx StepContext, log *logrus.Logger, batch *Batch) error
}
