package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MaganghubConfig MaganghubConfig
	CronSchedule    string
	Headless        bool
}

type MaganghubConfig struct {
	Username string
	Password string
}

func Load() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	cronSchedule := os.Getenv("CRON_SCHEDULE")
	if cronSchedule == "" {
		cronSchedule = "0 8 * * 1-5" // Default: 8am on weekdays
	}

	headless := os.Getenv("HEADLESS") != "false" // Default: true

	cfg := &Config{
		MaganghubConfig: MaganghubConfig{
			Username: os.Getenv("MAGANGHUB_USERNAME"),
			Password: os.Getenv("MAGANGHUB_PASSWORD"),
		},
		CronSchedule: cronSchedule,
		Headless:     headless,
	}

	if cfg.MaganghubConfig.Username == "" || cfg.MaganghubConfig.Password == "" {
		log.Fatal("MAGANGHUB_USERNAME or MAGANGHUB_PASSWORD is not set")
	}

	return cfg
}
