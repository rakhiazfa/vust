package vust

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type ModulePrefixHook struct{}

func (h *ModulePrefixHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *ModulePrefixHook) Fire(entry *logrus.Entry) error {
	entry.Message = fmt.Sprintf("[%s] - %s", "Vust", entry.Message)
	return nil
}
