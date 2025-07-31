package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Security SecurityConfig `mapstructure:"security"`
}

type ServerConfig struct {
    Port            int    `mapstructure:"port"`
    UpdateInterval  int    `mapstructure:"update_interval"`
    AllowedOrigins []string `mapstructure:"allowed_origins"`
}

type SecurityConfig struct {
    JWTSecret     string   `mapstructure:"jwt_secret"`
    AllowedIPs    []string `mapstructure:"allowed_ips"`
    SSLCertFile   string   `mapstructure:"ssl_cert_file"`
    SSLKeyFile    string   `mapstructure:"ssl_key_file"`
}

func Load() (*Config, error) {
    viper.SetDefault("server.port", 8443)
    viper.SetDefault("server.update_interval", 5)
    viper.SetDefault("server.allowed_origins", []string{"http://localhost:3000"})

    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("/etc/server-monitor/")

    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, err
        }
    }

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }

    return &config, nil
}
