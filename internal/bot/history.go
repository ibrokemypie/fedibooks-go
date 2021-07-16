package bot

import (
	"fmt"
	"log"
	"regexp"
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

func GetNewStatuses(history *History, historyFilePath string, instanceURL, accessToken string, learnFromCW bool) {
	botUser, err := fedi.GetCurrentUser(instanceURL, accessToken)
	if err != nil {
		log.Fatal(err)
	}
	followedUsers, err := fedi.GetUserFollowing(botUser, instanceURL, accessToken)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range history.LastStatus {
		fmt.Println("initial user: " + k + ", last status: " + v)
	}

	for _, user := range followedUsers {
		sinceID := history.LastStatus[user.ID]
		if len(sinceID) <= 0 {
			sinceID = "0"
		}

		for retrievedStatuses, err := fedi.GetUserStatuses(user, sinceID, instanceURL, accessToken); len(retrievedStatuses) > 0; retrievedStatuses, err = fedi.GetUserStatuses(user, sinceID, instanceURL, accessToken) {
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println("getting statuses from user: " + user.ID + ", from id: " + sinceID)

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
