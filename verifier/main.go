package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/iden3/go-circuits"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	auth "github.com/iden3/go-iden3-auth"
	"github.com/iden3/go-iden3-auth/loaders"
	"github.com/iden3/go-iden3-auth/pubsignals"
	"github.com/iden3/go-iden3-auth/state"
	"github.com/iden3/iden3comm/protocol"
)

func main() {
	Authenticate()
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
	rURL := "http://localhost:8080/"
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
				URL:  "http://localhost:8000/student-age.json-ld",
				Type: "AgeCredential",
			},
		},
	}

	request.Body.Scope = append(request.Body.Scope, mtpProofRequest)

	jsonBytes, _ := json.Marshal(request)
	sendRequestWallet(jsonBytes, request)
}

func Verify(response protocol.AuthorizationResponseMessage, request protocol.AuthorizationRequestMessage) {

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
		return
	}

	fmt.Println("Successfully authenticated")
}
