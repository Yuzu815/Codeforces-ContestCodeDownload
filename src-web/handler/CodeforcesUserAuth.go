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
		context.Redirect(http.StatusSeeOther, "?err=logErr")
		cores.LogServer.Errorln("Login fail. Please check your username and password.")
	} else {
		cores.LogServer.WithFields(logrus.Fields{
			"ApiKey":   userData.ApiKey,
			"Username": userData.Username,
		}).Info("Have access to user information.")
		contextCopy := context.Copy()
		go fetchContestData(userData, contextCopy)
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
func fetchContestData(userData model.CodeforcesUserModel, context *gin.Context) {
	contestID := 380042
	result := cores.MissionCall(contestID, userData)
	cores.LogServer.WithFields(logrus.Fields{
		"contestID":  contestID,
		"jsonResult": result,
	}).Info("Source code and record correspondence information has been obtained from codeforces.")
	cores.CodeforcesContestResult.Store(cores.RandomTaskName, result)
}
