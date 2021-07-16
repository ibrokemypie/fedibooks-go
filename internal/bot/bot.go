package bot

import (
	"time"

	"github.com/ibrokemypie/fedibooks-go/internal/fedi"
	"github.com/spf13/viper"
)

func InitBot() {
	instanceURL := viper.GetString("instance.instance_url")
	accessToken := viper.GetString("instance.access_token")
	getPostInterval := viper.GetInt("get_posts_interval")
	makePostInterval := viper.GetInt("make_post_interval")
	learnFromCW := viper.GetBool("learn_from_cw")
	historyFilePath := viper.GetString("history.file_path")
	maxStoredStatuses := viper.GetInt("history.length")

	history := LoadFromGob(historyFilePath)

	go GetStatusesLoop(&history, historyFilePath, instanceURL, accessToken, getPostInterval, learnFromCW, maxStoredStatuses)

	go PostQuotesLoop(&history.Statuses, instanceURL, accessToken, makePostInterval)

	select {}
}

func GetStatusesLoop(history *History, historyFilePath string, instanceURL, accessToken string, interval int, learnFromCW bool, maxStoredStatuses int) {
	for {
		GetNewStatuses(history, historyFilePath, instanceURL, accessToken, learnFromCW, maxStoredStatuses)
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}

func PostQuotesLoop(statuses *map[string]HistoryStatus, instanceURL, accessToken string, interval int) {
	for {
		quote := GenQuote(statuses)
		fedi.PostStatus(quote, "unlisted", fedi.Status{}, instanceURL, accessToken)
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}
