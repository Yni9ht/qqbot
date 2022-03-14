package log

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	DEBUG = "Debug"
	INFO  = "Info"
	WARN  = "Warning"
	ERROR = "Error"
)

func Debug(v ...interface{}) {
	printLog(DEBUG, v...)
}

func Info(v ...interface{}) {
	printLog(INFO, v...)
}

func Warn(v ...interface{}) {
	printLog(WARN, v...)
}

func Error(v ...interface{}) {
	printLog(ERROR, v...)
}

func Debugf(formatter string, v ...interface{}) {
	Debug(fmt.Sprintf(formatter, v...))
}

func Infof(formatter string, v ...interface{}) {
	Info(fmt.Sprintf(formatter, v...))
}

func Warnf(formatter string, v ...interface{}) {
	Warn(fmt.Sprintf(formatter, v...))
}

func Errorf(formatter string, v ...interface{}) {
	Error(fmt.Sprintf(formatter, v...))
}

func printLog(level string, v ...interface{}) {
	pc, file, line, _ := runtime.Caller(3)
	file = filepath.Base(file)
	funcName := strings.TrimPrefix(filepath.Ext(runtime.FuncForPC(pc).Name()), ".")

	logFormat := "[%s] %s %s:%d:%s " + fmt.Sprint(v...) + "\n"
	date := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf(logFormat, level, date, file, line, funcName)
}
