package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	noAccount := true
	if _, err := os.Stat("./account.json"); err == nil {

		identity := FromFile("./account.json")
		fmt.Println(identity)
		noAccount = false
	}

	router := gin.Default()

	if noAccount {
		router.POST("/api/v1/generate", generateAccount)
	}

	router.Run("localhost:8080")
}

func generateAccount(c *gin.Context) {
	account := NewAccount()
	file, _ := json.MarshalIndent(account, "", "	")
	_ = ioutil.WriteFile("account.json", file, 0644)
	c.IndentedJSON(http.StatusOK, account)
}
