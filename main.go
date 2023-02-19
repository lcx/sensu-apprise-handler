package main

import (
	"fmt"
	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/sensu-community/sensu-plugin-sdk/sensu"
  "encoding/json"
  "net/http"
  "bytes"
  "errors"
)

type HandlerConfig struct {
	sensu.PluginConfig
	AppriseWebhookUrl string
	AppriseKey        string
	AppriseTags       string
	SensuBaseUrl		  string
}

type Notification struct {
  Tag   string `json:"tag"`
  Title string `json:"title"`
  Body  string `json:"body"`
}

const (
	webHookUrl = "webhook-url"
	key        = "key"
	tags       = "tags"
	baseUrl		 = "base-url"
)

var (
	config = HandlerConfig{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-apprise-handler",
			Short:    "The Sensu Go Apprise handler for notifying via Apprise",
			Timeout:  10,
			Keyspace: "sensu.io/plugins/sensu-apprise-handler/config",
		},
	}

	AppriseConfigOptions = []*sensu.PluginConfigOption{
		{
			Path:      webHookUrl,
			Env:       "SENSU_APPRISE_WEBHOOK_URL",
			Argument:  webHookUrl,
			Shorthand: "w",
			Default:   "",
			Usage:     "The webhook url to send messages to, defaults to value of SENSU_APPRISE_WEBHOOK_URL env variable",
			Value:     &config.AppriseWebhookUrl,
		},
		{
			Path:      key,
			Env:       "SENSU_APPRISE_KEY",
			Argument:  key,
			Shorthand: "k",
			Usage:     "The Apprise Key to post to",
			Value:     &config.AppriseKey,
		},
		{
			Path:      tags,
			Env:       "SENSU_APPRISE_TAGS",
			Argument:  tags,
			Shorthand: "t",
			Default:   "all",
			Usage:     "The tags used for the notification",
			Value:     &config.AppriseTags,
		},
		{
			Path:      baseUrl,
			Env:       "SENSU_BASE_URL",
			Argument:  baseUrl,
			Shorthand: "b",
			Default:   "http://sensu-go.example.com:3000/c/~/n/",
			Usage:     "A URL to the sensu host, used for links in messages",
			Value:     &config.SensuBaseUrl,
		},
	}
)

func main() {
	goHandler := sensu.NewGoHandler(&config.PluginConfig, AppriseConfigOptions, checkArgs, sendMessage)
	goHandler.Execute()
}

func checkArgs(_ *corev2.Event) error {
	if len(config.AppriseWebhookUrl) == 0 {
		return fmt.Errorf("--webhook-url or SENSU_APPRISE_WEBHOOK_URL environment variable is required")
	}

	return nil
}

func formattedEventAction(event *corev2.Event) string {
	switch event.Check.Status {
	case 0:
		return "RESOLVED"
	default:
		return "ALERT"
	}
}

func messageStatus(event *corev2.Event) string {
	switch event.Check.Status {
	case 0:
		return "Resolved"
	case 1:
		return "Warning"
	case 2:
		return "Critical"
	default:
		return "Unknown"
	}
}

func messageTitle(event *corev2.Event) string {
  title := formattedEventAction(event)+":"+ 
           event.Entity.Name+
					"|"+event.Check.Name+"> "+
					messageStatus(event)

  return title
}

func messageText(event *corev2.Event) string {
  message :=  event.Entity.Name+
							"|"+event.Check.Namespace+"|"+
							event.Check.Name+"> "+
							event.Check.Output

  return message
}

func sendMessage(event *corev2.Event) error {
  url := config.AppriseWebhookUrl + "/notify/" + config.AppriseKey 
  body := &Notification{
      Tag:   config.AppriseTags,
      Title: messageTitle(event),
      Body:  messageText(event),
  }

  payloadBuf := new(bytes.Buffer)
  json.NewEncoder(payloadBuf).Encode(body)
  req, _ := http.NewRequest("POST", url, payloadBuf)
  req.Header.Set("Content-Type", "application/json; charset=UTF-8")
  client := &http.Client{}
  res, err := client.Do(req)
  if err != nil {
      return err
  }

  defer res.Body.Close()

  if res.Status != "200 OK" {
    return errors.New("response Status: " + res.Status)
  }


  return nil
}
