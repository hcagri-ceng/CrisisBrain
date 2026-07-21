package anomaly

import "encoding/json"

type EarthquakeDetector struct {
	magnitude float64
	depth     float64
}

func NewEarthquakeDetector(magnitude, depth float64) *EarthquakeDetector {
	return &EarthquakeDetector{
		magnitude: magnitude,
		depth:     depth,
	}
}

func (e *EarthquakeDetector) DetectData(data json.RawMessage) (bool, error) {
	type payload struct {
		Mag   float64 `json:"mag"`
		Depth float64 `json:"depth"`
	}
	p := payload{}
	err := json.Unmarshal(data, &p)
	if err != nil {
		return false, err
	}
	if p.Mag >= e.magnitude && p.Depth <= e.depth {
		return true, nil
	}
	return false, nil

}
