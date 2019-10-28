package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/hsanjuan/go-dcc"
)

// LoadConfig parses a configuration file and returns a Config object.
func LoadConfig(path string) (dcc.Config, error) {
	conf, err := ioutil.ReadFile(path)
	if err != nil {
		return dcc.Config{}, err
	}
	var cfg dcc.Config
	err = json.Unmarshal(conf, &cfg)
	if err != nil {
		return dcc.Config{}, err
	}
	log.Printf("Loaded configuration for %d locomotive(s)", len(cfg.Locomotives))
	return cfg, nil
}

// SaveConfig stores a Config object in the given path.
func SaveConfig(cfg dcc.Config, path string) error {
	pretty, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, pretty, 0644)
	return err
}
