package worker

import (
	"github.com/go-co-op/gocron"
	"time"
)

var scheduler *gocron.Scheduler

func init() {
	scheduler = gocron.NewScheduler(time.Now().Location())
}

func AddTask(durationAsSec int, jobFun interface{}, params ...interface{}) error {

	_, err := scheduler.Every(durationAsSec).Seconds().Do(jobFun, params...)
	if err != nil {
		return err
	}

	return nil
}

func Run() {
	scheduler.StartAsync()
}
