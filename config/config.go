package config

import (
	"os"
	"strconv"
)

// ApplicationConfig contains the actual
// application configuration obtained from the
// environment variables
type ApplicationConfig struct {
	VaultAddress    string
	ConfigFilePath  string
	SlackWebhookURL string
	Insecure        bool
	Debug           bool
	TTLThreshold    int64
	TTLIncrement    int64
}

// Get creates a new ApplicationConfig struct with
// the actual application configuration
// obtained from the environment variables
func Get() (*ApplicationConfig, error) {
	cfg := ApplicationConfig{}

	// parse environment variables if needed
	insecure, err := strconv.ParseBool(os.Getenv("INSECURE"))
	if err != nil {
		insecure = false
	}

	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		insecure = false
	}

	ttlThreshold, err := strconv.ParseInt(os.Getenv("TTL_THRESHOLD"), 10, 64)
	if err != nil {
		ttlThreshold = 15206400 // ~6 months
	}

	ttlIncrement, err := strconv.ParseInt(os.Getenv("TTL_INCREMENT"), 10, 64)
	if err != nil {
		ttlIncrement = 5184000 // 60 days
	}

	// get relevant environment variables
	cfg.VaultAddress = os.Getenv("VAULT_ADDRESS")
	cfg.ConfigFilePath = os.Getenv("CONFIG_FILE_PATH")
	cfg.SlackWebhookURL = os.Getenv("SLACK_WEBHOOK_URL")
	cfg.Insecure = insecure // defaults to false
	cfg.Debug = debug       // defaults to false
	cfg.TTLThreshold = ttlThreshold
	cfg.TTLIncrement = ttlIncrement
	return &cfg, nil
}
