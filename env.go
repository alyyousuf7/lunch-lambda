package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/alyyousuf7/lunch-lambda/notification"
	"github.com/pkg/errors"
)

func EnvCOGS() (COGS, error) {
	username := os.Getenv("COGS_USERNAME")
	if username == "" {
		return COGS{}, fmt.Errorf("environment variable COGS_USERNAME is missing")
	}

	password := os.Getenv("COGS_PASSWORD")
	if password == "" {
		return COGS{}, fmt.Errorf("environment variable COGS_PASSWORD is missing")
	}

	return COGS{username, password}, nil
}

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
