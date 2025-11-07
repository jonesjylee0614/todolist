package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Undo     UndoConfig
	CORS     CORSConfig
}

type AppConfig struct {
	Host string
	Port int
	Env  string
}

func (a AppConfig) Addr() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

type DatabaseConfig struct {
	DSN             string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

type UndoConfig struct {
	TTL time.Duration
}

type CORSConfig struct {
	AllowOrigins []string
	AllowMethods []string
	AllowHeaders []string
}

func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}

	if cfg.Undo.TTL == 0 {
		cfg.Undo.TTL = 5 * time.Second
	}

	if cfg.Database.ConnMaxLifetime == 0 {
		cfg.Database.ConnMaxLifetime = 15 * time.Minute
	}

	if len(cfg.CORS.AllowMethods) == 0 {
		cfg.CORS.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	}

	if len(cfg.CORS.AllowHeaders) == 0 {
		cfg.CORS.AllowHeaders = []string{"Content-Type", "Authorization", "X-Requested-With"}
	}

	return cfg, nil
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.host", "0.0.0.0")
	v.SetDefault("app.port", 8081)
	v.SetDefault("app.env", "development")

	v.SetDefault("database.dsn", "root:Jz@szM982io@tcp(localhost:3306)/todolist?charset=utf8mb4&parseTime=True&loc=Local")
	v.SetDefault("database.maxIdleConns", 5)
	v.SetDefault("database.maxOpenConns", 20)
	v.SetDefault("database.connMaxLifetime", "15m")

	v.SetDefault("undo.ttl", "5s")

	v.SetDefault("cors.allowOrigins", []string{"*"})
}
