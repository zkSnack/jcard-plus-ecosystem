package walletSDK

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"time"

	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/keccak256"
)

type ClaimAPI struct {
	SubjectID          string   `json:"subject_id"`
	ClaimSchema        string   `json:"claim_schema"`
	ClaimSchemaHashHex string   `json:"claim_schema_hash_hex"`
	CredentialType     string   `json:"credential_type"`
	IndexSlotA         *big.Int `json:"index_slot_a"`
	IndexSlotB         *big.Int `json:"index_slot_b"`
	ValueSlotA         *big.Int `json:"value_slot_a"`
	ValueSlotB         *big.Int `json:"value_slot_b"`
	ExpirationDate     int64    `json:"expiration_date"`
}

func CreateIden3ClaimFromAPI(claim ClaimAPI) *core.Claim {
	var schema core.SchemaHash
	if claim.ClaimSchemaHashHex != "" {
		schema, _ = core.NewSchemaHashFromHex(claim.ClaimSchemaHashHex)
	} else if claim.ClaimSchema != "" {
		schema, _ = core.NewSchemaHashFromHex(GetHashFromClaimSchemaURL(claim.ClaimSchema, claim.CredentialType))
	} else {
		log.Fatal("No schema hash or schema file provided")
	}

	var options []core.Option

	options = append(options,
		core.WithExpirationDate(time.Unix(claim.ExpirationDate, 0)),
		core.WithRevocationNonce(rand.Uint64()),
		core.WithIndexDataInts(claim.IndexSlotA, claim.IndexSlotB),
		core.WithValueDataInts(claim.ValueSlotA, claim.ValueSlotB))

	if claim.SubjectID != "" {
		id, _ := core.IDFromString(claim.SubjectID)
		options = append(options, core.WithIndexID(id))
	}

	iden3Claim, _ := core.NewClaim(schema, options...)
	return iden3Claim
}

func GetHashFromClaimSchemaURL(schemaHashURL string, credentialType string) string {
	// create tmp claim schema file
	claimSchema, err := ioutil.TempFile("tmp", "claim-schema-*.json-ld")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(claimSchema.Name())
	DownloadFileFromURL(schemaHashURL, claimSchema)
	err = claimSchema.Close()
	return GetHashFromClaimSchema(claimSchema.Name(), credentialType)
}

func GetHashFromClaimSchema(file string, credentialType string) string {
	schemaBytes, _ := os.ReadFile(file)

	var sHash core.SchemaHash
	h := keccak256.Hash(schemaBytes, []byte(credentialType))
	copy(sHash[:], h[len(h)-16:])
	sHashHex, _ := sHash.MarshalText()
	fmt.Printf("Got %s hash for the schema %s\n", string(sHashHex), file)
	return string(sHashHex)
}

// DownloadFileFromURL TODO: Put validation on this file path.
func DownloadFileFromURL(fullURL string, claimSchemaFile *os.File) {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	// Put content on file
	resp, err := client.Get(fullURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	size, err := io.Copy(claimSchemaFile, resp.Body)
	defer claimSchemaFile.Close()

	fmt.Printf("Downloaded a file %s with size %d\n", claimSchemaFile.Name(), size)
}
