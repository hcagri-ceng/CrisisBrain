package anomaly

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"strconv"
)

type FireDetector struct {
	minConfidence float64
	minFrp        float64
}

func NewFireDetector(confidence, frp float64) *FireDetector {
	return &FireDetector{

		minConfidence: confidence,
		minFrp:        frp,
	}
}

func (f *FireDetector) DetectData(data json.RawMessage) (bool, error) {
	// 1. Veri byte[] olarak gelir, bunu bir string okuyucuya çevirmeliyiz.
	// Go'da CSV okumak için "encoding/csv" ve "bytes" paketlerini import etmelisin.
	reader := csv.NewReader(bytes.NewReader(data))

	// Bütün CSV'yi satır satır okur (records bir [][]string olur)
	records, err := reader.ReadAll()
	if err != nil {
		return false, err
	}

	// 2. Satır satır gez
	for i, row := range records {
		// İlk satır başlıktır (latitude, longitude...), onu atlıyoruz.
		if i == 0 {
			continue
		}

		// CSV'de her sütun bir string'dir. İhtiyacımız olanları parse ediyoruz.
		// row[0] -> Latitude, row[1] -> Longitude, row[9] -> Confidence, row[12] -> FRP

		lat, errLat := strconv.ParseFloat(row[0], 64)
		lon, errLon := strconv.ParseFloat(row[1], 64)
		frp, errFrp := strconv.ParseFloat(row[12], 64)
		confidence := row[9] // Confidence artık harf dönüyor (n, l, h)

		// Veriler bozuksa (sayıya çevrilemiyorsa) o satırı atla, sistemi çökertme (Defansif Programlama)
		if errLat != nil || errLon != nil || errFrp != nil {
			continue
		}

		// 3. Kural Motoru (Business Logic)
		// Türkiye BBox: Enlem(36-42), Boylam(26-45)
		if lat >= 36.0 && lat <= 42.0 && lon >= 26.0 && lon <= 45.0 {
			// Güvenilirlik kontrolü: Sadece Nominal ('n') veya High ('h') olanları kabul et. Low ('l') olanları ele.
			// frp kontrolü: Yapılandırmadan gelen (f.minFrp) eşiğinden büyük mü?
			if (confidence == "h" || confidence == "n") && frp > f.minFrp {
				return true, nil // Tehlike tespit edildi! (İlk bulduğunda sistemi alarma geçirir)
			}
		}
	}

	// Bütün liste tarandı ve tehlikeli bir yangın bulunamadı.
	return false, nil
}
