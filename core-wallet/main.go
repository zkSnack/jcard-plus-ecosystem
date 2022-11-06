package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iden3/iden3comm/protocol"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	noAccount := true
	var account *Account
	if _, err := os.Stat("./account.json"); err == nil {

		account = LoadAccountFromFile("./account.json")
		noAccount = false
	}

	router := gin.Default()

	if noAccount {
		router.POST("/api/v1/generate", generateAccount)
	}

	router.POST("/api/v1/addClaim", addClaim(account))
	router.POST("/api/v1/requestProof", requestProof(account))
	router.POST("/api/v1/claim", getClaims(account))

	router.Run("localhost:8080")
}

func generateAccount(c *gin.Context) {
	account := NewAccount()
	file, _ := json.MarshalIndent(account, "", "	")
	_ = ioutil.WriteFile("account.json", file, 0644)
	c.IndentedJSON(http.StatusOK, account)
}

func addClaim(account *Account) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var newClaim ClaimAPI

		if err := c.BindJSON(&newClaim); err != nil {
			log.Println("Error while parsing claim JSON object. Err: ", err)
			return
		}

		inputJSON := toJSON(account.addClaim(newClaim))
		/*proof, err := GenerateZkProof("compiled-circuits/stateTransition", inputJSON)
		if err != nil {
			log.Fatal("Something went wrong", err)
		}*/
		c.IndentedJSON(http.StatusCreated, inputJSON)
	}

	return gin.HandlerFunc(fn)
}

func requestProof(account *Account) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var request protocol.AuthorizationRequestMessage

		if err := c.BindJSON(&request); err != nil {
			fmt.Println("Error while parsing request JSON object. Err: ", err)
			c.IndentedJSON(http.StatusBadRequest, err)
		} else {
			inputBytes, err := account.GenerateProof(request)
			if err != nil {
				c.IndentedJSON(http.StatusBadRequest, err)
			} else {
				inputJSON := toJSON(inputBytes)
				/*proof, err := GenerateZkProof("compiled-circuits/stateTransition", inputJSON)
				if err != nil {
					log.Fatal("Something went wrong", err)
				}*/
				c.IndentedJSON(http.StatusCreated, inputJSON)
			}
		}
	}

	return gin.HandlerFunc(fn)
}

func getClaims(account *Account) gin.HandlerFunc {
	fn := func(c *gin.Context) {

	}

	return gin.HandlerFunc(fn)
}
