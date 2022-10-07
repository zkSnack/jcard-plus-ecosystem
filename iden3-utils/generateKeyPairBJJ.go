package main

import (
	"fmt"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"math/big"
)

func toHexInt(n *big.Int) string {
	return fmt.Sprintf("%x", n)
}

func main() {
	// generate babyJubjub private key randomly
	babyJubjubPrivKey := babyjub.NewRandPrivKey()
	babyJubjubPrivKeyScaler := babyjub.SkToBigInt(&babyJubjubPrivKey)
	fmt.Println("Private Key: ", toHexInt(babyJubjubPrivKeyScaler))

	// generate public key from private key
	babyJubjubPubKey := babyJubjubPrivKey.Public()
	fmt.Println("Public Key: ", babyJubjubPubKey)
}
