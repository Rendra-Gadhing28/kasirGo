package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName              string
	AppEnv               string
	Port                 string
	AppURL               string
	CORSAllowedOrigins   string
	DatabaseURL          string
	DBHost               string
	DBPort               string
	DBUser               string
	DBPassword           string
	DBName               string
	DBSSL                string
	JWTSecret            string
	JWTAccessExpireMins  int
	JWTRefreshExpireDays int
	MidtransServerKey    string
	MidtransClientKey    string
	MidtransIsProduction bool
	DefaultTaxRate       float64
	DefaultCurrency      string
	StoreName            string
	StoreAddress         string
	StorePhone           string
}

var AppConfig *Config

func LoadConfig() *Config {
	// Try loading .env from multiple standard locations
	_ = godotenv.Load()
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("../../.env")

	jwtAccessMins, _ := strconv.Atoi(getEnv("JWT_ACCESS_EXPIRE_MINUTES", "1440"))
	jwtRefreshDays, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRE_DAYS", "7"))
	taxRate, _ := strconv.ParseFloat(getEnv("DEFAULT_TAX_RATE", "11"), 64)
	isProd, _ := strconv.ParseBool(getEnv("MIDTRANS_IS_PRODUCTION", "false"))

	cfg := &Config{
		AppName:              getEnv("APP_NAME", "KasirPro"),
		AppEnv:               getEnv("APP_ENV", "development"),
		Port:                 getEnv("PORT", "8080"),
		AppURL:               getEnv("APP_URL", "https://fish-warming-logos-lots.trycloudflare.com"),
		CORSAllowedOrigins:   getEnv("CORS_ALLOWED_ORIGINS", "https://fish-warming-logos-lots.trycloudflare.com,http://localhost:5173,http://localhost:3000,http://localhost"),
		DatabaseURL:          getEnv("DATABASE_URL", ""),
		DBHost:               getEnv("DB_HOST", "127.0.0.1"),
		DBPort:               getEnv("DB_PORT", "3306"),
		DBUser:               getEnv("DB_USER", "root"),
		DBPassword:           getEnv("DB_PASSWORD", "rootpassword"),
		DBName:               getEnv("DB_NAME", "kasirpro"),
		DBSSL:                getEnv("DB_SSL", "false"),
		JWTSecret:            getEnv("JWT_SECRET", "super_secret_jwt_key_kasirpro_256bit_minimum_length_required!"),
		JWTAccessExpireMins:  jwtAccessMins,
		JWTRefreshExpireDays: jwtRefreshDays,
		MidtransServerKey:    getEnv("MIDTRANS_SERVER_KEY", "SB-Mid-server-simulator-mock-key-12345"),
		MidtransClientKey:    getEnv("MIDTRANS_CLIENT_KEY", "SB-Mid-client-simulator-mock-key-12345"),
		MidtransIsProduction: isProd,
		DefaultTaxRate:       taxRate,
		DefaultCurrency:      getEnv("DEFAULT_CURRENCY", "IDR"),
		StoreName:            getEnv("STORE_NAME", "KasirPro Modern Market"),
		StoreAddress:         getEnv("STORE_ADDRESS", "Jl. Sudirman No. 123, Jakarta Selatan"),
		StorePhone:           getEnv("STORE_PHONE", "081234567890"),
	}

	AppConfig = cfg
	return cfg
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}
