package walletSDK

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func GetConfig(filename string) (*Config, error) {
	yfile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.New("Failed to open the config file.")
	}

	config := new(Config)
	err = yaml.Unmarshal(yfile, config)
	if err != nil {
		return nil, errors.New("Failed to unmarshal the yaml file.")
	}

	return config, nil
}

func GetIdentity(filename string) (*Identity, error) {
	var identity *Identity
	if _, err := os.Stat(filename); err == nil {
		identity, err = LoadIdentityFromFile(filename)
		if err != nil {
			log.Fatalf("Failed to load identity from the File. Err %s", err)
		}
		log.Println("Account loaded from saved file: account.json")
	} else {
		identity, err = GenerateAccount()
		if err != nil {
			log.Fatalf("Error %s. Failed to create new identity. Aborting...", err)
		}
	}

	return identity, nil
}

func GenerateAccount() (*Identity, error) {
	if identity, err := NewIdentity(); err == nil {
		err = DumpIdentity(identity)
		if err != nil {
			return nil, err
		}
		return identity, nil
	} else {
		return nil, errors.New("Failed to create new identity")
	}
}

func DumpIdentity(identity *Identity) error {
	file, err := json.MarshalIndent(identity, "", "	")
	if err != nil {
		return errors.New("Failed to json MarshalIdent identity struct")
	}
	err = ioutil.WriteFile("account.json", file, 0644)
	if err != nil {
		return errors.New("Failed to write identity state to the file")
	}
	log.Println("Account.json updated to latest identity state")

	return nil
}
