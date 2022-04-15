package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var SugarLogger *zap.SugaredLogger

func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	logg := zap.New(core)
	SugarLogger = logg.Sugar()
}

func getEncoder() zapcore.Encoder {
	encodeConfig := zap.NewProductionEncoderConfig()
	encodeConfig.LevelKey = "level"
	encodeConfig.MessageKey = "msg"
	encodeConfig.TimeKey = "time"
	encodeConfig.NameKey = "name"
	encodeConfig.CallerKey = "caller"
	encodeConfig.FunctionKey = "func"

	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encodeConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encodeConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./log/jsonrpc.log",
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     10,
		Compress:   true,
	}
	return zapcore.AddSync(lumberJackLogger)
}
