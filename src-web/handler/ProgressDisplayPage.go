package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ProgressDisplay(context *gin.Context) {
	context.HTML(http.StatusOK, "ProgressDisplayPage.gohtml", gin.H{})
}
