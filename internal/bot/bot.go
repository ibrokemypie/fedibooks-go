package bot

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

func InitBot() {
	instanceURL := viper.GetString("instance.instance_url")
	accessToken := viper.GetString("instance.access_token")

	historyFilePath := "./history.gob"

	history := LoadFromGob(historyFilePath)

	GetNewStatuses(&history, historyFilePath, instanceURL, accessToken)

	fmt.Println("stored status count: " + strconv.Itoa(len(history.Statuses)))
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

	encoder := gob.NewEncoder(historyFile)
	encoder.Encode(history)
}
