package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func QueryPage(context *gin.Context) {
	context.HTML(http.StatusOK, "querypage.gohtml", gin.H{})
}
