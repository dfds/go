package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitializeLogger(debug bool, logLevel string) {
	var logConf zap.Config
	if debug {
		logConf = zap.NewDevelopmentConfig()
	} else {
		logConf = zap.NewProductionConfig()
	}

	level, err := zapcore.ParseLevel(logLevel)
	if err != nil {
		fmt.Println(err)
		level = zapcore.InfoLevel
	}

	logConf.Level = zap.NewAtomicLevelAt(level)
	logConf.EncoderConfig.FunctionKey = "func"
	Logger, _ = logConf.Build(zap.AddStacktrace(zapcore.ErrorLevel))
	Logger.Info(fmt.Sprintf("Logging enabled, log level set to %s", Logger.Level().String()))
}
