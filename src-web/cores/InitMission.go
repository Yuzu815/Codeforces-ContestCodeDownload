package cores

import (
	"Codeforces-ContestCodeDownload/src-web/logMode"
	"Codeforces-ContestCodeDownload/src-web/model"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"sync"
)

const RESULT_IS_END_FLAG = "RESULT_IS_END_FLAG"

// MissionProgressMap 用于存UID对应的进度
var MissionProgressMap = sync.Map{}

// TaskMessageChan TODO F: 将每一ID映射为一个通道，需实现定期删除
var TaskMessageChan = sync.Map{}

func initRandomUIDLogServer() string {
	randomUID := getRandomStringHex(16)
	logMode.InitLogServer(randomUID)
	logMode.GetLogMap(randomUID).Infoln("Successful mission creation, UID: " + randomUID)
	return randomUID
}

func initTempFileDir(randomUID string) {
	err := os.Mkdir("./temp/"+randomUID, 0750)
	if err != nil {
		logMode.GetLogMap(randomUID).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Infoln("Error in initTempFileDir...")
	}
}

func initMessageChan(randomUID string) {
	TaskMessageChan.Store(randomUID, make(chan string, 100))
}

func initProcessInterface(randomUID string) {
	MissionProgressMap.Store(randomUID, 0.0)
}

// MissionInitiated 只用来初始化当前的任务
func MissionInitiated() string {
	randomUID := initRandomUIDLogServer()
	initTempFileDir(randomUID)
	initMessageChan(randomUID)
	initProcessInterface(randomUID)
	return randomUID
}

// MissionCall TODO F: 作为任务启动的接口，返回值格式需进行一定地修改，预期添加文件名。
func MissionCall(contestID int, randomUID string, info model.CodeforcesUserModel) {
	action := "/contest.status?"
	actionParameter := "contestId=" + strconv.Itoa(contestID)
	signedURL := getSignedURL(info.ApiKey, info.ApiSecret, action, actionParameter)
	logMode.GetLogMap(randomUID).WithFields(logrus.Fields{
		"signed URL": signedURL,
	}).Infoln("generate signed URL.")
	//TODO F: 对返回的loginRespond进行检查
	httpClient, _ := GetCodeforcesHttpClient(info.Username, info.Password, randomUID)
	// TODO E: 某些时候CSRF匹配失败时，程序会崩溃，需做特殊处理，此處直接返回簡單處理下。
	if httpClient == nil {
		logMode.GetLogMap(randomUID).Errorln("httpClient is NULL. Please check network or website availability.")
	}
	getAllAcceptSubmissionData(signedURL, randomUID, httpClient)
}
