package issuerSDK

import (
	"context"
	"log"
	"math/big"

	circuits "github.com/iden3/go-circuits"
	merkletree "github.com/iden3/go-merkletree-sql/v2"

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

func NewIssuer() *Issuer {
	config, _ := walletsdk.GetConfig("../holder/config.yaml")
	identity, _ := walletsdk.GetIdentity("../issuer/account.json")

	issuer := Issuer{
		Identity:     identity,
		Config:       config,
		IssuedClaims: make(map[string]circuits.Claim),
	}

	return &issuer
}

func (i *Issuer) IssueClaim(claim walletsdk.ClaimAPI) *circuits.Claim {
	// Get core claim from Claim API
	claimToAdd := walletsdk.CreateIden3ClaimFromAPI(claim)
	err := i.Identity.AddClaim(claim, i.Config)
	if err != nil {
		log.Fatalf("Error %s. Failed to add claim. Aborting...", err)
	}
	err = walletsdk.DumpIdentity(i.Identity)
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
