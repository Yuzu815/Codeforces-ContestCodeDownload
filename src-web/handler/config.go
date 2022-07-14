package handler

import (
	"Codeforces-ContestCodeDownload/src-web/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
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
	userData := decryptUserData(encryptedUserData, "123")
	fmt.Println(userData)

	//TODO F: 添加空值校验和账号密码校验（抓取登陆返回值），暂不对API KEY校验。
	f, err := os.Create("file.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file
	defer f.Close()
	_, err = f.WriteString(userData.ApiKey + "\n")
	if err != nil {
		log.Fatal(err)
	}
}
