package handler

import (
	"Codeforces-ContestCodeDownload/src-web/cores"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var httpUpgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WS_realTimeInfo(context *gin.Context) {
	ws, _ := httpUpgrade.Upgrade(context.Writer, context.Request, nil)
	defer ws.Close()
	for {
		var resultMessage string
		messageType, messageContext, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		if string(messageContext) == "TEST_CONNECTIVITY" {
			resultMessage = "[CONNECT] " + cores.RandomTaskName
		} else {
			//TODO F: 换用更加优雅的方式实现
			missionMapLogRef, OK := cores.TaskMessageChan.Load(cores.RandomTaskName)
			for OK == false {
				missionMapLogRef, OK = cores.TaskMessageChan.Load(cores.RandomTaskName)
				continue
			}
			//TODO F: 在通道中内容全被取出后，发送TEST_CONNECTIVITY测试不会返回结果，疑似因为此处被阻塞，应添加一个等待上限
			resultMessage = <-missionMapLogRef.(chan string)
		}
		ws.WriteMessage(messageType, []byte(resultMessage))
	}
}
