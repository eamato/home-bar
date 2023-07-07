package configs

import (
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"home-bar/internal"
	"os"
	"strconv"
	"strings"
)

const (
	EVNameConfigFile = "CONFIG_FILE"

	EVNameServerHost = "SERVER_HOST"
	EVNameServerPort = "SERVER_PORT"

	EVNameAccessTokenExpiryHour  = "TOKEN_ACCESS_TOKEN_EXPIRY_HOUR"
	EVNameRefreshTokenExpiryHour = "TOKEN_REFRESH_TOKEN_EXPIRY_HOUR"
	EVNameAccessTokenSecret      = "TOKEN_ACCESS_TOKEN_SECRET"
	EVNameRefreshTokenSecret     = "TOKEN_REFRESH_TOKEN_SECRET"

	EVNameDatabaseHost           = "DATABASE_HOST"
	EVNameDatabasePort           = "DATABASE_PORT"
	EVNameDatabaseUser           = "DATABASE_USER"
	EVNameDatabasePassword       = "DATABASE_PASSWORD"
	EVNameDatabaseName           = "DATABASE_NAME"
	EVNameDatabaseDumpToolPath   = "DATABASE_DUMP_TOOL_PATH"
	EVNameDatabaseBackupToolPath = "DATABASE_BACKUP_TOOL_PATH"

	EVNameOAuthWebClientID     = "OAUTH_WEB_CLIENT_ID"
	EVNameOAuthWebClientSecret = "OAUTH_WEB_CLIENT_SECRET"
	EVNameOAuthRedirectURL     = "OAUTH_REDIRECT_URL"
)

// Config struct holds all configuration things from config file
type Config struct {
	ServerConfig   ServeConfig
	DatabaseConfig DatabaseConfig
	TokenConfig    TokenConfig
	OAuthConfig    *oauth2.Config
}

type ServeConfig struct {
	Host string
	Port string
}

type DatabaseConfig struct {
	Host           string
	Port           string
	User           string
	Password       string
	Name           string
	DumpToolPath   string
	BackupToolPath string
}

type TokenConfig struct {
	AccessTokenExpiryHour  int
	RefreshTokenExpiryHour int
	AccessTokenSecret      string
	RefreshTokenSecret     string
}

func NewConfig() *Config {
	configFilePath := os.Getenv(EVNameConfigFile)

	if configFilePath == "" {
		internal.PrintWarning("%s not set, load .env.dev", EVNameConfigFile)

		configFilePath = "./.env.dev"
	}

	internal.PrintMessage("%s loading %s", EVNameConfigFile, configFilePath)

	env, err := godotenv.Read(configFilePath)
	if err != nil {
		internal.PrintFatal(fmt.Sprintf("failed to read %s", configFilePath), err)
	}

	return &Config{
		ServerConfig: ServeConfig{
			Host: getStringValue(env, EVNameServerHost),
			Port: getStringValue(env, EVNameServerPort),
		},
		DatabaseConfig: DatabaseConfig{
			Host:           getStringValue(env, EVNameDatabaseHost),
			Port:           getStringValue(env, EVNameDatabasePort),
			User:           getStringValue(env, EVNameDatabaseUser),
			Password:       getStringValue(env, EVNameDatabasePassword),
			Name:           getStringValue(env, EVNameDatabaseName),
			DumpToolPath:   getStringValue(env, EVNameDatabaseDumpToolPath),
			BackupToolPath: getStringValue(env, EVNameDatabaseBackupToolPath),
		},
		TokenConfig: TokenConfig{
			AccessTokenExpiryHour:  getIntValue(env, EVNameAccessTokenExpiryHour),
			RefreshTokenExpiryHour: getIntValue(env, EVNameRefreshTokenExpiryHour),
			AccessTokenSecret:      getStringValue(env, EVNameAccessTokenSecret),
			RefreshTokenSecret:     getStringValue(env, EVNameRefreshTokenSecret),
		},
		OAuthConfig: newOAuthConfig(env),
	}
}

func newOAuthConfig(env map[string]string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     getStringValue(env, EVNameOAuthWebClientID),
		ClientSecret: getStringValue(env, EVNameOAuthWebClientSecret),
		RedirectURL:  getStringValue(env, EVNameOAuthRedirectURL),
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
}

func printConfigError(configKey string) {
	internal.PrintWarning("neither config file nor environment variables contains config key %s", configKey)
}

func getStringValue(configFile map[string]string, configKey string) string {
	var result string

	value, ok := os.LookupEnv(configKey)
	if ok {
		result = value
	} else {
		result = configFile[configKey]
	}

	if result == "" {
		printConfigError(configKey)
	}

	return result
}

//nolint:unused
func getStringArrayValue(configFile map[string]string, configKey string) []string {
	return strings.Split(getStringValue(configFile, configKey), ",")
}

func getBoolValue(configFile map[string]string, configKey string) bool {
	var result bool
	var err error

	value, ok := os.LookupEnv(configKey)
	if ok {
		result, err = strconv.ParseBool(value)
	} else {
		result, err = strconv.ParseBool(configFile[configKey])
	}

	if err != nil {
		printConfigError(configKey)
	}

	return result
}

func getIntValue(configFile map[string]string, configKey string) int {
	var result int
	var err error

	value, ok := os.LookupEnv(configKey)
	if ok {
		result, err = strconv.Atoi(value)
	} else {
		result, err = strconv.Atoi(configFile[configKey])
	}

	if err != nil {
		printConfigError(configKey)
	}

	return result
}
