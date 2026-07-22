package main

import (
	"CrisisBrain/internal/domain/anomaly"
	"CrisisBrain/internal/repository/postgres"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../../.env")
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

	earthquakeDetector := anomaly.NewEarthquakeDetector(5.0, 10.0)
	fireDetector := anomaly.NewFireDetector(100.0, 50.0)
	weatherDetector := anomaly.NewWeatherDetector(20.0, 15.0)
	fmt.Println(anomalyService, earthquakeDetector, fireDetector, weatherDetector)

}
