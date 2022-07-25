package handler

import (
	"github.com/gin-gonic/gin"
)

func ResultDownload(context *gin.Context) {
	fileUID := context.Param("TASK_UID")
	context.File("temp/" + fileUID + ".zip")
}
