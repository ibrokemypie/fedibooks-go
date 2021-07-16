package bot

import (
	"fmt"
	"time"

	"github.com/ibrokemypie/fedibooks-go/internal/fedi"
	"github.com/spf13/viper"
)

func InitBot() {
	instanceURL := viper.GetString("instance.url")
	accessToken := viper.GetString("instance.access_token")
	getPostInterval := viper.GetInt("get_posts_interval")
	makePostInterval := viper.GetInt("make_post_interval")
	learnFromCW := viper.GetBool("learn_from_cw")
	historyFilePath := viper.GetString("history.file_path")
	maxStoredStatuses := viper.GetInt("history.max_length")
	postVisibility := viper.GetString("post_visibility")

	history := LoadFromGob(historyFilePath)

	go GetStatusesLoop(history, historyFilePath, instanceURL, accessToken, getPostInterval, learnFromCW, maxStoredStatuses)

	go PostQuotesLoop(history, instanceURL, accessToken, makePostInterval, postVisibility)

	go HandleReplies(history, instanceURL, accessToken, postVisibility)

	select {}
}

func HandleReplies(history *History, instanceURL, accessToken, postVisibility string) {
	notificationChannel := make(chan fedi.Notification)
	go fedi.NotificationStream(notificationChannel, instanceURL, accessToken)

	for {
		notification := <-notificationChannel
		botUser, err := fedi.GetCurrentUser(instanceURL, accessToken)
		if err != nil {
			fmt.Println(err)
			continue
		}
		followedUsers, err := fedi.GetUserFollowing(botUser, instanceURL, accessToken)
		if err != nil {
			fmt.Println(err)
			continue
		}
		quote := GenQuote(history, followedUsers)
		err = fedi.PostStatus(quote, postVisibility, notification.Status, instanceURL, accessToken)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func GetStatusesLoop(history *History, historyFilePath string, instanceURL, accessToken string, interval int, learnFromCW bool, maxStoredStatuses int) {
	for {
		GetNewStatuses(history, historyFilePath, instanceURL, accessToken, learnFromCW, maxStoredStatuses)
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}

func PostQuotesLoop(history *History, instanceURL, accessToken string, interval int, postVisibility string) {
	for {
		botUser, err := fedi.GetCurrentUser(instanceURL, accessToken)
		if err != nil {
			fmt.Println(err)
			continue
		}
		followedUsers, err := fedi.GetUserFollowing(botUser, instanceURL, accessToken)
		if err != nil {
			fmt.Println(err)
			continue
		}
		quote := GenQuote(history, followedUsers)
		err = fedi.PostStatus(quote, postVisibility, fedi.Status{}, instanceURL, accessToken)
		if err != nil {
			fmt.Println(err)
			continue
		}
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}
