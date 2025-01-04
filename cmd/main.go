package main

import (
	"context"
	"os"

	"github.com/YattaDeSune/Game-of-life/internal/app"
)

func main() {
	ctx := context.Background()
	os.Exit(mainWithExitCode(ctx))
}

func mainWithExitCode(ctx context.Context) int {
	cfg := app.Config{
		Width:  15,
		Height: 15,
	}
	app := app.New(cfg)

	return app.Run(ctx)
}
