package vust

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type JobPrefixHook struct {
	Name string
}

func (h *JobPrefixHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *JobPrefixHook) Fire(entry *logrus.Entry) error {
	entry.Message = fmt.Sprintf("[%s] - %s", h.Name, entry.Message)
	return nil
}
