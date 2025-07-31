package logger

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
	"time"
)

func TestLogrusOptimize(t *testing.T) {
	originalOutput := logrus.StandardLogger().Out
	originalLevel := logrus.StandardLogger().Level
	originalFormatter := logrus.StandardLogger().Formatter
	originalHooks := logrus.StandardLogger().Hooks
	defer func() {
		logrus.SetOutput(originalOutput)
		logrus.SetLevel(originalLevel)
		logrus.SetFormatter(originalFormatter)
		logrus.StandardLogger().ReplaceHooks(originalHooks) // 훅 복원
	}()

	testDuration := 1 * time.Millisecond
	logrus.SetOutput(io.Discard)
	logrus.SetLevel(logrus.InfoLevel)
	logrus.AddHook(&waitHook{testDuration * 2})

	t.Run("Debugf not evaluates", func(t *testing.T) {
		done := make(chan struct{})
		go func() {
			Debugf(context.Background(), "test %d", 1)
			close(done)
		}()

		select {
		case <-done:
			// hook evaluation skipped
		case <-time.After(testDuration):
			assert.Fail(t, "should not take long time")
		}
	})

	t.Run("Errorf evaluates", func(t *testing.T) {
		done := make(chan struct{})
		go func() {
			Errorf(context.Background(), "test %d", 1)
			close(done)
		}()

		select {
		case <-done:
			assert.Fail(t, "should take longer time")
		case <-time.After(testDuration):
			// evaluating hook. takes longer time
		}
	})

	t.Run("Error evaluates", func(t *testing.T) {
		done := make(chan struct{})
		go func() {
			Error(context.Background(), errors.New("test"))
			close(done)
		}()

		select {
		case <-done:
			assert.Fail(t, "should take longer time")
		case <-time.After(testDuration):
			// evaluating hook. takes longer time
		}
	})
}

type waitHook struct{ duration time.Duration }

func (c *waitHook) Levels() []logrus.Level { return logrus.AllLevels }
func (c *waitHook) Fire(_ *logrus.Entry) error {
	time.Sleep(c.duration)
	return nil
}
