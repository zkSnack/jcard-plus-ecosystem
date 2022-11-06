package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type IssueClaimsBody struct {
	AuthToken string `json:"authToken"`
	HolderID  string `json:"holderID"`
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

	claims := IssueClaims(jsonBody.AuthToken, jsonBody.HolderID)
	c.IndentedJSON(http.StatusOK, claims)
}
