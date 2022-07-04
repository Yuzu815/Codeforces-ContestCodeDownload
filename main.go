package main

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func readKeyFile() (string, string, string, string) {
	//TODO E: 文件读入检测报错
	bytes, _ := ioutil.ReadFile("api.key")
	fileString := string(bytes)
	//TODO E: json文件读取可能没有对应值报错
	jsonResult := gjson.GetMany(fileString, "apiKey", "apiSecret", "username", "password")
	return jsonResult[0].String(), jsonResult[1].String(), jsonResult[2].String(), jsonResult[3].String()
}

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

func matchCsrfString(htmlString string) string {
	//TODO F: 或许有更优秀的匹配CSRF的方案
	regexCsrfFirst, _ := regexp.Compile(`<meta name="X-Csrf-Token" content="([\da-f]*)"`)
	matchStringFirst := regexCsrfFirst.FindString(htmlString)
	regexCsrfSecond, _ := regexp.Compile(`"([\da-f]*)"`)
	matchStringSecond := regexCsrfSecond.FindString(matchStringFirst)
	//TODO E: 匹配失败的错误处理未完成，可能会导致数组越界
	return matchStringSecond[1 : len(matchStringSecond)-1]
}

func getCodeforcesHttpClient(username, password string) *http.Client {
	cookiejarValue, _ := cookiejar.New(nil)
	//此处可以设置Fiddler抓包地址，便于抓包
	//DEBUG_PROXY_URL, _ := url.Parse("http://127.0.0.1:8866")
	codeforcesHttpClient := &http.Client{
		Jar: cookiejarValue,
		/*
			Transport: &http.Transport{
				Proxy: http.ProxyURL(DEBUG_PROXY_URL),
			},
		*/
	}
	//TODO E: 异常处理未编写，暂不清楚会有什么异常
	getCsrfRequest, _ := http.NewRequest("GET", "https://codeforces.com/enter?back=%2F", nil)
	getCsrfRequest.Header.Add("Host", "codeforces.com")
	getCsrfRequest.Header.Add("User-Agent", "Golang-FetchCode")
	//TODO E: 异常处理未编写，可能发送GET请求时Codeforces崩溃，无法使用
	getCsrfRequestRespond, _ := codeforcesHttpClient.Do(getCsrfRequest)
	//TODO E: 关于IO的异常处理未编写
	includedCsrfBodyData, _ := ioutil.ReadAll(getCsrfRequestRespond.Body)
	csrfValue := matchCsrfString(string(includedCsrfBodyData))
	postValue := url.Values{
		"csrf_token":    {csrfValue},
		"action":        {"enter"},
		"ftaa":          {getRandomStringHex(18)},
		"bfaa":          {getRandomStringHex(32)},
		"handleOrEmail": {username},
		"password":      {password},
		"_tta":          {"200"},
	}
	getLoginCookieRequest, _ := http.NewRequest("POST", "https://codeforces.com/enter?back=%2F", strings.NewReader(postValue.Encode()))
	getLoginCookieRequest.Header.Add("Host", "codeforces.com")
	getLoginCookieRequest.Header.Add("User-Agent", "Golang-FetchCode")
	getLoginCookieRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//TODO E: 此处未进行错误处理
	codeforcesHttpClient.Do(getLoginCookieRequest)
	return codeforcesHttpClient
}

func fetchSubmissionCode(submissionURL string, manageClient *http.Client) string {
	getSourceCodeRequest, _ := http.NewRequest("GET", submissionURL, nil)
	sourceHtmlRespond, _ := manageClient.Do(getSourceCodeRequest)
	sourceHtmlReader, _ := goquery.NewDocumentFromReader(sourceHtmlRespond.Body)
	matchSourceCode := sourceHtmlReader.Find("#program-source-text").Text()
	//fmt.Println(matchSourceCode)
	return matchSourceCode
}

/*
# Return Information
ID := result.id
CID := result.contestId
PID := result.problem.index
PNAME := result.problem.name
CNAME := result.author.members.[name(Maybe NULL)/handle]
LANG := result.programmingLanguage

# Maybe
fileName := PID-PNAME-CNAME-LANG(CID#ID)
*/
type allNeedInformationStruct struct {
	ID    int64
	CID   int64
	PID   string
	PNAME string
	CNAME string
	LANG  string
}

func saveSourceCodeToFile(sourceCode string, infoCode allNeedInformationStruct) {
	//TODO F: 暂时只支持比较主流的语言，其他语言用Other表示，并用txt作为后缀存储
	sufferName := ".txt"abbrLANG := "Other"
	if strings.Contains(infoCode.LANG, "C++") || strings.Contains(infoCode.LANG, "G++") || strings.Contains(infoCode.LANG, "Clang") {
		sufferName = ".cpp"
		abbrLANG = "C++"
	} else if strings.Contains(infoCode.LANG, "Java") {
		sufferName = ".java"
		abbrLANG = "Java"
	} else if strings.Contains(infoCode.LANG, "Python") || strings.Contains(infoCode.LANG, "Pypi"){
		sufferName = ".py"
		abbrLANG = "Python"
	} else if strings.Contains(infoCode.LANG, "GCC") {
		sufferName = ".c"
		abbrLANG = "C"
	}
	//TODO E: 生成的文件名可能在某些操作系统是不合法的，需要设置一个排除表，对其中元素特殊编码
	fileName := fmt.Sprintf("%s-%s-%s-%s(%d#%d)%s", infoCode.PID, infoCode.PNAME, infoCode.CNAME, abbrLANG, infoCode.CID, infoCode.ID, sufferName)
	ioutil.WriteFile(fileName, []byte(sourceCode), 0664)
}

/*
# Condition
result.author.participantType = "CONTESTANT"
result.verdict = "OK",
*/
func getAllAcceptSubmissionData(signedURL string, manageClient *http.Client) []allNeedInformationStruct {
	//TODO E: 对Get的异常处理
	apiData, _ := http.Get(signedURL)
	//TODO E: 对response读取的异常处理
	apiBytes, _ := ioutil.ReadAll(apiData.Body)
	apiJsonString := string(apiBytes)
	//TODO F: 寻找到gjson是否能支持多条件匹配，目前采用的是取两个Result的交集
	allContestantResult := gjson.Get(apiJsonString, `result.#(author.participantType="CONTESTANT")#.id`)
	allVerdictOKResult := gjson.Get(apiJsonString, `result.#(verdict="OK")#.id`)
	allAcceptSubmissionID := intersectGjsonResult(allContestantResult, allVerdictOKResult)
	var allNeedInformation []allNeedInformationStruct
	for id, s := range allAcceptSubmissionID {
		infoForID := gjson.Get(apiJsonString, `result.#(id=`+string(s)+`)#`)
		var temp allNeedInformationStruct
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
		fmt.Printf("「DEBUG」 CID:%d ID:%d NAME:%s LANG:%s\n", temp.CID, temp.ID, temp.CNAME, temp.LANG)
		var tempSourceCode string
		if temp.CID > 100000 {
			tempSourceCode = fetchSubmissionCode(fmt.Sprintf(`https://codeforces.com/gym/%d/submission/%d`, temp.CID, temp.ID), manageClient)
		} else {
			tempSourceCode = fetchSubmissionCode(fmt.Sprintf(`https://codeforces.com/contest/%d/submission/%d`, temp.CID, temp.ID), manageClient)
		}
		saveSourceCodeToFile(tempSourceCode, temp)
		allNeedInformation = append(allNeedInformation, temp)
		fmt.Printf("「DEBUG」 Loading: %d/%d\n", id+1, len(allAcceptSubmissionID))
	}
	return allNeedInformation
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