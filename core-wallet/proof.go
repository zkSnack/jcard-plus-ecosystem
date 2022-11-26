package walletSDK

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/iden3/go-circuits"
	"github.com/iden3/go-iden3-auth/loaders"
	"github.com/iden3/go-iden3-auth/pubsignals"
	"github.com/iden3/go-rapidsnark/types"
	"github.com/iden3/go-schema-processor/processor"

	jsonSuite "github.com/iden3/go-schema-processor/json"
	jsonldSuite "github.com/iden3/go-schema-processor/json-ld"
	"github.com/pkg/errors"
)

const (
	jsonldExt string = "json-ld"
	jsonExt   string = "json"
)

// ZKInputs are inputs for proof generation
type ZKInputs map[string]interface{}

// GenerateZkProof executes snarkjs groth16prove function and returns proof only if it's valid
func GenerateZkProof(circuitPath string, inputs ZKInputs, config *Config) (*types.ZKProof, error) {
	if path.Clean(circuitPath) != circuitPath {
		return nil, errors.Errorf("Illegal circuitPath: %s", circuitPath)
	}

	// serialize inputs into json
	inputsJSON, err := json.Marshal(inputs)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to serialize inputs into json.")
	}

	// create tmf file for inputs
	inputFile, err := ioutil.TempFile("tmp", "input-*.json")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create tmf file for inputs.")
	}
	defer os.Remove(inputFile.Name())

	// write json inputs into tmp file
	_, err = inputFile.Write(inputsJSON)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to write json inputs into tmp file.")
	}
	err = inputFile.Close()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to close json inputs tmp file.")
	}

	// create tmp witness file
	wtnsFile, err := ioutil.TempFile("tmp", "witness-*.wtns")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tmp witness file.")
	}
	defer os.Remove(wtnsFile.Name())
	err = wtnsFile.Close()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to close tmp witness file.")
	}

	// calculate witness
	wtnsCmd := exec.Command("node", config.Circuits.JS+"generate_witness.js", circuitPath+"/circuit.wasm", inputFile.Name(), wtnsFile.Name())
	wtnsOut, err := wtnsCmd.CombinedOutput()
	if err != nil {
		log.Println("Failed to calculate witness", "wtnsOut", string(wtnsOut))
		return nil, errors.Wrap(err, "Failed to calculate witness.")
	}
	log.Println("-- witness calculate completed --")

	// Proof verification

	// create tmp proof file
	proofFile, err := ioutil.TempFile("tmp", "proof-*.json")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create tmp proof file.")
	}
	defer os.Remove(proofFile.Name())
	err = proofFile.Close()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to close tmp proof file.")
	}

	// create tmp public file
	publicFile, err := ioutil.TempFile("tmp", "public-*.json")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create tmp public file.")
	}
	defer os.Remove(publicFile.Name())
	err = publicFile.Close()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to close tmp public file.")
	}

	// generate proof
	err = generateProof(circuitPath, wtnsFile, proofFile, publicFile)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to generate proof.")
	}
	log.Println("-- groth16 proof generated --")

	// verify proof
	err = verifyProof(circuitPath, proofFile, publicFile)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to verify proof.")
	}
	log.Println("-- groth16 proof verified --")

	var proof types.ProofData
	var pubSignals []string

	// read generated public signals
	publicJSON, err := os.ReadFile(publicFile.Name())
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read generated public signals.")
	}

	err = json.Unmarshal(publicJSON, &pubSignals)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal public signals.")
	}
	// read generated proof
	proofJSON, err := os.ReadFile(proofFile.Name())
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read generated proof.")
	}

	err = json.Unmarshal(proofJSON, &proof)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal generated proof.")
	}

	return &types.ZKProof{Proof: &proof, PubSignals: pubSignals}, nil
}

// VerifyZkProof executes snarkjs verify function and returns if proof is valid
func VerifyZkProof(circuitPath string, zkp *types.ZKProof) error {
	if path.Clean(circuitPath) != circuitPath {
		return errors.New("Illegal circuitPath")
	}

	// create tmp proof file
	proofFile, err := ioutil.TempFile("tmp", "proof-*.json")
	if err != nil {
		return errors.Wrap(err, "Failed to create tmp proof file")
	}
	defer os.Remove(proofFile.Name())

	// create tmp public file
	publicFile, err := ioutil.TempFile("tmp", "public-*.json")
	if err != nil {
		return errors.Wrap(err, "Failed to create tmp public file")
	}
	defer os.Remove(publicFile.Name())

	// serialize proof into json
	proofJSON, err := json.Marshal(zkp.Proof)
	if err != nil {
		return errors.Wrap(err, "Failed to serialize proof into json")
	}

	// serialize public signals into json
	publicJSON, err := json.Marshal(zkp.PubSignals)
	if err != nil {
		return errors.Wrap(err, "Failed to serialize public signals into json")
	}

	// write json proof into tmp file
	_, err = proofFile.Write(proofJSON)
	if err != nil {
		return errors.Wrap(err, "Failed to write json proof into tmp file")
	}
	err = proofFile.Close()
	if err != nil {
		return errors.Wrap(err, "Failed to close json proof tmp file")
	}

	// write json public signals into tmp file
	_, err = publicFile.Write(publicJSON)
	if err != nil {
		return errors.Wrap(err, "Failed to write json public signals into tmp file.")
	}
	err = publicFile.Close()
	if err != nil {
		return errors.Wrap(err, "Failed to close json public signals tmp file.")
	}

	// verify proof
	err = verifyProof(circuitPath, proofFile, publicFile)
	if err != nil {
		return errors.Wrap(err, "Failed to verify proof.")
	}
	log.Println("-- groth16 proof verified --")

	return nil
}

