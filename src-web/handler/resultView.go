package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ResultPage(context *gin.Context) {
	context.HTML(http.StatusOK, "index.gohtml", gin.H{
		"title":      "Result Page",
		"resultBody": context.PostForm("CodeforcesResult"),
	})
}
