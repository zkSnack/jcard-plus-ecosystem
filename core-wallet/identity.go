package main

import (
	"context"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-merkletree-sql"
	"github.com/iden3/go-merkletree-sql/db/memory"
	"log"
)

type Identity struct {
	Clt *merkletree.MerkleTree
	Ret *merkletree.MerkleTree
	Rot *merkletree.MerkleTree
}

func (identity *Identity) Init(clt *merkletree.MerkleTree, ret *merkletree.MerkleTree, rot *merkletree.MerkleTree) {
	identity.Clt = clt
	identity.Ret = ret
	identity.Rot = rot
}

func (identity *Identity) FromAccount(account *Account) {
	ctx := context.Background()
	clt, _ := merkletree.NewMerkleTree(ctx, memory.NewMemoryStorage(), 32)
	ret, _ := merkletree.NewMerkleTree(ctx, memory.NewMemoryStorage(), 32)
	rot, _ := merkletree.NewMerkleTree(ctx, memory.NewMemoryStorage(), 32)

	hIndex, hValue, _ := account.AuthClaim.HiHv()

	clt.Add(ctx, hIndex, hValue)

	state, _ := merkletree.HashElems(
		clt.Root().BigInt(),
		ret.Root().BigInt(),
		rot.Root().BigInt())

	id, _ := core.IdGenesisFromIdenState(core.TypeDefault, state.BigInt())

	if id.String() != account.ID.String() {
		log.Fatal("ID differs while recreating it from json file.", id, account.ID)
	}

	for i := 0; i < len(account.Claims); i++ {
		hIndex, hValue, _ := account.Claims[i].HiHv()
		clt.Add(ctx, hIndex, hValue)
	}

	state, _ = merkletree.HashElems(
		clt.Root().BigInt(),
		ret.Root().BigInt(),
		rot.Root().BigInt())

	if state.String() != account.IDS.String() {
		log.Fatal("IDS differs while recreating it from json file.")
	}

	identity.Init(clt, ret, rot)
}

func FromFile(file string) *Identity {
	account := new(Account)
	account.LoadAccount(file)
	identity := new(Identity)
	identity.FromAccount(account)
	return identity
}
