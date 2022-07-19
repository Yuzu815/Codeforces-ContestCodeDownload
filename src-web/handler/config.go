package handler

import (
	"Codeforces-ContestCodeDownload/src-web/cores"
	"Codeforces-ContestCodeDownload/src-web/model"
	"fmt"
	"github.com/gin-gonic/gin"
)

// decryptUserData TODO F: 添加加密安全传输信息
func decryptUserData(encryptedInformation model.CodeforcesUserModel, decryptedKey string) model.CodeforcesUserModel {
	return encryptedInformation
}

func SaveCodeforcesConfig(context *gin.Context) {
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
	result := cores.MissionInitiated(381185, userData)
	context.Set("CodeforcesResult", result)
	fmt.Println(context.Value("CodeforcesResult"))
	//TODO F: 重定向到Result界面，并尝试上下文传值
	context.Request.URL.Path = "/test2"
	//router.HandleContext(context)
}
