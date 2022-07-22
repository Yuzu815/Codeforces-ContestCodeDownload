package handler

import (
	"Codeforces-ContestCodeDownload/src-web/cores"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetTask TODO F: 使用更为优雅的方式来传递参数，而不是字符串互相转化。后期添加多用户支持。
func GetTask(context *gin.Context) {
	val, _ := cores.PROCESS.Load(cores.RandomTaskName)
	if val == nil {
		context.JSON(http.StatusOK, gin.H{
			"error": "null",
		})
	} else {
		taskProgress := strconv.FormatFloat(val.(float64), 'f', -1, 64)
		context.JSON(http.StatusOK, gin.H{
			"UID":         cores.RandomTaskName,
			"taskProcess": taskProgress,
		})
	}
}
