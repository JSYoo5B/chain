package logger

import (
	"context"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.AddHook(&runnerNameHook{})
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	logrus.WithContext(ctx).Debugf(format, args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	logrus.WithContext(ctx).Errorf(format, args...)
}

func Error(ctx context.Context, err error) {
	logrus.WithContext(ctx).Error(err)
}

type runnerNameHook struct{}

func (c *runnerNameHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (c *runnerNameHook) Fire(entry *logrus.Entry) error {
	if entry.Context == nil {
		return nil
	}

	if runnerName, ok := RunnerNameFromContext(entry.Context); ok {
		entry.Data["runnerName"] = runnerName
	}

	return nil
}
