package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	ccc "github.com/heqzha/goutils/concurrency"
	"github.com/Sirupsen/logrus"
)

var (
	nowDate    string
	rotateLock sync.RWMutex

	logger *Logger
)

const (
	LOG_LEVEL_DEBUG = logrus.DebugLevel
	LOG_LEVEL_INFO  = logrus.InfoLevel
	LOG_LEVEL_WARN  = logrus.WarnLevel
	LOG_LEVEL_ERROR = logrus.ErrorLevel
	LOG_LEVEL_FATAL = logrus.FatalLevel
	LOG_LEVEL_PANIC = logrus.PanicLevel
)

func LogConfig(path string, level logrus.Level) {
	nowDate = time.Now().Format("2006-01-02")
	logger = newLogger(path, level)
	ccc.TaskRunPeriodic(rotateLog, "LoggerConfig", 5 * time.Second)
}

func LogDebug(typ, msg string) error {
	return LogDebugM(map[string]interface{}{
		"type": typ,
		"msg":  msg,
	})
}

func LogDebugM(m map[string]interface{}) error {
	if logger == nil {
		return fmt.Errorf("Configure logger before using it")
	}
	logger.Debug(m)
	return nil
}

func LogInfo(typ, msg string) error {
	return LogInfoM(map[string]interface{}{
		"type": typ,
		"msg":  msg,
	})
}

func LogInfoM(m map[string]interface{}) error {
	if logger == nil {
		return fmt.Errorf("Configure logger before using it")
	}
	logger.Info(m)
	return nil
}

func LogWarn(typ, msg string) error {
	return LogWarnM(map[string]interface{}{
		"type": typ,
		"msg":  msg,
	})
}

func LogWarnM(m map[string]interface{}) error {
	if logger == nil {
		return fmt.Errorf("Configure logger before using it")
	}
	logger.Warn(m)
	return nil
}

func LogError(typ, msg string) error {
	return LogErrorM(map[string]interface{}{
		"type": typ,
		"msg":  msg,
	})
}

func LogErrorM(m map[string]interface{}) error {
	if logger == nil {
		return fmt.Errorf("Configure logger before using it")
	}
	logger.Error(m)
	return nil
}

func LogFatal(typ, msg string) error {
	return LogFatalM(map[string]interface{}{
		"type": typ,
		"msg":  msg,
	})
}

func LogFatalM(m map[string]interface{}) error {
	if logger == nil {
		return fmt.Errorf("Configure logger before using it")
	}
	logger.Fatal(m)
	return nil
}

func LogPanic(typ, msg string) error {
	return LogPanicM(map[string]interface{}{
		"type": typ,
		"msg":  msg,
	})
}

func LogPanicM(m map[string]interface{}) error {
	if logger == nil {
		return fmt.Errorf("Configure logger before using it")
	}
	logger.Panic(m)
	return nil
}

type Logger struct {
	*logrus.Logger
	path string
	base string
	file *os.File
}

func (logger *Logger) init() {
	err := os.MkdirAll(logger.path, 0777)
	if err != nil {
		log.Fatalf("Failed to init log path %s: %s", logger.path, err.Error())
		return
	}

	p := filepath.Join(logger.path, fmt.Sprintf("%s.log.%s", logger.base, nowDate))
	f, err := os.OpenFile(p, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
	if err != nil {
		log.Printf("Failed to init [%s] log: %s", logger.base, err.Error())
		os.Exit(1)
	}

	logger.file = f
	logger.Out = f
	logger.Formatter = new(logrus.JSONFormatter)
}

func (logger *Logger) rotate() {
	logger.file.Close()
	p := filepath.Join(logger.path, fmt.Sprintf("%s.log.%s", logger.base, nowDate))
	f, err := os.OpenFile(p, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
	if err != nil {
		logger.Fatalf("Failed to init [%s] log: %s", logger.base, err.Error())
	}

	logger.Out = f
	logger.Formatter = new(logrus.JSONFormatter)
}

func (logger *Logger) newLog(level logrus.Level, data map[string]interface{}) {
	rotateLock.RLock()
	defer rotateLock.RUnlock()

	fields := logrus.Fields(data)
	switch level {
	case LOG_LEVEL_DEBUG:
		logger.WithFields(fields).Debug("")
	case LOG_LEVEL_INFO:
		logger.WithFields(fields).Info("")
	case LOG_LEVEL_WARN:
		logger.WithFields(fields).Warn("")
	case LOG_LEVEL_ERROR:
		logger.WithFields(fields).Error("")
	case LOG_LEVEL_FATAL:
		logger.WithFields(fields).Fatal("")
	case LOG_LEVEL_PANIC:
		logger.WithFields(fields).Panic("")
	default:
		logger.WithFields(fields).Info("")
	}
}

func (logger *Logger) Debug(data map[string]interface{}) {
	logger.newLog(LOG_LEVEL_DEBUG, data)
}

func (logger *Logger) Warn(data map[string]interface{}) {
	logger.newLog(LOG_LEVEL_WARN, data)
}

func (logger *Logger) Info(data map[string]interface{}) {
	logger.newLog(LOG_LEVEL_INFO, data)
}

func (logger *Logger) Error(data map[string]interface{}) {
	logger.newLog(LOG_LEVEL_ERROR, data)
}

func (logger *Logger) Fatal(data map[string]interface{}) {
	logger.newLog(LOG_LEVEL_FATAL, data)
}

func (logger *Logger) Panic(data map[string]interface{}) {
	logger.newLog(LOG_LEVEL_PANIC, data)
}

func (logger *Logger) SetLevel(level logrus.Level) {
	logger.Level = level
}

func newLogger(path string, level logrus.Level) *Logger {
	var base string
	switch level {
	case LOG_LEVEL_DEBUG:
		base = "DEBUG"
	case LOG_LEVEL_INFO:
		base = "INFO"
	case LOG_LEVEL_WARN:
		base = "WARN"
	case LOG_LEVEL_ERROR:
		base = "ERROR"
	case LOG_LEVEL_FATAL:
		base = "FATAL"
	case LOG_LEVEL_PANIC:
		base = "PANIC"
	default:
		base = "UNKNOW"
	}
	l := &Logger{new(logrus.Logger), path, base, nil}
	l.init()
	l.SetLevel(level)

	return l
}

func rotateLog() time.Duration {
	_nowDate := time.Now().Format("2006-01-02")
	if _nowDate == nowDate {
		now := time.Now()
		return time.Duration(23-now.Hour())*time.Hour + time.Duration(59-now.Minute())*time.Minute + time.Duration(60-now.Second())*time.Second
	}

	rotateLock.Lock()
	defer rotateLock.Unlock()

	nowDate = _nowDate
	logger.rotate()
	now := time.Now()
	return time.Duration(23-now.Hour())*time.Hour + time.Duration(59-now.Minute())*time.Minute + time.Duration(60-now.Second())*time.Second
}
