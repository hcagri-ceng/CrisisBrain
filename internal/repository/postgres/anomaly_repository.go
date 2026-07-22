package postgres

import (
	"CrisisBrain/internal/domain/anomaly"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AnomalyRepository struct {
	p *pgxpool.Pool
}

func NewAnomalyRepository(p *pgxpool.Pool) *AnomalyRepository {
	return &AnomalyRepository{
		p: p,
	}
}

func (r *AnomalyRepository) Save(anomaly *anomaly.Anomaly) error {
	query := `INSERT INTO anomalies (source_name, severity_level, raw_sensor_data, detected_at) VALUES ($1, $2, $3, $4)`
	_, err := r.p.Exec(context.Background(), query, anomaly.SourceName, anomaly.SeverityLevel, anomaly.RawSensorData, anomaly.DetectedAt)
	return err
}
