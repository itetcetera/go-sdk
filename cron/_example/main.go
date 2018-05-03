package main

import (
	"context"
	"math/rand"
	"time"

	"github.com/blend/go-sdk/cron"
	"github.com/blend/go-sdk/logger"
)

type emptyJob struct {
	running bool
}

func (j *emptyJob) Timeout() time.Duration {
	return 2 * time.Second
}

func (j *emptyJob) Name() string {
	return "printJob"
}

func (j *emptyJob) Execute(ctx context.Context) error {
	j.running = true
	if rand.Int()%2 == 1 {
		time.Sleep(2000 * time.Millisecond)
	} else {
		time.Sleep(8000 * time.Millisecond)
	}
	j.running = false
	return nil
}

func (j *emptyJob) OnCancellation() {
	j.running = false
}

func (j *emptyJob) Status() string {
	if j.running {
		return "Request in progress"
	}
	return "Request idle."
}

func (j *emptyJob) Schedule() cron.Schedule {
	return cron.Immediately().Then(cron.Every(10 * time.Second))
}

func main() {
	jm := cron.New().WithLogger(logger.All())
	jm.LoadJob(&emptyJob{})
	jm.Start()

	for {
		status := jm.Status()
		for _, job := range status.Jobs {
			if task, hasTask := status.Tasks[job.Name]; hasTask {
				jm.Logger().Infof("job: %s state: running elapsed: %v", job.Name, cron.Since(task.StartTime))
			} else {
				jm.Logger().Infof("job: %s state: stopped", job.Name)
			}
		}

		time.Sleep(1000 * time.Millisecond)
	}
}