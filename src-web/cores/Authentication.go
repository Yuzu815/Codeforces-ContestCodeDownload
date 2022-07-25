package cores

import (
	"Codeforces-ContestCodeDownload/src-web/logserver"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)

func matchCsrfString(htmlString string) string {
	regexCsrfFirst, _ := regexp.Compile(`<meta name="X-Csrf-Token" content="([\da-f]*)"`)
	matchStringFirst := regexCsrfFirst.FindString(htmlString)
	regexCsrfSecond, _ := regexp.Compile(`"([\da-f]*)"`)
	matchStringSecond := regexCsrfSecond.FindString(matchStringFirst)
	return matchStringSecond[1 : len(matchStringSecond)-1]
}

func GetCodeforcesHttpClient(username, password, randomUID string) (*http.Client, *http.Response) {
	cookiejarValue, _ := cookiejar.New(nil)
	//TODO E: 添加网络检查，代理连接可能会失败，需处理
	//TODO F: 首頁表單添加伸縮，打開高級選項可以配置代理
	//AccelerateProxyUrl, _ := url.Parse("http://127.0.0.1:44444")
	codeforcesHttpClient := &http.Client{
		Jar: cookiejarValue,
		/*
			Transport: &http.Transport{
				Proxy: http.ProxyURL(AccelerateProxyUrl),
			},
		//*/
	}
	getCsrfRequest, _ := http.NewRequest("GET", "https://codeforces.com/enter?back=%2F", nil)
	getCsrfRequest.Header.Add("Host", "codeforces.com")
	getCsrfRequest.Header.Add("User-Agent", "Golang-FetchCode")
	getCsrfRequestRespond, err := codeforcesHttpClient.Do(getCsrfRequest)
	if err != nil {
		logserver.GetLogMap(randomUID).WithFields(logrus.Fields{
			"reason": err.Error(),
		}).Errorln("An error occurred while fetching the CSRF TOKEN.")
		return nil, nil
	}
	includedCsrfBodyData, _ := ioutil.ReadAll(getCsrfRequestRespond.Body)
	csrfValue := matchCsrfString(string(includedCsrfBodyData))
	logserver.GetLogMap(randomUID).WithFields(logrus.Fields{
		"CSRF Value": csrfValue,
	}).Infoln("Matched CSRF Value")
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
	response, err := codeforcesHttpClient.Do(getLoginCookieRequest)
	if err != nil {
		logserver.GetLogMap(randomUID).WithFields(logrus.Fields{
			"reason": err.Error(),
		}).Errorln("Error when sending a POST request to simulate a login.")
		return nil, response
	}
	return codeforcesHttpClient, response
}
