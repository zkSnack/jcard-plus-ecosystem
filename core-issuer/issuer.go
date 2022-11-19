package issuerSDK

import (
	"context"
	"errors"
	"math/big"

	circuits "github.com/iden3/go-circuits"
	merkletree "github.com/iden3/go-merkletree-sql/v2"

	"zkSnacks/walletSDK"
)

type Issuer struct {
	Config       *walletSDK.Config           `json:"config"`
	Identity     *walletSDK.Identity         `json:"identity"`
	IssuedClaims map[string][]circuits.Claim `json:"issued_claims"`
}

const (
	CLAIM_SCHEMA_ROOT_DIR = "../claim-schemas/"
)

func NewIssuer() *Issuer {
	config, _ := walletSDK.GetConfig("../holder/config.yaml")
	identity, _ := walletSDK.GetIdentity("../issuer/account.json")

	issuer := Issuer{
		Identity:     identity,
		Config:       config,
		IssuedClaims: make(map[string][]circuits.Claim),
	}

	return &issuer
}

func (i *Issuer) IssueClaim(claim walletSDK.ClaimAPI) (*circuits.Claim, error) {
	// Get core claim from Claim API
	claimToAdd := walletSDK.CreateIden3ClaimFromAPI(claim)

	err := i.Identity.AddClaim(claim, i.Config)
	if err != nil {
		return nil, errors.New("Failed to add claim.")
	}

	err = walletSDK.DumpIdentity(i.Identity)
	if err != nil {
		return nil, errors.New("Failed to dump file.")
	}

	hIndexClaim, hValueClaim, _ := claimToAdd.HiHv()
	claimHash, err := merkletree.HashElems(hIndexClaim, hValueClaim)
	if err != nil {
		return nil, errors.New("Failed to to hash claim.")
	}

	// Generate proof of claim
	claimProof, _, err := i.Identity.Clt.GenerateProof(context.Background(), hIndexClaim, i.Identity.Clt.Root())
	if err != nil {
		return nil, errors.New("Failed to generate MTP for claim.")
	}

	claimRevNonce := new(big.Int).SetUint64(claimToAdd.GetRevocationNonce())
	proofNotRevoke, _, err := i.Identity.Ret.GenerateProof(context.Background(), claimRevNonce, i.Identity.Ret.Root())
	if err != nil {
		return nil, errors.New("Failed to generate revocation MTP for claim")
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

	// Assign claim to associate holder
	i.IssuedClaims[claim.SubjectID] = append(i.IssuedClaims[claim.SubjectID], signedClaim)

	return &signedClaim, nil
}

func (i *Issuer) GetIssuedClaims(holderID string) []circuits.Claim {
	return i.IssuedClaims[holderID]
}
