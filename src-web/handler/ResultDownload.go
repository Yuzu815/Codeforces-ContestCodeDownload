package handler

import (
	"Codeforces-ContestCodeDownload/src-web/logserver"
	"github.com/gin-gonic/gin"
)

func ResultDownload(context *gin.Context) {
	fileUID := context.Param("UID")
	context.File("temp/" + fileUID + ".zip")
	logserver.GetLogMap(fileUID).Infoln(fileUID + "zip Downing...")
}
