package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iden3/go-circuits"

	auth "github.com/iden3/go-iden3-auth"
	"github.com/iden3/go-iden3-auth/loaders"
	"github.com/iden3/go-iden3-auth/pubsignals"
	"github.com/iden3/go-iden3-auth/state"
	"github.com/iden3/iden3comm/protocol"
)

var idToQueryInfo map[string]pubsignals.Query
var sessionIDToVerificationReqMap map[uint64]protocol.AuthorizationRequestMessage

const VerifierID = "1125GJqgw6YEsKFwj63GY87MMxPL9kwDKxPUiwMLNZ"

var VerifierHost = "https://a488-205-215-243-16.ngrok.io"

const CallbackURL = "/api/v1/callback"

// Currently fixing the query directly in the code, will be changed to a dynamic query in the future
func Init() {
	idToQueryInfo = make(map[string]pubsignals.Query)
	sessionIDToVerificationReqMap = make(map[uint64]protocol.AuthorizationRequestMessage)
	idToQueryInfo["1"] = pubsignals.Query{
		AllowedIssuers: []string{"*"},
		Req: map[string]interface{}{
			"birthDay": map[string]interface{}{
				"$lt": 20100101,
			},
		},
		Schema: protocol.Schema{
			URL:  "https://raw.githubusercontent.com/pratik1998/jcard-plus-schema-holder/master/claim-schemas/student-age.json-ld",
			Type: "AgeCredential",
		},
	}
	VerifierHost = os.Getenv("DOMAIN_NAME")
}

func main() {
	// Authenticate()

	// Populate the map with the query
	Init()

	router := gin.Default()
	router.LoadHTMLGlob("./static/index.html")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.Static("/static", "./static")
	router.GET("/api/v1/sign-in", generateQR())
	router.GET("/api/v1/viewQuery", viewQuery())
	router.GET("/api/v1/requestVerificationQuery", requestVerificationQuery())
	router.POST("/api/v1/callback", authenticateCallback())

	router.Run("0.0.0.0:9090")
}

func generateQR() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		responseData := map[string]interface{}{
			"url": VerifierHost + "/api/v1/requestVerificationQuery?queryId=1",
		}
		c.IndentedJSON(http.StatusOK, responseData)
	}
	return gin.HandlerFunc(fn)
}

func viewQuery() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		query := c.Request.URL.Query()
		id := query.Get("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "id is required",
			})
			return
		}
		queryInfo, ok := idToQueryInfo[id]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "id not found",
			})
			return
		}
		responseData := map[string]interface{}{
			"query":      queryInfo,
			"verifierId": VerifierID,
		}
		c.IndentedJSON(http.StatusOK, responseData)
	}
	return gin.HandlerFunc(fn)
}

func requestVerificationQuery() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		query := c.Request.URL.Query()
		queryId := query.Get("queryId")
		senderId := query.Get("senderId")
		if queryId == "" || senderId == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "queryId and senderId query parameter is required",
			})
			return
		}
		queryInfo, ok := idToQueryInfo[queryId]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "queryId not found",
			})
			return
		}
		sessionID := rand.Uint64()
		callbackUri := fmt.Sprintf("%s%s?sessionId=", VerifierHost, CallbackURL) + strconv.FormatUint(sessionID, 10)

		var request protocol.AuthorizationRequestMessage
		// Generate request for basic authentication
		request = auth.CreateAuthorizationRequestWithMessage("To do adult stuff", "Age is above 18", VerifierID, callbackUri)
		request.To = senderId

		// Add request for a specific proof
		var mtpProofRequest protocol.ZeroKnowledgeProofRequest
		mtpProofRequest.ID = 1
		mtpProofRequest.CircuitID = string(circuits.AtomicQuerySigCircuitID)
		mtpProofRequest.Rules = map[string]interface{}{
			"query": queryInfo,
		}

		request.Body.Scope = append(request.Body.Scope, mtpProofRequest)

		// Add request to map serve it later when the callback is received
		sessionIDToVerificationReqMap[sessionID] = request

		c.IndentedJSON(http.StatusOK, request)
	}
	return gin.HandlerFunc(fn)
}

