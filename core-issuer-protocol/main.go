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

	issuer := NewIssuer()

	router := gin.Default()

	router.POST("/api/v1/issueClaim", issueClaims(issuer))

	router.Run("localhost:8090")
}

func issueClaims(issuer *Issuer) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var body IssueClaimsBody
		c.BindJSON(&body)

		ageClaimAPI := generateAgeClaim(body.Token, body.ID)
		issuer.IssueClaim(*ageClaimAPI)
		claims := issuer.GetIssuedClaims()
		c.IndentedJSON(http.StatusOK, claims)
	}
	return gin.HandlerFunc(fn)
}
