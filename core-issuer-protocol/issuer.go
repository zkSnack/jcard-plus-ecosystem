package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-iden3-crypto/poseidon"
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
	AuthClaim  *core.Claim        `json:"authClaim"`
	Claims     []*core.Claim      `json:"claims"`
	Identity   *Identity          `json:"identity"`
}

type Identity struct {
	PrivateKey babyjub.PrivateKey     `json:"private_key"`
	ID         *core.ID               `json:"id"`
	Clt        *merkletree.MerkleTree `json:"clt"`
	Ret        *merkletree.MerkleTree `json:"ret"`
	Rot        *merkletree.MerkleTree `json:"rot"`
	AuthClaim  *core.Claim            `json:"authClaim"`
}

type Issuer struct {
	ID        *core.ID           `json:"id"`
	IDS       *merkletree.Hash   `json:"identity_state"`
	PublicKey *babyjub.PublicKey `json:"public_key"`
	Claims    []*core.Claim      `json:"claims"`
	Identity  *Identity          `json:"identity"`
}

type Claim struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

func generateAuthClaim(babyJubjubPubKey *babyjub.PublicKey) *core.Claim {
	authSchemaHashHex := generateHashFromClaimSchemaFile("auth.json-ld", "AuthBJJCredential")
	authSchemaHash, _ := core.NewSchemaHashFromHex(authSchemaHashHex)

	// Add revocation nonce. Used to invalidate the claim. Update it to random number once finish testing.
	revNonce := uint64(1)
	// revNonce := rand.Uint64()

	authClaim, _ := core.NewClaim(authSchemaHash,
		core.WithIndexDataInts(babyJubjubPubKey.X, babyJubjubPubKey.Y),
		core.WithRevocationNonce(revNonce))

	authClaimToMarshal, _ := json.Marshal(authClaim)
	fmt.Println("Auth Claim:", string(authClaimToMarshal))
	return authClaim
}

func generateIssuerIdentity(ctx context.Context, privateKey babyjub.PrivateKey, authClaim *core.Claim) *Identity {
	// Claim Merkle Tree
	clt, _ := merkletree.NewMerkleTree(ctx, memory.NewMemoryStorage(), 32)

	// Revocation Merkle Tree
	ret, _ := merkletree.NewMerkleTree(ctx, memory.NewMemoryStorage(), 32)

	// Roots Merkle Tree
	rot, _ := merkletree.NewMerkleTree(ctx, memory.NewMemoryStorage(), 32)

	hIndex, hValue, _ := authClaim.HiHv()

	// Add Auth Claim to Claim Merkle Tree
	clt.Add(ctx, hIndex, hValue)

	// idenState, _ := core.IdenState(clt, ret, rot)
	idenState, _ := merkletree.HashElems(
		clt.Root().BigInt(),
		ret.Root().BigInt(),
		rot.Root().BigInt())
	id, _ := core.IdGenesisFromIdenState(core.TypeDefault, idenState.BigInt())

	issuerIdentity := &Identity{
		PrivateKey: privateKey,
		ID:         id,
		Clt:        clt,
		Ret:        ret,
		Rot:        rot,
		AuthClaim:  authClaim,
	}

	// Print Roots of Merkle Trees
	fmt.Println("Genesis ID:", issuerIdentity.ID)
	fmt.Println("Claim Merkle Tree Root:", issuerIdentity.Clt.Root().BigInt())
	fmt.Println("Revocation Merkle Tree Root:", issuerIdentity.Ret.Root().BigInt())
	fmt.Println("Roots Merkle Tree Root:", issuerIdentity.Rot.Root().BigInt())
	return issuerIdentity
}

func generateAgeClaim(babyJubjubPubKey *babyjub.PublicKey) *core.Claim {
	claimSchemaHashHex := generateHashFromClaimSchemaFile("student-age.json-ld", "AgeCredential")
	claimSchemaHash, _ := core.NewSchemaHashFromHex(claimSchemaHashHex)

	// Why is this needed? Is it any use on the wallet or verifier side?
	subjectId, _ := core.IDFromString("113TCVw5KMeMp99Qdvub9Mssfz7krL9jWNvbdB7Fd2")

	// Add revocation nonce. Used to invalidate the claim. Update it to random number once finish testing.
	revNonce := uint64(7)
	// revNonce := rand.Uint64()

	birthday := big.NewInt(19960424)

	claim, _ := core.NewClaim(claimSchemaHash,
		core.WithIndexDataInts(birthday, big.NewInt(0)),
		core.WithRevocationNonce(revNonce),
		core.WithIndexID(subjectId))

	claimToMarshal, _ := json.Marshal(claim)
	fmt.Println("Auth Claim:", string(claimToMarshal))
	return claim
}

func (identity *Identity) offerClaim(claim *core.Claim) {
	// TO-D0: Add logic to offer a claim to the user
}

func (identity *Identity) revokeClaim(claim *core.Claim) {
	// TO-D0: Add logic to revoke a claim
}

func (identity *Identity) IssueClaimBySignature(claim *core.Claim) {
	// TO-D0: Add logic to issue a claim
	claimIndex, claimValue := claim.RawSlots()
	indexHash, _ := poseidon.Hash(core.ElemBytesToInts(claimIndex[:]))
	valueHash, _ := poseidon.Hash(core.ElemBytesToInts(claimValue[:]))

	// Poseidon Hash the indexHash and the valueHash together to get the claimHash
	claimHash, _ := merkletree.HashElems(indexHash, valueHash)

	// Sign the claimHash with the private key of the issuer
	claimSignature := identity.PrivateKey.SignPoseidon(claimHash.BigInt())

	fmt.Println("Claim Signature:", claimSignature)
}

func IssueClaim(holderID string) []Claim {
	var claims = []Claim{
		{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
		{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
	}
	return claims
}

// func main() {

// 	// 1. BabyJubJub key - Generate a new key pair randomly

// 	// TODO: update it to random number once finish testing.
// 	babyJubjubPrivKeyString := "0x8a2e1766a7f4851b6d27d313b7c4b7b271772763eb33466c50671f3e8597c658"
// 	babyJubjubPrivKeyInByte, _ := utils.HexDecode(babyJubjubPrivKeyString)
// 	var babyJubjubPrivKey babyjub.PrivateKey
// 	copy(babyJubjubPrivKey[:], babyJubjubPrivKeyInByte)
// 	// babyJubjubPrivKey := babyjub.NewRandPrivKey()
// 	fmt.Println("Private Key: ", babyJubjubPrivKey)

// 	// generate public key from private key
// 	babyJubjubPubKey := babyJubjubPrivKey.Public()

// 	// print public key
// 	fmt.Println("BabyJubJub Public Key:", babyJubjubPubKey)

// 	// 2. Create an Identity

// 	// 2.1 Create an Auth Claim
// 	authClaim := generateAuthClaim(babyJubjubPubKey)

// 	// 2.2 Generate an identity
// 	ctx := context.Background()
// 	issuerIdentity := generateIssuerIdentity(ctx, babyJubjubPrivKey, authClaim)

// 	ageClaim := generateAgeClaim(babyJubjubPubKey)
// 	// Offer a claim to the user
// 	issuerIdentity.offerClaim(ageClaim)

// 	// Issue a claim
// 	issuerIdentity.IssueClaimBySignature(ageClaim)
// 	// identity, _ := core.NewIdentity(babyJubjubPrivKey, authClaim)

// }
