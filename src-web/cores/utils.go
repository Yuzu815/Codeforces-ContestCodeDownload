package cores

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
)

func getRandomStringHex(strLen int) string {
	if strLen <= 0 {
		return string([]byte{})
	}
	var need int
	if strLen&1 == 0 {
		need = strLen
	} else {
		need = strLen + 1
	}
	size := need / 2
	dst := make([]byte, need)
	src := dst[size:]
	if _, err := rand.Read(src[:]); err != nil {
		return string([]byte{})
	}
	hex.Encode(dst, src)
	return string(dst[:strLen])
}

func parseJsonFiles(infoForID gjson.Result) informationStruct {
	var temp informationStruct
	temp.ID = infoForID.Get(`0.id`).Int()
	temp.CID = infoForID.Get(`0.contestId`).Int()
	temp.PID = infoForID.Get(`0.problem.index`).String()
	temp.PNAME = infoForID.Get(`0.problem.name`).String()
	if infoForID.Get(`0.author.members`).Int() == 1 {
		temp.CNAME = infoForID.Get(`0.author.members.0.handle`).String()
	} else {
		temp.CNAME = infoForID.Get(`0.author.members.0.name`).String()
	}
	temp.LANG = infoForID.Get(`0.programmingLanguage`).String()
	return temp
}

func getAllAcceptSubmissionID(apiJsonString string) []string {
	allContestantResult := gjson.Get(apiJsonString, `result.#(author.participantType="CONTESTANT")#.id`)
	allVerdictOKResult := gjson.Get(apiJsonString, `result.#(verdict="OK")#.id`)
	return intersectGjsonResult(allContestantResult, allVerdictOKResult)
}

func getAPIJsonString(signedURL string) string {
	apiData, err := http.Get(signedURL)
	if err != nil {
		logServer.WithFields(logrus.Fields{
			"reason": err.Error(),
		}).Errorln("An error occurred while getting API URL.")
		return ""
	}
	apiBytes, _ := ioutil.ReadAll(apiData.Body)
	apiJsonString := string(apiBytes)
	return apiJsonString
}
