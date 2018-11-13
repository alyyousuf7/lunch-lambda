package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"

	"github.com/alyyousuf7/lunch-lambda/notification"
	"github.com/pkg/errors"
)

// EnvCOGS returns a COGS instance with username and password set using
// environment variables
// Password can also be base64 encoded value in environment variable
func EnvCOGS() (COGS, error) {
	username := os.Getenv("COGS_USERNAME")
	if username == "" {
		return COGS{}, fmt.Errorf("environment variable COGS_USERNAME is missing")
	}

	passwordb64 := os.Getenv("COGS_PASSWORD")
	if passwordb64 == "" {
		return COGS{}, fmt.Errorf("environment variable COGS_PASSWORD is missing")
	}

	password, err := base64.StdEncoding.DecodeString(passwordb64)
	if err != nil {
		// probably the input is not base64 at all, so let's use it directly
		return COGS{username, passwordb64}, nil
	}

	// remove the last character if it is a linefeed
	if password[len(password)-1] == '\n' {
		password = password[:len(password)-1]
	}

	return COGS{username, string(password)}, nil
}

// EnvNotifiers returns relevant Notifier array using environment variable
func EnvNotifiers() ([]notification.Notifier, error) {
	clientNum, err := strconv.Atoi(os.Getenv("CLIENT_NUM"))
	if err != nil || clientNum < 1 {
		return nil, fmt.Errorf("invalid environment variable CLIENT_NUM value")
	}

	notifiers := []notification.Notifier{}
	for i := 1; i <= clientNum; i++ {
		notifier, err := EnvNotifier(i)
		if err != nil {
			return nil, errors.Wrapf(err, "Client %d", i)
		}

		notifiers = append(notifiers, notifier)
	}

	return notifiers, nil
}

// EnvNotifier returns Notifier instance using environment variable
func EnvNotifier(index int) (notification.Notifier, error) {
	envPrefix := fmt.Sprintf("CLIENT_%d", index)

	switch os.Getenv(envPrefix) {
	case "console":
		return &notification.Console{}, nil

	case "pushbullet":
		accessToken := os.Getenv(fmt.Sprintf("%s_TOKEN", envPrefix))

		if accessToken == "" {
			return nil, fmt.Errorf("token not provided")
		}

		return &notification.PushBullet{
			AccessToken: accessToken,
		}, nil

	default:
		return nil, fmt.Errorf("invalid notification service")
	}
}
