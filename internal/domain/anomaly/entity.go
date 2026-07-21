package anomaly

import (
	"encoding/json"
	"time"
)

type Anomaly struct {
	ID            int32           `json:"id"`
	SourceName    string          `json:"source_name"`
	SeverityLevel float64         `json:"severity_level"`
	RawSensorData json.RawMessage `json:"raw_sensor_data"`
	DetectedAt    time.Time       `json:"detected_at"`
}
