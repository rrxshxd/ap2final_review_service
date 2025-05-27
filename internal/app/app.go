package app

import (
	grpcserver "ap2final_review_service/internal/adapter/grpc"
	mongorepo "ap2final_review_service/internal/adapter/mongo"
	"ap2final_review_service/internal/config"
	"ap2final_review_service/internal/usecase"
	"context"
	"github.com/sorawaslocked/ap2final_base/pkg/logger"
	mongocfg "github.com/sorawaslocked/ap2final_base/pkg/mongo"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const serviceName = "review service"

type App struct {
	grpcServer *grpcserver.Server
	log        *slog.Logger
}

func New(
	ctx context.Context,
	cfg *config.Config,
	log *slog.Logger,
) (*App, error) {
	const op = "App.New"

	newLog := log.With(slog.String("op", op))
	newLog.Info("starting service", slog.String("service", serviceName))

	newLog.Info("connecting to mongo database", slog.String("uri", cfg.Mongo.URI))

	db, err := mongocfg.NewDB(ctx, cfg.Mongo)
	if err != nil {
		newLog.Error("error connecting to mongo database", logger.Err(err))
		return nil, err
	}

	reviewRepo := mongorepo.NewReview(db.Connection)

	reviewUseCase := usecase.NewReviewUseCase(reviewRepo, log)

	grpcServer := grpcserver.New(cfg.Server.GRPC, log, reviewUseCase)

	return &App{
		grpcServer: grpcServer,
		log:        log,
	}, nil
}

func (a *App) stop() {
	a.grpcServer.Stop()
}

func (a *App) Run() {
	a.grpcServer.MustRun()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	s := <-shutdownCh

	a.log.Info("received system shutdown signal", slog.Any("signal", s.String()))
	a.log.Info("stopping the application")
	a.stop()
	a.log.Info("graceful shutdown complete")
}
