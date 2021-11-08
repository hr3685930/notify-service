package server

import (
	"context"
	"notify-service/api/proto/pb"
	"notify-service/internal/errs"
)

type notify struct {
}

func NewNotify() *notify {
	return &notify{}
}

func (n *notify) SendEmail(ctx context.Context, in *proto.EmailRequest) (*proto.EmailResponse, error) {
	//_, _ = user.Repo.GetAll(ctx)
	return nil, errs.RpcValidationFailed("hahha")
}
func (n *notify) SendWechatOfficialMsg(context.Context, *proto.OfficialMsgRequest) (*proto.TestReq, error) {
	return nil, nil
}
