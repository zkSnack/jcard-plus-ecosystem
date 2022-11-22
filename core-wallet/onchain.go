package walletSDK

import (
	"math/big"
	state "zkSnacks/walletSDK/contracts"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	types2 "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-rapidsnark/types"
	"github.com/pkg/errors"
)

func getClient(config *Config) (*ethclient.Client, error) {
	client, err := ethclient.Dial(config.Web3.URL)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to init client")
	}
	return client, nil
}

func TransitState(config *Config, id *core.ID, proof *types.ZKProof) (*types2.Transaction, error) {
	client, err := getClient(config)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create Eth Client")
	}

	address := common.HexToAddress(config.Web3.StateTransition)
	instance, err := state.NewState(address, client)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to instance of State Contract")
	}
	pub, err := stringsToArrayBigInt(proof.PubSignals)
	if err != nil {
		return nil, errors.Wrap(err, "Failed while converting pubsignal to BigInt")
	}
	var a [2]*big.Int
	for index, val := range proof.Proof.A[:2] {
		tmp, err := stringToBigInt(val)
		if err != nil {
			return nil, errors.Wrap(err, "Failed while converting Proof.A to BigInt")
		}
		a[index] = tmp
	}

	var b [2][2]*big.Int
	for index1, val1 := range proof.Proof.B[:2] {
		for index2, val2 := range val1[:2] {
			tmp, err := stringToBigInt(val2)
			if err != nil {
				return nil, errors.Wrap(err, "Failed while converting Proof.B to BigInt")
			}
			b[index1][1-index2] = tmp
		}
	}
	var c [2]*big.Int
	for index, val := range proof.Proof.C[:2] {
		tmp, err := stringToBigInt(val)
		if err != nil {
			return nil, errors.Wrap(err, "Failed while converting Proof.C to BigInt")
		}
		c[index] = tmp
	}

	privateKey, err := crypto.HexToECDSA(config.Web3.PrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "Failed while deriving privateKey")
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, config.Web3.ChainID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed while creating NewKeyedTransactorWithChainID")
	}
	tx, err := instance.TransitState(auth, pub[0], pub[1], pub[2], len(pub[3].Bytes()) != 0, a, b, c)
	if err != nil {
		return nil, errors.Wrap(err, "Failed while calling TransitState func")
	}
	return tx, nil
}

func GetCurrentState(config *Config, id *core.ID) (*big.Int, error) {
	client, err := getClient(config)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create Eth Client")
	}

	address := common.HexToAddress(config.Web3.StateTransition)
	instance, err := state.NewState(address, client)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create NewState")
	}
	currentState, err := instance.GetState(nil, id.BigInt())
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get State from smart contract")
	}
	return currentState, nil
}
