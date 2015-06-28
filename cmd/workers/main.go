package main

import (
	"os"

	"github.com/carlescere/scheduler"
	"github.com/mitsuse/workers"
	"github.com/mitsuse/workers/test"
)

func main() {
	run(
		test.New(
			"test worker",
			os.Getenv("SLACK_TOKEN"),
			os.Getenv("SLACK_CHANNEL_TEST"),
		),
		scheduler.Every(10).Seconds(),
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
