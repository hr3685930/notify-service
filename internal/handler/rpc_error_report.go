package handler

import (
	"fmt"
	"github.com/ddliu/go-httpclient"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"notify-service/configs"
)

var dontReportGrpcCode = []codes.Code{
	codes.NotFound,
}

func GRPCErrorReport(md metadata.MD, req interface{}, stack string, status *status.Status) {
	isDontReport := false
	for _, value := range dontReportGrpcCode {
		if value == status.Code() {
			isDontReport = true
		}
	}
	errUrl := configs.ENV.App.ErrReport
	if errUrl != "" && !isDontReport {
		request := map[string]interface{}{
			"header": md,
			"params": req,
		}

		app := map[string]string{
			"name":        configs.ENV.App.Name,
			"environment": configs.ENV.App.Env,
		}

		exception := map[string]interface{}{
			"code":  status.Code().String(),
			"trace": stack,
		}

		option := map[string]interface{}{
			"error_type": "error",
			"app":        app,
			"exception":  exception,
			"request":    request,
		}

		fmt.Print(option)

		go func() {
			_, _ = httpclient.Begin().PostJson(errUrl, option)
		}()
	}
}
