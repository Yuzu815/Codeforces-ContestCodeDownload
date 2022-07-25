package cores

import (
	"Codeforces-ContestCodeDownload/src-web/logserver"
	"Codeforces-ContestCodeDownload/src-web/model"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func saveSourceCodeToFile(sourceCode, randomUID string, infoCode model.InformationStruct) {
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
	} else if strings.Contains(infoCode.LANG, "GCC") || strings.Contains(infoCode.LANG, "C11") {
		sufferName = ".c"
		abbrLANG = "C"
	}
	invalidChar := "<|>\\/:\"*?"
	fileName := fmt.Sprintf("%s-%s-%s-%s(%d#%d)%s", infoCode.PID, infoCode.PNAME, infoCode.CNAME, abbrLANG, infoCode.CID, infoCode.ID, sufferName)
	if strings.ContainsAny(fileName, invalidChar) {
		invalidRegexp := regexp.MustCompile(`[` + invalidChar + `]`)
		newFileName := invalidRegexp.ReplaceAllString(fileName, "")
		logserver.GetLogMap(randomUID).WithFields(logrus.Fields{
			"oldFileName": fileName,
			"newFileName": newFileName,
		}).Warnln("The constructed file name may contain illegal characters that are not allowed by the operating system.")
		fileName = newFileName
	}
	err := ioutil.WriteFile("./temp/"+randomUID+"/"+fileName, []byte(sourceCode), 0664)
	if err != nil {
		logserver.GetLogMap(randomUID).WithFields(logrus.Fields{
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
func getAllAcceptSubmissionData(signedURL, randomUID string, manageClient *http.Client) []model.InformationStruct {
	apiJsonString := getAPIJsonString(signedURL, randomUID)
	allAcceptSubmissionID := getAllAcceptSubmissionID(apiJsonString)
	if len(allAcceptSubmissionID) == 0 {
		logserver.GetLogMap(randomUID).WithFields(logrus.Fields{
			"signedURL": signedURL,
		}).Errorln("The list of obtained submission records is empty.")
	}
	var allNeedInformation []model.InformationStruct
	missionLogMapRef, _ := TaskMessageChan.Load(randomUID)
	for index, submissionID := range allAcceptSubmissionID {
		infoForID := gjson.Get(apiJsonString, `result.#(id=`+string(submissionID)+`)#`)
		infoStruct := parseJsonFiles(infoForID)
		logserver.GetLogMap(randomUID).WithFields(logrus.Fields{
			"CID":   infoStruct.CID,
			"ID":    infoStruct.ID,
			"CNAME": infoStruct.CNAME,
			"LANG":  infoStruct.LANG,
		}).Infoln("Fetching this source...")
		missionLogMapRef.(chan string) <- realtimeLogStr(infoStruct)
		var tempSourceCode string
		if infoStruct.CID > 100000 {
			tempSourceCode = fetchSubmissionCode(fmt.Sprintf(`https://codeforces.com/gym/%d/submission/%d`, infoStruct.CID, infoStruct.ID), randomUID, manageClient)
		} else {
			tempSourceCode = fetchSubmissionCode(fmt.Sprintf(`https://codeforces.com/contest/%d/submission/%d`, infoStruct.CID, infoStruct.ID), randomUID, manageClient)
		}
		saveSourceCodeToFile(tempSourceCode, randomUID, infoStruct)
		allNeedInformation = append(allNeedInformation, infoStruct)
		logserver.GetLogMap(randomUID).WithFields(logrus.Fields{
			"Index":        index,
			"SubmissionID": submissionID,
		}).Infoln("Have fetched this source...")
		//TODO F: 此处使用全局变量来计算，后期需修正
		MissionProgressMap.Store(randomUID, float64(index+1)/float64(len(allAcceptSubmissionID)))
	}
	missionLogMapRef.(chan string) <- RESULT_IS_END_FLAG
	ZipCompress("./temp/"+randomUID, "./temp/"+randomUID, randomUID)
	return allNeedInformation
}

func realtimeLogStr(jsonResult model.InformationStruct) string {
	return fmt.Sprintf("[Info] CID:%d ID:%d CNAME:%s LANG:%s\n", jsonResult.CID, jsonResult.ID, jsonResult.CNAME, jsonResult.LANG)
}
