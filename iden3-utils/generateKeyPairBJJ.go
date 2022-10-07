package main

import (
	"fmt"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-iden3-crypto/utils"
)

func main() {
	// generate babyJubjub private key randomly
	babyJubjubPrivKey := babyjub.NewRandPrivKey()
	babyJubjubPrivKeyString := utils.HexEncode(babyJubjubPrivKey[:])
	fmt.Println("Private Key: ", babyJubjubPrivKeyString)

	// generate public key from private key
	babyJubjubPubKey := babyJubjubPrivKey.Public()
	fmt.Println("Public Key: ", babyJubjubPubKey)
}
