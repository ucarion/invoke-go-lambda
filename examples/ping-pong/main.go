package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	Ping bool `json:"ping"`
	Pong bool `json:"pong"`
}

func main() {
	lambda.Start(func(ctx context.Context, event Event) (Event, error) {
		if event.Ping {
			return Event{Pong: true}, nil
		}

		return Event{Ping: true}, nil
	})
}
