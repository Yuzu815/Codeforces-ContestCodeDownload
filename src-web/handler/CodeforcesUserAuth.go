package handler

import (
	"Codeforces-ContestCodeDownload/src-web/cores"
	"Codeforces-ContestCodeDownload/src-web/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

// CodeforcesUserAuth TODO F: 验证通过后，重定向到/result页面，实时显示抓取情况，并展示进度条，否则重定向到/error页面
func CodeforcesUserAuth(context *gin.Context) {
	cores.MissionInitiated()
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
	if checkLoginStatus(cores.GetCodeforcesHttpClient(userData.Username, userData.Password)) == false {
		context.Redirect(http.StatusMovedPermanently, "?err=logErr")
		cores.LogServer.Errorln("Login fail. Please check your username and password.")
	} else {
		cores.LogServer.WithFields(logrus.Fields{
			"ApiKey":   userData.ApiKey,
			"Username": userData.Username,
		}).Info("Have access to user information.")
		go fetchContestData(userData, context)
		context.Redirect(http.StatusMovedPermanently, "/result")
	}
}

// checkLoginStatus TODO F: 或许将这一检查写成client, error会更合适。
func checkLoginStatus(client *http.Client, response *http.Response) bool {
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

// fetchContestData TODO F: 目前采用抽离函数并发执行的方式，未知*gin.Context是否会受影响
func fetchContestData(userData model.CodeforcesUserModel, context *gin.Context) {
	contestID := 381185
	result := cores.MissionStart(contestID, userData)
	cores.LogServer.WithFields(logrus.Fields{
		"contestID":  contestID,
		"jsonResult": result,
	}).Info("Source code and record correspondence information has been obtained from codeforces.")
	context.Set("CodeforcesResult", result)
}
