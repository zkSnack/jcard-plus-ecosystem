package walletSDK

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"time"

	"github.com/iden3/go-circuits"
	"github.com/iden3/go-iden3-auth/loaders"
	"github.com/iden3/go-iden3-auth/pubsignals"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/go-merkletree-sql/v2/db/memory"
	"github.com/iden3/iden3comm/packers"
	"github.com/iden3/iden3comm/protocol"
	"github.com/pkg/errors"
)

type Identity struct {
	ID             *core.ID                  `json:"id"`
	IDS            *merkletree.Hash          `json:"identity_state"`
	PrivateKey     babyjub.PrivateKey        `json:"private_key"`
	AuthClaim      *core.Claim               `json:"authClaim"`
	Claims         []*core.Claim             `json:"claims"`
	Clt            *merkletree.MerkleTree    `json:"clt"`
	Ret            *merkletree.MerkleTree    `json:"ret"`
	Rot            *merkletree.MerkleTree    `json:"rot"`
	ReceivedClaims map[string]circuits.Claim `json:"received_claims"`
}

type Config struct {
	Issuer struct {
		URL string `yaml:"url"`
		ID  string `yaml:"id"`
	} `yaml:"issuer"`
	Circuits struct {
		Path string `yaml:"path"`
		JS   string `yaml:"js"`
	} `yaml:"circuits"`
	Web3 struct {
		StateTransition string   `yaml:"stateTransition"`
		URL             string   `yaml:"url"`
		PrivateKey      string   `yaml:"privateKey"`
		ChainID         *big.Int `yaml:"chainID"`
	} `yaml:"web3"`
}

func NewIdentity() (*Identity, error) {
	babyJubjubPrivKey := babyjub.NewRandPrivKey()
	babyJubjubPubKey := babyJubjubPrivKey.Public()

	ctx := context.Background()

	authSchemaHash, _ := core.NewSchemaHashFromHex("ca938857241db9451ea329256b9c06e5")
	authClaim, _ := core.NewClaim(authSchemaHash,
		core.WithIndexDataInts(babyJubjubPubKey.X, babyJubjubPubKey.Y),
		core.WithRevocationNonce(rand.Uint64()))

	clt, _ := merkletree.NewMerkleTree(ctx, memory.NewMemoryStorage(), 32)
	ret, _ := merkletree.NewMerkleTree(ctx, memory.NewMemoryStorage(), 32)
	rot, _ := merkletree.NewMerkleTree(ctx, memory.NewMemoryStorage(), 32)

	// Get the Index of the claim and the Value of the authClaim
	hIndex, hValue, _ := authClaim.HiHv()

	err := clt.Add(ctx, hIndex, hValue)
	if err != nil {
		return nil, errors.Wrap(err, "Error while adding AuthClaim to the Clt during NewIdentity Generation")
	}

	state, _ := merkletree.HashElems(
		clt.Root().BigInt(),
		ret.Root().BigInt(),
		rot.Root().BigInt())

	id, _ := core.IdGenesisFromIdenState(core.TypeDefault, state.BigInt())

	identity := Identity{
		ID:             id,
		IDS:            state,
		PrivateKey:     babyJubjubPrivKey,
		AuthClaim:      authClaim,
		Claims:         make([]*core.Claim, 0),
		Clt:            clt,
		Ret:            ret,
		Rot:            rot,
		ReceivedClaims: make(map[string]circuits.Claim),
	}
	return &identity, nil
}

