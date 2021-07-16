package config

import (
	"fmt"
	"os"

	"github.com/ibrokemypie/magickbot/pkg/auth"
	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("make_post_interval", 30)
	viper.SetDefault("get_posts_interval", 30)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			pwd, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			fmt.Println("Config file not found, creating one at " + pwd)

			// attempt to write a new config file
			if err := viper.SafeWriteConfig(); err != nil {
				panic(err)
			}
		} else {
			// Config file was found but another error was produced
			panic(err)
		}
	}

	if !viper.IsSet("instance.instance_url") || !viper.IsSet("instance.access_token") {
		instanceURL, accessToken := auth.Authorize()

		viper.Set("instance.instance_url", instanceURL)
		viper.Set("instance.access_token", accessToken)
	}

	viper.WriteConfig()
}
