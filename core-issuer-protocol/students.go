package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
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

// TO-DO: Add error handling
func LoadStudentInfo() {
	idToStudentInfo = make(map[string]Student)
	var students []Student

	content, err := ioutil.ReadFile("../data/students.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	_ = json.Unmarshal(content, &students)

	for _, student := range students {
		idToStudentInfo[student.Token] = student
	}

}

func getStudentInfoByToken(token string) Student {
	return idToStudentInfo[token]
}
