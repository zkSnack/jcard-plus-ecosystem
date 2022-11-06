package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/iden3/go-circuits"
	"github.com/iden3/go-iden3-auth/loaders"
	"github.com/iden3/go-iden3-auth/pubsignals"
	"github.com/iden3/go-schema-processor/processor"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"os/exec"
	"path"
	"strings"

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

// ZKProof is structure that represents SnarkJS library result of proof generation
type ZKProof struct {
	A        []string   `json:"pi_a"`
	B        [][]string `json:"pi_b"`
	C        []string   `json:"pi_c"`
	Protocol string     `json:"protocol"`
	Curve    string     `json:"curve"`
}

// FullProof is ZKP proof with public signals
type FullProof struct {
	Proof      *ZKProof `json:"proof"`
	PubSignals []string `json:"pub_signals"`
}

// GenerateZkProof executes snarkjs groth16prove function and returns proof only if it's valid
func GenerateZkProof(circuitPath string, inputs ZKInputs) (*FullProof, error) {

	if path.Clean(circuitPath) != circuitPath {
		return nil, fmt.Errorf("illegal circuitPath")
	}

	// serialize inputs into json
	inputsJSON, err := json.Marshal(inputs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to serialize inputs into json")
	}

	// create tmf file for inputs
	inputFile, err := ioutil.TempFile("tmp", "input-*.json")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tmf file for inputs")
	}
	defer os.Remove(inputFile.Name())

	// write json inputs into tmp file
	_, err = inputFile.Write(inputsJSON)
	if err != nil {
		return nil, errors.Wrap(err, "failed to write json inputs into tmp file")
	}
	err = inputFile.Close()
	if err != nil {
		return nil, errors.Wrap(err, "failed to close json inputs tmp file")
	}

	// create tmp witness file
	wtnsFile, err := ioutil.TempFile("tmp", "witness-*.wtns")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tmp witness file")
	}
	defer os.Remove(wtnsFile.Name())
	err = wtnsFile.Close()
	if err != nil {
		return nil, errors.Wrap(err, "failed to close tmp witness file")
	}

	// calculate witness
	wtnsCmd := exec.Command("node", "js/generate_witness.js", circuitPath+"/circuit.wasm", inputFile.Name(), wtnsFile.Name())
	wtnsOut, err := wtnsCmd.CombinedOutput()
	if err != nil {
		log.Println("failed to calculate witness", "wtnsOut", string(wtnsOut))
		return nil, errors.Wrap(err, "failed to calculate witness")
	}
	log.Println("-- witness calculate completed --")

	// create tmp proof file
	proofFile, err := ioutil.TempFile("tmp", "proof-*.json")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tmp proof file")
	}
	defer os.Remove(proofFile.Name())
	err = proofFile.Close()
	if err != nil {
		return nil, errors.Wrap(err, "failed to close tmp proof file")
	}

	// create tmp public file
	publicFile, err := ioutil.TempFile("tmp", "public-*.json")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tmp public file")
	}
	defer os.Remove(publicFile.Name())
	err = publicFile.Close()
	if err != nil {
		return nil, errors.Wrap(err, "failed to close tmp public file")
	}

	// generate proof
	var execCommandName string
	var execCommandParams []string
	execCommandName = "snarkjs"
	execCommandParams = append(execCommandParams, "groth16", "prove")
	execCommandParams = append(execCommandParams, circuitPath+"/circuit_final.zkey", wtnsFile.Name(), proofFile.Name(), publicFile.Name())
	proveCmd := exec.Command(execCommandName, execCommandParams...)
	log.Println("used prover: %s", execCommandName)
	proveOut, err := proveCmd.CombinedOutput()
	if err != nil {
		log.Println("failed to generate proof", "proveOut", string(proveOut))
		return nil, errors.Wrap(err, "failed to generate proof")
	}
	log.Println("-- groth16 prove completed --")

	// verify proof
	verifyCmd := exec.Command("snarkjs", "groth16", "verify", circuitPath+"/verification_key.json", publicFile.Name(), proofFile.Name())
	verifyOut, err := verifyCmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrap(err, "failed to verify proof")
	}
	log.Println("-- groth16 verify -- snarkjs result %s", strings.TrimSpace(string(verifyOut)))

	if !strings.Contains(string(verifyOut), "OK!") {
		return nil, errors.New("invalid proof")
	}

	var proof ZKProof
	var pubSignals []string

	// read generated public signals
	publicJSON, err := os.ReadFile(publicFile.Name())
	if err != nil {
		return nil, errors.Wrap(err, "failed to read generated public signals")
	}

	err = json.Unmarshal(publicJSON, &pubSignals)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal public signals")
	}
	// read generated proof
	proofJSON, err := os.ReadFile(proofFile.Name())
	if err != nil {
		return nil, errors.Wrap(err, "failed to read generated proof")
	}

	err = json.Unmarshal(proofJSON, &proof)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal generated proof")
	}

	return &FullProof{Proof: &proof, PubSignals: pubSignals}, nil
}

