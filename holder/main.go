package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/utils"
	"github.com/iden3/iden3comm/protocol"
	"github.com/pkg/errors"

	"zkSnacks/walletSDK"
)

type ClaimResponseBody struct {
	ID               string                 `json:"id"`
	SchemaURL        string                 `json:"schemaURL"`
	CredentialType   string                 `json:"credentialType"`
	Expiration       *time.Time             `json:"expiration"`
	Updatable        bool                   `json:"updatable"`
	Version          uint32                 `json:"version"`
	RevocationNonce  uint64                 `json:"revocationNonce"`
	RevocationStatus bool                   `json:"revocationStatus"`
	IssuerID         string                 `json:"issuerID"`
	ClaimData        map[string]interface{} `json:"claimData"`
	ClaimRawData     *core.Claim            `json:"claimRawData"`
}

func main() {
	config, _ := walletSDK.GetConfig("./config.yaml")
	identity, _ := walletSDK.GetIdentity("./account.json")

	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	router.Use(cors.New(corsConfig)) // TODO: Remove this in production and allow only specific domains
	router.POST("/api/v1/addClaim", addClaim(identity, config))
	router.POST("/api/v1/requestProof", requestProof(identity, config))
	router.POST("/api/v1/fetchClaimsFromIssuer", fetchClaimsFromIssuer(identity, config)) // Why this endpoint is POST?
	router.GET("/api/v1/getClaims", getClaims(identity, config))
	router.GET("/api/v1/getAccount", getAccount(identity))
	router.GET("/api/v1/getCurrentState", getCurrentState(config, identity))
	router.GET("/api/v1/getAccountInfo", getAccountInfo(identity))

	router.Run("localhost:8080")
}

func getAccount(identity *walletSDK.Identity) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, identity)
	}
	return gin.HandlerFunc(fn)
}

func getAccountInfo(identity *walletSDK.Identity) gin.HandlerFunc {
	fn := func(c *gin.Context) {

		responseData := map[string]interface{}{
			"id":            identity.ID.String(),
			"identityState": identity.IDS,
			"privateKey":    utils.HexEncode(identity.PrivateKey[:]),
		}
		c.IndentedJSON(http.StatusOK, responseData)
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

func sendRequestToIssuerToGetClaims(identity *walletSDK.Identity, config *walletSDK.Config) ([]walletSDK.Iden3CredentialClaimBody, error) {
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
	var claims []walletSDK.Iden3CredentialClaimBody
	err = json.Unmarshal(body, &claims)
	if err != nil {
		return nil, errors.Wrap(err, "Error unmarshaling data from response.")
	}
	return claims, nil
}

func fetchClaimsFromIssuer(identity *walletSDK.Identity, config *walletSDK.Config) gin.HandlerFunc {
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

func convertIden3CredClaimBodyToResponse(claims []walletSDK.Iden3CredentialClaimBody) []ClaimResponseBody {
	var claimsResponse []ClaimResponseBody
	for _, claim := range claims {
		claimsResponse = append(claimsResponse, ClaimResponseBody{
			ID:               claim.Iden3credential.ID,
			SchemaURL:        claim.Iden3credential.CredentialSchema.ID,
			CredentialType:   claim.Iden3credential.CredentialSchema.Type,
			Expiration:       claim.Iden3credential.Expiration,
			Updatable:        claim.Iden3credential.Updatable,
			Version:          claim.Iden3credential.Version,
			RevocationNonce:  claim.Iden3credential.RevNonce,
			RevocationStatus: false, // TODO: get revocation status from Issuer API
			IssuerID:         claim.Data.IssuerID.String(),
			ClaimData:        claim.Iden3credential.CredentialSubject,
			ClaimRawData:     claim.Data.Claim,
		})
	}
	return claimsResponse
}

func getClaims(identity *walletSDK.Identity, config *walletSDK.Config) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		receivedClaims := convertIden3CredClaimBodyToResponse(identity.GetStoredClaims())
		responseData := map[string]interface{}{
			"claims": receivedClaims,
		}
		c.IndentedJSON(http.StatusOK, responseData)
	}
	return gin.HandlerFunc(fn)
}
