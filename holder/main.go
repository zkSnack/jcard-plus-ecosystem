package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/cors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iden3/go-iden3-auth/pubsignals"
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

type ProofRequest struct {
	ProofRequestData protocol.AuthorizationRequestMessage `json:"proofRequestData"`
	Status           string                               `json:"status"`
	TimeStamp        time.Time                            `json:"timeStamp"`
}

type ProofRequestQueryBody struct {
	AllowedIssuers []string               `json:"allowedIssuers"`
	SchemaURL      string                 `json:"schemaURL"`
	SchemaHash     string                 `json:"schemaHash"`
	CredentialType string                 `json:"credentialType"`
	Data           map[string]interface{} `json:"data"`
}

type ProofRequestResponseBody struct {
	ID          string                  `json:"id"`
	From        string                  `json:"from"`
	To          string                  `json:"to"`
	Message     string                  `json:"message"`
	Reason      string                  `json:"reason"`
	CallbackURL string                  `json:"callbackURL"`
	Status      string                  `json:"status"`
	TimeStamp   time.Time               `json:"timeStamp"`
	QueryData   []ProofRequestQueryBody `json:"queryData"`
}

var proofRequests []ProofRequest

func main() {
	config, _ := walletSDK.GetConfig("./config.yaml")
	identity, _ := walletSDK.GetIdentity("./account.json")

	router := gin.Default()
	if gin.Mode() == gin.DebugMode {
		corsConfig := cors.DefaultConfig()
		corsConfig.AllowOrigins = []string{"http://localhost:3000"}
		router.Use(cors.New(corsConfig))
	} else if gin.Mode() == gin.ReleaseMode {
		router.LoadHTMLGlob(config.UI.HtmlDir)
		router.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.html", nil)
		})
		router.Static("/static", config.UI.StaticDir)
	}

	router.POST("/api/v1/addClaim", addClaim(identity, config))
	router.POST("/api/v1/requestProof", requestProof(identity, config))
	router.POST("/api/v1/fetchClaimsFromIssuer", fetchClaimsFromIssuer(identity, config)) // Why this endpoint is POST?
	router.GET("/api/v1/getClaims", getClaims(identity, config))
	router.GET("/api/v1/getAccount", getAccount(identity))
	router.GET("/api/v1/getCurrentState", getCurrentState(config, identity))
	router.GET("/api/v1/getAccountInfo", getAccountInfo(identity))
	router.POST("/api/v1/addProofRequest", addProofRequest(identity))
	router.GET("/api/v1/getProofRequests", getProofRequests())
	router.GET("/api/v1/acceptProofRequest", acceptProofRequest(identity, config))

	router.Run("0.0.0.0:8080")
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
		"token": config.Issuer.Token,
	})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(config.Issuer.URL+"/api/v1/issueClaim", "application/json", responseBody)
	if err != nil {
		return nil, errors.Wrap(err, "Failed while submitting claim to the issuer.")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Failed while submitting claim to the issuer.")
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
	claimsResponse := make([]ClaimResponseBody, 0)
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

// TODO: add error handling
func getProofRequestToResponse() []ProofRequestResponseBody {
	var proofRequestResponse []ProofRequestResponseBody
	for _, proofRequest := range proofRequests {
		var proofQueries []ProofRequestQueryBody
		for _, queryReq := range proofRequest.ProofRequestData.Body.Scope {
			rules := queryReq.Rules
			jsonStr, _ := json.Marshal(rules["query"])
			var query pubsignals.Query
			if err := json.Unmarshal(jsonStr, &query); err != nil {
				return nil
			}
			proofQueries = append(proofQueries, ProofRequestQueryBody{
				AllowedIssuers: query.AllowedIssuers,
				SchemaURL:      query.Schema.URL,
				SchemaHash:     walletSDK.GetHashFromClaimSchemaURL(query.Schema.URL, query.Schema.Type),
				CredentialType: query.Schema.Type,
				Data:           query.Req,
			})
		}
		proofRequestResponse = append(proofRequestResponse, ProofRequestResponseBody{
			ID:          proofRequest.ProofRequestData.ID,
			From:        proofRequest.ProofRequestData.From,
			To:          proofRequest.ProofRequestData.To,
			Message:     proofRequest.ProofRequestData.Body.Message,
			Reason:      proofRequest.ProofRequestData.Body.Reason,
			CallbackURL: proofRequest.ProofRequestData.Body.CallbackURL,
			TimeStamp:   proofRequest.TimeStamp,
			Status:      proofRequest.Status,
			QueryData:   proofQueries,
		})
	}
	return proofRequestResponse
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

func addProofRequest(identity *walletSDK.Identity) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var url struct {
			URL string `json:"url"`
		}

		if err := c.BindJSON(&url); err != nil {
			fmt.Println("Error while getting url from JSON object. Err: ", err)
			c.IndentedJSON(http.StatusInternalServerError, err)
		} else {
			resp, err := http.Get(url.URL + "&senderId=" + identity.ID.String())
			if err != nil {
				fmt.Println("Failed while getting query from verifier. Err: ", err)
				c.IndentedJSON(http.StatusInternalServerError, err)
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("Failed while getting query from verifier. Err: ", err)
				c.IndentedJSON(http.StatusInternalServerError, err)
			}

			var request protocol.AuthorizationRequestMessage

			err = json.Unmarshal(body, &request)
			if err != nil {
				fmt.Println("Error while parsing AuthorizationRequestMessage JSON object. Err: ", err)
				c.IndentedJSON(http.StatusInternalServerError, err)
			} else {
				newProofRequest := ProofRequest{
					ProofRequestData: request,
					TimeStamp:        time.Now(),
					Status:           "pending", // Initial status of the proof request
				}
				proofRequests = append(proofRequests, newProofRequest)
				resp := map[string]interface{}{
					"status": "success",
				}
				c.IndentedJSON(http.StatusCreated, resp)
			}
		}
	}
	return gin.HandlerFunc(fn)
}

