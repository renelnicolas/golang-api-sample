package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// DatabaseConfig :
type DatabaseConfig struct {
	Host      string
	Port      int
	User      string
	Pwd       string
	Schema    string
	Connector string
}

// AmqConfig :
type AmqConfig struct {
	Host      string
	Port      int
	User      string
	Pwd       string
	Connector string
}

// DomainConfig :
type DomainConfig struct {
	Schema string
	Host   string
}

// Config :
type Config struct {
	Database  DatabaseConfig
	Amq       AmqConfig
	Domain    DomainConfig
	DebugMode bool
	Env       string
}

var (
	config *Config
)

// GetConfig :
func GetConfig() *Config {
	return config
}

// New : New returns a new Config struct
func New(path string) *Config {
	if err := godotenv.Load(path); nil != err {
		log.Print("No .env file found")
	}

	config = &Config{
		Database: DatabaseConfig{
			Host:   getEnv("DB_HOST", "localhost"),
			Port:   getEnvAsInt("DB_PORT", 3306),
			User:   getEnv("DB_USER", "root"),
			Pwd:    getEnv("DB_PWD", "root"),
			Schema: getEnv("DB_SCHEMA", "optimiads"),
		},
		Amq: AmqConfig{
			Host: getEnv("AMQ_HOST", "localhost"),
			Port: getEnvAsInt("AMQ_PORT", 5672),
			User: getEnv("AMQ_USER", "guest"),
			Pwd:  getEnv("AMQ_PWD", "guest"),
		},
		Domain: DomainConfig{
			Schema: getEnv("HTTP_SCHEMA", "http"),
			Host:   getEnv("HTTP_DOMAIN", "ohmytech.local"),
		},
		DebugMode: getEnvAsBool("DEBUG_MODE", true),
		Env:       getEnv("ENV", "dev"),
	}

	toDatabaseConnector(&config.Database)
	toAmqConnector(&config.Amq)

	return config
}

// toDatabaseConnector :
func toDatabaseConnector(db *DatabaseConfig) {
	db.Connector = db.User + `:` + db.Pwd + `@tcp(` + db.Host + `:` + strconv.Itoa(db.Port) + `)/` + db.Schema
}

// toAmqConnector :
func toAmqConnector(amq *AmqConfig) {
	amq.Connector = `amqp://` + amq.User + `:` + amq.Pwd + `@` + amq.Host + `:` + strconv.Itoa(amq.Port) + `/`
}

// getEnv : Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

// getEnvAsInt : Simple helper function to read an environment variable into integer or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); nil == err {
		return value
	}

	return defaultVal
}

// getEnvAsBool : Helper to read an environment variable into a bool or return default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); nil == err {
		return val
	}

	return defaultVal
}

// getEnvAsSlice : Helper to read an environment variable into a string slice or return default value
func getEnvAsSlice(name string, defaultVal []string, sep string) []string {
	valStr := getEnv(name, "")

	if valStr == "" {
		return defaultVal
	}

	val := strings.Split(valStr, sep)

	return val
}
