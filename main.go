package main

import (
	"context"
	"fmt"
	"github.com/aaronjan/hunch"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	proto "notify-service/api/proto/pb"
	"notify-service/configs"
	"notify-service/internal/commands"
	"notify-service/internal/handler"
	"notify-service/internal/provider"
	"notify-service/internal/server"
	"notify-service/pkg/command"
	"notify-service/pkg/config"
	"notify-service/pkg/http/gin"
	zap "notify-service/pkg/log"
	"notify-service/pkg/rpc"
	"time"
)

func main() {
	ctx := context.Background()
	_, err := hunch.Waterfall(
		ctx,
		func(ctx context.Context, n interface{}) (interface{}, error) {
			// config
			return nil, config.Load(&configs.ENV)
		},
		func(ctx context.Context, n interface{}) (interface{}, error) {
			return hunch.All(
				ctx,
				func(ctx context.Context) (interface{}, error) {
					path := "./storage/log/"
					filename := "app.log"
					return nil, zap.NewLog(path, filename).Init()
				},
				func(ctx context.Context) (interface{}, error) {
					return hunch.Retry(ctx, 0, func(c context.Context) (interface{}, error) {
						err := config.Drive(configs.ENV.Cache, configs.ENV.App)
						if err != nil {
							fmt.Println("缓存重连中...", err)
							time.Sleep(time.Second * 2)
						}
						return nil, err
					})
				},
				func(ctx context.Context) (interface{}, error) {
					return hunch.Retry(ctx, 0, func(c context.Context) (interface{}, error) {
						err := config.Drive(configs.ENV.Queue, configs.ENV.App)
						if err != nil {
							fmt.Println("队列重连中...", err)
							time.Sleep(time.Second * 2)
						}
						return nil, err
					})
				},
				func(ctx context.Context) (interface{}, error) {
					return hunch.Retry(ctx, 0, func(c context.Context) (interface{}, error) {
						err := config.Drive(configs.ENV.DataBase, configs.ENV.App)
						if err != nil {
							fmt.Println("数据库重连中...", err)
							time.Sleep(time.Second * 2)
						}
						return nil, err
					})
				},
			)
		},
		func(ctx context.Context, n interface{}) (interface{}, error) {
			provider.Register()
			return nil, nil
		},
		func(ctx context.Context, n interface{}) (interface{}, error) {
			command.NewCommand(commands.Commands).Init()
			return nil, nil
		},
		func(ctx context.Context, n interface{}) (interface{}, error) {
			err := gin.NewHTTPServer(true).HTTP(":8080", handler.Route)
			return nil, err
		},
		func(ctx context.Context, n interface{}) (interface{}, error) {
			opts := []grpc.ServerOption{
				grpc.UnaryInterceptor(rpc.CustomErrInterceptor(handler.GRPCErrorReport)),
			}
			err := rpc.NewGrpc(":8081").Register(opts, func(s *grpc.Server) {
				proto.RegisterNotifyServer(s, server.NewNotify())
				reflection.Register(s)
			})
			return nil, err
		},
	)
	if err != nil {
		panic(err)
	}
}
