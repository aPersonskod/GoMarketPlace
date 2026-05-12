package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost               string
	Port                     string
	JwtSecret                string
	BuyServiceAddressDev     string
	OrderServiceAddressDev   string
	ProductServiceAddressDev string
	UserServiceAddressDev    string
}

var Env = initConfig()

func initConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	return &Config{
		PublicHost:               os.Getenv("PUBLIC_HOST"),
		Port:                     os.Getenv("PORT"),
		JwtSecret:                os.Getenv("JWTSECRET"),
		BuyServiceAddressDev:     os.Getenv("BUY_SERVICE_ADDRESS_DEV"),
		OrderServiceAddressDev:   os.Getenv("ORDER_SERVICE_ADDRESS_DEV"),
		ProductServiceAddressDev: os.Getenv("PRODUCT_SERVICE_ADDRESS_DEV"),
		UserServiceAddressDev:    os.Getenv("USER_SERVICE_ADDRESS_DEV"),
	}
}
