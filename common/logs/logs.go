package logs

import (
	"common/config"
	log "github.com/charmbracelet/log"
	"os"
	"time"
)

var logger *log.Logger

func InitLog(appName string) {
	logger = log.New(os.Stderr)
	if config.Conf.Log.Level == "DEBUG" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	logger.SetPrefix(appName)
	logger.SetReportTimestamp(true)
	logger.SetTimeFormat(time.Layout)
}

func Warn(format string, values ...any) {
	if len(values) == 0 {
		logger.Warn(format)
	} else {
		logger.Warnf(format, values...)
	}

}
func Error(format string, values ...any) {
	if len(values) == 0 {
		logger.Error(format)
	} else {
		logger.Errorf(format, values...)
	}
}
func Fatal(format string, values ...any) {
	if len(values) == 0 {
		logger.Fatal(format)
	} else {
		logger.Fatalf(format, values...)
	}
}
func Info(format string, values ...any) {
	if len(values) == 0 {
		logger.Info(format)
	} else {
		logger.Infof(format, values...)
	}

}
