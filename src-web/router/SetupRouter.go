package router

import (
	"Codeforces-ContestCodeDownload/src-web/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	//TODO F: 添加/error，尝试调用接口将错误信息返回到页面上
	//TODO F: 添加/process，为前端提供轮询接口，返回代码抓取进度，并在前端中并使用进度条显示
	//TODO F: 添加/work，显示进度条以及实时抓取状态
	//TODO F: 添加/download，提供一个下载链接，下载代码压缩包。链接需要经过签名，签名后60分钟内下载有效，过期自动删除文件保护隐私。
	//TODO F: 添加CodeforcesAPI不可用时的检测
	router := gin.Default()
	router.Static("/statics", "./statics")
	router.LoadHTMLGlob("templates/*")
	router.Any("/", handler.IndexPage)
	//router.Any("/result", handler.ResultPage)
	//TODO F: 优化掉Set，Value的传递方法
	router.POST("/auth", handler.CodeforcesUserAuth)
	router.Any("/process", handler.GetProgressInterface)
	//router.Any("/result", handler.ProgressDisplay, handler.ResultPage)
	router.Any("/result", handler.ProgressDisplay)
	router.Any("/download", handler.ResultPage)
	router.Any("/download/:TASK_UID", handler.ResultDownload)
	router.Any("/realtime_ws", handler.WS_realTimeInfo)
	return router
}
