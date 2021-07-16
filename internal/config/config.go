package config

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/ibrokemypie/magickbot/pkg/auth"
	"github.com/spf13/viper"
)

func LoadConfig(configFile string) {
	viper.SetConfigName(strings.TrimSuffix(filepath.Base(configFile), filepath.Ext(configFile)))
	viper.SetConfigType(strings.TrimPrefix(filepath.Ext(configFile), "."))
	viper.AddConfigPath(filepath.Dir(configFile))

	viper.SetDefault("post.make_interval", 30)
	viper.SetDefault("post.max_words", 30)
	viper.SetDefault("post.visibility", "unlisted")
	viper.SetDefault("history.get_interval", 30)
	viper.SetDefault("history.learn_from_cw", false)
	viper.SetDefault("history.file_path", "./history.gob")
	viper.SetDefault("history.max_length", 100000)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found, creating one at " + filepath.Dir(configFile))

			// attempt to write a new config file
			if err := viper.SafeWriteConfig(); err != nil {
				log.Fatal(err)
			}
		} else {
			// Config file was found but another error was produced
			log.Fatal(err)
		}
	}

	if !viper.IsSet("instance.url") || !viper.IsSet("instance.access_token") {
		instanceURL, accessToken := auth.Authorize()

		viper.Set("instance.url", instanceURL)
		viper.Set("instance.access_token", accessToken)
	}

	viper.WriteConfig()
}
