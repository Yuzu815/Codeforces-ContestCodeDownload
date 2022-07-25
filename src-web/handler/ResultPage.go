package handler

import (
	"Codeforces-ContestCodeDownload/src-web/cores"
	"Codeforces-ContestCodeDownload/src-web/logMode"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func ResultPage(context *gin.Context) {
	checkTaskProcess(context.Copy())
	UID, _ := context.Cookie("UID")
	context.HTML(http.StatusOK, "ResultPage.gohtml", gin.H{
		"title":      "Result Page",
		"MissionUID": UID,
	})
}

func checkTaskProcess(context *gin.Context) {
	UID, _ := context.Cookie("UID")
	logMode.GetLogMap(UID).Infoln("/result checking process...")
	processVal, processOk := cores.MissionProgressMap.Load(UID)
	for processOk == false || processVal.(float64) < 1.0 {
		time.Sleep(time.Second)
		processVal, processOk = cores.MissionProgressMap.Load(UID)
	}
	logMode.GetLogMap(UID).Infoln("/result checking process over...")
}
