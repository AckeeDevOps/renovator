package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/AckeeDevOps/renovator/client"
	"github.com/AckeeDevOps/renovator/config"
	"github.com/AckeeDevOps/renovator/notifier"
)

func main() {
	log.Println("Renovator started ...")

	// create configuration struct
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("could not create configuration struct: %s", err)
	}

	// show configuration if needed
	if cfg.Debug {
		log.Println(cfg)
	}

	// create Valt client
	vaultClient, err := client.NewClient(cfg.VaultAddress, cfg.Insecure, cfg.Debug)
	if err != nil {
		log.Fatalf("could not initialize Vault client: %s", err)
	}

	// create registry
	registry := notifier.NewRegistry()

	// read tokens from JSON
	var tokenConfig TokenConfig
	jsonFile, err := os.Open(cfg.ConfigFilePath)
	if err != nil {
		log.Fatalf("could not open config file: %s", err)
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("could not read config file: %s", err)
	}

	err = json.Unmarshal(byteValue, &tokenConfig)
	if err != nil {
		log.Fatalf("could not unmarshal config file: %s", err)
	}

	// go through all tokens
	for _, t := range tokenConfig.Tokens {
		status, err := vaultClient.CheckOrRenew(t.Token, cfg.TTLThreshold, cfg.TTLIncrement)
		if err != nil {
			message := fmt.Sprintf("Could not renew; display name: %s; error: %s", t.Name, err.Error())
			registry.AddStatus(t.Token, false, message)
			continue
		}

		// renewed token
		days := status.TTL / 60 / 60 / 24
		message := fmt.Sprintf("new/current TTL is %d (%d days); display name: %s", status.TTL, days, t.Name)
		registry.AddStatus(t.Token, true, message)
	}

	err = notifier.NotifySlack(registry, cfg.SlackWebhookURL)
	if err != nil {
		log.Fatalf("could not send Slack message: %s", err)
	}

	log.Println("done.")
}
