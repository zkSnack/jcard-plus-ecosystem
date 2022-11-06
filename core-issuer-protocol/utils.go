package main

import (
	"fmt"
	"os"

	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/keccak256"
)

const CLAIM_SCHEMA_ROOT_DIR = "../claim-schemas/"
const CLAIM_SCHEMA_VOCAB_ROOT_DIR = "../claim-schemas-vocab/"

func generateHashFromClaimSchemaFile(schemaFileName string, credentialType string) string {
	// Check for path injection vulnerbility
	schemaBytes, _ := os.ReadFile(CLAIM_SCHEMA_ROOT_DIR + schemaFileName)

	var sHash core.SchemaHash
	h := keccak256.Hash(schemaBytes, []byte(credentialType))
	copy(sHash[:], h[len(h)-16:])
	sHashHex, _ := sHash.MarshalText()

	fmt.Println("Schema File:", schemaFileName)
	fmt.Println("Schema Credential Type:", credentialType)
	fmt.Println("Schema Hash:", string(sHashHex))

	return string(sHashHex)
}
