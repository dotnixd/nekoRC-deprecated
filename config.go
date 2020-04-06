package main

import (
	"io/ioutil"

	"github.com/goccy/go-yaml"
)

// Config declares configuration file
type Config struct {
	Distribution string `yaml:"Distribution"` // Name of distribution
	CTL          struct {
		ConnectionTimeout uint `yaml:"ConnectionTimeout"` // I/O timeout beetween nekoRC and nekoCTL
	} `yaml:"CTL"` // nekoCTL-interaction settings
}

// Load loads config
func (c *Config) Load() {
	content, err := ioutil.ReadFile("/etc/nekoRC/config.neko.yml")
	_check(err, true)
	_check(yaml.Unmarshal(content, c), true)
}
