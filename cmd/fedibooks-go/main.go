package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/ibrokemypie/fedibooks-go/internal/bot"
	"github.com/ibrokemypie/fedibooks-go/internal/config"
)

func main() {
	configFile := flag.String("c", "./config.yaml", "Path and name of config file")
	flag.Parse()
	configFilePath, err := filepath.Abs(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	config.LoadConfig(configFilePath)
	bot.InitBot()
}
