package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/YattaDeSune/Game-of-life/http/server/handler"
	"github.com/YattaDeSune/Game-of-life/internal/service"
	"go.uber.org/zap"
)

func new(ctx context.Context,
	logger *zap.Logger,
	lifeService service.LifeService,
) (http.Handler, error) {
	muxHandler, err := handler.New(ctx, lifeService)
	if err != nil {
		return nil, fmt.Errorf("handler initialization error: %w", err)
	}

	muxHandler = handler.Decorate(muxHandler, loggingMiddleware(logger))

	return muxHandler, nil
}

func Run(
	ctx context.Context,
	logger *zap.Logger,
	height, width int,
) (func(context.Context) error, error) {
	lifeService, err := service.New(height, width)
	if err != nil {
		return nil, err
	}

	muxHandler, err := new(ctx, logger, *lifeService)
	if err != nil {
		return nil, err
	}

	srv := &http.Server{Addr: ":8081", Handler: muxHandler}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("ListenAndServe",
				zap.String("err", err.Error()))
		}
	}()

	return srv.Shutdown, nil
}

func loggingMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			duration := time.Since(start)
			logger.Info("HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Duration("duration", duration),
			)
		})
	}
}
