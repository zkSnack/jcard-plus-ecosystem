package main

import (
	"fmt"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"math/big"
)

func main() {
	babyJubjubPrivKey := "priv_key"
	babyJubjubPrivKeyBigInt := new(big.Int)
	babyJubjubPrivKeyBigInt.SetString(babyJubjubPrivKey, 16)

	babyJubjubPrivKeyScaler := babyjub.NewPrivKeyScalar(babyJubjubPrivKeyBigInt)

	// generate public key from private key
	babyJubjubPubKey := babyJubjubPrivKeyScaler.Public()
	fmt.Println("Public Key: ", babyJubjubPubKey)
}
