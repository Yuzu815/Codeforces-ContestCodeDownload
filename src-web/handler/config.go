package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func saveCodeforcesConfig(context *gin.Context) {
	apiKey := context.PostForm("apiKey")
	apiSecret := context.PostForm("apiSecret")
	usernameOrEmail := context.PostForm("usernameOrEmail")
	password := context.PostForm("password")
	fmt.Println(apiKey, apiSecret, usernameOrEmail, password)
}
