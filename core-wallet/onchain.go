package walletsdk

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-rapidsnark/types"
	"github.com/pkg/errors"
	"log"
	"math/big"
	state "zkSnacks/walletsdk/contracts"
)

func getClient(config *Config) (*ethclient.Client, error) {
	client, err := ethclient.Dial(config.Web3.URL)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to init client")
	}
	return client, nil
}

func TransitState(config *Config, id *core.ID, proof *types.ZKProof) {
	client, err := getClient(config)
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress(config.Web3.StateTransition)
	instance, err := state.NewState(address, client)
	if err != nil {
		log.Fatal(err)
	}
	oldState, err := instance.GetStateInfoById(nil, id.BigInt())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(oldState)
	pub, err := stringsToArrayBigInt(proof.PubSignals)
	if err != nil {
		log.Fatal(err)
	}
	var a [2]*big.Int
	for index, val := range proof.Proof.A[:2] {
		tmp, err := stringToBigInt(val)
		if err != nil {
			log.Fatal(err)
		}
		a[index] = tmp
	}

	var b [2][2]*big.Int
	for index1, val1 := range proof.Proof.B[:2] {
		for index2, val2 := range val1[:2] {
			tmp, err := stringToBigInt(val2)
			if err != nil {
				log.Fatal(err)
			}
			b[index1][1-index2] = tmp
		}
	}
	var c [2]*big.Int
	for index, val := range proof.Proof.C[:2] {
		tmp, err := stringToBigInt(val)
		if err != nil {
			log.Fatal(err)
		}
		c[index] = tmp
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pub, a, b, c)

	privateKey, err := crypto.HexToECDSA(config.Web3.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, config.Web3.ChainID)
	if err != nil {
		log.Fatal(err)
	}
	tx, err := instance.TransitState(auth, pub[0], pub[1], pub[2], len(pub[3].Bytes()) != 0, a, b, c)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tx sent: %s", tx.Hash().Hex())
}
