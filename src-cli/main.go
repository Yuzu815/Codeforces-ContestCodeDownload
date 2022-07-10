package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var logServer = logrus.New()

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
type informationStruct struct {
	ID    int64
	CID   int64
	PID   string
	PNAME string
	CNAME string
	LANG  string
}

func saveSourceCodeToFile(sourceCode string, infoCode informationStruct) {
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
		logServer.WithFields(logrus.Fields{
			"oldFileName": fileName,
			"newFileName": newFileName,
		}).Warnln("The constructed file name may contain illegal characters that are not allowed by the operating system.")
		fileName = newFileName
	}
	err := ioutil.WriteFile(fileName, []byte(sourceCode), 0664)
	if err != nil {
		logServer.WithFields(logrus.Fields{
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
func getAllAcceptSubmissionData(signedURL string, manageClient *http.Client) []informationStruct {
	apiJsonString := getAPIJsonString(signedURL)
	allAcceptSubmissionID := getAllAcceptSubmissionID(apiJsonString)
	if len(allAcceptSubmissionID) == 0 {
		logServer.WithFields(logrus.Fields{
			"signedURL": signedURL,
		}).Errorln("The list of obtained submission records is empty.")
	}
	var allNeedInformation []informationStruct
	for index, submissionID := range allAcceptSubmissionID {
		infoForID := gjson.Get(apiJsonString, `result.#(id=`+string(submissionID)+`)#`)
		temp := parseJsonFiles(infoForID)
		logServer.WithFields(logrus.Fields{
			"CID":   temp.CID,
			"ID":    temp.ID,
			"CNAME": temp.CNAME,
			"LANG":  temp.LANG,
		}).Infoln("Fetching this source...")
		fmt.Printf("「DEBUG」 CID:%d ID:%d NAME:%s LANG:%s\n", temp.CID, temp.ID, temp.CNAME, temp.LANG)
		var tempSourceCode string
		if temp.CID > 100000 {
			tempSourceCode = fetchSubmissionCode(fmt.Sprintf(`https://codeforces.com/gym/%d/submission/%d`, temp.CID, temp.ID), manageClient)
		} else {
			tempSourceCode = fetchSubmissionCode(fmt.Sprintf(`https://codeforces.com/contest/%d/submission/%d`, temp.CID, temp.ID), manageClient)
		}
		saveSourceCodeToFile(tempSourceCode, temp)
		allNeedInformation = append(allNeedInformation, temp)
		logServer.WithFields(logrus.Fields{
			"Index":        index,
			"SubmissionID": submissionID,
		}).Infoln("Have fetched this source...")
		fmt.Printf("「DEBUG」 Loading: %d/%d\n", index+1, len(allAcceptSubmissionID))
	}
	return allNeedInformation
}

func readKeyFile() (string, string, string, string) {
	bytes, err := ioutil.ReadFile("api.key")
	if err != nil {
		logrus.Errorln("Exception when trying to read api.key file.")
		return "", "", "", ""
	}
	fileString := string(bytes)
	jsonResult := gjson.GetMany(fileString, "apiKey", "apiSecret", "username", "password")
	return jsonResult[0].String(), jsonResult[1].String(), jsonResult[2].String(), jsonResult[3].String()
}

func initLogServer() {
	//logrus.SetLevel(logrus.TraceLevel)
	logServer.SetLevel(logrus.InfoLevel)
	logServer.SetFormatter(&logrus.JSONFormatter{})
	logfile, _ := os.OpenFile("./logrus.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if 1 == 0 {
		logWriters := []io.Writer{
			logfile,
			os.Stdout,
		}
		fileAndStdoutWriter := io.MultiWriter(logWriters...)
		logServer.SetOutput(fileAndStdoutWriter)
	} else {
		logServer.SetOutput(logfile)
	}
}

func main() {
	initLogServer()
	apiKey, apiSecret, username, password := readKeyFile()
	if apiKey == "" || apiSecret == "" || username == "" || password == "" {
		logServer.Errorln("Some of the parameters read from api.key are empty.")
	}
	var contestID int
	fmt.Print("Please enter contest ID: ")
	_, err := fmt.Scanf("%d", &contestID)
	if err != nil {
		logServer.Errorln("Exception when reading in the contest ID.")
	}
	action := "/contest.status?"
	actionParameter := "contestId=" + strconv.Itoa(contestID)
	signedURL := getSignedURL(apiKey, apiSecret, action, actionParameter)
	getAllAcceptSubmissionData(signedURL, getCodeforcesHttpClient(username, password))
}
