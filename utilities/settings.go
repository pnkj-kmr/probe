package util

import (
	"log"

	"github.com/spf13/viper"
)

type _yaml struct {
	Name      string `mapstructure:"name"`
	Port      int    `mapstructure:"port"`
	Debug     bool   `mapstructure:"debug"`
	SecretKey []byte `mapstructure:"secret_key"`
	Expiry    int    `mapstructure:"expiry"`
}

var Settings _yaml

func init() {
	// Set config file name and path
	viper.SetConfigName("probe")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	// Unmarshal config into the struct
	if err := viper.Unmarshal(&Settings); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	if len(Settings.SecretKey) == 0 {
		// 32 bytes - key
		// openssl rand -hex 32
		// openssl rand -base64 32
		Settings.SecretKey = []byte("probe_ilBTO6toliukgkjtdfKJRbiu")
	}

	if Settings.Expiry == 0 {
		Settings.Expiry = 1
	}

	log.Println("Loaded struct:", Settings)
}
