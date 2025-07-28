package chain

import (
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.AddHook(&runnerNameHook{})
}

type runnerNameHook struct{}

func (c *runnerNameHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (c *runnerNameHook) Fire(entry *logrus.Entry) error {
	if entry.Context == nil {
		return nil
	}

	if runnerName, ok := entry.Context.Value(runnerNameKey{}).(string); ok {
		entry.Data["runnerName"] = runnerName
	}

	return nil
}
