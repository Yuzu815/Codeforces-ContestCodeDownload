package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func ResultDownload(context *gin.Context) {
	fileName := context.Param("TASK_UID")
	fmt.Println("DEBUG: " + fileName)
	context.File("temp/" + fileName + ".zip")
}
