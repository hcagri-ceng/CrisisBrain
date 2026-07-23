package main

import (
	"CrisisBrain/internal/domain/anomaly"
	"CrisisBrain/internal/repository/postgres"
	"CrisisBrain/internal/worker"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Failed to parse DB config: %v", err)
	}
	// Config settings
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("Bağlantı havuzu oluşturulamadı: %v", err)
	}
	defer pool.Close()
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("Veritabanına erişilemedi: %v", err)
	}
	fmt.Println("Veritabanı bağlantı havuzu başarıyla başlatıldı!")
	//

	anomalyRepo := postgres.NewAnomalyRepository(pool)
	anomalyService := anomaly.NewAnomalyService(anomalyRepo)

	earthquakeDetector := anomaly.NewEarthquakeDetector(1.0, 1000.0)
	fireDetector := anomaly.NewFireDetector(100.0, 50.0)
	weatherDetector := anomaly.NewWeatherDetector(20.0, 15.0)
	fmt.Println(anomalyService, earthquakeDetector, fireDetector, weatherDetector)

	//
	kd := os.Getenv("KANDILLI_API_URL")
	if kd == "" {
		log.Fatalf("KANDILLI_API_URL ortam değişkeni bulunamadı!")
	}
	kandilliWorker := worker.NewKandilliWorker(anomalyService, earthquakeDetector, kd)
	go kandilliWorker.Start()
	nk := os.Getenv("NASA_FIRMS_API_KEY")
	if nk == "" {
		log.Fatalf("NASA_FIRMS_API_KEY ortam değişkeni bulunamadı!")
	}
	nasaWorker := worker.NewNasaWorker(anomalyService, fireDetector, nk)
	go nasaWorker.Start()

	select {}
}
