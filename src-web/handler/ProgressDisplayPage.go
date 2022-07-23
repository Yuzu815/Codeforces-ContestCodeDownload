package handler

import (
	"Codeforces-ContestCodeDownload/src-web/cores"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func ProgressDisplay(context *gin.Context) {
	context.HTML(http.StatusOK, "ProgressDisplayPage.gohtml", gin.H{})
	go RedirectDownloadPage(context)
}

func RedirectDownloadPage(context *gin.Context) {
	val, ok := cores.PROCESS.Load(cores.RandomTaskName)
	for ok == false || val.(float64) < 1.0 || context.Value("CodeforcesResult") == nil {
		//TODO F: 后期进行修正，不考虑硬写入时间的方案
		time.Sleep(time.Second)
		val, ok = cores.PROCESS.Load(cores.RandomTaskName)
	}
	context.Redirect(http.StatusMovedPermanently, "/download")
}
