package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/YattaDeSune/Game-of-life/http/server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Width  int
	Height int
}

type App struct {
	Cfg Config
}

func New(cfg Config) *App {
	return &App{
		Cfg: cfg,
	}
}

func (a *App) Run(ctx context.Context) int {
	logger := setupLogger()

	shutDownFunc, err := server.Run(ctx, logger, a.Cfg.Height, a.Cfg.Width)
	if err != nil {
		logger.Error(err.Error())

		return 1 // Возвращаем код для регистрации ошибки системой
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-c
	cancel()
	shutDownFunc(ctx)

	return 0

}

func setupLogger() *zap.Logger {
	config := zap.NewProductionConfig()

	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	logger, err := config.Build()
	if err != nil {
		fmt.Printf("Ошибка настройки логгера: %v\n", err)
	}

	return logger
}
