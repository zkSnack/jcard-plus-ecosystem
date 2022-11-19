package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math/big"
	"time"
	"zkSnacks/issuerSDK"
	"zkSnacks/walletSDK"
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

var idToStudentInfo map[string]*Student

func loadStudentInfo() {
	idToStudentInfo = make(map[string]*Student)
	var students []Student

	content, err := ioutil.ReadFile("../data/students.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	err = json.Unmarshal(content, &students)
	if err != nil {
		log.Fatal("Error in unmarshaling student json data: ", err)
	}

	for _, student := range students {
		idToStudentInfo[student.Token] = &student
	}

	log.Printf("Successfully loaded %d students data\n", len(students))
}

func getStudentInfoBy(token string) (*Student, error) {
	if idToStudentInfo == nil {
		loadStudentInfo()
	}
	if idToStudentInfo == nil {
		return nil, errors.New("Failed to load student info")
	} else if idToStudentInfo[token] == nil {
		return nil, errors.New("Failed to find student associated with token")
	}
	return idToStudentInfo[token], nil
}

func generateAgeClaim(holderID string, token string) (*walletSDK.ClaimAPI, error) {
	studentInfo, err := getStudentInfoBy(token)
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