// VerifyZkProof executes snarkjs verify function and returns if proof is valid
func VerifyZkProof(circuitPath string, zkp *FullProof) error {

	if path.Clean(circuitPath) != circuitPath {
		return fmt.Errorf("illegal circuitPath")
	}

	// create tmp proof file
	proofFile, err := ioutil.TempFile("tmp", "proof-*.json")
	if err != nil {
		return errors.Wrap(err, "failed to create tmp proof file")
	}
	defer os.Remove(proofFile.Name())

	// create tmp public file
	publicFile, err := ioutil.TempFile("tmp", "public-*.json")
	if err != nil {
		return errors.Wrap(err, "failed to create tmp public file")
	}
	defer os.Remove(publicFile.Name())

	// serialize proof into json
	proofJSON, err := json.Marshal(zkp.Proof)
	if err != nil {
		return errors.Wrap(err, "failed to serialize proof into json")
	}

	// serialize public signals into json
	publicJSON, err := json.Marshal(zkp.PubSignals)
	if err != nil {
		return errors.Wrap(err, "failed to serialize public signals into json")
	}

	// write json proof into tmp file
	_, err = proofFile.Write(proofJSON)
	if err != nil {
		return errors.Wrap(err, "failed to write json proof into tmp file")
	}
	err = proofFile.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close json proof tmp file")
	}

	// write json public signals into tmp file
	_, err = publicFile.Write(publicJSON)
	if err != nil {
		return errors.Wrap(err, "failed to write json public signals into tmp file")
	}
	err = publicFile.Close()
	if err != nil {
		return errors.Wrap(err, "failed to close json public signals tmp file")
	}

	// verify proof
	verifyCmd := exec.Command("snarkjs", "groth16", "verify", circuitPath+"/verification_key.json", publicFile.Name(), proofFile.Name())
	verifyOut, err := verifyCmd.CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "failed to verify proof")
	}
	log.Println("-- groth16 verify -- snarkjs result %s", strings.TrimSpace(string(verifyOut)))

	if !strings.Contains(string(verifyOut), "OK!") {
		return errors.New("invalid proof")
	}

	return nil
}

func ValidateAndGetCircuitsQuery(q pubsignals.Query, ctx context.Context, loader loaders.SchemaLoader) (*circuits.Query, error) {

	schemaBytes, ext, err := loader.Load(ctx, q.Schema)
	if err != nil {
		return nil, errors.Wrap(err, "can't load schema for request query")
	}

	pr, err := prepareProcessor(q.Schema.Type, ext)
	if err != nil {
		return nil, errors.Wrap(err, "can't prepare processor for request query")
	}
	queryReq, err := parseRequest(q.Req, schemaBytes, pr)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse request query")
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
		return nil, errors.Errorf(
			"process suite for schema format %s is not supported", ext)
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
			return nil, 0, errors.New("query operator is not supported")
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
		return "", nil, errors.New("multiple requests not supported")
	}

	for field, body := range req {
		fieldName = field
		var ok bool
		fieldPredicate, ok = body.(map[string]interface{})
		if !ok {
			return "", nil, errors.New("failed cast type map[string]interface")
		}
		if len(fieldPredicate) > 1 {
			return "", nil, errors.New("multiple predicates for one field not supported")
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
		return nil, errors.Errorf("unsupported values type %T", v)
	}

	return values, nil
}
