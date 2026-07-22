package anomaly

import (
	"encoding/json"
	"time"
)

type AnomalyRepository interface {
	Save(anomaly *Anomaly) error
}

type AnomalyService struct {
	repo AnomalyRepository
}

func NewAnomalyService(repo AnomalyRepository) *AnomalyService {
	return &AnomalyService{
		repo: repo,
	}
}

func (s *AnomalyService) ProcessSensorData(detector AnomalyDetector, sourceName string, data json.RawMessage) error {
	isAnomaly, err := detector.DetectData(data)
	if err != nil {
		return err
	}
	if isAnomaly {
		newAnomaly := &Anomaly{
			SourceName:    sourceName,
			SeverityLevel: 0, // Set an appropriate default or calculate based on the anomaly
			RawSensorData: data,
			DetectedAt:    time.Now(),
		}
		return s.repo.Save(newAnomaly)
	}
	return nil
}
