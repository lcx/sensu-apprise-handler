# sensu-apprise-handler
## Overview
The Senso Go Apprise Handler is a [Sensu Event Handler](https://docs.sensu.io/sensu-go/latest/reference/handlers/#how-do-sensu-handlers-work) for sending incident notifications to [Apprise](https://github.com/caronc/apprise).

## Usage examples
### Help Output

```
The Sensu Go Apprise handler for notifying via Apprise

Usage:
  sensu-apprise-handler [flags]
  sensu-apprise-handler [command]

Available Commands:
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -b, --base-url string      A URL to the sensu host, used for links in messages (default "http://sensu-go.example.com:3000/c/~/n/")
  -h, --help                 help for sensu-apprise-handler
  -k, --key string           The Apprise Key to post to
  -t, --tags string          The tags used for the notification (default "all")
  -w, --webhook-url string   The webhook url to send messages to, defaults to value of SENSU_APPRISE_WEBHOOK_URL env variable
```

