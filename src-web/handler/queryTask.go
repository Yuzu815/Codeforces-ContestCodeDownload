package handler

import (
	"Codeforces-ContestCodeDownload/src-web/cores"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func QueryPage(context *gin.Context) {
	context.HTML(http.StatusOK, "querypage.gohtml", gin.H{})
	for true {
		now := cores.PROCESS[cores.RandomTaskName]
		if now >= 1 {
			break
		}
	}
	fmt.Println(cores.PROCESS[cores.RandomTaskName])
	context.Next()
}
