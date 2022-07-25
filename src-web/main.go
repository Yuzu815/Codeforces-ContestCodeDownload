package main

import (
	"Codeforces-ContestCodeDownload/src-web/router"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	RouterServer := router.SetupRouter()
	port := portArg()
	fmt.Println("[Listen] PORT" + port)
	_ = RouterServer.Run(port)
}

func portArg() string {
	args := os.Args
	argNum := len(args)
	if argNum == 2 {
		argVal, err := strconv.Atoi(args[1])
		if err == nil && argVal >= 1 && argVal <= 65535 {
			return ":" + args[1]
		}
	}
	return ":8080"
}