func authenticateCallback() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		query := c.Request.URL.Query()
		sessionID := query.Get("sessionId")
		if sessionID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "sessionId is required",
			})
			return
		}
		var response protocol.AuthorizationResponseMessage
		if err := c.BindJSON(&response); err != nil {
			fmt.Println("Error while parsing AuthorizationResponseMessage JSON object. Err: ", err)
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		sessionIDInt, err := strconv.ParseUint(sessionID, 10, 64)
		if err != nil {
			fmt.Println("Error while parsing sessionId. Err: ", err)
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		request, ok := sessionIDToVerificationReqMap[sessionIDInt]
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "sessionId not found",
			})
			return
		}
		verified := Verify(response, request)
		status := "failed"
		if verified {
			status = "success"
		}

		resp := map[string]interface{}{
			"status": status,
		}
		c.IndentedJSON(http.StatusCreated, resp)
	}
	return gin.HandlerFunc(fn)
}

func sendRequestWallet(postBody []byte, request protocol.AuthorizationRequestMessage) {
	requestBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("http://localhost:8080/api/v1/requestProof", "application/json", requestBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var res protocol.AuthorizationResponseMessage
	err = json.Unmarshal(respBody, &res)
	if err != nil {
		log.Fatalf("Failed to marshal response to Auth Type")
	}
	Verify(res, request)
}

func Authenticate() {
	// Audience is verifier id
	rURL := "http://localhost:8080"
	sessionID := 1
	CallbackURL := "/api/callback"
	Audience := "1125GJqgw6YEsKFwj63GY87MMxPL9kwDKxPUiwMLNZ"

	uri := fmt.Sprintf("%s%s?sessionId=%s", rURL, CallbackURL, strconv.Itoa(sessionID))

	var request protocol.AuthorizationRequestMessage

	// Generate request for basic authentication
	request = auth.CreateAuthorizationRequestWithMessage("test flow", "message to sign", Audience, uri)

	request.ID = "7f38a193-0918-4a48-9fac-36adfdb8b542"
	request.ThreadID = "7f38a193-0918-4a48-9fac-36adfdb8b542"

	// Add request for a specific proof
	var mtpProofRequest protocol.ZeroKnowledgeProofRequest
	mtpProofRequest.ID = 1
	mtpProofRequest.CircuitID = string(circuits.AtomicQuerySigCircuitID)
	mtpProofRequest.Rules = map[string]interface{}{
		"query": pubsignals.Query{
			AllowedIssuers: []string{"*"},
			Req: map[string]interface{}{
				"birthDay": map[string]interface{}{
					"$lt": 20100101,
				},
			},
			Schema: protocol.Schema{
				URL:  "https://raw.githubusercontent.com/pratik1998/jcard-plus-schema-holder/master/claim-schemas/student-age.json-ld",
				Type: "AgeCredential",
			},
		},
	}

	request.Body.Scope = append(request.Body.Scope, mtpProofRequest)

	jsonBytes, _ := json.Marshal(request)
	sendRequestWallet(jsonBytes, request)
}

func Verify(response protocol.AuthorizationResponseMessage, request protocol.AuthorizationRequestMessage) bool {

	// Add Polygon RPC node endpoint - needed to read on-chain state
	ethURL := "https://rpc-mumbai.matic.today"

	// Add identity state contract address
	contractAddress := "0x87B36cE5393D4ea6EEf3eb7b1ca6aAd7ae295D4F"

	// Locate the directory that contains circuit's verification keys
	keyDIR := "./keys"

	// load the verifcation key
	var verificationKeyloader = &loaders.FSKeyLoader{Dir: keyDIR}
	resolver := state.ETHResolver{
		RPCUrl:   ethURL,
		Contract: contractAddress,
	}

	// EXECUTE VERIFICATION
	verifier := auth.NewVerifier(verificationKeyloader, loaders.DefaultSchemaLoader{IpfsURL: "ipfs.io"}, resolver)
	err := verifier.VerifyAuthResponse(context.Background(), response, request)
	if err != nil {
		log.Fatalf("Failed to verify %s", err)
		return false
	}

	fmt.Println("Successfully authenticated")
	return true
}
