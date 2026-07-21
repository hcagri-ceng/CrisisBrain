package anomaly

import "encoding/json"

type FireDetector struct {
	fireRadiativePower float64
	brightness         float64
}

func NewFiredetector(fireRadiativePower, brightness float64) *FireDetector {
	return &FireDetector{
		fireRadiativePower: fireRadiativePower,
		brightness:         brightness,
	}
}

func (f *FireDetector) DetectData(data json.RawMessage) (bool, error) {
	type payload struct {
		FireRadiativePower float64 `json:"frp"`
		Brightness         float64 `json:"brightness"`
	}
	s := payload{}
	err := json.Unmarshal(data, &s)
	if err != nil {
		return false, err
	}
	if s.FireRadiativePower >= f.fireRadiativePower && s.Brightness >= f.brightness {
		return true, nil
	}
	return false, nil
}
