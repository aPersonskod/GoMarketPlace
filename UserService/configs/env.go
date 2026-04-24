package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost     string
	Port           string
	DbUser         string
	DbPassword     string
	DbName         string
	JwtSecret      string
	ExpirationMins int
}

var Env = initConfig()

func initConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	expirationMins, err := strconv.Atoi(os.Getenv("EXPIRATIONMINS"))
	if err != nil {
		expirationMins = 30
	}

	return &Config{
		PublicHost:     os.Getenv("PUBLIC_HOST"),
		Port:           os.Getenv("PORT"),
		DbUser:         os.Getenv("DBUSER"),
		DbPassword:     os.Getenv("DBPASSWORD"),
		DbName:         os.Getenv("DBNAME"),
		JwtSecret:      os.Getenv("JWTSECRET"),
		ExpirationMins: expirationMins,
	}
}
