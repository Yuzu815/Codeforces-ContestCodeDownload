package logMode

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
)

// LogServerMap TODO F: 稍晚将日志部分抽离出来，并日志中应能返回对应的函数名
var LogServerMap = sync.Map{}

func InitLogServer(randomTaskUID string) {
	logServer := logrus.New()
	//logrus.SetLevel(logrus.TraceLevel)
	logServer.SetLevel(logrus.InfoLevel)
	logServer.SetFormatter(&logrus.JSONFormatter{})
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
	LogServerMap.Store(randomTaskUID, logServer)
}

func GetLogMap(randomTaskID string) *logrus.Logger {
	logServer, OK := LogServerMap.Load(randomTaskID)
	if OK == false {
		return nil
	} else {
		return logServer.(*logrus.Logger)
	}
}

func GetLogMapBool(randomTaskID any, isOK bool) *logrus.Logger {
	if isOK == false {
		return nil
	}
	logServer, OK := LogServerMap.Load(randomTaskID)
	if OK == false {
		return nil
	} else {
		return logServer.(*logrus.Logger)
	}
}
