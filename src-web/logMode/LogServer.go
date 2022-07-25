package logMode

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
)

// logServerMap TODO F: 稍晚将日志部分抽离出来，并日志中应能返回对应的函数名
var logServerMap = sync.Map{}

func InitLogServer(randomTaskUID string) {
	logServer := logrus.New()
	logServer.SetLevel(logrus.InfoLevel)
	logServer.SetFormatter(&logrus.JSONFormatter{})
	logServer.SetReportCaller(true)
	logfile, _ := os.OpenFile("./log/"+randomTaskUID+".log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if gin.Mode() == gin.DebugMode {
		logWriters := []io.Writer{
			logfile,
			os.Stdout,
		}
		fileAndStdoutWriter := io.MultiWriter(logWriters...)
		logServer.SetOutput(fileAndStdoutWriter)
	} else {
		logServer.SetOutput(logfile)
	}
	logServerMap.Store(randomTaskUID, logServer)
}

func GetLogMap(randomTaskID string) *logrus.Logger {
	logServer, OK := logServerMap.Load(randomTaskID)
	if OK == false {
		return nil
	} else {
		return logServer.(*logrus.Logger)
	}
}

func GetLogMapWithBool(randomTaskID any, isOK bool) *logrus.Logger {
	if isOK == false {
		return nil
	}
	logServer, OK := logServerMap.Load(randomTaskID)
	if OK == false {
		return nil
	} else {
		return logServer.(*logrus.Logger)
	}
}
