package main

import (
	"fmt"
	"go.uber.org/zap"
	"id-maker/config"
	"id-maker/internal/controller/http/router"
	v1 "id-maker/internal/controller/http/v1"
	"id-maker/internal/controller/rpc"
	"id-maker/internal/initialize"
	"id-maker/internal/usecase"
	"id-maker/internal/usecase/repo"
	"id-maker/pkg/grpcserver"
	"id-maker/pkg/httpserver"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 1、加载配置文件
	if err := initialize.Load_Config(); err != nil {
		fmt.Printf("initialize config setting failed, err:%v \n", err)
		return
	}
	// 2、初始化日志
	lg, err := initialize.Init_Log(config.Conf.LogConfig, config.Conf.AppConfig.Mode)
	if err != nil {
		fmt.Printf("initialize logger failed, err:%v \n", err)
		return
	}
	defer lg.Sync()
	lg.Debug("logger initialize success...")
	// 3、初始化MySQL连接
	sqlSession, err := initialize.Init_Mysql(config.Conf.MySqlConfig)
	if err != nil {
		fmt.Printf("initialize mysql link failed, err:%v \n", err)
		return
	}
	// 4、初始化Redis连接
	//if err := initialize.Init_Redis(config.Conf.RedisConfig); err != nil {
	//	fmt.Printf("initialize redis link failed, err:%v \n", err)
	//	return
	//}
	//defer initialize.Close_Redis()
	// 5、注册路由

	// 业务处理的注册，当业务增多的时候可以另起一个包来处理
	segmentUsecae := usecase.New(repo.New(sqlSession))
	v1.RegisterRouteSrv(segmentUsecae, lg)

	router := router.InitRouter()

	// 6、启动服务(优雅开关）
	// http 服务
	httpServer := httpserver.New(router, httpserver.Port(config.Conf.AppConfig.Port))
	// gRPC 服务
	grpcServer := grpcserver.New(grpcserver.Port(config.Conf.GrpcConfig.Port))
	rpc.NewRouter(segmentUsecae, lg)
	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞

	select {
	case <-quit:
		lg.Info("Shutdown Server ...")
	case err = <-httpServer.Notify():
		lg.Error("app - Run - httpServer.Notify", zap.Error(err))
	case err = <-grpcServer.Notify():
		lg.Error("app - Run - grpcServer.Notify", zap.Error(err))
	}
	if err := httpServer.Shutdown(); err != nil {
		lg.Error("Server Shutdown: ", zap.Error(err))
	}
	grpcServer.Shutdown()

	//	srv := &http.Server{
	//		Addr:    fmt.Sprintf(":%d", config.Conf.AppConfig.Port),
	//		Handler: router,
	//	}
	//
	//	go func() {
	//		// 开启一个goroutine启动服务
	//		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	//			lg.Fatal("listen: %s\n", zap.Error(err))
	//		}
	//	}()
	//
	//	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	//	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	//	// kill 默认会发送 syscall.SIGTERM 信号
	//	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	//	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	//	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	//	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	//	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	//	lg.Info("Shutdown Server ...")
	//	// 创建一个5秒超时的context
	//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//	defer cancel()
	//	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	//	if err := srv.Shutdown(ctx); err != nil {
	//		lg.Error("Server Shutdown: ", zap.Error(err))
	//	}
	//
	//	lg.Info("Server exiting")
}
