package anomaly

import (
	"encoding/json"
)

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
		Title string  `json:"title"`
	}
	type apiResponse struct {
		Result []payload `json:"result"`
	}
	var response apiResponse
	err := json.Unmarshal(data, &response)
	if err != nil {
		return false, err
	}

	if len(response.Result) == 0 {
		return false, nil // Ortada deprem yok, tehlike yok.
	}

	if response.Result[0].Mag >= e.magnitude && response.Result[0].Depth <= e.depth {
		return true, nil
	}
	return false, nil

}
