package worker

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"CrisisBrain/internal/domain/anomaly"
)

type KandilliWorker struct {
	kandilliURL string
	as          *anomaly.AnomalyService
	ed          *anomaly.EarthquakeDetector
}

func NewKandilliWorker(as *anomaly.AnomalyService, ed *anomaly.EarthquakeDetector, kd string) *KandilliWorker {
	return &KandilliWorker{
		as:          as,
		ed:          ed,
		kandilliURL: kd,
	}
}

func (kw *KandilliWorker) Start() {
	log.Println("Kandilli Worker başlatıldı. Her 10 saniyede bir veri kontrol edilecek...")

	client := &http.Client{
		Timeout: 10 * time.Second, // Zaman aşımı süresi
	}
	// Güvenli aralık (10 saniye)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop() // Goroutine kapandığında memory leak olmaması için ticker'ı durdururuz.

	for {
		select {

		case <-ticker.C:
			// Döngünün bu adımı için bir anonim fonksiyon başlatıyoruz ve anında () ile çağırıyoruz.
			func() {
				log.Println("[KANDİLLİ] Yeni sensör verisi alınıyor (Kandilli)...")
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel() // ARTIK GÜVENLİ! Çünkü bu fonksiyon bittiğinde (yani bu döngü adımı bittiğinde) defer çalışır.

				req, err := http.NewRequestWithContext(ctx, "GET", kw.kandilliURL, nil)
				if err != nil {
					log.Println("İstek oluşturulamadı:", err)
					return // DİKKAT! Buradaki return işçiyi öldürmez, sadece bu anonim fonksiyondan çıkar. Döngü devam eder!
				}

				resp, err := client.Do(req)
				if err != nil {
					log.Println("İstek hatası:", err)
					return // Anında çık, defer sayesinde cancel() otomatik çalışır.
				}
				defer resp.Body.Close()
				kandilliData, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Println("Yanıt Okunamadı")
					return
				}
				// ARTIK GÜVENLİ! Fonksiyon bittiğinde body kesinlikle kapanacak.
				log.Println("Veri başarıyla çekildi, uzunluk:", len(kandilliData))

				// 2. Veriyi orkestra şefine (Service) yolla
				err = kw.as.ProcessSensorData(kw.ed, "Kandilli API", kandilliData)
				if err != nil {
					log.Printf("[KANDİLLİ] Veri işlenirken hata oluştu: %v\n", err)
				} else {
					log.Println("[KANDİLLİ] Veri başarıyla işlendi ve anomali analizi tamamlandı.")
				}

			}()

		}

	}
}
