# SubFlow - Özellikler ve Yetenekler

**Son Güncelleme:** 20 Ocak 2026  
**Versiyon:** 1.0.0  
**Mimar:** Muhammet Ali Büyük | [alibuyuk.net](https://alibuyuk.net)

---

## ✅ Mevcut Özellikler (v1.0.0)

### 🧮 AIA G702/G703 Hesaplama Motoru

| Özellik | Durum | Açıklama |
|---------|-------|----------|
| BigInt Aritmetik | ✅ | Tüm para birimleri kuruş cinsinden (int64) - IEEE 754 float hatası yok |
| AIA G702 Satırları | ✅ | Line 1-8 tam uyumlu hesaplama |
| G703 Özet | ✅ | Continuation Sheet desteği |
| Retainage Hesabı | ✅ | Yüzde bazlı teminat kesintisi |
| Çoklu Para Birimi | ✅ | TRY, USD, EUR desteği |

**Hesaplanan Değerler:**
- Original Contract Sum (Orijinal Sözleşme)
- Net Change by Change Orders (Değişiklik Emirleri)
- Contract Sum to Date (Güncel Sözleşme Toplamı)
- Total Completed & Stored (Tamamlanan + Malzeme)
- Retainage (Teminat Kesintisi)
- Total Earned Less Retainage (Net Kazanç)
- Less Previous Certificates (Önceki Ödemeler)
- **Current Payment Due (Ödenecek Tutar)**
- Balance to Finish (Kalan İş)

---

### 🏗️ Clean Architecture Backend

| Bileşen | Durum | Teknoloji |
|---------|-------|-----------|
| API Server | ✅ | Go Fiber v2 |
| Entity Models | ✅ | Project, Transaction, User, Tenant |
| Repository | ✅ | PostgreSQL + In-Memory |
| Services | ✅ | Calculator, Ledger, WorkerPool |
| Middleware | ✅ | Auth, CORS, Rate Limit, Security Headers |

---

### 📊 Rapor Oluşturma

| Özellik | Durum | Format |
|---------|-------|--------|
| HTML Rapor | ✅ | AIA G702 standardında |
| Print-Ready | ✅ | CSS @media print |
| Türkçe Format | ✅ | Binlik ayraç, virgül decimal |
| Progress Bar | ✅ | Görsel tamamlanma göstergesi |

---

### 💾 Veritabanı

| Özellik | Durum | Teknoloji |
|---------|-------|-----------|
| PostgreSQL Schema | ✅ | 6 tablo |
| İmmutable Ledger | ✅ | UPDATE/DELETE trigger koruması |
| Multi-Tenant | ✅ | Tenant izolasyonu |
| SQLC Queries | ✅ | Type-safe SQL |
| Goose Migrations | ✅ | Versiyon kontrolü |

---

### 🖥️ Frontend

| Özellik | Durum | Teknoloji |
|---------|-------|-----------|
| React + TypeScript | ✅ | Vite build |
| TanStack Query | ✅ | Server state management |
| Zustand Store | ✅ | Client state (theme, auth) |
| Tailwind CSS | ✅ | Utility-first styling |
| Dark Mode | ✅ | Tema desteği |
| Dashboard | ✅ | Stat cards, grafikler |
| Projects Page | ✅ | Tablo, arama, filtre |
| Calculator Page | ✅ | AIA hesaplayıcı UI |

---

### 🔧 DevOps

| Özellik | Durum | Teknoloji |
|---------|-------|-----------|
| Dockerfile | ✅ | Multi-stage (~20MB) |
| Docker Compose | ✅ | PostgreSQL + Redis + API |
| Makefile | ✅ | 20+ komut |
| GitHub Actions | ✅ | CI/CD pipeline |
| Swagger Docs | ✅ | OpenAPI spec |

---

## 🚧 Planlanan Özellikler (Roadmap)

### v1.1.0 - PDF Engine
- [ ] Maroto PDF kütüphanesi entegrasyonu
- [ ] AIA G702 PDF şablonu
- [ ] G703 Continuation Sheet PDF
- [ ] Toplu PDF üretimi (Worker Pool)

### v1.2.0 - Authentication
- [ ] JWT token authentication
- [ ] Role-Based Access Control (RBAC)
- [ ] OAuth2 / SSO desteği
- [ ] Password hashing (argon2)

### v1.3.0 - Advanced Features
- [ ] Change Order yönetimi
- [ ] Subcontractor sözleşmeleri
- [ ] E-fatura entegrasyonu
- [ ] Webhook notifications
- [ ] Audit log görüntüleme

### v1.4.0 - Reports & Analytics
- [ ] Dashboard grafikleri (Recharts)
- [ ] Excel export
- [ ] Karşılaştırmalı raporlar
- [ ] KPI metrikleri

### v2.0.0 - Enterprise
- [ ] Kubernetes deployment
- [ ] Multi-region support
- [ ] Real-time collaboration
- [ ] Mobile app (React Native)

---

## 📁 Proje Yapısı

```
SubFlow/
├── cmd/api/                 # Entry point (main.go)
├── internal/
│   ├── core/
│   │   ├── entity/          # Domain models (5 dosya)
│   │   └── service/         # Business logic (3 dosya)
│   ├── adapter/
│   │   ├── handler/         # HTTP handlers (3 dosya)
│   │   └── repository/      # Data access (5 dosya)
│   └── pkg/                 # Utilities (logger)
├── migrations/              # SQL schemas (2 dosya)
├── docs/                    # Swagger
├── tools/                   # Report generator
└── web/                     # React frontend
    └── src/
        ├── pages/           # 3 sayfa
        ├── components/      # Layout
        ├── store/           # Zustand
        └── lib/             # API client
```

---

## 🔢 İstatistikler

| Metrik | Değer |
|--------|-------|
| Toplam Go Dosyası | ~20 |
| Frontend Dosyası | ~15 |
| SQL Dosyası | 6 |
| Toplam Satır | ~5,000+ |
| Test Coverage | Unit tests ✅ |

---

© 2026 Muhammet Ali Büyük | [alibuyuk.net](https://alibuyuk.net)
