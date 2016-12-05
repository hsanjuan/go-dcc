package dcc

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var DefaultConfigPath = "~/.config/go-dcc/config"

type DCCConfig struct {
	Locomotives []*Locomotive `json:"locomotives"`
}

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
	log.Printf("Loaded configuration for %d locomotive(s)", len(dccConfig.Locomotives))
	return &dccConfig, nil
}
