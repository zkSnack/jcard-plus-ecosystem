package main

import (
	"context"
	"fmt"
	"math/big"

	circuits "github.com/iden3/go-circuits"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-iden3-crypto/utils"
	merkletree "github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/go-merkletree-sql/v2/db/memory"
)

type Identity struct {
	PrivateKey      babyjub.PrivateKey     `json:"private_key"`
	ID              *core.ID               `json:"id"`
	Clt             *merkletree.MerkleTree `json:"clt"`
	Ret             *merkletree.MerkleTree `json:"ret"`
	Rot             *merkletree.MerkleTree `json:"rot"`
	AuthClaim       *core.Claim            `json:"authClaim"`
	AuthState       *circuits.TreeState    `json:"authState"`
	AuthMTPProof    *merkletree.Proof      `json:"authMTPProof"`
	AuthNonRevProof *merkletree.Proof      `json:"authNonRevProof"`
	IDS             *merkletree.Hash       `json:"identity_state"`
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

	authTreeState := circuits.TreeState{
		State:          idenState,
		ClaimsRoot:     clt.Root(),
		RevocationRoot: &merkletree.HashZero,
		RootOfRoots:    &merkletree.HashZero,
	}

	authMtpProof, _, _ := clt.GenerateProof(ctx, hIndex, clt.Root())

	authClaimRevNonce := new(big.Int).SetUint64(authClaim.GetRevocationNonce())
	authNonRevProof, _, _ := ret.GenerateProof(ctx, authClaimRevNonce, ret.Root())

	issuerIdentity := &Identity{
		PrivateKey:      privateKey,
		ID:              id,
		Clt:             clt,
		Ret:             ret,
		Rot:             rot,
		AuthClaim:       authClaim,
		AuthState:       &authTreeState,
		AuthMTPProof:    authMtpProof,
		AuthNonRevProof: authNonRevProof,
		IDS:             idenState,
	}

	// Print Roots of Merkle Trees
	fmt.Println("Genesis ID:", issuerIdentity.ID)
	return issuerIdentity
}

// TO-DO: Use logic from Iden3-SDK instead of this
func generateAgeClaim(issuerIdentity *Identity, holderID string) *circuits.Claim {
	ctx := context.Background()
	claimSchemaHashHex := generateHashFromClaimSchemaFile("student-age.json-ld", "AgeCredential")
	claimSchemaHash, _ := core.NewSchemaHashFromHex(claimSchemaHashHex)

	// Why is this needed? Is it any use on the wallet or verifier side?
	subjectId, _ := core.IDFromString("113TCVw5KMeMp99Qdvub9Mssfz7krL9jWNvbdB7Fd2")

	// Add revocation nonce. Used to invalidate the claim. Update it to random number once finish testing.
	revNonce := uint64(7)
	// revNonce := rand.Uint64()

	birthday := big.NewInt(19960424)

	ageClaim, _ := core.NewClaim(claimSchemaHash,
		core.WithIndexDataInts(birthday, big.NewInt(0)),
		core.WithRevocationNonce(revNonce),
		core.WithIndexID(subjectId))

	hIndexAgeClaim, hValueageClaim, _ := ageClaim.HiHv()
	claimHash, _ := merkletree.HashElems(hIndexAgeClaim, hValueageClaim)

	claimSignature := issuerIdentity.PrivateKey.SignPoseidon(claimHash.BigInt())

	// Add Age Claim to Claim Merkle Tree
	issuerIdentity.Clt.Add(ctx, hIndexAgeClaim, hValueageClaim)

	// Generate Proof of Claim
	ageClaimProof, _, _ := issuerIdentity.Clt.GenerateProof(ctx, hIndexAgeClaim, issuerIdentity.Clt.Root())

	// Generate Revocation Proof
	claimRevNonce := new(big.Int).SetUint64(ageClaim.GetRevocationNonce())
	proofNotRevoke, _, _ := issuerIdentity.Ret.GenerateProof(ctx, claimRevNonce, issuerIdentity.Ret.Root())

	idsAfterClaimAdd, _ := merkletree.HashElems(
		issuerIdentity.Clt.Root().BigInt(),
		issuerIdentity.Ret.Root().BigInt(),
		issuerIdentity.Rot.Root().BigInt())

	issuerIdentity.IDS = idsAfterClaimAdd

	issuerStateAfterClaimAdd := circuits.TreeState{
		State:          idsAfterClaimAdd,
		ClaimsRoot:     issuerIdentity.Clt.Root(),
		RevocationRoot: issuerIdentity.Ret.Root(),
		RootOfRoots:    issuerIdentity.Rot.Root(),
	}

	claimIssuerSignature := circuits.BJJSignatureProof{
		IssuerID:           issuerIdentity.ID,
		IssuerTreeState:    *issuerIdentity.AuthState,
		IssuerAuthClaimMTP: issuerIdentity.AuthMTPProof,
		Signature:          claimSignature,
		IssuerAuthClaim:    issuerIdentity.AuthClaim,
		IssuerAuthNonRevProof: circuits.ClaimNonRevStatus{
			TreeState: *issuerIdentity.AuthState,
			Proof:     issuerIdentity.AuthNonRevProof,
		},
	}

	holderAgeClaim := circuits.Claim{
		Claim:     ageClaim,
		Proof:     ageClaimProof,
		TreeState: issuerStateAfterClaimAdd,
		IssuerID:  issuerIdentity.ID,
		NonRevProof: &circuits.ClaimNonRevStatus{
			TreeState: issuerStateAfterClaimAdd,
			Proof:     proofNotRevoke,
		},
		SignatureProof: claimIssuerSignature,
	}

	// fmt.Println("Age Claim:")
	// claimMarshelText, _ := json.MarshalIndent(holderAgeClaim, "", "\t")
	// fmt.Println(string(claimMarshelText))
	return &holderAgeClaim
}

func IssueClaims(holderID string) []circuits.Claim {

	babyJubjubPrivKeyString := "0x8a2e1766a7f4851b6d27d313b7c4b7b271772763eb33466c50671f3e8597c658"
	babyJubjubPrivKeyInByte, _ := utils.HexDecode(babyJubjubPrivKeyString)
	var babyJubjubPrivKey babyjub.PrivateKey
	copy(babyJubjubPrivKey[:], babyJubjubPrivKeyInByte)
	// babyJubjubPrivKey := babyjub.NewRandPrivKey()
	fmt.Println("Private Key: ", babyJubjubPrivKey)

	// generate public key from private key
	babyJubjubPubKey := babyJubjubPrivKey.Public()

	// print public key
	fmt.Println("BabyJubJub Public Key:", babyJubjubPubKey)

	// 2. Create an Identity

	// 2.1 Create an Auth Claim
	authClaim := generateAuthClaim(babyJubjubPubKey)

	// 2.2 Generate an identity
	ctx := context.Background()
	issuerIdentity := generateIssuerIdentity(ctx, babyJubjubPrivKey, authClaim)

	var claims []circuits.Claim
	claims = append(claims, *generateAgeClaim(issuerIdentity, "113TCVw5KMeMp99Qdvub9Mssfz7krL9jWNvbdB7Fd2"))

	return claims
}
