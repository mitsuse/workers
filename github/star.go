package github

import (
	"os"

	"golang.org/x/oauth2"

	api "github.com/google/go-github/github"
	"github.com/mitsuse/workers"
)

type startCollector struct {
	name   string
	client *api.Client
}

func NewStarCollector(name, token string) workers.Worker {
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
	}

	return w
}

func (w *startCollector) Name() string {
	return w.name
}

func (w *startCollector) Work() {
	// TODO: Implement.
}
