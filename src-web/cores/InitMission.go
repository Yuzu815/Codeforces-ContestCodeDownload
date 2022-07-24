package cores

import (
	"Codeforces-ContestCodeDownload/src-web/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strconv"
	"sync"
)

// LogServer TODO F: 稍晚将日志部分抽离出来，并日志中应能返回对应的函数名
var LogServer = logrus.New()
var RandomTaskName = "RandomTaskName"
var PROCESS = sync.Map{}

// CodeforcesContestResult TODO F: 后期使用结构体封装，改造为定时删除数据
var CodeforcesContestResult = sync.Map{}

func initLogServer() {
	//logrus.SetLevel(logrus.TraceLevel)
	LogServer.SetLevel(logrus.InfoLevel)
	LogServer.SetFormatter(&logrus.JSONFormatter{})
	logfile, _ := os.OpenFile("./log/"+RandomTaskName+".log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if gin.Mode() == gin.DebugMode {
		logWriters := []io.Writer{
			logfile,
			os.Stdout,
		}
		fileAndStdoutWriter := io.MultiWriter(logWriters...)
		LogServer.SetOutput(fileAndStdoutWriter)
	} else {
		LogServer.SetOutput(logfile)
	}
}

func initRandomUID() {
	RandomTaskName = getRandomStringHex(16)
}

func initTempFileDir() {
	err := os.Mkdir("./temp/"+RandomTaskName, 0750)
	if err != nil {
		LogServer.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Infoln("Error in initTempFileDir...")
	}
}

// MissionInitiated 只用来初始化当前的任务
func MissionInitiated() {
	initRandomUID()
	initLogServer()
	initTempFileDir()
	PROCESS.Store(RandomTaskName, 0.0)
}

// MissionCall TODO F: 作为任务启动的接口，返回值格式需进行一定的修改，预期添加文件名。
func MissionCall(contestID int, info model.CodeforcesUserModel) []InformationStruct {
	action := "/contest.status?"
	actionParameter := "contestId=" + strconv.Itoa(contestID)
	signedURL := getSignedURL(info.ApiKey, info.ApiSecret, action, actionParameter)
	//TODO F: 对返回的loginRespond进行检查
	httpClient, _ := GetCodeforcesHttpClient(info.Username, info.Password)
	return getAllAcceptSubmissionData(signedURL, httpClient)
}