func getProofRequests() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		responseData := map[string]interface{}{
			"proofRequests": getProofRequestToResponse(),
		}
		c.IndentedJSON(http.StatusCreated, responseData)
	}
	return gin.HandlerFunc(fn)
}

func sendProofToVerifier(verfierCallbackURL string, request protocol.AuthorizationResponseMessage) error {
	postBody, _ := json.Marshal(request)
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(verfierCallbackURL, "application/json", responseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return errors.Wrap(err, "Error unmarshaling data from response.")
	}
	if responseData["status"] == "failed" {
		return errors.Wrap(err, "Verifier rejected the proof.")
	}
	return nil
}

func acceptProofRequest(identity *walletSDK.Identity, config *walletSDK.Config) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		query := c.Request.URL.Query()
		requestId := query.Get("requestId")
		if requestId == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "requestId is required",
			})
			return
		}

		for i, proofRequest := range proofRequests {
			if proofRequest.ProofRequestData.ID == requestId {
				if resp, err := identity.ProofRequest(proofRequest.ProofRequestData, config); err == nil {
					err := sendProofToVerifier(proofRequest.ProofRequestData.Body.CallbackURL, *resp)
					if err != nil {
						c.IndentedJSON(http.StatusBadRequest, gin.H{
							"error": err,
						})
						return
					}
					proofRequests[i].Status = "accepted"
					resp := map[string]interface{}{
						"status": "success",
					}
					c.IndentedJSON(http.StatusCreated, resp)
				} else {
					log.Printf("Failed to process proof request. Err %s\n", err)
					c.IndentedJSON(http.StatusInternalServerError, "Something went wrong! Failed to generate proof")
				}
			}
		}
	}
	return gin.HandlerFunc(fn)
}
