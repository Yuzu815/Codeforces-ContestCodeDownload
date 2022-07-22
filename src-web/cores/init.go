package cores

import (
	"Codeforces-ContestCodeDownload/src-web/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

// LogServer TODO F: 稍晚将日志部分抽离出来
var LogServer = logrus.New()
var RandomTaskName = "RandomTaskName"
var PROCESS = sync.Map{}

// InformationStruct
/*
# Return Information
ID := result.id
CID := result.contestId
PID := result.problem.index
PNAME := result.problem.name
CNAME := result.author.members.[name(Maybe NULL)/handle]
LANG := result.programmingLanguage

fileName := PID-PNAME-CNAME-LANG(CID#ID)
*/
type InformationStruct struct {
	ID    int64
	CID   int64
	PID   string
	PNAME string
	CNAME string
	LANG  string
	//TODO F: 对部分缩写进行重构
}

func saveSourceCodeToFile(sourceCode string, infoCode InformationStruct) {
	sufferName := ".txt"
	abbrLANG := "Other"
	if strings.Contains(infoCode.LANG, "C++") || strings.Contains(infoCode.LANG, "G++") || strings.Contains(infoCode.LANG, "Clang") {
		sufferName = ".cpp"
		abbrLANG = "C++"
	} else if strings.Contains(infoCode.LANG, "Java") {
		sufferName = ".java"
		abbrLANG = "Java"
	} else if strings.Contains(infoCode.LANG, "Python") || strings.Contains(infoCode.LANG, "PyPy") {
		sufferName = ".py"
		abbrLANG = "Python"
	} else if strings.Contains(infoCode.LANG, "GCC") {
		sufferName = ".c"
		abbrLANG = "C"
	}
	invalidChar := "<|>\\/:\"*?"
	fileName := fmt.Sprintf("%s-%s-%s-%s(%d#%d)%s", infoCode.PID, infoCode.PNAME, infoCode.CNAME, abbrLANG, infoCode.CID, infoCode.ID, sufferName)
	if strings.ContainsAny(fileName, invalidChar) {
		invalidRegexp := regexp.MustCompile(`[` + invalidChar + `]`)
		newFileName := invalidRegexp.ReplaceAllString(fileName, "")
		LogServer.WithFields(logrus.Fields{
			"oldFileName": fileName,
			"newFileName": newFileName,
		}).Warnln("The constructed file name may contain illegal characters that are not allowed by the operating system.")
		fileName = newFileName
	}
	err := ioutil.WriteFile("./temp/"+RandomTaskName+"/"+fileName, []byte(sourceCode), 0664)
	if err != nil {
		LogServer.WithFields(logrus.Fields{
			"reason":   err.Error(),
			"fileName": fileName,
		}).Errorln("An error occurred while saving the fetched code to a file.")
	}
}

/*
# Condition
result.author.participantType = "CONTESTANT"
result.verdict = "OK",
*/
func getAllAcceptSubmissionData(signedURL string, manageClient *http.Client) []InformationStruct {
	apiJsonString := getAPIJsonString(signedURL)
	allAcceptSubmissionID := getAllAcceptSubmissionID(apiJsonString)
	if len(allAcceptSubmissionID) == 0 {
		LogServer.WithFields(logrus.Fields{
			"signedURL": signedURL,
		}).Errorln("The list of obtained submission records is empty.")
	}
	var allNeedInformation []InformationStruct
	for index, submissionID := range allAcceptSubmissionID {
		infoForID := gjson.Get(apiJsonString, `result.#(id=`+string(submissionID)+`)#`)
		temp := parseJsonFiles(infoForID)
		LogServer.WithFields(logrus.Fields{
			"CID":   temp.CID,
			"ID":    temp.ID,
			"CNAME": temp.CNAME,
			"LANG":  temp.LANG,
		}).Infoln("Fetching this source...")
		var tempSourceCode string
		if temp.CID > 100000 {
			tempSourceCode = fetchSubmissionCode(fmt.Sprintf(`https://codeforces.com/gym/%d/submission/%d`, temp.CID, temp.ID), manageClient)
		} else {
			tempSourceCode = fetchSubmissionCode(fmt.Sprintf(`https://codeforces.com/contest/%d/submission/%d`, temp.CID, temp.ID), manageClient)
		}
		saveSourceCodeToFile(tempSourceCode, temp)
		allNeedInformation = append(allNeedInformation, temp)
		LogServer.WithFields(logrus.Fields{
			"Index":        index,
			"SubmissionID": submissionID,
		}).Infoln("Have fetched this source...")
		//TODO F: 此处使用全局变量来计算，后期需修正
		PROCESS.Store(RandomTaskName, float64(index+1)/float64(len(allAcceptSubmissionID)))
	}
	ZipCompress("./temp/"+RandomTaskName, RandomTaskName)
	return allNeedInformation
}

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

// MissionStart TODO F: 作为任务启动的接口，返回值格式需进行一定的修改，预期添加文件名。
func MissionStart(contestID int, info model.CodeforcesUserModel) []InformationStruct {
	action := "/contest.status?"
	actionParameter := "contestId=" + strconv.Itoa(contestID)
	signedURL := getSignedURL(info.ApiKey, info.ApiSecret, action, actionParameter)
	//TODO F: 对返回的loginRespond进行检查
	httpClient, _ := GetCodeforcesHttpClient(info.Username, info.Password)
	return getAllAcceptSubmissionData(signedURL, httpClient)
}
