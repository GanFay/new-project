package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ganfay/split-core/internal/config"
	"github.com/ganfay/split-core/internal/delivery/telegram"
	"github.com/ganfay/split-core/internal/pkg/logger"
	"github.com/ganfay/split-core/internal/repository/postgres"
	"github.com/ganfay/split-core/internal/repository/redisRepository"
	"github.com/ganfay/split-core/internal/usecase"

	"github.com/redis/go-redis/v9"

	tele "gopkg.in/telebot.v4"
)

func main() {
	ctx := context.Background()

	cfg := config.LoadConfig()

	logger.SetupLogger(cfg.Env)

	settings := tele.Settings{
		Token:  cfg.BotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	pool, err := postgres.NewPostgresPool(ctx, cfg.Postgres.URL())
	if err != nil {
		slog.Error("Error connecting to database", "err", err)
		os.Exit(1)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr(),
		Password: cfg.Redis.Password,
	})

	userRepository := postgres.NewUserRepository(pool)
	fundRepository := postgres.NewFundRepository(pool)
	purchaseRepository := postgres.NewPurchaseRepository(pool)
	StateRepository := redisRepository.NewRepository(rdb)

	fundUC := usecase.NewFundUsecase(fundRepository, purchaseRepository)
	userUC := usecase.NewUserUsecase(userRepository)
	stateUC := usecase.NewStateUsecase(StateRepository)

	h := telegram.NewBotHandler(fundUC, userUC, stateUC)

	b, err := tele.NewBot(settings)
	if err != nil {
		slog.Error("Error creating bot", "err", err)
		os.Exit(1)
	}
	h.SetupRegister(b)

	go func() {
		slog.Info("Starting bot", "version", cfg.BotVersion, "env", cfg.Env)
		b.Start()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sign := <-quit
	slog.Info("Stopping application...", "signal", sign.String())

	b.Stop()

	pool.Close()
	err = rdb.Close()
	if err != nil {
		panic("Failed to close the redis database: " + err.Error())
	}
	slog.Info("Application stopped gracefully.")
}
