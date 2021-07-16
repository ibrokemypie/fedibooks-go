package bot

import (
	"encoding/gob"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
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

func GenQuote(statuses *map[string]HistoryStatus) string {
	rand.Seed(time.Now().UnixNano())
	c := NewChain(2)

	for _, s := range *statuses {
		c.Build(strings.NewReader(s.Text))
	}
	text := c.Generate(20) // Generate text.
	return text
}

func LoadFromGob(historyFilePath string) History {
	historyFile, err := os.OpenFile(historyFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer historyFile.Close()

	var history History
	decoder := gob.NewDecoder(historyFile)

	err = decoder.Decode(&history)
	if err != nil {
		return History{LastStatus: make(map[string]string), Statuses: make(map[string]HistoryStatus)}
	}

	return history
}

func SaveToGob(history *History, historyFilePath string) {
	historyFile, err := os.OpenFile(historyFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer historyFile.Close()

	// backup the original file before modifying it
	historyFileBak, err := os.Create(historyFilePath + ".bak")
	if err != nil {
		log.Fatal(err)
	}
	defer historyFileBak.Close()

	_, err = io.Copy(historyFileBak, historyFile)
	if err != nil {
		log.Fatal(err)
	}

	err = historyFileBak.Sync()
	if err != nil {
		log.Fatal(err)
	}

	encoder := gob.NewEncoder(historyFile)
	encoder.Encode(history)
}
