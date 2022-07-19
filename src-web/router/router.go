package router

import (
	"Codeforces-ContestCodeDownload/src-web/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	//TODO F: 添加/error，尝试调用接口将错误信息返回到页面上
	//TODO F: 添加/process，前端轮询代码抓取进度并使用进度条显示
	//TODO F: 添加/work，显示进度条以及实时抓取状态
	//TODO F: 添加/download，提供一个下载链接，链接需要经过签名，签名后60分钟内下载有效，过期自动删除文件保护隐私
	router := gin.Default()
	router.Static("/statics", "./statics")
	router.LoadHTMLGlob("templates/*")
	router.Any("/", handler.IndexPage)
	router.POST("/auth", handler.SaveCodeforcesConfig)
	router.Any("/result", handler.ResultPage)
	return router
}
