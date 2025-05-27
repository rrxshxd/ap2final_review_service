package main

import (
	"ap2final_review_service/internal/app"
	"ap2final_review_service/internal/config"
	"context"
	"github.com/sorawaslocked/ap2final_base/pkg/logger"
)

func main() {
	ctx := context.Background()

	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info("initializing application")
	application, err := app.New(ctx, cfg, log)
	if err != nil {
		log.Error("failed to initialize application")
		return
	}

	log.Info("running application")
	application.Run()
}
