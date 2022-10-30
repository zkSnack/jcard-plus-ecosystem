package main

import (
	"encoding/json"
	"log"
)

func toJSON(inputJSON []byte) map[string]interface{} {
	var jsonBlob map[string]interface{}
	if err := json.Unmarshal(inputJSON, &jsonBlob); err != nil {
		log.Fatal("Failed while converting to JSON format.")
	}
	return jsonBlob
}
