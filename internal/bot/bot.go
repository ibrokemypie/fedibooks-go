package bot

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ibrokemypie/fedibooks-go/internal/fedi"
	"github.com/spf13/viper"
)

func GetFollowingStatuses() {
	instanceURL := viper.GetString("instance.instance_url")
	accessToken := viper.GetString("instance.access_token")

	botUser, err := fedi.GetCurrentUser(instanceURL, accessToken)
	if err != nil {
		log.Fatal(err)
	}

	followedUsers, err := fedi.GetUserFollowing(botUser, instanceURL, accessToken)
	if err != nil {
		log.Fatal(err)
	}

	for _, user := range followedUsers {
		userStatuses, err := fedi.GetUserStatuses(user, "0", instanceURL, accessToken)
		if err != nil {
			log.Fatal(err)
		}

		for _, status := range userStatuses {
			statusText, err := cleanStatus(status.Content)
			if err != nil {
				log.Fatal(err)
			}
			if len(statusText) > 0 {
				fmt.Println("Author: " + user.Username + ", content: " + statusText)
			}
		}
	}
}

// Cleans the status HTML into a manageable string
// Taken from https://github.com/Lynnesbian/mstdn-ebooks/blob/master/functions.py
func cleanStatus(status string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(status))
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
	mastodonMentionRegex := regexp.MustCompile(`https://([^/]+)/(@[^\\s]+)`)
	pleromaMentionRegex := regexp.MustCompile(`https://([^/]+)/users/([^\\s/]+)`)

	text = mastodonMentionRegex.ReplaceAllString(text, "$2@$1")
	text = pleromaMentionRegex.ReplaceAllString(text, "$2@$1")
	text = strings.TrimSpace(text)

	return text, nil
}
