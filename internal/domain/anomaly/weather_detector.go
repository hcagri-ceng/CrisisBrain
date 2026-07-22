package anomaly

import "encoding/json"

type WeatherherDetector struct {
	precipitation float64
	windSpeed     float64
}

func NewWeatherDetector(precipitation, windSpeed float64) *WeatherherDetector {
	return &WeatherherDetector{
		precipitation: precipitation,
		windSpeed:     windSpeed,
	}
}

func (w *WeatherherDetector) DetectData(data json.RawMessage) (bool, error) {
	type payload struct {
		Precipitation float64 `json:"precipitation"`
		WindSpeed     float64 `json:"windSpeed"`
	}
	s := payload{}
	err := json.Unmarshal(data, &s)
	if err != nil {
		return false, err
	}
	if s.Precipitation >= w.precipitation || s.WindSpeed >= w.windSpeed {
		return true, nil
	}
	return false, nil
}
