package main

import (
	"os"

	"github.com/carlescere/scheduler"
	"github.com/mitsuse/workers/workers/test"
)

func main() {
	testWorker := test.New(
		"test worker",
		os.Getenv("SLACK_TOKEN"),
		os.Getenv("SLACK_CHANNEL_TEST"),
	)
	scheduler.Every(10).Seconds().Run(testWorker.Work)

	sleep()
}

func sleep() {
	c := make(chan struct{})
	<-c
}
