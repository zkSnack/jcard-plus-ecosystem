package main

import (
	"context"
	"encoding/json"
	"github.com/iden3/go-iden3-auth/loaders"
	"github.com/iden3/iden3comm/protocol"
	"github.com/pkg/errors"
	"math/big"

	"github.com/iden3/go-iden3-auth/pubsignals"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/go-merkletree-sql/v2/db/memory"
	"io/ioutil"
	"log"
	"math/rand"
)

type Account struct {
	ID         *core.ID           `json:"id"`
	IDS        *merkletree.Hash   `json:"identity_state"`
	PrivateKey babyjub.PrivateKey `json:"private_key"`
	PublicKey  *babyjub.PublicKey `json:"public_key"`
	CltRoot    *merkletree.Hash   `json:"clt"`
	RetRoot    *merkletree.Hash   `json:"ret"`
	RotRoot    *merkletree.Hash   `json:"rot"`
	AuthClaim  *core.Claim        `json:"authClaim"`
	Claims     []*core.Claim      `json:"claims"`
	Identity   *Identity          `json:"identity"`
}

type Config struct {
	Issuer Issuer `yaml:"issuer"`
}

type Issuer struct {
	URL string `yaml:"url"`
	ID  string `yaml:"id"`
}

func NewAccount() *Account {
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

	clt.Add(ctx, hIndex, hValue)

	state, _ := merkletree.HashElems(
		clt.Root().BigInt(),
		ret.Root().BigInt(),
		rot.Root().BigInt())

	id, _ := core.IdGenesisFromIdenState(core.TypeDefault, state.BigInt())

	account := new(Account)
	account.ID = id
	account.IDS = state
	account.PrivateKey = babyJubjubPrivKey
	account.PublicKey = babyJubjubPubKey
	account.CltRoot = clt.Root()
	account.RetRoot = ret.Root()
	account.RotRoot = rot.Root()
	account.AuthClaim = authClaim
	identity := NewIdentity(babyJubjubPrivKey, id, clt, ret, rot, authClaim)
	account.Identity = identity
	return account
}

func LoadAccountFromFile(file string) *Account {
	account := new(Account)
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	err = json.Unmarshal(content, account)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	identity := FromFileData(account)
	account.Identity = identity
	return account
}

func (account *Account) addClaim(claim ClaimAPI) []byte {
	inputJSON, newClaim := account.Identity.addClaim(claim)

	account.IDS = account.Identity.GetIDS()
	account.CltRoot = account.Identity.Clt.Root()
	account.RetRoot = account.Identity.Ret.Root()
	account.RotRoot = account.Identity.Rot.Root()
	account.Claims = append(account.Claims, newClaim)
	return inputJSON
}

func (account *Account) GenerateProof(request protocol.AuthorizationRequestMessage) ([]byte, error) {
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
	return account.Identity.GenerateProof(challenge, *parsedQuery, query.Schema)
}
