package config

import (
	"encoding/json"
	"os"
)

const (
	cfgFile = ".gatorconfig.json"
)

type Config struct {
	DbUrl			string		`json:"db_url"`
	CurrentUser		string		`json:"current_user_name"`
}

func Read() (Config, error) {
	cfg, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}
	
	data, err := os.ReadFile(cfg)
	if err != nil {
		return Config{}, err
	}

	dbCfg := Config{}
	err = json.Unmarshal(data, &dbCfg)
	if err != nil {
		return Config{}, err
	}

	return dbCfg, nil
}

func (cfg *Config) SetUser (current string) error {
	cfg.CurrentUser = current
	return writeConfig(cfg)
}

// Helper Functions
func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return home + "/" + cfgFile, nil
}

func writeConfig(cfg *Config) error {
	cfgPath, err := getConfigPath()
	if err != nil {
		return err
	}

	file, err := os.Create(cfgPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	return encoder.Encode(cfg)
}
