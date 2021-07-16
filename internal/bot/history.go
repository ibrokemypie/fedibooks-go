package bot

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ibrokemypie/fedibooks-go/internal/fedi"
)

type History struct {
	LastStatus map[string]string
	Statuses   map[string]HistoryStatus
}

type HistoryStatus struct {
	AuthorID string
	Text     string
}

func GetNewStatuses(history *History, historyFilePath string, instanceURL, accessToken string, learnFromCW bool, maxStoredStatuses int) {
	botUser, err := fedi.GetCurrentUser(instanceURL, accessToken)
	if err != nil {
		log.Fatal(err)
	}
	followedUsers, err := fedi.GetUserFollowing(botUser, instanceURL, accessToken)
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range followedUsers {
		sinceID := history.LastStatus[user.ID]
		if len(sinceID) <= 0 {
			sinceID = "0"
		}

		fmt.Println("Getting new statuses for user: " + user.ID)

		for retrievedStatuses, err := fedi.GetUserStatuses(user, sinceID, instanceURL, accessToken); len(retrievedStatuses) > 0; retrievedStatuses, err = fedi.GetUserStatuses(user, sinceID, instanceURL, accessToken) {
			if err != nil {
				fmt.Println(err)
				break
			}

			for _, status := range retrievedStatuses {
				cleanedContent, err := cleanStatus(status.Content)
				if err != nil {
					fmt.Println(err)
				}
				if len(cleanedContent) > 0 {
					if learnFromCW || (!learnFromCW && len(status.CW) == 0 && !status.Sensitive) {
						history.Statuses[status.ID] = HistoryStatus{AuthorID: status.Account.ID, Text: cleanedContent}
					}
				}
				sinceID = status.ID
				history.LastStatus[user.ID] = sinceID
			}

			SaveToGob(history, historyFilePath)
		}
	}

	if len(history.Statuses) > maxStoredStatuses {
		postsToClean := len(history.Statuses) - maxStoredStatuses
		fmt.Println("History has reached length " + strconv.Itoa(len(history.Statuses)) + ", removing " + strconv.Itoa(postsToClean) + " oldest statuses.")
		keys := make([]string, 0, len(history.Statuses))
		for k := range history.Statuses {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for i := 0; i <= postsToClean; i++ {
			delete(history.Statuses, keys[i])
		}
		SaveToGob(history, historyFilePath)
	}

	fmt.Println("Finished retrieving statuses.")
}

// Cleans the status HTML into a manageable string
// Taken from https://github.com/Lynnesbian/mstdn-ebooks/blob/master/functions.py
func cleanStatus(content string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return "", err
	}

	doc.Find("br").Each(func(i int, s *goquery.Selection) {
		s.ReplaceWithHtml("\n")
	})

	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		s.ReplaceWithHtml("\n")
	})

	doc.Find("a.hashtag").Each(func(i int, s *goquery.Selection) {
		s.ReplaceWithHtml(s.Text())
	})

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok {
			s.ReplaceWithHtml(href)
		} else {
			s.Remove()
		}
	})

	text := doc.Text()
	mastodonMentionRegex := regexp.MustCompile(`https:\/\/([^/]+)\/(@[^\s]+)`)
	pleromaMentionRegex := regexp.MustCompile(`https:\/\/([^/]+)\/users\/([^\s/]+)`)

	// replace mentions with fake text only mentions
	// zero width space after the username to preventactually mentioning someone
	text = mastodonMentionRegex.ReplaceAllString(text, "$2\u200B@$1")
	text = pleromaMentionRegex.ReplaceAllString(text, "@$2\u200B@$1")
	text = strings.TrimSpace(text)

	return text, nil
}

func LoadFromGob(historyFilePath string) *History {
	historyFile, err := os.OpenFile(historyFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer historyFile.Close()

	var history *History
	decoder := gob.NewDecoder(historyFile)

	err = decoder.Decode(&history)
	if err != nil {
		history = &History{LastStatus: make(map[string]string), Statuses: make(map[string]HistoryStatus)}
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
