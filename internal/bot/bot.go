package bot

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ibrokemypie/fedibooks-go/internal/fedi"
	"github.com/spf13/viper"
)

func InitBot() {
	instanceURL := viper.GetString("instance.url")
	accessToken := viper.GetString("instance.access_token")
	makePostInterval := viper.GetInt("post.make_interval")
	postVisibility := viper.GetString("post.visibility")
	maxWords := viper.GetInt("post.max_words")
	getPostInterval := viper.GetInt("history.get_interval")
	learnFromCW := viper.GetBool("history.learn_from_cw")
	historyFilePath := viper.GetString("history.file_path")
	maxStoredStatuses := viper.GetInt("history.max_length")

	history := LoadFromGob(historyFilePath)

	for k, v := range history.LastStatus {
		fmt.Println("Last status for user " + k + ": " + v)
	}

	rand.Seed(time.Now().UnixNano())

	go GetStatusesLoop(history, historyFilePath, instanceURL, accessToken, getPostInterval, learnFromCW, maxStoredStatuses)

	go PostQuotesLoop(history, instanceURL, accessToken, makePostInterval, postVisibility, maxWords)

	go HandleReplies(history, instanceURL, accessToken, postVisibility, maxWords)

	select {}
}

func HandleReplies(history *History, instanceURL, accessToken, postVisibility string, maxWords int) {
	for {
		notificationChannel := make(chan fedi.Notification)
		go fedi.NotificationStream(notificationChannel, instanceURL, accessToken)

		for notification := range notificationChannel {
			if notification.Type == "lost connection" {
				fmt.Println("Websocket connection closed. Reopening")
				break
			}

			if notification.Type == "mention" {
				botUser, err := fedi.GetCurrentUser(instanceURL, accessToken)
				if err != nil {
					fmt.Println(err)
					continue
				}

				// dont reply to yourself
				if notification.Account.ID == botUser.ID {
					continue
				}

				followedUsers, err := fedi.GetUserFollowing(botUser, instanceURL, accessToken)
				if err != nil {
					fmt.Println(err)
					continue
				}

				// mention the person we are replying to
				replyMention := "@" + notification.Status.Account.Acct
				quote := replyMention + " " + GenQuote(history, followedUsers, maxWords)

				err = fedi.PostStatus(quote, postVisibility, notification.Status.ID, "false", instanceURL, accessToken)
				if err != nil {
					fmt.Println(err)
					continue
				}
			}
		}
	}
}

func GetStatusesLoop(history *History, historyFilePath string, instanceURL, accessToken string, interval int, learnFromCW bool, maxStoredStatuses int) {
	for {
		GetNewStatuses(history, historyFilePath, instanceURL, accessToken, learnFromCW, maxStoredStatuses)
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}

func PostQuotesLoop(history *History, instanceURL, accessToken string, interval int, postVisibility string, maxWords int) {
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
		quote := GenQuote(history, followedUsers, maxWords)
		err = fedi.PostStatus(quote, postVisibility, "", "false", instanceURL, accessToken)
		if err != nil {
			fmt.Println(err)
			continue
		}
		time.Sleep(time.Duration(interval) * time.Minute)
	}
}
