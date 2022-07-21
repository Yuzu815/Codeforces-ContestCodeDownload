package cores

import (
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

func GetCodeforcesHttpClient(username, password string) (*http.Client, *http.Response) {
	cookiejarValue, _ := cookiejar.New(nil)
	//Fiddler DEBUG PROXY ADDRESS
	//TODO E: 添加网络检查，代理连接可能会失败，需处理
	//DEBUG_PROXY_URL, _ := url.Parse("http://127.0.0.1:8866")
	ACCELERATE_PROXY_URL, _ := url.Parse("http://127.0.0.1:44444")
	codeforcesHttpClient := &http.Client{
		Jar: cookiejarValue,
		///*
		Transport: &http.Transport{
			Proxy: http.ProxyURL(ACCELERATE_PROXY_URL),
			//Proxy: http.ProxyURL(DEBUG_PROXY_URL),
		},
		//*/
	}
	getCsrfRequest, _ := http.NewRequest("GET", "https://codeforces.com/enter?back=%2F", nil)
	getCsrfRequest.Header.Add("Host", "codeforces.com")
	getCsrfRequest.Header.Add("User-Agent", "Golang-FetchCode")
	getCsrfRequestRespond, err := codeforcesHttpClient.Do(getCsrfRequest)
	if err != nil {
		LogServer.WithFields(logrus.Fields{
			"reason": err.Error(),
		}).Errorln("An error occurred while fetching the CSRF TOKEN.")
		return nil, nil
	}
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
	loginRespond, err := codeforcesHttpClient.Do(getLoginCookieRequest)
	if err != nil {
		LogServer.WithFields(logrus.Fields{
			"reason": err.Error(),
		}).Errorln("Error when sending a POST request to simulate a login.")
		return nil, nil
	}
	return codeforcesHttpClient, loginRespond
}
