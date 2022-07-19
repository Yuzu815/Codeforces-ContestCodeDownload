package main

import (
	"Codeforces-ContestCodeDownload/src-web/router"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.DebugMode)
	router := router.SetupRouter()
	_ = router.Run()
}
