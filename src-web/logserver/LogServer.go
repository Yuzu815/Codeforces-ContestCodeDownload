package logserver

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
)

// logServerMap TODO F: 稍晚将日志部分抽离出来，并日志中应能返回对应的函数名
var logServerMap = sync.Map{}

func InitLogServer(randomUID string) {
	logServer := logrus.New()
	logServer.SetReportCaller(true)
	logServer.SetFormatter(&logrus.JSONFormatter{})
	logfile, _ := os.OpenFile("./log/"+randomUID+".log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	logCollection, _ := os.OpenFile("./log/logCollection.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	prefixUIDHook := &prefixHook{UID: randomUID}
	logServer.AddHook(prefixUIDHook)
	if gin.Mode() == gin.DebugMode {
		logWriters := []io.Writer{
			logfile,
			logCollection,
			os.Stdout,
		}
		fileAndStdoutWriter := io.MultiWriter(logWriters...)
		logServer.SetOutput(fileAndStdoutWriter)
	} else {
		logWriters := []io.Writer{
			logfile,
			logCollection,
		}
		fileAndStdoutWriter := io.MultiWriter(logWriters...)
		logServer.SetOutput(fileAndStdoutWriter)
	}
	logServerMap.Store(randomUID, logServer)
}

func GetLogMap(randomUID string) *logrus.Logger {
	logServer, OK := logServerMap.Load(randomUID)
	if OK == false {
		return nil
	} else {
		return logServer.(*logrus.Logger)
	}
}

func GetLogMapWithBool(randomUID any, isOK bool) *logrus.Logger {
	if isOK == false {
		return nil
	}
	logServer, OK := logServerMap.Load(randomUID)
	if OK == false {
		return nil
	} else {
		return logServer.(*logrus.Logger)
	}
}
