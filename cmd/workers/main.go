package main

import (
	"os"

	"github.com/carlescere/scheduler"
	"github.com/mitsuse/workers"
	"github.com/mitsuse/workers/github"
)

func main() {
	run(
		github.NewStarCollector(
			"GitHub Star Collector",
			os.Getenv("GITHUB_TOKEN"),
		),
		scheduler.Every(15).Minutes(),
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
