package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func parseErrInfo(context *gin.Context) string {
	err := context.Query("err")
	var returnStr string
	if err == "logErr" {
		returnStr = "Login failed. Please check the account and password are correct!"
	}
	return returnStr
}

func IndexPage(context *gin.Context) {
	context.HTML(http.StatusOK, "Index.gohtml", gin.H{
		"title": "Login Page",
		"error": parseErrInfo(context.Copy()),
	})
}
