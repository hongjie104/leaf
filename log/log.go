package log

import (
	"fmt"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger Logger
var Logger *zap.SugaredLogger

// New New
func New() *zap.SugaredLogger {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
}

func getLogWriter() zapcore.WriteSyncer {
	now := time.Now()
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fmt.Sprintf("./logs/%04d-%02d-%02d- %02d:%02d:%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()),
		MaxSize:    10, // 在进行切割之前，日志文件的最大大小 以MB为单位
		MaxBackups: 5,  // 保留旧文件的最大个数
		MaxAge:     30, // 保留旧文件的最大天数
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// Debug Debug
func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

// Debugf Debugf
func Debugf(template string, args ...interface{}) {
	Logger.Debugf(template, args...)
}

// Info Info
func Info(args ...interface{}) {
	Logger.Info(args...)
}

// Infof Infof
func Infof(template string, args ...interface{}) {
	Logger.Infof(template, args...)
}

// Warn Warn
func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

// Warnf Warnf
func Warnf(template string, args ...interface{}) {
	Logger.Warnf(template, args...)
}

// Error Error
func Error(args ...interface{}) {
	Logger.Error(args...)
}

// Errorf Errorf
func Errorf(template string, args ...interface{}) {
	Logger.Errorf(template, args...)
}

// DPanic DPanic
func DPanic(args ...interface{}) {
	Logger.DPanic(args...)
}

// DPanicf DPanicf
func DPanicf(template string, args ...interface{}) {
	Logger.DPanicf(template, args...)
}

// Panic Panic
func Panic(args ...interface{}) {
	Logger.Panic(args...)
}

// Panicf Panicf
func Panicf(template string, args ...interface{}) {
	Logger.Panicf(template, args...)
}

// Fatal Fatal
func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

// Fatalf Fatalf
func Fatalf(template string, args ...interface{}) {
	Logger.Fatalf(template, args...)
}
