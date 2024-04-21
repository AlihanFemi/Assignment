package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Database Database `json:"database"`
}
type Database struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
	SSL      string `json:"sslmode"`
}

func Parse() (Config, error) {
	env := os.Getenv("APP_ENV")
	filename := fmt.Sprintf("config.%v.json", env)
	file, err := os.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = json.Unmarshal(file, &config)
	return config, err
}

func (db Database) BuildConnectionString() string {
	return fmt.Sprintf(
		"user=%v dbname=%v password=%v sslmode=%v",
		db.Username,
		db.DBName,
		db.Password,
		db.SSL)
}
