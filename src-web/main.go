package main

import (
	"Codeforces-ContestCodeDownload/src-web/router"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)
	RouterServer := router.SetupRouter()
	_ = RouterServer.Run(":8080")
}
