package handler

import (
	"Codeforces-ContestCodeDownload/src-web/cores"
	"Codeforces-ContestCodeDownload/src-web/logMode"
	"Codeforces-ContestCodeDownload/src-web/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// CodeforcesUserAuth POST /auth
func CodeforcesUserAuth(context *gin.Context) {
	randomUID := cores.MissionInitiated()
	cookieDomain := getSiteDomain(context.Request.Header.Get("Referer"))
	context.SetCookie("UID", randomUID, 86400, "/", cookieDomain, false, true)
	userData := decryptUserData(context.Copy(), randomUID, nil)
	if checkLoginStatus(cores.GetCodeforcesHttpClient(userData.Username, userData.Password, randomUID)) == false {
		context.Redirect(http.StatusSeeOther, "?err=logErr")
		logMode.GetLogMap(randomUID).Errorln("Login fail. Please check your username and password.")
	} else {
		logMode.GetLogMap(randomUID).WithFields(logrus.Fields{
			"ApiKey":   userData.ApiKey,
			"Username": userData.Username,
		}).Infoln("Have access to user information.")
		go fetchContestData(userData, randomUID)
		context.Redirect(http.StatusSeeOther, "/result")
	}
}

// checkLoginStatus TODO F: 或许将这一检查写成client, error会更合适。
func checkLoginStatus(client *http.Client, response *http.Response) bool {
	if response == nil || client == nil {
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
func decryptUserData(context *gin.Context, UID string, decryptedKey any) model.CodeforcesUserModel {
	if decryptedKey == nil {
		logMode.GetLogMap(UID).Infoln("Encryption and decryption not enabled")
	}
	encryptedApiKey := context.PostForm("apiKey")
	encryptedApiSecret := context.PostForm("apiSecret")
	encryptedUsername := context.PostForm("usernameOrEmail")
	encryptedPassword := context.PostForm("password")
	encryptedContestID := context.PostForm("CID")
	encryptedUserData := model.CodeforcesUserModel{
		ApiKey:    encryptedApiKey,
		ApiSecret: encryptedApiSecret,
		Username:  encryptedUsername,
		Password:  encryptedPassword,
		ContestID: encryptedContestID,
	}
	//TODO F: 添加空值校验
	return encryptedUserData
}

// fetchContestData TODO F: 已改用全局MAP方案，后期需实现将每一context都绑定上一个randomID
func fetchContestData(userData model.CodeforcesUserModel, randomUID string) {
	contestID, err := strconv.Atoi(userData.ContestID)
	if err != nil {
		logMode.GetLogMap(randomUID).Errorln(err.Error())
	}
	cores.MissionCall(contestID, randomUID, userData)
	logMode.GetLogMap(randomUID).WithFields(logrus.Fields{
		"contestID": contestID,
	}).Infoln("Source code and record correspondence information has been obtained from codeforces.")
}

func getSiteDomain(urlStr string) string {
	matchString := strings.Split(urlStr, "//")[1]
	matchStringAgain := strings.Split(matchString, "/")[0]
	return matchStringAgain
}