func LoadIdentityFromFile(file string) (*Identity, error) {
	identity := new(Identity)
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, identity)
	if err != nil {
		return nil, errors.Wrap(err, "Error during Unmarshal of identity file")
	}

	ctx := context.Background()
	identity.Clt, _ = merkletree.NewMerkleTree(ctx, memory.NewMemoryStorage(), 32)
	identity.Ret, _ = merkletree.NewMerkleTree(ctx, memory.NewMemoryStorage(), 32)
	identity.Rot, _ = merkletree.NewMerkleTree(ctx, memory.NewMemoryStorage(), 32)

	hIndex, hValue, _ := identity.AuthClaim.HiHv()

	err = identity.Clt.Add(ctx, hIndex, hValue)
	if err != nil {
		return nil, errors.Wrap(err, "Failed while adding auth claim from JSON File")
	}

	state := identity.GetIDS()

	id, _ := core.IdGenesisFromIdenState(core.TypeDefault, state.BigInt())

	// Check if the generated ID from Genesis state is same as we have in the file.
	if id.String() != identity.ID.String() {
		return nil, errors.Errorf("ID differs while recreating it from json file. Generated id is %s but in file it is %s", id.String(), identity.ID.String())
	}

	// TODO: Handle Ret Tree
	for _, claim := range identity.Claims {
		// Before updating the claims tree, add the claims tree root at Genesis state to the Roots tree.
		err := identity.Rot.Add(ctx, identity.Clt.Root().BigInt(), big.NewInt(0))
		if err != nil {
			return nil, errors.Wrap(err, "Error while adding the root of the Clt to the Rot In the load file func.")
		}
		hIndex, hValue, _ := claim.HiHv()
		err = identity.Clt.Add(ctx, hIndex, hValue)
		if err != nil {
			return nil, errors.Wrap(err, "Failed while adding claims to the Clt from JSON File")
		}
	}

	if identity.GetIDS().String() != identity.IDS.String() {
		return nil, errors.Errorf("IDS differs while recreating it from json file. Generated IDS is %s but in file it is %s", identity.GetIDS().String(), identity.IDS.String())
	}
	return identity, nil
}

func (identity *Identity) AddClaim(claim ClaimAPI, config *Config) error {
	ctx := context.Background()

	authClaim := identity.AuthClaim
	authClaimIndex, _ := authClaim.HIndex()
	authClaimRevNonce := authClaim.GetRevocationNonce()

	// 1. Generate Merkle Tree Proof for authClaim at Genesis State
	authMTPProof, _, _ := identity.Clt.GenerateProof(ctx, authClaimIndex, identity.Clt.Root())

	// 2. Generate the Non-Revocation Merkle tree proof for the authClaim at Genesis State
	authNonRevMTPProof, _, _ := identity.Ret.GenerateProof(ctx, big.NewInt(int64(authClaimRevNonce)), identity.Ret.Root())

	oldState := identity.GetIDS()
	isOldStateGenesis, _ := identity.IsAtGenesisState()
	oldTreeState := identity.GetTreeState()

	// Before updating the claims tree, add the claims tree root at Genesis state to the Roots tree.
	err := identity.Rot.Add(ctx, identity.Clt.Root().BigInt(), big.NewInt(0))
	if err != nil {
		return errors.Wrap(err, "Error while adding the root of the Clt to the Rot")
	}

	claimToAdd := CreateIden3ClaimFromAPI(claim)
	hIndex, hValue, _ := claimToAdd.HiHv()

	err = identity.Clt.Add(ctx, hIndex, hValue)
	if err != nil {
		return errors.Wrap(err, "Error while adding the new claim to Clt")
	}
	// Add the claim to our array
	identity.Claims = append(identity.Claims, claimToAdd)

	// Fetch the new Identity State
	newState := identity.GetIDS()
	identity.IDS = newState

	// Sign a message (hash of the old state + the new state) using your private key
	hashOldAndNewStates, _ := poseidon.Hash([]*big.Int{oldState.BigInt(), newState.BigInt()})

	signature := identity.PrivateKey.SignPoseidon(hashOldAndNewStates)

	// Generate state transition inputs
	stateTransitionInputs := circuits.StateTransitionInputs{
		ID:                identity.ID,
		OldTreeState:      oldTreeState,
		NewState:          newState,
		IsOldStateGenesis: isOldStateGenesis,
		AuthClaim: circuits.Claim{
			Claim: authClaim,
			Proof: authMTPProof,
			NonRevProof: &circuits.ClaimNonRevStatus{
				Proof: authNonRevMTPProof,
			},
		},
		Signature: signature,
	}

	// Perform marshalling of the state transition inputs
	inputBytes, _ := stateTransitionInputs.InputsMarshal()
	proof, err := GenerateZkProof(config.Circuits.Path+"stateTransition", toJSON(inputBytes), config)
	if err != nil {
		return errors.Wrap(err, "Error while creating proof using snarkJS")
	}
	transaction, err := TransitState(config, identity.ID, proof)
	if err != nil {
		return errors.Wrap(err, "Errored while submitting transaction to blockchain")
	}
	log.Printf("Add Claim successful. Submitted to change state on blockchain with txID: %s\n", transaction.Hash().String())
	return nil
}

func (identity *Identity) AddClaimsFromIssuer(claims []circuits.Claim) error {
	// TODO: Better key for looking up Claims
	for _, claim := range claims {
		schemaHash, _ := claim.Claim.GetSchemaHash().MarshalText()
		identity.ReceivedClaims[string(schemaHash)] = claim
	}
	return nil
}

