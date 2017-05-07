package config

import (
	"log"
	"os"
)

// Config represents configuration information
type Config struct {
	BotID      string
	SlovnikURL string
	WebhookURL string
}

const (
	envBotID      = "SLOVNIK_BOT_ID"
	envAPIURL     = "SLOVNIK_API_URL"
	envWebhookURL = "SLOVNIK_WEBHOOK_URL"
)

// Setup configuration
func Setup() *Config {
	botID, ok := os.LookupEnv(envBotID)
	if !ok {
		log.Panic(envBotID + " is not set!")
	}

	slovnikURL, ok := os.LookupEnv(envAPIURL)

	if !ok {
		log.Panic(envAPIURL + " is not set!")
	}

	webhookURL := os.Getenv(envWebhookURL)

	return &Config{
		BotID:      botID,
		SlovnikURL: slovnikURL,
		WebhookURL: webhookURL,
	}
}
