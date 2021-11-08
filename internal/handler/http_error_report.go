package handler

import (
	"github.com/ddliu/go-httpclient"
	"github.com/gin-gonic/gin"
	"net/http"
	"notify-service/configs"
)

var dontReportHTTPCode = []int{
	http.StatusUnauthorized,
	http.StatusForbidden,
	http.StatusMethodNotAllowed,
	http.StatusUnsupportedMediaType,
	http.StatusUnprocessableEntity,
}

func HTTPErrorReport(HTTPCode int, response gin.H, stack string, c *gin.Context) {
	isDontReport := false
	for _, value := range dontReportHTTPCode {
		if value == HTTPCode {
			isDontReport = true
		}
	}
	errUrl := configs.ENV.App.ErrReport
	if errUrl != "" && !isDontReport {
		body, _ := c.Get("jsonBody")

		request := map[string]interface{}{
			"url":     c.Request.Host + c.Request.RequestURI,
			"method":  c.Request.Method,
			"headers": c.Request.Header,
			"params":  body.(string),
		}

		app := map[string]string{
			"name":        configs.ENV.App.Name,
			"environment": configs.ENV.App.Env,
		}

		exception := map[string]interface{}{
			"code":  HTTPCode,
			"trace": stack,
		}

		option := map[string]interface{}{
			"error_type": "api_error",
			"app":        app,
			"exception":  exception,
			"request":    request,
		}

		go func() {
			_, _ = httpclient.Begin().PostJson(errUrl, option)
		}()
	}
}
