package main

import (
	"log"
	"time"

	"github.com/BurntSushi/toml"
)

type config struct {
	Account     []account
	RenewBefore duration `toml:"renew_before"`
	AutoRenew   bool     `toml:"auto_renew"`
	Report      string
	Smtp        smtpCfg
}

type duration struct {
	time.Duration
}

type smtpCfg struct {
	Username   string
	Password   string
	Hostname   string
	Port       int
	Recipients []string
}

func (d *duration) UnmarshalText(text []byte) (err error) {
	d.Duration, err = time.ParseDuration(string(text))
	return
}

func loadConfig() (c *config) {
	if _, err := toml.DecodeFile(confFile, &c); err != nil {
		log.Fatal(err)
	}
	return
}
