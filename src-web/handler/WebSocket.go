package handler

import (
	"Codeforces-ContestCodeDownload/src-web/cores"
	"Codeforces-ContestCodeDownload/src-web/logMode"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var httpUpgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketRealTimeInfo(context *gin.Context) {
	UID, _ := context.Cookie("UID")
	ws, _ := httpUpgrade.Upgrade(context.Writer, context.Request, nil)
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			logMode.GetLogMap(UID).Errorln(err.Error())
		}
	}(ws)
	for {
		var resultMessage string
		missionMapLogRef, OK := cores.TaskMessageChan.Load(UID)
		for OK == false {
			missionMapLogRef, OK = cores.TaskMessageChan.Load(UID)
			continue
		}
		for {
			if len(missionMapLogRef.(chan string)) > 0 {
				resultMessage, _ = <-missionMapLogRef.(chan string)
				err := ws.WriteMessage(websocket.TextMessage, []byte(resultMessage))
				if err != nil {
					logMode.GetLogMap(UID).Errorln(err.Error())
					return
				}
			} else {
				break
			}
		}
	}
}
