package cores

import (
	"Codeforces-ContestCodeDownload/src-web/logserver"
	"crypto/sha512"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"net/http"
	"strconv"
	"time"
)

/*
If your key is xxx, secret is yyy, chosen rand is 123456, and you want to access method contest.hacks for contest 566,
you should compose request like this:
https://codeforces.com/api/contest.hacks?contestId=566&apiKey=xxx&time=1656689340&apiSig=123456<hash>,
where <hash> is sha512Hex(123456/contest.hacks?apiKey=xxx&contestId=566&time=1656689340#yyy)
Note: First six characters of the apiSig parameter can be arbitrary.
*/
func getSignedURL(apiKey, apiSecret, action, actionParameter string) string {
	nowTime := strconv.FormatInt(time.Now().Unix(), 10)
	magicStr := getRandomStringHex(6)
	hashRaw := fmt.Sprintf("%s%sapiKey=%s&%s&time=%s#%s", magicStr, action, apiKey, actionParameter, nowTime, apiSecret)
	sha512Bytes := sha512.Sum512([]byte(hashRaw))
	sha512String := fmt.Sprintf("%x", sha512Bytes)
	apiSig := magicStr + sha512String
	signedURL := fmt.Sprintf("https://codeforces.com/api%s%s&apiKey=%s&time=%s&apiSig=%s", action, actionParameter, apiKey, nowTime, apiSig)
	return signedURL
}

func intersectGjsonResult(resultA gjson.Result, resultB gjson.Result) []string {
	var result []string
	resultA.ForEach(func(_, elementA gjson.Result) bool {
		isCount := false
		resultB.ForEach(func(_, elementB gjson.Result) bool {
			if elementA.String() == elementB.String() {
				isCount = true
			}
			return true
		})
		if isCount {
			result = append(result, elementA.String())
		}
		return true
	})
	return result
}

func fetchSubmissionCode(submissionURL, randomUID string, manageClient *http.Client) string {
	getSourceCodeRequest, _ := http.NewRequest("GET", submissionURL, nil)
	sourceHtmlRespond, err := manageClient.Do(getSourceCodeRequest)
	if err != nil {
		logserver.GetLogMap(randomUID).WithFields(logrus.Fields{
			"reason": err.Error(),
		}).Errorln("An error occurred while fetching the submission.")
	}
	sourceHtmlReader, _ := goquery.NewDocumentFromReader(sourceHtmlRespond.Body)
	matchSourceCode := sourceHtmlReader.Find("#program-source-text").Text()
	return matchSourceCode
}