func (identity *Identity) ProofRequest(request protocol.AuthorizationRequestMessage, config *Config) (*protocol.AuthorizationResponseMessage, error) {
	rules := request.Body.Scope[0].Rules
	jsonStr, err := json.Marshal(rules["query"])
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert query into jsonStr")
	}
	var query pubsignals.Query
	if err := json.Unmarshal(jsonStr, &query); err != nil {
		return nil, errors.Wrap(err, "Failed to typecast rule to pubsignals.Query")
	}
	parsedQuery, _ := ValidateAndGetCircuitsQuery(query, context.Background(), loaders.DefaultSchemaLoader{IpfsURL: ""})

	challenge := new(big.Int).SetInt64(1)
	schemaHash := GetHashFromClaimSchemaURL(query.Schema.URL, query.Schema.Type)

	// TODO: Get Dynamic circuit name from proof request
	circuitName := circuits.AtomicQuerySigCircuitID
	if val, ok := identity.ReceivedClaims[schemaHash]; ok {
		atomicInputs := circuits.AtomicQuerySigInputs{
			ID:               identity.ID,
			AuthClaim:        identity.GetUserAuthClaim(),
			Challenge:        challenge,
			Signature:        identity.PrivateKey.SignPoseidon(challenge),
			CurrentTimeStamp: time.Now().Unix(),
			Claim:            val,
			Query:            *parsedQuery,
		}
		inputBytes, err := atomicInputs.InputsMarshal()
		if err != nil {
			return nil, errors.Wrapf(err, "Error during marshalling of %s circuit inputs", circuitName)
		}
		proof, err := GenerateZkProof(config.Circuits.Path+"credentialAtomicQuerySig", toJSON(inputBytes), config)
		if err != nil {
			return nil, errors.Wrap(err, "Error while generating proof using snarkJS")
		}
		resp := protocol.AuthorizationResponseMessage{
			ID:       request.ID,
			Typ:      packers.MediaTypePlainMessage,
			Type:     protocol.AuthorizationResponseMessageType,
			ThreadID: request.ThreadID,
			Body: protocol.AuthorizationMessageResponseBody{
				Message: request.Body.Message,
				Scope: []protocol.ZeroKnowledgeProofResponse{
					{
						ID:        1,
						CircuitID: string(circuits.AtomicQuerySigCircuitID),
						ZKProof:   *proof,
					},
				},
			},
			From: identity.ID.String(),
			To:   request.From,
		}
		return &resp, nil
	} else {
		return nil, errors.New("Requested claim does not exists in the wallet.")
	}
}

func (identity *Identity) GetTreeState() circuits.TreeState {
	return circuits.TreeState{
		State:          identity.GetIDS(),
		ClaimsRoot:     identity.Clt.Root(),
		RevocationRoot: identity.Ret.Root(),
		RootOfRoots:    identity.Rot.Root(),
	}
}

func (identity *Identity) GetIDS() *merkletree.Hash {
	state, _ := merkletree.HashElems(
		identity.Clt.Root().BigInt(),
		identity.Ret.Root().BigInt(),
		identity.Rot.Root().BigInt())
	return state
}

func (identity *Identity) GetUserAuthClaim() circuits.Claim {
	ctx := context.Background()

	authClaim := identity.AuthClaim
	authClaimIndex, _ := authClaim.HIndex()
	authClaimRevNonce := authClaim.GetRevocationNonce()

	authMTPProof, _, _ := identity.Clt.GenerateProof(ctx, authClaimIndex, identity.Clt.Root())

	authNonRevMTPProof, _, _ := identity.Ret.GenerateProof(ctx, big.NewInt(int64(authClaimRevNonce)), identity.Ret.Root())

	inputsAuthClaim := circuits.Claim{
		//Schema:    authClaim.Schema,
		Claim:     identity.AuthClaim,
		Proof:     authMTPProof,
		TreeState: identity.GetTreeState(),
		NonRevProof: &circuits.ClaimNonRevStatus{
			TreeState: identity.GetTreeState(),
			Proof:     authNonRevMTPProof,
		},
	}
	return inputsAuthClaim
}

// IsAtGenesisState TODO: Implement this function
func (identity *Identity) IsAtGenesisState() (bool, error) {
	ans, err := checkGenesisStateID(identity.ID.BigInt(), identity.IDS.BigInt())
	if err != nil {
		return false, errors.Wrap(err, "Failed to check if state is genesis")
	}
	return ans, nil
}
