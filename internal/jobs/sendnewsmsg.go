package jobs

import (
	"errors"
	"fmt"
	"notify-service/pkg/queue"
)

type WechatNewsMsg struct {
	Num int `json:"num"`
}

func (w *WechatNewsMsg) Handler() (queueErr *queue.Error) {
	defer func() {
		if err := recover(); err != nil {
			queueErr = queue.Err(fmt.Errorf("error send news msg: %+v", err))
		}
	}()

	fmt.Println("w111", w.Num)
	return queue.Err(errors.New("hahaha"))

}
