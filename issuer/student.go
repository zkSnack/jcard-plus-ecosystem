package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"strconv"
	"time"
	"zkSnacks/walletSDK"

	core "github.com/iden3/go-iden3-core"

	"github.com/gofrs/uuid"
	verifiable "github.com/iden3/go-schema-processor/verifiable"

	"github.com/pkg/errors"
)

type Student struct {
	JHED_ID   string `json:"jhed_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Gender    string `json:"gender"`
	BirthDate string `json:"birth_date"`
	Degree    string `json:"degree"`
	Program   string `json:"program"`
	Token     string `json:"token"`
}

type CredentialSchema struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

var idToStudentInfo map[string]Student

func loadStudentInfo(config *walletSDK.Config) error {
	// Read database
	content, err := ioutil.ReadFile(config.Issuer.DataDirectory + "/students.json")
	if err != nil {
		return err
	}

	// Serialize data
	var students []Student
	err = json.Unmarshal(content, &students)
	if err != nil {
		return err
	}

	// Init data store
	idToStudentInfo = make(map[string]Student)
	for _, student := range students {
		idToStudentInfo[student.Token] = student
	}

	log.Printf("Successfully loaded %d students data\n", len(students))
	return nil
}

func getStudentInfoByToken(token string) (*Student, error) {
	if val, ok := idToStudentInfo[token]; ok {
		return &val, nil
	}
	return nil, errors.New("Failed to find student associated with token.")
}

func generateAgeClaim(config *walletSDK.Config, holderID string, token string) (*walletSDK.ClaimAPI, error) {
	studentInfo, err := getStudentInfoByToken(token)
	if err != nil {
		return nil, err
	}

	birthday := new(big.Int)
	birthday.SetString(studentInfo.BirthDate, 10)

	claimAPI := walletSDK.ClaimAPI{
		SubjectID:      holderID,
		ClaimSchema:    "https://raw.githubusercontent.com/pratik1998/jcard-plus-schema-holder/master/claim-schemas/student-age.json-ld",
		CredentialType: "AgeCredential",
		IndexSlotA:     birthday,
		ExpirationDate: time.Now().Unix() + 31536000,
	}

	return &claimAPI, nil
}

func generateAgeClaimV2(config *walletSDK.Config, holderID string, token string) (*verifiable.Iden3Credential, error) {
	id, err := core.IDFromString(holderID)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}
	did, err := core.NewDID(id.String(), core.WithNetwork("polygon", "mumbai"))
	if err != nil {
		return nil, err
	}
	fmt.Print("DID: ", did.String())
	studentInfo, err := getStudentInfoByToken(token)
	if err != nil {
		return nil, err
	}
	iden3credentialAPI := verifiable.Iden3Credential{}
	iden3credentialAPI.ID = uuid.Must(uuid.NewV4()).String()
	iden3credentialAPI.Context = []string{
		"https://raw.githubusercontent.com/iden3/claim-schema-vocab/main/schemas/json-ld/iden3credential.json-ld",
		"https://raw.githubusercontent.com/pratik1998/jcard-plus-schema-holder/master/claim-schemas/student-age.json-ld",
	}
	iden3credentialAPI.Type = []string{
		"Iden3Credential",
	}
	expiration_date := time.Unix(time.Now().Unix()+31536000, 0)
	iden3credentialAPI.Expiration = &expiration_date
	iden3credentialAPI.Updatable = false
	iden3credentialAPI.Version = 1
	iden3credentialAPI.RevNonce = rand.Uint64()
	iden3credentialAPI.CredentialSubject = map[string]interface{}{
		"birthDay": studentInfo.BirthDate,
		"id":       holderID,
		"type":     "AgeCredential",
	}
	iden3credentialAPI.CredentialStatus = &verifiable.CredentialStatus{
		ID:   "http://localhost:8080/api/v1/claims/revocation/status/" + strconv.Itoa(int(iden3credentialAPI.RevNonce)),
		Type: "SparseMerkleTreeProof", // Should be using constants but too lazy do work now
	}
	iden3credentialAPI.CredentialSchema.ID = "https://raw.githubusercontent.com/pratik1998/jcard-plus-schema-holder/master/claim-schemas/student-age.json-ld"
	iden3credentialAPI.CredentialSchema.Type = "AgeCredential"

	return &iden3credentialAPI, nil
}
