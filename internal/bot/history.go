package bot

import (
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/ibrokemypie/fedibooks-go/internal/fedi"
	"github.com/spf13/viper"
)

func GetFollowingStatuses() []fedi.Status {
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

	var retrievedStatuses []fedi.Status

	for _, user := range followedUsers {
		userStatuses, err := fedi.GetUserStatuses(user, "0", instanceURL, accessToken)
		if err != nil {
			log.Fatal(err)
		}

		for _, status := range userStatuses {
			status, err := cleanStatus(status)
			if err != nil {
				log.Fatal(err)
			}
			if len(status.Text) > 0 {
				retrievedStatuses = append(retrievedStatuses, status)
			}
		}
	}

	return retrievedStatuses
}

// Cleans the status HTML into a manageable string
// Taken from https://github.com/Lynnesbian/mstdn-ebooks/blob/master/functions.py
func cleanStatus(status fedi.Status) (fedi.Status, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(status.Content))
	if err != nil {
		return fedi.Status{}, err
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

	text = mastodonMentionRegex.ReplaceAllString(text, "@$2@$1")
	text = pleromaMentionRegex.ReplaceAllString(text, "@$2@$1")
	text = strings.TrimSpace(text)

	status.Text = text

	return status, nil
}
