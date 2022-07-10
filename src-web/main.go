package main

import "Codeforces-ContestCodeDownload/src-web/router"

func main() {
	router := router.SetupRouter()
	_ = router.Run()
}
