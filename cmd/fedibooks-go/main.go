package main

import (
	"github.com/ibrokemypie/fedibooks-go/internal/bot"
	"github.com/ibrokemypie/fedibooks-go/internal/config"
)

func main() {
	config.LoadConfig()
	bot.InitBot()
}
