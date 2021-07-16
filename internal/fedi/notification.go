package fedi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type NotificationEvent struct {
	Event   string      `json:"event"`
	Payload interface{} `json:"payload"`
}

type TypedEvent struct {
	Event   string `json:"event"`
	Payload string `json:"payload"`
}

type Notification struct {
	ID      string  `json:"id"`
	Type    string  `json:"type"`
	Account Account `json:"account"`
	Status  Status  `json:"status"`
}

// ClearNotifications - Clear all notifications
func ClearNotifications(instanceURL, accessToken string) error {
	u, err := url.Parse(instanceURL + "/api/v1/notifications/clear")
	if err != nil {
		return err
	}

	_, err = resty.New().R().
		SetAuthToken(accessToken).
		Get(u.String())
	if err != nil {
		return err
	}

	return nil
}

func NotificationStream(notificationChannel chan Notification, instanceURL, accessToken string) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	u, err := url.Parse(instanceURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	if u.Scheme == "https" {
		u.Scheme = "wss"
	} else if u.Scheme == "http" {
		u.Scheme = "ws"
	}

	u.Path = "/api/v1/streaming"

	q := u.Query()
	q.Set("stream", "user:notification")
	q.Set("access_token", accessToken)
	u.RawQuery = q.Encode()

	c, _, err := websocket.Dial(ctx, u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "the sky is falling")

	for {
		var typedEvent TypedEvent
		err = wsjson.Read(ctx, c, &typedEvent)
		if err != nil {
			fmt.Println("event")
			fmt.Println(err)
			continue
		}

		if typedEvent.Event == "notification" {
			var notification Notification
			err = json.Unmarshal([]byte(typedEvent.Payload), &notification)
			if err != nil {
				fmt.Println(err)
				continue
			}

			notificationChannel <- notification
		}
	}
}
