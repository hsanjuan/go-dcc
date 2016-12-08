package dcc

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// DCCConfig allows to store configuration settings to initialize go-dcc.
type DCCConfig struct {
	Locomotives []*Locomotive `json:"locomotives"`
}

// LoadConfig parses a configuration file and returns a DCCConfig object.
func LoadConfig(path string) (*DCCConfig, error) {
	conf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var dccConfig DCCConfig
	err = json.Unmarshal(conf, &dccConfig)
	if err != nil {
		return nil, err
	}
	log.Printf("XLoaded configuration for %d locomotive(s)", len(dccConfig.Locomotives))
	return &dccConfig, nil
}