func ValidateAndGetCircuitsQuery(q pubsignals.Query, ctx context.Context, loader loaders.SchemaLoader) (*circuits.Query, error) {
	schemaBytes, ext, err := loader.Load(ctx, q.Schema)
	if err != nil {
		return nil, errors.Wrap(err, "Can't load schema for request query.")
	}

	pr, err := prepareProcessor(q.Schema.Type, ext)
	if err != nil {
		return nil, errors.Wrap(err, "Can't prepare processor for request query.")
	}
	queryReq, err := parseRequest(q.Req, schemaBytes, pr)
	if err != nil {
		return nil, errors.Wrap(err, "Can't parse request query.")
	}

	return queryReq, nil
}

func prepareProcessor(claimType, ext string) (*processor.Processor, error) {
	pr := &processor.Processor{}
	var parser processor.Parser
	switch ext {
	case jsonExt:
		parser = jsonSuite.Parser{ParsingStrategy: processor.OneFieldPerSlotStrategy}
	case jsonldExt:
		parser = jsonldSuite.Parser{ClaimType: claimType, ParsingStrategy: processor.OneFieldPerSlotStrategy}
	default:
		return nil, errors.Errorf("Process suite for schema format %s is not supported.", ext)
	}

	return processor.InitProcessorOptions(pr, processor.WithParser(parser)), nil
}

func parseRequest(req map[string]interface{}, schema []byte, pr *processor.Processor) (*circuits.Query, error) {
	if req == nil {
		return &circuits.Query{
			SlotIndex: 0,
			Values:    nil,
			Operator:  circuits.NOOP,
		}, nil
	}

	fieldName, fieldPredicate, err := extractQueryFields(req)
	if err != nil {
		return nil, err
	}

	values, operator, err := parseFieldPredicate(fieldPredicate, err)
	if err != nil {
		return nil, err
	}

	slotIndex, err := pr.GetFieldSlotIndex(fieldName, schema)
	if err != nil {
		return nil, err
	}

	return &circuits.Query{SlotIndex: slotIndex, Values: values, Operator: operator}, nil
}

func parseFieldPredicate(fieldPredicate map[string]interface{}, err error) ([]*big.Int, int, error) {
	var values []*big.Int
	var operator int
	for op, v := range fieldPredicate {

		var ok bool
		operator, ok = circuits.QueryOperators[op]
		if !ok {
			return nil, 0, errors.New("Query operator is not supported.")
		}

		values, err = getValuesAsArray(v)
		if err != nil {
			return nil, 0, err
		}

		// only one predicate for field is supported
		break
	}
	return values, operator, err
}

func extractQueryFields(req map[string]interface{}) (fieldName string, fieldPredicate map[string]interface{}, err error) {
	if len(req) > 1 {
		return "", nil, errors.New("Multiple requests not supported.")
	}

	for field, body := range req {
		fieldName = field
		var ok bool
		fieldPredicate, ok = body.(map[string]interface{})
		if !ok {
			return "", nil, errors.New("Failed cast type map[string]interface.")
		}
		if len(fieldPredicate) > 1 {
			return "", nil, errors.New("Multiple predicates for one field not supported.")
		}
		break
	}
	return fieldName, fieldPredicate, nil
}

func getValuesAsArray(v interface{}) ([]*big.Int, error) {
	var values []*big.Int

	switch value := v.(type) {
	case float64:
		values = make([]*big.Int, 1)
		values[0] = new(big.Int).SetInt64(int64(value))
	case []interface{}:
		values = make([]*big.Int, len(value))
		for i, item := range value {
			values[i] = new(big.Int).SetInt64(int64(item.(float64)))
		}
	default:
		return nil, errors.Errorf("Unsupported values type %T.", v)
	}

	return values, nil
}

func generateProof(circuitPath string, wtnsFile *os.File, proofFile *os.File, publicFile *os.File) error {
	var prog string
	var args []string

	prog = "snarkjs"
	args = append(args, "groth16", "prove")
	args = append(args, circuitPath+"/circuit_final.zkey", wtnsFile.Name(), proofFile.Name(), publicFile.Name())

	cmd := exec.Command(prog, args...)
	proveOut, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "Failed to generate proof: proveOut %s", string(proveOut))
	}

	return nil
}

func verifyProof(circuitPath string, proofFile *os.File, publicFile *os.File) error {
	cmd := exec.Command("snarkjs", "groth16", "verify", circuitPath+"/verification_key.json", publicFile.Name(), proofFile.Name())

	verifyOut, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "Failed to verify proof.")
	}
	log.Printf("snarkjs result %s", strings.TrimSpace(string(verifyOut)))

	if !strings.Contains(string(verifyOut), "OK!") {
		return errors.New("Invalid proof.")
	}

	return nil
}
