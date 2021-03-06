package main

import (
	"reflect"
	"strings"

	"github.com/alyyousuf7/lunch-lambda/notification"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pkg/errors"
)

// Do fetches lunch and pushes notification to all services
func Do() error {
	cogs, err := EnvCOGS()
	if err != nil {
		return err
	}

	notifiers, err := EnvNotifiers()
	if err != nil {
		return err
	}

	lunch, err := cogs.Lunch()
	if err != nil {
		return errors.Wrap(err, "COGS")
	}

	for i, notifier := range notifiers {
		if err := notifier.Notify("Lunch", strings.Join(lunch, "\n")); err != nil {
			return errors.Wrapf(err, "Client #%d", i+1)
		}
	}

	return nil
}

// Handler handles AWS Lambda requests
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if err := Do(); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func onlyConsole() bool {
	notifiers, err := EnvNotifiers()
	if err != nil {
		panic(err)
	}

	for _, n := range notifiers {
		if reflect.TypeOf(n) != reflect.TypeOf(&notification.Console{}) {
			return false
		}
	}

	return true
}

func main() {
	if onlyConsole() {
		err := Do()
		if err != nil {
			panic(err)
		}
	} else {
		lambda.Start(Handler)
	}
}
