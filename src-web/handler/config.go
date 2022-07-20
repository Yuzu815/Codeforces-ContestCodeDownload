package handler

import (
	"Codeforces-ContestCodeDownload/src-web/cores"
	"Codeforces-ContestCodeDownload/src-web/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// decryptUserData TODO F: 添加加密安全传输信息
func decryptUserData(encryptedInformation model.CodeforcesUserModel, decryptedKey string) model.CodeforcesUserModel {
	return encryptedInformation
}

// TODO F: 验证通过后，重定向到/result页面，实时显示抓取情况，并展示进度条，否则重定向到/error页面
func CodeforcesUserAuth(context *gin.Context) {
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
	//TODO F: 添加空值校验和账号密码校验（抓取登陆返回值），暂不对API KEY校验。
	userData := decryptUserData(encryptedUserData, "123")
	cores.LogServer.WithFields(logrus.Fields{
		"ApiKey":   userData.ApiKey,
		"Username": userData.Username,
	}).Info("Have access to user information.")
	contestID := 381185
	result := cores.MissionInitiated(contestID, userData)
	cores.LogServer.WithFields(logrus.Fields{
		"contestID":  contestID,
		"jsonResult": result,
	}).Info("Source code and record correspondence information has been obtained from codeforces.")
	context.Set("CodeforcesResult", result)
	context.Next()
}
