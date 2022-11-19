package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/big"
	"time"
	"zkSnacks/issuerSDK"
	"zkSnacks/walletSDK"

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

var idToStudentInfo map[string]Student

func loadStudentInfo() error {
	// Read database
	content, err := ioutil.ReadFile("../data/students.json")
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
	if idToStudentInfo == nil {
		loadStudentInfo()
	}
	if val, ok := idToStudentInfo[token]; ok {
		return &val, nil
	}
	return nil, errors.New("Fail to find student associated with token")
}

func generateAgeClaim(holderID string, token string) (*walletSDK.ClaimAPI, error) {
	studentInfo, err := getStudentInfoByToken(token)
	if err != nil {
		return nil, err
	}

	claimSchemaHashHex := walletSDK.GetHashFromClaimSchema(issuerSDK.CLAIM_SCHEMA_ROOT_DIR+"student-age.json-ld", "AgeCredential")

	birthday := new(big.Int)
	birthday.SetString(studentInfo.BirthDate, 10)

	claimAPI := walletSDK.ClaimAPI{
		SubjectID:          holderID,
		ClaimSchemaHashHex: claimSchemaHashHex,
		IndexSlotA:         birthday,
		ExpirationDate:     time.Now().Unix() + 31536000,
	}

	return &claimAPI, nil
}
