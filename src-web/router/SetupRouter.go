package router

import (
	"Codeforces-ContestCodeDownload/src-web/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	//TODO F: 添加/download，提供一个下载链接，下载代码压缩包。链接需要经过签名，签名后60分钟内下载有效，过期自动删除文件保护隐私。
	//TODO F: 添加CodeforcesAPI不可用时的检测
	router := gin.Default()
	router.Static("/statics", "./statics")
	router.LoadHTMLGlob("templates/*")
	router.Any("/", handler.IndexPage)
	router.POST("/auth", handler.CodeforcesUserAuth)
	//TODO F: 为Context添加超时功能，若某一页面被人挂起可能会导致通道缓冲区满阻塞崩溃
	router.Any("/process", handler.GetProgressInterface)
	//router.Any("/result", handler.ProgressDisplay, handler.ResultPage)
	router.Any("/result", handler.ProgressDisplay)
	router.Any("/download", handler.ResultPage)
	router.Any("/download/:UID", handler.ResultDownload)
	router.Any("/result/realtime_ws", handler.WebSocketRealTimeInfo)
	return router
}
