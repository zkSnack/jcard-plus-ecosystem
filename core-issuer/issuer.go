package issuerSDK

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"

	circuits "github.com/iden3/go-circuits"
	core "github.com/iden3/go-iden3-core"
	merkletree "github.com/iden3/go-merkletree-sql/v2"
	jsonld "github.com/iden3/go-schema-processor/json-ld"
	"github.com/iden3/go-schema-processor/loaders"
	"github.com/iden3/go-schema-processor/processor"
	verifiable "github.com/iden3/go-schema-processor/verifiable"

	"zkSnacks/walletSDK"
)

type Issuer struct {
	Config       *walletSDK.Config                               `json:"config"`
	Identity     *walletSDK.Identity                             `json:"identity"`
	IssuedClaims map[string][]walletSDK.Iden3CredentialClaimBody `json:"issued_claims"`
}

func NewIssuer() *Issuer {
	config, _ := walletSDK.GetConfig("./config.yaml")
	identity, _ := walletSDK.GetIdentity("./account.json")

	issuer := Issuer{
		Config:       config,
		Identity:     identity,
		IssuedClaims: make(map[string][]walletSDK.Iden3CredentialClaimBody),
	}

	return &issuer
}

func (i *Issuer) getClaimToAddV2(iden3credentialAPI verifiable.Iden3Credential) (*core.Claim, error) {
	loader := loaders.HTTP{URL: iden3credentialAPI.CredentialSchema.ID}
	credType := iden3credentialAPI.CredentialSubject["type"].(string)
	// subjectDID, err := core.ParseDID(iden3credentialAPI.CredentialSubject["id"].(string))
	// if err != nil {
	// 	return nil, err
	// }
	parser := jsonld.Parser{ClaimType: credType,
		ParsingStrategy: processor.OneFieldPerSlotStrategy}
	schemaBytes, _, err := loader.Load(context.Background())
	if err != nil {
		return nil, errors.New("Failed to add claim.")
	}
	// Careful: This will remove some fields from the iden3credentialAPI.CredentialSubject
	// https://github.com/iden3/go-schema-processor/blob/main/json-ld/parser.go#L69
	// But for our purpose we don't need them
	claimToAdd, err := parser.ParseClaim(&iden3credentialAPI, schemaBytes)
	if err != nil {
		fmt.Println("Error parsing claim: ", err)
		return nil, errors.New("Failed to parse claim.")
	}
	id, err := claimToAdd.GetID()
	fmt.Println("Claim to add: ", id.String())
	if err != nil {
		return nil, errors.New("Failed to get ID from claim.")
	}

	// if subjectDID.ID != id {
	// 	return nil, errors.New("ID from claim and credential subject do not match.")
	// }
	return claimToAdd, nil
}

func (i *Issuer) IssueClaim(iden3credentialAPI verifiable.Iden3Credential) (*circuits.Claim, error) {
	// Get core claim from Claim API
	claimToAdd, err := i.getClaimToAddV2(iden3credentialAPI)
	if err != nil {
		log.Printf("Failed to add claim: %v\n", err)
		return nil, err
	}

	id, err := claimToAdd.GetID()
	if err != nil {
		fmt.Println("Subject ID not provided. Currently self claim not supported.")
		return nil, err
	}
	subjectID := id.String()

	// TO-DO: Add Core claim should return infromation about blockchain transaction
	// https://github.com/iden3/go-schema-processor/blob/main/verifiable/proof.go#L21
	// It will help to show information about the transaction in the UI
	err = i.Identity.AddCoreClaim(claimToAdd, i.Config)
	if err != nil {
		log.Println("Error while adding claim to identity", err)
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
	claimFullData := walletSDK.Iden3CredentialClaimBody{Iden3credential: iden3credentialAPI, Data: signedClaim}

	i.IssuedClaims[subjectID] = append(i.IssuedClaims[subjectID], claimFullData)

	return &signedClaim, nil
}

// func (i *Issuer) IssueClaim(claim walletSDK.ClaimAPI) (*circuits.Claim, error) {
// 	// Get core claim from Claim API
// 	claimToAdd := walletSDK.CreateIden3ClaimFromAPI(claim)

// 	err := i.Identity.AddClaim(claim, i.Config)
// 	if err != nil {
// 		log.Println("Error while adding claim to identity", err)
// 		return nil, errors.New("Failed to add claim.")
// 	}

// 	err = walletSDK.DumpIdentity(i.Identity)
// 	if err != nil {
// 		return nil, errors.New("Failed to dump file.")
// 	}

// 	hIndexClaim, hValueClaim, _ := claimToAdd.HiHv()
// 	claimHash, err := merkletree.HashElems(hIndexClaim, hValueClaim)
// 	if err != nil {
// 		return nil, errors.New("Failed to to hash claim.")
// 	}

// 	// Generate proof of claim
// 	claimProof, _, err := i.Identity.Clt.GenerateProof(context.Background(), hIndexClaim, i.Identity.Clt.Root())
// 	if err != nil {
// 		return nil, errors.New("Failed to generate MTP for claim.")
// 	}

// 	claimRevNonce := new(big.Int).SetUint64(claimToAdd.GetRevocationNonce())
// 	proofNotRevoke, _, err := i.Identity.Ret.GenerateProof(context.Background(), claimRevNonce, i.Identity.Ret.Root())
// 	if err != nil {
// 		return nil, errors.New("Failed to generate revocation MTP for claim")
// 	}

// 	// Sign claim
// 	claimSignature := i.Identity.PrivateKey.SignPoseidon(claimHash.BigInt())

// 	// Generate circuit.Claim
// 	issuerAuthClaimMTP := i.Identity.GetUserAuthClaim()
// 	currentTreeState := i.Identity.GetTreeState()

// 	claimIssuerSignature := circuits.BJJSignatureProof{
// 		IssuerID:              i.Identity.ID,
// 		IssuerTreeState:       issuerAuthClaimMTP.TreeState,
// 		IssuerAuthClaimMTP:    issuerAuthClaimMTP.Proof,
// 		Signature:             claimSignature,
// 		IssuerAuthClaim:       issuerAuthClaimMTP.Claim,
// 		IssuerAuthNonRevProof: *issuerAuthClaimMTP.NonRevProof,
// 	}

// 	signedClaim := circuits.Claim{
// 		Claim:     claimToAdd,
// 		Proof:     claimProof,
// 		TreeState: currentTreeState,
// 		IssuerID:  i.Identity.ID,
// 		NonRevProof: &circuits.ClaimNonRevStatus{
// 			TreeState: currentTreeState,
// 			Proof:     proofNotRevoke,
// 		},
// 		SignatureProof: claimIssuerSignature,
// 	}

// 	// Assign claim to associate holder
// 	claimFullData := walletSDK.ClaimBody{Header: claim, Data: signedClaim}

// 	i.IssuedClaims[claim.SubjectID] = append(i.IssuedClaims[claim.SubjectID], claimFullData)

// 	return &signedClaim, nil
// }

func (i *Issuer) GetIssuedClaims(holderID string) []walletSDK.Iden3CredentialClaimBody {
	return i.IssuedClaims[holderID]
}
