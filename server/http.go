package server

import (
	"context"
	"fmt"
	"github.com/openai-chatGPT-server/config"
	. "github.com/openai-chatGPT-server/util"
	"github.com/openai-chatGPT-server/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

)

type HTTPServer struct {
	router *router.Router
	config *config.Configuration
}

func NewHTTPServer(router *router.Router, config *config.Configuration) *HTTPServer {
	return &HTTPServer{router: router, config: config}
}

func (r *HTTPServer) Run() error {

	Logger.Warnf("server name:%v", r.config.Server.Name)
	Logger.Warnf("port:%v", r.config.Server.Port)
	srv := &http.Server{
		Addr:           fmt.Sprintf("0.0.0.0:%d", r.config.Server.Port),
		Handler:        r.router.LoadRouters(),
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   40 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			Logger.Errorf("server err %v", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	Logger.Warn("server stop")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		Logger.Errorf("server shutdown err %v", err)
		return err
	}

	Logger.Warn("server quit")

	return nil

}
