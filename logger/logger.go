package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
	"time"
)

func init() {
	rotateOptions := []rotatelogs.Option{
		rotatelogs.WithRotationTime(time.Hour * 24),
	}
	rotateOptions = append(rotateOptions, rotatelogs.WithMaxAge(time.Hour*24*7))
	w, err := rotatelogs.New(path.Join("./data/logs", "%Y-%m-%d.log"), rotateOptions...)
	if err != nil {
		log.Errorf("rotatelogs init err: %v", err)
		panic(err)
	}
	log.SetReportCaller(true)
	consoleFormatter := LogFormat{EnableColor: true, CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
		return frame.Function, path.Base(frame.File)
	}}
	fileFormatter := LogFormat{EnableColor: false, CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
		return frame.Function, path.Base(frame.File)
	}}
	var logLevel string
	if os.Getenv("DEBUG") == "true" {
		log.SetLevel(log.DebugLevel)
	} else {
		logLevel = "info"
	}
	log.AddHook(NewLocalHook(w, consoleFormatter, fileFormatter, GetLogLevel(logLevel)...))
}
