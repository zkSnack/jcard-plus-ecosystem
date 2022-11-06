package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type IssueClaimsBody struct {
	Token string `json:"token"`
	ID    string `json:"id"`
}

func main() {
	LoadStudentInfo()
	router := gin.Default()

	router.POST("/api/v1/issueClaim", issueClaim)

	router.Run("localhost:8090")
}

func issueClaim(c *gin.Context) {

	var jsonBody IssueClaimsBody
	c.BindJSON(&jsonBody)

	claims := IssueClaims(jsonBody.Token, jsonBody.ID)
	c.IndentedJSON(http.StatusOK, claims)
}
