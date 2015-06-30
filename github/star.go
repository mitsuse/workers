package github

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/oauth2"

	api "github.com/google/go-github/github"
	"github.com/mitsuse/workers"
	"github.com/mitsuse/workers/notifiers"
)

type startCollector struct {
	name      string
	client    *api.Client
	notifier  notifiers.Notifier
	last      time.Time
	firstWork bool
}

func NewStarCollector(name, token string, notifier notifiers.Notifier) workers.Worker {
	w := &startCollector{
		name: name,
		client: api.NewClient(
			oauth2.NewClient(
				oauth2.NoContext,
				oauth2.StaticTokenSource(
					&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
				),
			),
		),
		notifier:  notifier,
		firstWork: true,
	}

	return w
}

func (w *startCollector) Name() string {
	return w.name
}

func (w *startCollector) Work() {
	// TODO: Exclusive access.
	name, err := w.getLoginName()
	if err != nil {
		workers.Log(w, err)
		return
	}

	last := w.getLast()

	for r := range w.watchEvents(name, last) {
		text := fmt.Sprintf(
			"%s starred %s/%s - %s",
			*r.event.Actor.Login,
			*r.event.Repo.Owner.Login,
			*r.event.Repo.Name,
			*r.event.Repo.HTMLURL,
		)
		w.notifier.Notify(text)
	}
}

func (w *startCollector) getLast() time.Time {
	if w.firstWork {
		w.last = time.Now()
	}

	return w.last
}

func (w *startCollector) getLoginName() (string, error) {
	user, _, err := w.client.Users.Get("")
	return *user.Login, err
}

func (w *startCollector) watchEvents(
	loginName string,
	last time.Time,
) <-chan *EventResponse {
	responseChan := make(chan *EventResponse)

	go func() {
		_, lastPage, err := w.getEvents(loginName, 0)
		if err != nil {
			close(responseChan)
			return
		}

		for page := 0; page <= lastPage; page++ {
			eventSeq, _, err := w.getEvents(loginName, page)
			if err != nil {
				responseChan <- &EventResponse{event: nil, err: err}
				close(responseChan)
				return
			}

			for _, event := range eventSeq {
				if last.Unix()-(*event.CreatedAt).Unix() > 0 {
					close(responseChan)
					return
				}

				if *event.Type != "WatchEvent" {
					continue
				}
				responseChan <- &EventResponse{event: &event, err: nil}
			}
		}

		close(responseChan)
	}()

	return responseChan
}

func (w *startCollector) getEvents(
	loginName string,
	page int,
) (eventSeq []api.Event, lastPage int, err error) {
	options := &api.ListOptions{Page: page, PerPage: 10}

	eventSeq, response, err := w.client.Activity.ListEventsRecievedByUser(
		loginName,
		true,
		options,
	)

	if err != nil {
		return nil, 0, err
	}

	return eventSeq, response.LastPage, nil
}

type EventResponse struct {
	event *api.Event
	err   error
}

func (r *EventResponse) Event() *api.Event {
	return r.event
}

func (r *EventResponse) Error() error {
	return r.err
}
