package commands

import (
	"github.com/urfave/cli"
	"notify-service/internal/jobs"
	"notify-service/pkg/queue"
)

func Queue(c *cli.Context) {
	ch := make(chan int)
	go queue.NewConsumer("wechat", &jobs.WechatNewsMsg{}, 3, 3, 0)
	go queue.NewConsumer("wechat", &jobs.WechatTextMsg{}, 0, 1, 2)
	<-ch
}
