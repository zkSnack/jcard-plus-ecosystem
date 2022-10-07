package main

import (
	"fmt"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-iden3-crypto/utils"
)

func main() {
	// This is example privateKey, please do not use it as your own.
	input := "0x50c4174f8e9cef5dffc2c905bc66fda886453b9f06794b56131ef5e16c396715"
	babyJubjubPrivKeyInByte, _ := utils.HexDecode(input)
	var babyJubjubPrivKey babyjub.PrivateKey
	copy(babyJubjubPrivKey[:], babyJubjubPrivKeyInByte)

	// generate public key from private key
	babyJubjubPubKey := babyJubjubPrivKey.Public()
	fmt.Println("Public Key: ", babyJubjubPubKey)
	// Should output: bb153d95621a7aaede6f8a07f7ddfcfa05c70f16b2b772582301a2908242ea88
}
