package vust

import (
	"github.com/sirupsen/logrus"
)

type ItemProcessor interface {
	Process(ctx StepContext, log *logrus.Logger, item any) (any, error)
}
