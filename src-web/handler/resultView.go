package handler

import (
	"Codeforces-ContestCodeDownload/src-web/cores"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ResultPage(context *gin.Context) {
	temp := context.Value("CodeforcesResult").([]cores.InformationStruct)
	//TODO E: json错误解析处理
	aLittleJson, _ := json.Marshal(temp)
	//TODO F: 后端提供结构，实时返回进度，实现进度条&日志返回
	context.HTML(http.StatusOK, "result.gohtml", gin.H{
		"title":      "Result Page",
		"resultBody": string(aLittleJson),
	})
}
