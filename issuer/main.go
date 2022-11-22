package main

import (
	"net/http"
	"zkSnacks/issuerSDK"
	"zkSnacks/walletSDK"

	"github.com/gin-gonic/gin"
)

type IssueClaimsBody struct {
	Token string `json:"token"`
	ID    string `json:"id"`
}

func main() {
	loadStudentInfo()

	jhuIssuer := issuerSDK.NewIssuer()

	router := gin.Default()
	router.POST("/api/v1/issueClaim", issueClaim(jhuIssuer))
	router.GET("/api/v1/getCurrentState", getCurrentState(jhuIssuer.Config, jhuIssuer.Identity))

	router.Run("localhost:8090")
}

func issueClaim(jhuIssuer *issuerSDK.Issuer) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var body IssueClaimsBody
		c.BindJSON(&body)

		ageClaimAPI, err := generateAgeClaim(jhuIssuer.Config, body.ID, body.Token)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		} else {
			_, err := jhuIssuer.IssueClaim(*ageClaimAPI)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			}

			claims := jhuIssuer.GetIssuedClaims(ageClaimAPI.SubjectID)
			c.IndentedJSON(http.StatusOK, claims)
		}
	}
	return gin.HandlerFunc(fn)
}

func getCurrentState(config *walletSDK.Config, identity *walletSDK.Identity) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		state, err := walletSDK.GetCurrentState(config, identity.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		} else {
			c.IndentedJSON(http.StatusOK, state)
		}
	}
	return gin.HandlerFunc(fn)
}
