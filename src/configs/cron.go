package configs

import (
	"fmt"
	"go-vieon-content-crawling-api/src/utils"
	"os"

	"github.com/robfig/cron/v3"
)

type cronClient struct {
	instance *cron.Cron
}

func NewCronClient() *cronClient {
	return &cronClient{instance: nil}
}

func (instance *cronClient) Instance() *cron.Cron {
	if instance.instance == nil {
		instance.instance = cron.New()
	}
	return instance.instance
}

func (c *cronClient) Schedule() {
	expression := os.Getenv("CRON_JOB_EXPRESSION")
	c.instance.AddFunc(expression, func() {
		if err := utils.NewNetUtil().ScheduleCrawlActive(); err != nil {
			fmt.Println("Waiting for try again after 1 minutes... ")
			utils.NewNetUtil().ScheduleCrawlActive()
			return
		} else {
			fmt.Println("Cron job has been started!")
		}
	})
	c.instance.Start()
}
