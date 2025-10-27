package vust

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type StepPrefixHook struct {
	Name string
}

func (h *StepPrefixHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *StepPrefixHook) Fire(entry *logrus.Entry) error {
	entry.Message = fmt.Sprintf("[%s] - %s", h.Name, entry.Message)
	return nil
}
