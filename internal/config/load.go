package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("TSCHED")

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return nil, fmt.Errorf("config: read file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config: unmarshal: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("config: validate: %w", err)
	}

	return &cfg, nil
}

func setDefaults() {
	viper.SetDefault("app.web_dir", "web")

	viper.SetDefault("http.port", 7540)
	viper.SetDefault("http.read_timeout", 5*time.Second)
	viper.SetDefault("http.write_timeout", 10*time.Second)
	viper.SetDefault("http.idle_timeout", 15*time.Second)

	viper.SetDefault("db.file", "scheduler.db")
}

func validate(cfg *Config) error {
	if cfg.DB.File == "" {
		return fmt.Errorf("db.file is required")
	}
	return nil
}
