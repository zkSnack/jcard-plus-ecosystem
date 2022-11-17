package issuerSDK

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
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

func NewIssuer() *Issuer {
	config, _ := walletSDK.GetConfig("./config.yaml")
	identity, _ := walletSDK.GetIdentity("./account.json")

	issuer := Issuer{
		Config:       config,
		Identity:     identity,
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
