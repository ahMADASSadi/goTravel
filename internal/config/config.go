package config

import "path/filepath"

type Config struct {
	AppName  string
	HTTPPort uint
	DBDriver string
	DBSource string
}

var Cfg *Config

func LoadConfig() *Config {
	dbPath, _ := filepath.Abs("./ticketings.sqlite3")

	return &Config{
		AppName:  "Ticketing",
		HTTPPort: 8000,
		DBDriver: "sqlite3",
		DBSource: dbPath,
	}
}
