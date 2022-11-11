package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/big"
	"os"

	yaml "gopkg.in/yaml.v2"

	circuits "github.com/iden3/go-circuits"
	merkletree "github.com/iden3/go-merkletree-sql/v2"
	"github.com/pkg/errors"

	"zkSnacks/walletsdk"
)

type Issuer struct {
	Config       *walletsdk.Config         `json:"config"`
	Identity     *walletsdk.Identity       `json:"identity"`
	IssuedClaims map[string]circuits.Claim `json:"issued_claims"`
}

const (
	CLAIM_SCHEMA_ROOT_DIR = "../claim-schemas/"
)

// Duplicate code from holder. Didn't want to update HTTP endpoint for my needs
// TO-DO: Use logic from holder package instead of this? Might move code to utils package?
func readConfig(filename string) (*walletsdk.Config, error) {
	yfile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to open the config file.")
	}

	config := new(walletsdk.Config)
	err = yaml.Unmarshal(yfile, config)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal the yaml file.")
	}
	return config, nil
}

// Duplicate code from holder.  Didn't want to update HTTP endpoint for my needs
func generateAccount() (*walletsdk.Identity, error) {
	if identity, err := walletsdk.NewIdentity(); err == nil {
		err = dumpIdentity(identity)
		if err != nil {
			return nil, err
		}
		return identity, nil
	} else {
		return nil, errors.Wrap(err, "Failed to create new identity")
	}
}

// Duplicate code from holder.  Didn't want to update HTTP endpoint for my needs
func dumpIdentity(identity *walletsdk.Identity) error {
	file, err := json.MarshalIndent(identity, "", "	")
	if err != nil {
		return errors.Wrap(err, "Failed to json MarshalIdent identity struct")
	}
	err = ioutil.WriteFile("account.json", file, 0644)
	if err != nil {
		return errors.Wrap(err, "Failed to write identity state to the file")
	}
	log.Println("Account.json updated to latest identity state")
	return nil
}

func NewIssuer() *Issuer {
	config, err := readConfig("../holder/config.yaml")
	if err != nil {
		log.Fatalf("Failed while loading config file. Error %s", err)
	}
	var identity *walletsdk.Identity
	if _, err := os.Stat("./account.json"); err == nil {
		identity, err = walletsdk.LoadIdentityFromFile("./account.json")
		if err != nil {
			log.Fatalf("Failed to load identity from the File. Err %s", err)
		}
		log.Println("Account loaded from saved file: account.json")
	} else {
		identity, err = generateAccount()
		if err != nil {
			log.Fatalf("Error %s. Failed to create new identity. Aborting...", err)
		}
	}
	return &Issuer{
		Identity:     identity,
		Config:       config,
		IssuedClaims: make(map[string]circuits.Claim),
	}
}

func (i *Issuer) IssueClaim(claim walletsdk.ClaimAPI) *circuits.Claim {
	// Get core claim from Claim API
	claimToAdd := walletsdk.CreateIden3ClaimFromAPI(claim)
	err := i.Identity.AddClaim(claim, i.Config)
	if err != nil {
		log.Fatalf("Error %s. Failed to add claim. Aborting...", err)
	}
	err = dumpIdentity(i.Identity)
	if err != nil {
		log.Fatalf("Error %s. Failed to dump File. Aborting...", err)
	}
	hIndexClaim, hValueClaim, _ := claimToAdd.HiHv()
	claimHash, err := merkletree.HashElems(hIndexClaim, hValueClaim)
	if err != nil {
		log.Fatalf("Error %s. Failed to hash claim. Aborting...", err)
	}

	// Generate proof of claim
	claimProof, _, err := i.Identity.Clt.GenerateProof(context.Background(), hIndexClaim, i.Identity.Clt.Root())
	if err != nil {
		log.Fatalf("Error %s. Failed to generate MTP for claim. Aborting...", err)
	}

	claimRevNonce := new(big.Int).SetUint64(claimToAdd.GetRevocationNonce())
	proofNotRevoke, _, err := i.Identity.Ret.GenerateProof(context.Background(), claimRevNonce, i.Identity.Ret.Root())
	if err != nil {
		log.Fatalf("Error %s. Failed to generate revocation MTP for claim. Aborting...", err)
	}

	// Sign claim
	claimSignature := i.Identity.PrivateKey.SignPoseidon(claimHash.BigInt())

	// Generate circuit.Claim
	issuerAuthClaimMTP := i.Identity.GetUserAuthClaim()
	currentTreeState := i.Identity.GetTreeState()

	claimIssuerSignature := circuits.BJJSignatureProof{
		IssuerID:              i.Identity.ID,
		IssuerTreeState:       issuerAuthClaimMTP.TreeState,
		IssuerAuthClaimMTP:    issuerAuthClaimMTP.Proof,
		Signature:             claimSignature,
		IssuerAuthClaim:       issuerAuthClaimMTP.Claim,
		IssuerAuthNonRevProof: *issuerAuthClaimMTP.NonRevProof,
	}

	signedClaim := circuits.Claim{
		Claim:     claimToAdd,
		Proof:     claimProof,
		TreeState: currentTreeState,
		IssuerID:  i.Identity.ID,
		NonRevProof: &circuits.ClaimNonRevStatus{
			TreeState: currentTreeState,
			Proof:     proofNotRevoke,
		},
		SignatureProof: claimIssuerSignature,
	}

	// Use better key than claim hash
	i.IssuedClaims[claim.ClaimSchemaHashHex] = signedClaim

	return &signedClaim
}

// Returns all the claims
// TO-DO: Only return claims associated with particular holder
func (i *Issuer) GetIssuedClaims() []circuits.Claim {
	var claims []circuits.Claim
	for _, claim := range i.IssuedClaims {
		claims = append(claims, claim)
	}
	return claims
}
