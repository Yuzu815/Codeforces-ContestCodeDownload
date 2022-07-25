package handler

import (
	"Codeforces-ContestCodeDownload/src-web/logMode"
	"github.com/gin-gonic/gin"
)

func ResultDownload(context *gin.Context) {
	fileUID := context.Param("UID")
	context.File("temp/" + fileUID + ".zip")
	logMode.GetLogMap(fileUID).Infoln(fileUID + "zip Downing...")
}
