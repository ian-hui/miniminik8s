package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	// 1. init logger
	K8sLogger.Info("init logger success")
	K8sLogger.Sync()
}
