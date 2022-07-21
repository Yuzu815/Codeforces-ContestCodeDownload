package handler

import (
	"Codeforces-ContestCodeDownload/src-web/cores"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetTask TODO F: 使用更为优雅的方式来传递参数，而不是字符串互相转化。后期添加多用户支持。
func GetTask(context *gin.Context) {
	val := strconv.FormatFloat(cores.PROCESS[cores.RandomTaskName], 'f', -1, 64)
	context.String(http.StatusOK, val)
}
