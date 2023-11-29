package main

import "minik8s/logger"

var K8sLogger = logger.K8sLogger

func main() {
	// 1. init logger
	K8sLogger.Info("init logger success")
	K8sLogger.Sync()

}
