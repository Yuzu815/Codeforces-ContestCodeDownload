package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func SaveCodeforcesConfig(context *gin.Context) {
	apiKey := context.PostForm("apiKey")
	apiSecret := context.PostForm("apiSecret")
	usernameOrEmail := context.PostForm("usernameOrEmail")
	password := context.PostForm("password")
	fmt.Println(apiKey, apiSecret, usernameOrEmail, password)

	f, err := os.Create("file.txt")
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file
	defer f.Close()
	_, err = f.WriteString(apiKey + apiSecret + usernameOrEmail + password + "\n")
	if err != nil {
		log.Fatal(err)
	}
}
