package handler

import (
	"Codeforces-ContestCodeDownload/src-web/cores"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func ResultPage(context *gin.Context) {
	contextCopy := context.Copy()
	checkTaskProcess(contextCopy)
	UID, _ := context.Cookie("UID")
	//TODO F: 后端提供结构，实时返回进度，实现进度条&日志返回
	context.HTML(http.StatusOK, "ResultPage.gohtml", gin.H{
		"title":      "Result Page",
		"MissionUID": UID,
	})
}

func checkTaskProcess(context *gin.Context) {
	UID, _ := context.Cookie("UID")
	processVal, processOk := cores.PROCESS.Load(UID)
	for processOk == false || processVal.(float64) < 1.0 {
		//TODO F: 后期进行修正，不考虑硬写入时间的方案
		time.Sleep(time.Second)
		processVal, processOk = cores.PROCESS.Load(UID)
	}
}
