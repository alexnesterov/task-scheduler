// Package config
package config

import "time"

type Config struct {
	App  AppConfig
	HTTP HTTPConfig
	DB   DBConfig
}

type AppConfig struct {
	WebDir string `mapstructure:"web_dir"`
}

type HTTPConfig struct {
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type DBConfig struct {
	File string `mapstructure:"file"`
}
