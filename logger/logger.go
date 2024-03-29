package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var K8sLogger *zap.SugaredLogger

func init() {
	//初始化日志
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewJSONEncoder(config)
	defaultLogLevel := zapcore.DebugLevel // 设置 loglevel，debug表示所有日志都输出，info表示只输出info以上的日志

	logFile, _ := os.OpenFile("./logFile/minik8s_log.json", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 06666)
	// or os.Create()
	writer := zapcore.AddSync(logFile)

	logger := zap.New(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	K8sLogger = logger.Sugar()

}
