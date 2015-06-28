package test

import (
	"errors"
	"fmt"
	"os"

	"github.com/mitsuse/workers/workers"
	"github.com/nlopes/slack"
)

type workerImpl struct {
	name        string
	channelName string
	client      *slack.Slack
}

func New(name, token, channelName string) workers.Worker {
	w := &workerImpl{
		name:        name,
		channelName: channelName,
		client:      slack.New(token),
	}

	return w
}

func (w *workerImpl) Work() {
	channel, err := w.findChannel(w.channelName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s", w.name, err)
		return
	}

	message := slack.NewPostMessageParameters()
	message.Text = "test"
	message.Username = w.name

	w.client.PostMessage(channel.Id, message.Text, message)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s", w.name, err)
		return
	}
}

func (w *workerImpl) findChannel(name string) (channel slack.Channel, err error) {
	channelSeq, err := w.client.GetChannels(true)
	if err != nil {
		return channel, err
	}

	for _, c := range channelSeq {
		if c.Name == "test" {
			channel = c
			return channel, nil
		}
	}

	return channel, errors.New("no such channel")
}
