package jobs

import (
	"errors"
	"fmt"
	"notify-service/pkg/queue"
	"time"
)

type WechatTextMsg struct {
	Num int `json:"num"`
}

func (w *WechatTextMsg) Handler() (queueErr *queue.Error) {
	defer func() {
		if err := recover(); err != nil {
			queueErr = queue.Err(fmt.Errorf("error send text msg: %+v", err))
		}
	}()

	fmt.Println("s111", w.Num)
	time.Sleep(time.Second * 10)
	fmt.Println("e111", w.Num)

	return queue.Err(errors.New("hahaha"))

	fmt.Println("e111", w.Num)

	return
}
