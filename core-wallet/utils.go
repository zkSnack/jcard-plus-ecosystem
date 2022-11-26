package walletSDK

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-merkletree-sql/v2"
)

func toJSON(inputJSON []byte) map[string]interface{} {
	var jsonBlob map[string]interface{}
	if err := json.Unmarshal(inputJSON, &jsonBlob); err != nil {
		log.Fatal("Failed while converting to JSON format.")
	}
	return jsonBlob
}

func stringsToArrayBigInt(publicInputs []string) ([]*big.Int, error) {
	p := make([]*big.Int, 0, len(publicInputs))
	for _, s := range publicInputs {
		sb, err := stringToBigInt(s)
		if err != nil {
			return nil, err
		}
		p = append(p, sb)
	}
	return p, nil
}

func stringToBigInt(s string) (*big.Int, error) {
	base := 10
	if bytes.HasPrefix([]byte(s), []byte("0x")) {
		base = 16
		s = strings.TrimPrefix(s, "0x")
	}
	n, ok := new(big.Int).SetString(s, base)
	if !ok {
		return nil, fmt.Errorf("can not parse string to *big.Int: %s", s)
	}
	return n, nil
}

func checkGenesisStateID(id, state *big.Int) (bool, error) {

	stateHash, err := merkletree.NewHashFromBigInt(state)
	if err != nil {
		return false, err
	}

	IDFromState, err := core.IdGenesisFromIdenState(core.TypeDefault, stateHash.BigInt())
	if err != nil {
		return false, err
	}

	idBytes := merkletree.NewElemBytesFromBigInt(id)
	IDFromParam, err := core.IDFromBytes(idBytes[:31])
	if err != nil {
		return false, err
	}
	if IDFromState.String() != IDFromParam.String() {
		return false, nil
	}
	return true, nil
}

func envMapper(placeholderName string) string {
	split := strings.Split(placeholderName, ":-")
	defValue := ""
	if len(split) == 2 {
		placeholderName = split[0]
		defValue = split[1]
	}

	val, ok := os.LookupEnv(placeholderName)
	if !ok {
		return defValue
	}
	return val
}
