package router

import (
	"Codeforces-ContestCodeDownload/src-web/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	//router.Static("/statics", "./statics")
	router.LoadHTMLGlob("templates/*")
	router.Any("/", handler.Index)
	router.POST("/auth", handler.SaveCodeforcesConfig)
	return router
}
