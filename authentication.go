package main

import (
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

func getCodeforcesHttpClient(username, password string) *http.Client {
	cookiejarValue, _ := cookiejar.New(nil)
	//Fiddler DEBUG PROXY ADDRESS
	//DEBUG_PROXY_URL, _ := url.Parse("http://127.0.0.1:8866")
	codeforcesHttpClient := &http.Client{
		Jar: cookiejarValue,
		/*
			Transport: &http.Transport{
				Proxy: http.ProxyURL(DEBUG_PROXY_URL),
			},
		*/
	}
	getCsrfRequest, _ := http.NewRequest("GET", "https://codeforces.com/enter?back=%2F", nil)
	getCsrfRequest.Header.Add("Host", "codeforces.com")
	getCsrfRequest.Header.Add("User-Agent", "Golang-FetchCode")
	getCsrfRequestRespond, _ := codeforcesHttpClient.Do(getCsrfRequest)
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
	codeforcesHttpClient.Do(getLoginCookieRequest)
	return codeforcesHttpClient
}
