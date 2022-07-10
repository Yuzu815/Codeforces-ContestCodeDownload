package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

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
	fileName := fmt.Sprintf("%s-%s-%s-%s(%d#%d)%s", infoCode.PID, infoCode.PNAME, infoCode.CNAME, abbrLANG, infoCode.CID, infoCode.ID, sufferName)
	ioutil.WriteFile(fileName, []byte(sourceCode), 0664)
}

/*
# Condition
result.author.participantType = "CONTESTANT"
result.verdict = "OK",
*/
func getAllAcceptSubmissionData(signedURL string, manageClient *http.Client) []informationStruct {
	apiJsonString := getAPIJsonString(signedURL)
	allAcceptSubmissionID := getAllAcceptSubmissionID(apiJsonString)
	var allNeedInformation []informationStruct
	for index, submissionID := range allAcceptSubmissionID {
		infoForID := gjson.Get(apiJsonString, `result.#(id=`+string(submissionID)+`)#`)
		temp := parseJsonFiles(infoForID)
		fmt.Printf("「DEBUG」 CID:%d ID:%d NAME:%s LANG:%s\n", temp.CID, temp.ID, temp.CNAME, temp.LANG)
		var tempSourceCode string
		if temp.CID > 100000 {
			tempSourceCode = fetchSubmissionCode(fmt.Sprintf(`https://codeforces.com/gym/%d/submission/%d`, temp.CID, temp.ID), manageClient)
		} else {
			tempSourceCode = fetchSubmissionCode(fmt.Sprintf(`https://codeforces.com/contest/%d/submission/%d`, temp.CID, temp.ID), manageClient)
		}
		saveSourceCodeToFile(tempSourceCode, temp)
		allNeedInformation = append(allNeedInformation, temp)
		fmt.Printf("「DEBUG」 Loading: %d/%d\n", index+1, len(allAcceptSubmissionID))
	}
	return allNeedInformation
}

func readKeyFile() (string, string, string, string) {
	bytes, _ := ioutil.ReadFile("api.key")
	fileString := string(bytes)
	jsonResult := gjson.GetMany(fileString, "apiKey", "apiSecret", "username", "password")
	return jsonResult[0].String(), jsonResult[1].String(), jsonResult[2].String(), jsonResult[3].String()
}

func main() {
	var contestID int
	fmt.Print("Please enter contest ID: ")
	fmt.Scanf("%d", &contestID)
	apiKey, apiSecret, username, password := readKeyFile()
	action := "/contest.status?"
	actionParameter := "contestId=" + strconv.Itoa(contestID)
	signedURL := getSignedURL(apiKey, apiSecret, action, actionParameter)
	getAllAcceptSubmissionData(signedURL, getCodeforcesHttpClient(username, password))
}
