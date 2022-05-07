package http

import (
	"blockchain/pkg/config"
	"blockchain/pkg/logger"
	"blockchain/pkg/ossignal"
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
	"time"
)

type server struct {
	server  *http.Server
	options *options
	handler func(router chi.Router)
}

type options struct {
	traceName   string
	addr        string
	timeout     time.Duration
	middlewares []func(http.Handler) http.Handler
}

type Option func(*options)

func WithTraceName(name string) Option {
	return func(o *options) {
		o.traceName = name
	}
}

func WithFlag(f *flag) Option {
	return func(o *options) {
		o.addr = f.addr
		o.timeout = f.timeout
	}
}

func WithMiddlewares(middlewares ...func(http.Handler) http.Handler) Option {
	return func(o *options) {
		o.middlewares = append(o.middlewares, middlewares...)
	}
}

func NewServer(registerRoute func(router chi.Router)) *server {
	cfg := config.GetConfig().System
	server := &server{
		options: &options{
			traceName: "default-trace-name",
			addr:      cfg.Http.Addr,
			timeout:   cfg.Http.ContextTimeoutSec,
		},
	}

	loc, err := time.LoadLocation(cfg.TimeZone)
	if err != nil {
		logger.Logging.Error(fmt.Sprintf("error: %v", err))
	}
	logger.Logging.Info(fmt.Sprintf("current timezone: %s", loc))

	server.handler = registerRoute

	return server
}

func (srv *server) Start(opts ...Option) {
	for i := range opts {
		opts[i](srv.options)
	}

	mux := chi.NewRouter()
	if len(srv.options.middlewares) > 0 {
		mux.Use(srv.options.middlewares...)
	}

	registerHealth(mux)
	srv.handler(mux)

	srv.server = &http.Server{
		Addr:    srv.options.addr,
		Handler: mux,
	}

	go func() {
		logger.Logging.Info(fmt.Sprintf("start http server addr: %s", srv.options.addr))
		if err := srv.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Logging.Error(fmt.Sprintf("http serve error: %v", err))
		}
	}()
}

func (srv *server) Stop(signal os.Signal) {
	if signal == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), srv.options.timeout)
	defer cancel()
	if err := srv.server.Shutdown(ctx); err != nil {
		logger.Logging.Fatal(err.Error())
	}
	logger.Logging.Info(fmt.Sprintf("stopping http server... ExitCode: %d, Signal: %s", ossignal.GetExitCode(signal), signal.String()))
}
