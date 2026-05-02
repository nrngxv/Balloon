package config

import (
	"os"
	"strings"
	
	"github.com/knadh/koanf/providers/env"
	"github.com/go-playground/validator/v10"
	_"github.com/joho/godotenv/autoload"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

type Config struct{
	Primary			Primary			`koanf:"primary" validate:"required"`
	Server			ServerConfig	`koanf:"server" validate:"required"`
	Database		DatabaseConfig	`koanf:"database" validate:"required"`
	Env				string 			`koanf:"env" validate:"required"`
	Auth			AuthConfig		`koanf:"Auth" validate:"required"`
	Redis 			RedisConfig		`koand:"redis" validate:"required"`
}

type ServerConfig struct{
	Port               string 		`koanf:"port" validate:"required"`
	ReadTimeout        int	  		`koanf:"read_timeout" validate:"required"`
	WriteTimeout  	   int	  		`koanf:"write_timeout" validate:"required"`
	IdleTimeout        int    		`koanf:"idle_timeout" validate:"required"`
	CORSAllowedOrigins []string		`koanf:"cors_allowed_origins" validate:'"required"`
	Redis			   RedisConfig 	`koanf:"redis" validate"required`
}

type DatabaseConfig struct{
	Host				string	`koanf:"host" validate:"required"`
	Port				string	`koanf:"port" validate:"required"`
	User 				string	`koanf:"user" validate:"required"`
	Password			string	`koanf:"password"`
	Name				string	`koanf:"name" validate:"required"`
	SSLMode				string	`koanf:"ssl_mode" validate:"required"`
	MaxOpenConns		int		`koanf:"max_open_conns" validate:"required"`
	MaxIdleConns		int		`koanf:"max_idle_conns" validate:"required"`
	ConnMaxLifetime		int		`koanf:"conn_max_lifetime" validate:"required"`
	ConnMaxIdleLifetime int		`koanf:"conn_max_idle_time" validate:"required"`
}

// for managing redis instances
type RedisConfig struct{
	Address string `koanf:"address" validate:"required"`
}

type IntegrationConfig struct {
	ResendAPIKey string
}

// for managing authentication instances
type AuthConfig struct{
	Address string `koanf:"address" validate:"required"`
}

func LoadConfig() (*Config, error){
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	k := koanf.New(".")

	// Prefixing all the env variables with the name of the app "BOILERPLATE_"
	err := k.Load(env.Provider("BALLOON_", ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "BALLOON_"))
	}), nil)

	if err != nil {
		logger.Fatal().Err(err).Msg("Could not load initial env variables")
	}

	// for taking data from k to config
	mainConfig := &Config{}

	// Basically taking all the variables stored in k, and put it in config struct, which then transfer to the children structs
	err = k.Unmarshal("", mainConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("Could not unmarshall main config")
	}

	validate := validator.New()

	err = validate.Struct(mainConfig)
	
	if err != nil{
		logger.Fatal().Err(err).Msg("Could not validate Config Struct")
	}

	return mainConfig, nil
} 
