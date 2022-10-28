package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	merkletree "github.com/iden3/go-merkletree-sql"
	"github.com/iden3/go-merkletree-sql/db/memory"
)

type Account struct {
	ID         *core.ID           `json:"id"`
	IDS        *merkletree.Hash   `json:"identity_state"`
	PrivateKey babyjub.PrivateKey `json:"private_key"`
	PublicKey  *babyjub.PublicKey `json:"public_key"`
	CltRoot    *merkletree.Hash   `json:"clt"`
	RetRoot    *merkletree.Hash   `json:"ret"`
	RotRoot    *merkletree.Hash   `json:"rot"`
	Claims     []*core.Claim      `json:"claims"`
}

func main() {
	router := gin.Default()
	router.GET("/api/v1/account", loadAccount)
	router.POST("/api/v1/generate", generateAccount)

	router.Run("localhost:8080")
}

func loadAccount(c *gin.Context) {
	content, err := ioutil.ReadFile("./account.json")
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var account Account
	err = json.Unmarshal(content, &account)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
	c.IndentedJSON(http.StatusOK, account)
}

func generateAccount(c *gin.Context) {
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

	var claims []*core.Claim
	claims = append(claims, authClaim)
	var account Account
	account.ID = id
	account.IDS = state
	account.PrivateKey = babyJubjubPrivKey
	account.PublicKey = babyJubjubPubKey
	account.CltRoot = clt.Root()
	account.RetRoot = ret.Root()
	account.RotRoot = rot.Root()
	account.Claims = claims

	file, _ := json.MarshalIndent(account, "", "	")
	_ = ioutil.WriteFile("account.json", file, 0644)
	c.IndentedJSON(http.StatusOK, account)
}
