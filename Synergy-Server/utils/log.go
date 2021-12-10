package utils

import (
	"ARTeamViewService/conf"
	"ARTeamViewService/consts"
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func newLfsHook(logDir string, logName string, maxRemainCnt uint) logrus.Hook {
	writer, err := rotatelogs.New(
		logDir+"/%Y%m%d%H_"+logName,
		//// WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
		rotatelogs.WithLinkName(logName),
		///// WithRotationTime设置日志分割的时间，这里设置为一小时分割一次
		rotatelogs.WithRotationTime(time.Hour*1),

		// WithMaxAge和WithRotationCount二者只能设置一个，
		// WithMaxAge设置文件清理前的最长保存时间，
		// WithRotationCount设置文件清理前最多保存的个数。
		rotatelogs.WithMaxAge(-1),
		rotatelogs.WithRotationCount(maxRemainCnt),
		)
	if err != nil {
		_ = fmt.Errorf("config local file system for logger error :%v", err)
	}

	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		//logrus.TraceLevel: writer,
		logrus.DebugLevel: writer,
		logrus.InfoLevel: writer,
		logrus.WarnLevel: writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.TextFormatter{DisableColors:false})
	//}, &logrus.JSONFormatter{})

	return lfsHook
}

func GetLogger(confLog conf.Logger) (log *logrus.Logger) {
	log = logrus.New()
	logExpire := confLog.LogExpire
	if logExpire == 0 {
		logExpire = consts.DftLogExpire
	}
	log.AddHook(newLfsHook(confLog.LogDir, confLog.LogName, uint(logExpire*24)))

	log.SetOutput(os.Stdout)

	switch confLog.LogLevel {
	case 0:
		log.SetLevel(logrus.PanicLevel)
	case 1:
		log.SetLevel(logrus.FatalLevel)
	case 2:
		log.SetLevel(logrus.ErrorLevel)
	case 3:
		log.SetLevel(logrus.WarnLevel)
	case 4:
		log.SetLevel(logrus.InfoLevel)
	case 5:
		log.SetLevel(logrus.DebugLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}
	return log
}