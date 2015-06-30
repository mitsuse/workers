package main

import (
	"os"

	"github.com/carlescere/scheduler"
	"github.com/mitsuse/workers"
	"github.com/mitsuse/workers/github"
	"github.com/mitsuse/workers/notifiers"
)

func main() {
	run(
		github.NewStarCollector(
			"GitHub Star Collector",
			os.Getenv("GITHUB_TOKEN"),
			notifiers.New(
				"Star Collector",
				os.Getenv("SLACK_TOKEN"),
				os.Getenv("SLACK_CHANNEL_GITHUB_STARRED"),
			),
		),
		scheduler.Every(10).Minutes(),
	)

	wait()
}

func run(worker workers.Worker, job *scheduler.Job) {
	if _, err := job.Run(worker.Work); err != nil {
		workers.Log(worker, err)
		return
	}
}

func wait() {
	c := make(chan struct{})
	<-c
}
