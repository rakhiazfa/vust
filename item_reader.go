package vust

import (
	"github.com/sirupsen/logrus"
)

type ItemReader interface {
	Read(ctx StepContext, log *logrus.Logger) (any, error)
}
