package main

import (
	"context"
	"evo-campaign-engine/internal/config"
	"evo-campaign-engine/internal/domain"
	"evo-campaign-engine/internal/driver"
	"evo-campaign-engine/internal/engine"
	"evo-campaign-engine/internal/handler"
	pginfra "evo-campaign-engine/internal/infra/postgres"
	rdinfra "evo-campaign-engine/internal/infra/redis"
	"evo-campaign-engine/internal/repository"
	"evo-campaign-engine/internal/service"
	"evo-campaign-engine/internal/throttle"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	devMode := flag.Bool("dev", false, "development mode")
	flag.Parse()

	if *devMode {
		if err := godotenv.Load(".env"); err != nil {
			log.Printf("warning: .env not loaded: %v", err)
		}
	}

	cfg := config.Load()

	db := pginfra.Connect(&cfg.DB)
	rdb := rdinfra.Connect(&cfg.Redis)

	campaignRepo := repository.NewCampaignRepo(db)
	audienceRepo := repository.NewAudienceRepo(db)
	jobRepo := repository.NewSendJobRepo(db)
	senderRepo := repository.NewSenderRepo(db)
	suppressRepo := repository.NewSuppressionRepo(db)
	throttleRepo := repository.NewThrottleRepo(db)

	campaignSvc := service.NewCampaignService(
		campaignRepo, audienceRepo, jobRepo, senderRepo, suppressRepo, throttleRepo,
	)

	throttleEngine := throttle.New(rdb)
	rotator := engine.NewRotator(senderRepo)

	waDriver := driver.NewWhatsApp(cfg.EvolutionAPIURL, cfg.N8NWebhookURL)
	drivers := map[string]driver.ChannelDriver{
		"whatsapp":   waDriver,
		"api_whatsapp": waDriver,
	}

	jobCh := make(chan domain.SendJob, cfg.WorkerCount*10)
	scheduler := engine.NewScheduler(jobRepo, campaignRepo, cfg.SchedulerTick, 100, jobCh)
	workerDeps := engine.WorkerDeps{
		JobRepo:        jobRepo,
		CampaignRepo:   campaignRepo,
		SenderRepo:     senderRepo,
		SuppressRepo:   suppressRepo,
		ThrottleRepo:   throttleRepo,
		ThrottleEngine: throttleEngine,
		Rotator:        rotator,
		Drivers:        drivers,
		MaxRetries:     cfg.MaxRetries,
	}
	pool := engine.NewPool(scheduler, workerDeps, cfg.WorkerCount)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go pool.Start(ctx)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.RedirectTrailingSlash = false
	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health", "/ready"},
	}))

	healthHandler := handler.NewHealthHandler(db)
	healthHandler.RegisterRoutes(router)

	v1 := router.Group("/api/v1")
	{
		campaignHandler := handler.NewCampaignHandler(campaignSvc)
		campaignHandler.RegisterRoutes(v1)

		throttleHandler := handler.NewThrottleHandler(throttleRepo)
		throttleHandler.RegisterRoutes(v1)

		senderHandler := handler.NewSenderHandler(senderRepo)
		senderHandler.RegisterRoutes(v1)
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("evo-campaign-engine listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	if err := rdb.Close(); err != nil {
		log.Printf("redis close error: %v", err)
	}

	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}

	log.Println("evo-campaign-engine stopped")
}
