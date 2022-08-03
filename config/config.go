package config

import (
	"fmt"
	"hamza72x/bankify/util"
	"path"
	"path/filepath"

	"github.com/spf13/viper"
)

const DB_DRIVER = "postgres"

type Config struct {
	GIN_MODE    string `mapstructure:"GIN_MODE"`
	PORT        string `mapstructure:"PORT"`
	DB_HOST     string `mapstructure:"DB_HOST"`
	DB_PORT     string `mapstructure:"DB_PORT"`
	DB_USERNAME string `mapstructure:"DB_USERNAME"`
	DB_PASSWORD string `mapstructure:"DB_PASSWORD"`
	DB_NAME     string `mapstructure:"DB_NAME"`
}

func (c *Config) GetDBUrl() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.DB_USERNAME, c.DB_PASSWORD,
		c.DB_HOST, c.DB_PORT,
		c.DB_NAME,
	)
}

// @param fileName: just name, not full path
func LoadConfig(fileName string) (config Config, err error) {

	dir, err := filepath.Abs(".")
	util.LogFatal(err, "failed to get current directory")

	viper.SetConfigFile(path.Join(dir, fileName))

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
