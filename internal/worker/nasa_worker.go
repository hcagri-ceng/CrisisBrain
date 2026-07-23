package worker

import (
	"CrisisBrain/internal/domain/anomaly"
	"context"
	"io"
	"log"
	"net/http"
	"time"
)

type NasaWorker struct {
	NasaKey string
	as      *anomaly.AnomalyService
	fd      *anomaly.FireDetector
}

func NewNasaWorker(as *anomaly.AnomalyService, fd *anomaly.FireDetector, nasakey string) *NasaWorker {
	return &NasaWorker{
		NasaKey: nasakey,
		as:      as,
		fd:      fd,
	}

}

func (nw *NasaWorker) Start() {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			func() {
				log.Println("NASA FIRMS API VERİ ALINIYOR !")
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel() // ARTIK GÜVENLİ! Çünkü bu fonksiyon bittiğinde (yani bu döngü adımı bittiğinde) defer çalışır.

				req, err := http.NewRequestWithContext(ctx, "GET", nw.NasaKey, nil)
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
				nasaData, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Println("Yanıt Okunamadı")
					return
				}
				// ARTIK GÜVENLİ! Fonksiyon bittiğinde body kesinlikle kapanacak.
				log.Println("Veri başarıyla çekildi, uzunluk:", len(nasaData))

				// 2. Veriyi orkestra şefine (Service) yolla
				err = nw.as.ProcessSensorData(nw.fd, "NASA API", nasaData)
				if err != nil {
					log.Printf("[NASA FIRMS] Veri işlenirken hata oluştu: %v\n", err)
				} else {
					log.Println("[NASA FIRMS] Veri başarıyla işlendi ve anomali analizi tamamlandı.")
				}

			}()

		}

	}
}
