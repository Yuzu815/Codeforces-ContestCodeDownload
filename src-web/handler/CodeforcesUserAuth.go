package handler

import (
	"Codeforces-ContestCodeDownload/src-web/cores"
	"Codeforces-ContestCodeDownload/src-web/logMode"
	"Codeforces-ContestCodeDownload/src-web/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

// CodeforcesUserAuth TODO F: 验证通过后，重定向到/result页面，实时显示抓取情况，并展示进度条，否则重定向到/error页面
func CodeforcesUserAuth(context *gin.Context) {
	randomUID := cores.MissionInitiated()
	encryptedApiKey := context.PostForm("apiKey")
	encryptedApiSecret := context.PostForm("apiSecret")
	encryptedUsername := context.PostForm("usernameOrEmail")
	encryptedPassword := context.PostForm("password")
	encryptedUserData := model.CodeforcesUserModel{
		ApiKey:    encryptedApiKey,
		ApiSecret: encryptedApiSecret,
		Username:  encryptedUsername,
		Password:  encryptedPassword,
	}
	//TODO F: 添加空值校验。账号密码校验已添加，未对API KET做校验
	userData := decryptUserData(encryptedUserData, "123")
	if checkLoginStatus(cores.GetCodeforcesHttpClient(userData.Username, userData.Password, randomUID)) == false {
		context.Redirect(http.StatusSeeOther, "?err=logErr")
		logMode.GetLogMap(randomUID).Errorln("Login fail. Please check your username and password.")
	} else {
		logMode.GetLogMap(randomUID).WithFields(logrus.Fields{
			"ApiKey":   userData.ApiKey,
			"Username": userData.Username,
		}).Info("Have access to user information.")
		go fetchContestData(userData, randomUID)
		cookieDomain := getSiteDomain(context.Request.Header.Get("Referer"))
		context.SetCookie("UID", randomUID, 86400, "/", cookieDomain, false, true)
		context.Redirect(http.StatusSeeOther, "/result")
	}
}

// checkLoginStatus TODO F: 或许将这一检查写成client, error会更合适。
func checkLoginStatus(client *http.Client, response *http.Response) bool {
	if response == nil {
		return false
	}
	body, _ := ioutil.ReadAll(response.Body)
	if strings.Contains(string(body), "Invalid handle/email or password") ||
		strings.Contains(string(body), "Please, confirm email before entering the website.") {
		return false
	}
	return true
}

// decryptUserData TODO F: 添加加密安全传输信息
func decryptUserData(encryptedInformation model.CodeforcesUserModel, decryptedKey string) model.CodeforcesUserModel {
	return encryptedInformation
}

// fetchContestData TODO F: 已改用全局MAP方案，后期需实现将每一context都绑定上一个randomID
func fetchContestData(userData model.CodeforcesUserModel, randomUID string) {
	contestID := 380042
	result := cores.MissionCall(contestID, randomUID, userData)
	logMode.GetLogMap(randomUID).WithFields(logrus.Fields{
		"contestID":  contestID,
		"jsonResult": result,
	}).Info("Source code and record correspondence information has been obtained from codeforces.")
	cores.CodeforcesContestResult.Store(randomUID, result)
}

func getSiteDomain(urlStr string) string {
	matchString := strings.Split(urlStr, "//")[1]
	matchStringAgain := strings.Split(matchString, "/")[0]
	return matchStringAgain
}
