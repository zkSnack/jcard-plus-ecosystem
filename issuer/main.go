package main

import (
	"log"
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

		ageClaimAPI, err := generateAgeClaim(body.ID, body.Token)
		if err != nil {
			log.Println("Error when generating age claim: ", err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Error occurred while generating age claim"})
		} else {
			jhuIssuer.IssueClaim(*ageClaimAPI)
			claims := jhuIssuer.GetIssuedClaims()
			c.IndentedJSON(http.StatusOK, claims)
		}
	}
	return gin.HandlerFunc(fn)
}

func getCurrentState(config *walletSDK.Config, identity *walletSDK.Identity) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		state, err := walletSDK.GetCurrentState(config, identity.ID)
		if err != nil {
			log.Println("Error when getting current State: ", err)
			c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "Error occurred getting IDS from Blockchain"})
		} else {
			c.IndentedJSON(http.StatusOK, state)
		}
	}
	return gin.HandlerFunc(fn)
}
