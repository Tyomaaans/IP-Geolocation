package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ip-geo/internal/config"
	"ip-geo/internal/infrastructure/client/ipgeo"
	"ip-geo/internal/infrastructure/rabbitmq"
	"ip-geo/internal/infrastructure/redis"
	"ip-geo/internal/infrastructure/sqlite"
	"ip-geo/internal/ip"
	"ip-geo/internal/router"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	
	cfg := config.NewConfig()

	db, err := sqlite.NewSQLiteDB(cfg.DSN)
	if err != nil {
		log.Fatalf("failed to connect sqlite: %v", err)
	}

	ipGeoClient := ipgeo.NewIpGeoClient(cfg.IpGeo)
	redisClient := redis.NewRedisClient(cfg.REDISaddr, cfg.REDISpassword)
	rmq, _      := rabbitmq.NewRabbitMQClient(cfg.RabbitMQ, logger)

	ipRepo    := ip.NewIpRepository(db)
	ipSvc     := ip.NewIpService(ipRepo, ipGeoClient, redisClient, rmq.Channel())
	ipHandler := ip.NewIpHandler(ipSvc)

	r := routes.NewUserRouter(ipHandler)

	srv := &http.Server{
		Addr:         ":" + cfg.APPport,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		log.Printf("Server running on :%s", cfg.APPport)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	go func() {
		ctx := context.Background()
		if err := rmq.StartWorker(ctx, ipSvc); err != nil {
			logger.Error("rabbitmq: worker stopped with error", slog.Any("error", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	sqlDB, err := db.DB()
	if err == nil {
		if err := sqlDB.Close(); err != nil {
			log.Printf("DB close error: %v", err)
		}
	}

	log.Println("Server exited")
}