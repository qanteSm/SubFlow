<div align="center">
  <h1>🏗️ SubFlow</h1>
  <h3>Enterprise Construction Financial Ledger & Compliance Engine</h3>
  
  <p>
    <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version" />
    <img src="https://img.shields.io/badge/PostgreSQL-15+-4169E1?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL" />
    <img src="https://img.shields.io/badge/License-Proprietary-red?style=for-the-badge" alt="License" />
    <img src="https://img.shields.io/badge/Architecture-Clean-green?style=for-the-badge" alt="Architecture" />
  </p>
  
  <p>
    <strong>İnşaat sektörü için hakediş, teminat ve finansal uyum yönetimi.</strong><br/>
    AIA G702/G703 standardına uyumlu, BigInt tabanlı hassas hesaplama motoru.
  </p>
</div>

---

## ✨ Özellikler

| Özellik | Açıklama |
|---------|----------|
| 🔢 **BigInt Aritmetik** | IEEE 754 hataları olmadan hassas finansal hesaplamalar |
| 📒 **Çift Girişli Defter** | Değiştirilemez (immutable) muhasebe defteri yapısı |
| 🏛️ **AIA Uyumu** | G702/G703 standardına uygun hakediş hesaplaması |
| 🚀 **Yüksek Performans** | Goroutine tabanlı paralel PDF üretimi (100+/sn) |
| 🏢 **Multi-Tenant** | SaaS mimarisi ile müşteri izolasyonu |
| 📊 **Gerçek Zamanlı** | Anlık finansal görünüm ve raporlama |

---

## 🏛️ Mimari

```
┌──────────────────────────────────────────────────────┐
│                    API LAYER (Fiber)                  │
├──────────────────────────────────────────────────────┤
│               ADAPTER LAYER (Repository/PDF)          │
├──────────────────────────────────────────────────────┤
│                  CORE LAYER (Domain)                  │
│           Calculator │ Ledger │ Entities              │
├──────────────────────────────────────────────────────┤
│              INFRASTRUCTURE (PostgreSQL/Redis)        │
└──────────────────────────────────────────────────────┘
```

> Detaylı mimari için: [ARCHITECTURE.md](ARCHITECTURE.md)

---

## 🚀 Hızlı Başlangıç

### Gereksinimler

- Go 1.21+
- Docker & Docker Compose
- Make (opsiyonel)

### Kurulum

```bash
# 1. Projeyi klonla
git clone https://github.com/mabuyuk/subflow.git
cd subflow

# 2. Docker ile başlat
docker-compose up -d

# 3. Veya yerel olarak çalıştır
make run
```

### API Test

```bash
# Health check
curl http://localhost:3000/health

# Sistem bilgisi
curl http://localhost:3000/api/v1/system/version

# Finansal özet hesaplama
curl http://localhost:3000/api/v1/projects/123/financials/summary
```

---

## 📁 Proje Yapısı

```
subflow/
├── cmd/api/           # Uygulama giriş noktası
├── internal/
│   ├── core/          # İş mantığı (Domain Layer)
│   │   ├── entity/    # Veri modelleri
│   │   └── service/   # Hesaplama servisleri
│   ├── adapter/       # Dış dünya adaptörleri
│   │   ├── handler/   # HTTP kontrolcüleri
│   │   └── repository/# Veri erişimi
│   └── pkg/           # Yardımcı paketler
├── migrations/        # Veritabanı şemaları
├── web/               # React frontend
└── docker-compose.yml # Geliştirme ortamı
```

---

## 🧮 AIA Hesaplama Motoru

```go
// Örnek hakediş hesaplaması
result := calculator.Calculate(AIABillingInput{
    OriginalContractSum:   100000000, // ₺1,000,000.00
    CurrentWorkCompleted:   15000000, // ₺150,000.00
    LaborRetainageRate:         1000, // %10
})

// Sonuç
result.CurrentPaymentDue // Ödenecek tutar (kuruş cinsinden)
result.TotalRetainage    // Tutulan teminat
result.PercentComplete   // Tamamlanma yüzdesi
```

---

## 🛠️ Makefile Komutları

```bash
make help          # Yardım mesajı
make run           # Uygulamayı başlat
make test          # Testleri çalıştır
make lint          # Kod kalite kontrolü
make docker-build  # Docker imajı oluştur
make migrate-up    # Veritabanı migrasyonu
```

---

## 📡 API Endpoints

| Method | Endpoint | Açıklama |
|--------|----------|----------|
| `GET` | `/health` | Sağlık kontrolü |
| `GET` | `/api/v1/system/version` | Sistem bilgisi |
| `GET` | `/api/v1/projects` | Proje listesi |
| `GET` | `/api/v1/projects/:id/financials/summary` | Finansal özet |
| `POST` | `/api/v1/applications/:id/generate-pdf` | PDF oluştur |

---

## 🧪 Test

```bash
# Tüm testleri çalıştır
make test

# Coverage raporu
make test-coverage

# Benchmark
make bench
```

---

## 📄 Lisans

Bu proje **Source-Available** lisansı altındadır. Kod görüntülenebilir ve öğrenilebilir, ancak ticari kullanım için yazılı izin gereklidir.

Ticari lisans için: [iletisim@alibuyuk.net](mailto:iletisim@alibuyuk.net)

---

<div align="center">
  <p>
    <strong>Muhammet Ali Büyük</strong><br/>
    <a href="https://alibuyuk.net">alibuyuk.net</a> • 
    <a href="mailto:iletisim@alibuyuk.net">iletisim@alibuyuk.net</a>
  </p>
  
  <sub>© 2026 Tüm hakları saklıdır.</sub>
</div>
