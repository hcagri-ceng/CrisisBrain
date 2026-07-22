package anomaly

import "encoding/json"

type AnomalyDetector interface {
	DetectData(data json.RawMessage) (bool, error)
}
