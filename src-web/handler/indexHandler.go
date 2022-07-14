package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func IndexPage(context *gin.Context) {
	context.HTML(http.StatusOK, "index.gohtml", gin.H{
		"title": "Login Page",
	})
}
