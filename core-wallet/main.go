package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iden3/go-circuits"
	"github.com/iden3/go-rapidsnark/types"
	"github.com/iden3/iden3comm/packers"
	"github.com/iden3/iden3comm/protocol"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	config := readConfig("config.yaml")
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
	router.POST("/api/v1/getClaims", getClaims(account, config))

	router.Run("localhost:8080")
}

func readConfig(file string) Config {
	yfile, err := ioutil.ReadFile(file)

	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err2 := yaml.Unmarshal(yfile, &config)
	if err2 != nil {
		log.Fatal(err2)
	}
	return config
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
				proof, err := GenerateZkProof("compiled-circuits/credentialAtomicQuerySig", inputJSON)
				if err != nil {
					log.Fatal("Something went wrong", err)
				}
				resp := prepareProofRequestResponse(proof, request, account)
				c.IndentedJSON(http.StatusCreated, resp)
			}
		}
	}

	return gin.HandlerFunc(fn)
}

func prepareProofRequestResponse(proof *types.ZKProof, authReq protocol.AuthorizationRequestMessage, account *Account) protocol.AuthorizationResponseMessage {
	resp := protocol.AuthorizationResponseMessage{
		ID:       authReq.ID,
		Typ:      packers.MediaTypePlainMessage,
		Type:     protocol.AuthorizationResponseMessageType,
		ThreadID: authReq.ThreadID,
		Body: protocol.AuthorizationMessageResponseBody{
			Message: "test",
			Scope: []protocol.ZeroKnowledgeProofResponse{
				{
					ID:        1,
					CircuitID: string(circuits.AtomicQuerySigCircuitID),
					ZKProof:   *proof,
				},
			},
		},
		From: account.ID.String(),
		To:   authReq.From,
	}
	return resp
}

func sendRequestToIssuerToGetClaims(account *Account, config Config) ([]circuits.Claim, error) {
	postBody, _ := json.Marshal(map[string]string{
		"id":    account.ID.String(),
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
	fmt.Println(claims)
	return claims, nil
}

func getClaims(account *Account, config Config) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		claims, err := sendRequestToIssuerToGetClaims(account, config)
		if err != nil {
			fmt.Println(err)
			c.IndentedJSON(http.StatusCreated, err)
			return
		}
		account.Identity.addClaimsFromIssuer(claims)
		c.IndentedJSON(http.StatusCreated, account)
	}

	return gin.HandlerFunc(fn)
}
