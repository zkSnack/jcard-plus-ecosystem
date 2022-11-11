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

	jhuIssuer := NewIssuer()

	router := gin.Default()

	router.POST("/api/v1/issueClaim", issueClaims(jhuIssuer))

	router.Run("localhost:8090")
}

func issueClaims(jhuIssuer *Issuer) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var body IssueClaimsBody
		c.BindJSON(&body)

		ageClaimAPI := generateAgeClaim(body.Token, body.ID)
		jhuIssuer.IssueClaim(*ageClaimAPI)
		claims := jhuIssuer.GetIssuedClaims()
		c.IndentedJSON(http.StatusOK, claims)
	}
	return gin.HandlerFunc(fn)
}
