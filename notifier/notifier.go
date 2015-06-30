package notifier

import (
	"errors"

	"github.com/nlopes/slack"
)

type Notifier interface {
	Notify(text string) error
}

type slackNotifier struct {
	name        string
	channelName string
	client      *slack.Slack
}

func New(name, token, channelName string) Notifier {
	n := &slackNotifier{
		name:        name,
		channelName: channelName,
		client:      slack.New(token),
	}

	return n
}

func (n *slackNotifier) Notify(text string) error {
	channel, err := n.findChannel(n.channelName)
	if err != nil {
		return err
	}

	message := slack.NewPostMessageParameters()
	message.Text = text
	message.Username = n.name

	n.client.PostMessage(channel.Id, message.Text, message)
	if err != nil {
		return err
	}

	return nil
}

func (n *slackNotifier) findChannel(name string) (channel slack.Channel, err error) {
	channelSeq, err := n.client.GetChannels(true)
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
