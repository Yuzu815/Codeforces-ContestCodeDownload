package handler

import (
	"Codeforces-ContestCodeDownload/src-web/cores"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetProgressInterface
// TODO F: 使用更为优雅的方式来传递参数，而不是字符串互相转化。
// TODO F: 目前使用全局变量，后期添加多用户支持。
func GetProgressInterface(context *gin.Context) {
	UID, _ := context.Cookie("UID")
	val, _ := cores.MissionProgressMap.Load(UID)
	if val == nil {
		context.JSON(http.StatusOK, gin.H{
			"error": "null",
		})
	} else {
		missionProgress := strconv.FormatFloat((val.(float64))*100, 'f', -1, 64)
		context.JSON(http.StatusOK, gin.H{
			"UID":                    UID,
			"CurrentMissionProgress": missionProgress,
		})
	}
}
