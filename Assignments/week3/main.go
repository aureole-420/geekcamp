package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		time.Sleep(10 * time.Second) // 模拟长耗时请求
		c.String(http.StatusOK, "Hello")
	})

	srv := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	eg, egCtx := errgroup.WithContext(context.Background())

	// 启动server
	eg.Go(func() error {
		select {
		case <-egCtx.Done():
			return fmt.Errorf("Server initialization cancelled")
		default:
			return srv.ListenAndServe()
		}
	})

	// 监听sys interrupt信号
	sysCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	eg.Go(func() error {
		select {
		case <-sysCtx.Done():
			stop()
			return fmt.Errorf("Interrupt signal received")
		case <-egCtx.Done():
			return egCtx.Err()
		}
	})

	// 优雅关闭server
	eg.Go(func() error {
		<-egCtx.Done() //其实这里可以直接用sysCtx.done();所以errGroup并不是必须的.
		// The context is used to inform the server it has 5 seconds to finish
		// the request it is currently handling
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		return srv.Shutdown(ctx)
	})

	err := eg.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
