package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iden3/go-circuits"
	"github.com/iden3/iden3comm/protocol"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"zkSnacks/walletSDK"
)

func main() {
	config, _ := walletSDK.GetConfig("../holder/config.yaml")
	identity, _ := walletSDK.GetIdentity("../holder/account.json")

	router := gin.Default()
	router.POST("/api/v1/addClaim", addClaim(identity, config))
	router.POST("/api/v1/requestProof", requestProof(identity, config))
	router.POST("/api/v1/getClaims", getClaims(identity, config))
	router.GET("/api/v1/getAccount", getAccount(identity))
	router.GET("/api/v1/getCurrentState", getCurrentState(config, identity))

	router.Run("localhost:8080")
}

func readConfig(file string) (*walletSDK.Config, error) {
	yfile, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open the config file.")
	}

	config := new(walletSDK.Config)
	err = yaml.Unmarshal(yfile, config)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal the yaml file.")
	}
	return config, nil
}

func getAccount(identity *walletSDK.Identity) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, identity)
	}
	return gin.HandlerFunc(fn)
}

func getCurrentState(config *walletSDK.Config, identity *walletSDK.Identity) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		state, err := walletSDK.GetCurrentState(config, identity.ID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Something went wrong! Failed to get IDS from smart contract")
		} else {
			c.IndentedJSON(http.StatusOK, state)
		}
	}
	return gin.HandlerFunc(fn)
}

func addClaim(identity *walletSDK.Identity, config *walletSDK.Config) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var newClaim walletSDK.ClaimAPI

		if err := c.BindJSON(&newClaim); err != nil {
			c.IndentedJSON(http.StatusBadRequest, "Error while parsing claimAPI JSON object")
		} else {
			err := identity.AddClaim(newClaim, config)
			if err != nil {
				log.Printf("Failed to create new claim. Err %s\n", err)
				c.IndentedJSON(http.StatusInternalServerError, "Something went wrong! Failed to create new claim")
			}
			err = walletSDK.DumpIdentity(identity)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, "Something went wrong! Failed to update account file")
			} else {
				c.IndentedJSON(http.StatusCreated, identity)
			}
		}
	}
	return gin.HandlerFunc(fn)
}

func requestProof(identity *walletSDK.Identity, config *walletSDK.Config) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var request protocol.AuthorizationRequestMessage

		if err := c.BindJSON(&request); err != nil {
			fmt.Println("Error while parsing AuthorizationRequestMessage JSON object. Err: ", err)
			c.IndentedJSON(http.StatusInternalServerError, err)
		} else {
			if resp, err := identity.ProofRequest(request, config); err == nil {
				c.IndentedJSON(http.StatusCreated, resp)
			} else {
				log.Printf("Failed to process proof request. Err %s\n", err)
				c.IndentedJSON(http.StatusInternalServerError, "Something went wrong! Failed to generate proof")
			}
		}
	}
	return gin.HandlerFunc(fn)
}

func sendRequestToIssuerToGetClaims(identity *walletSDK.Identity, config *walletSDK.Config) ([]circuits.Claim, error) {
	postBody, _ := json.Marshal(map[string]string{
		"id":    identity.ID.String(),
		"token": "fe7d9c51-5dcf-46dd-8bbc-ae9a0b716ee3",
	})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(config.Issuer.URL+"/api/v1/issueClaim", "application/json", responseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var claims []circuits.Claim
	err = json.Unmarshal(body, &claims)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshaling data from response.")
	}
	return claims, nil
}

func getClaims(identity *walletSDK.Identity, config *walletSDK.Config) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		claims, err := sendRequestToIssuerToGetClaims(identity, config)
		if err != nil {
			fmt.Println(err)
			c.IndentedJSON(http.StatusCreated, err)
			return
		}
		if err := identity.AddClaimsFromIssuer(claims); err != nil {
			log.Printf("Error while adding issued claim to the wallet. Err %s\n", err)
			c.IndentedJSON(http.StatusInternalServerError, "Failed to add issued claim to the wallet")
		} else {
			err = walletSDK.DumpIdentity(identity)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, "Something went wrong! Failed to update account file")
			} else {
				c.IndentedJSON(http.StatusCreated, identity)
			}
		}
	}
	return gin.HandlerFunc(fn)
}
