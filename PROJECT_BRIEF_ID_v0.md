# Rencana Proyek: Universal MQTT Gateway Server (UM-Gateway)

## Ringkasan Eksekutif

**Nama Proyek**: Universal MQTT Gateway Server (UM-Gateway)
**Versi**: 2.0
**Status**: Proposal/Perencanaan
**Rilis Target**: Q2 2026
**Organisasi**: Divisi Teknologi Rafflesia Agro

---

## Riwayat Versi

| Versi | Tanggal | Perubahan | Penulis |
|-------|---------|-----------|--------|
| **v0.0** | 2026-01-09 | **Proposal Asli**: Konfigurasi berbasis JSON<br>• Konsep awal dengan definisi file JSON<br>• Koneksi broker MQTT eksternal diperlukan<br>• PostgreSQL/MongoDB dapat dikonfigurasi via JSON<br>• Fokus pada kasus penggunaan farm telemetry<br>• Anggaran: $283,200<br>• Konfigurasi dan deployment manual<br>• Skema JSON untuk definisi pub/sub | Tim Pengembangan |

---

## 1. Latar Belakang

### Kondisi Saat Ini
Server Go-MQTT yang ada (v1.0) dikembangkan sebagai solusi khusus untuk operasi peternakan ayam cerdas Rafflesia Agro. Meskipun sukses di domainnya, implementasi saat ini memiliki beberapa keterbatasan:

1. **Tight Coupling (Keterikatan Kuat)**: Logika bisnis di-hard-code untuk kasus telemetry pertanian
2. **Skema Tidak Fleksibel**: Skema database kaku dan memerlukan perubahan kode untuk modifikasi
3. **Desain Satu Kegunaan**: Tidak dapat mudah diadaptasi untuk skenario IoT lainnya (smart home, pemantauan industri, dll)
4. **Konfigurasi Manual**: Pemetaan hardware, tipe sensor, dan struktur data dikompilasi ke dalam kode
5. **Keterbatasan Reusabilitas**: Setiap deployment baru memerlukan modifikasi kode yang signifikan

### Peluang Pasar
Lanskap IoT berkembang pesat di berbagai domain:
- **Pertanian**: Smart farming, pemantauan ternak, otomatisasi greenhouse
- **Bangunan Cerdas**: Kontrol HVAC, manajemen energi, sistem keamanan
- **IoT Industri**: Pemantauan peralatan, pemeliharaan prediktif, kontrol kualitas
- **Kesehatan**: Pemantauan pasien, pelacakan peralatan medis
- **Retail**: Manajemen inventaris, analitik pelanggan, kontrol lingkungan

Setiap domain membutuhkan gateway MQTT-broker-ke-database, namun saat ini harus membangun solusi kustom dari awal.

---

## 2. Model Bisnis

### Proposition Nilai

**Untuk Integrator Sistem**:
- Mengurangi waktu pengembangan 60-80% untuk proyek IoT
- Menghilangkan pekerjaan integrasi broker MQTT yang berulang
- Fokus pada fitur spesifik domain daripada infrastruktur

**Untuk Pelanggan Akhir**:
- Time-to-market lebih cepat untuk solusi IoT
- Biaya pengembangan lebih rendah
- Infrastruktur lebih andal dan teruji
- Skalabilitas dan pemeliharaan yang mudah

**Untuk Rafflesia Agro**:
- Stream pendapatan baru melalui lisensi/SaaS
- Membangun kepemimpinan teknologi di ruang IoT
- Menggunakan kembali keahlian internal di berbagai industri
- Membangun ekosistem solusi IoT yang kompatibel

### Stream Pendapatan

1. **Open Source Core**: Versi dasar gratis dengan dukungan komunitas
2. **Enterprise Edition**: Fitur lanjutan, dukungan komersial, jaminan SLA
3. **Layanan Cloud**: Hosting UM-Gateway terkelola (berbasis langganan)
4. **Pengembangan Kustom**: Konsultasi untuk integrasi khusus
5. **Pelatihan & Sertifikasi**: Program pelatihan pengembang dan administrator

### Strategi Harga

| Edisi | Pasar Target | Model Harga | Fitur |
|-------|--------------|-------------|-------|
| Komunitas | Hobiis, startup | Gratis (Lisensi MIT) | Gateway MQTT-DB dasar, config JSON, dukungan komunitas |
| Profesional | UMKM, integrator | $99/bulan/instance | Fitur lanjutan, dukungan email, dashboard |
| Enterprise | Organisasi besar | Harga kustom | Ketersediaan tinggi, dukungan khusus, integrasi kustom |
| Cloud | Semua pasar | $0.10/perangkat/bulan | Terkelola penuh, auto-scaling, SLA 99.99% |

---

## 3. Target Pengguna

### Pengguna Utama

**1. Integrator Sistem IoT**
- Mengembangkan solusi IoT untuk banyak klien
- Membutuhkan infrastruktur yang dapat diandalkan dan dapat dikustomisasi
- Keahlian teknis: Tinggi
- Titik nyeri: Membangun integrasi MQTT yang sama berulang kali

**2. Tim Pengembangan Internal**
- Perusahaan yang membangun sistem IoT internal
- Membutuhkan solusi yang dapat dipelihara dan didokumentasikan
- Keahlian teknis: Sedang hingga Tinggi
- Titik nyeri: Sumber daya terbatas untuk pengembangan infrastruktur

**3. Perusahaan Startup**
- Membangun produk berbasis IoT
- Membutuhkan prototyping cepat dan skalabilitas
- Keahlian teknis: Bervariasi
- Titik nyeri: Tekanan time-to-market

**4. MSP (Managed Service Providers)**
- Menawarkan layanan manajemen IoT
- Membutuhkan kemampuan multi-tenant
- Keahlian teknis: Tinggi
- Titik nyeri: Mengelola deployment klien yang beragam

### Pengguna Sekunder

**5. Mahasiswa & Peneliti**
- Belajar konsep IoT
- Membangun proof-of-concept
- Kebutuhan: Dokumentasi yang jelas, contoh

---

## 4. Pernyataan Masalah

### Masalah Utama

**Masalah 1: Duplikasi Kode di Seluruh Proyek**
Setiap proyek IoT membutuhkan fungsionalitas gateway MQTT-ke-database yang serupa:
- Manajemen koneksi MQTT
- Parsing dan validasi pesan
- Persistensi data
- Eksposur endpoint API
Pendekatan saat ini: Membangun dari awal setiap kali

**Masalah 2: Model Data Tidak Fleksibel**
Domain IoT yang berbeda memiliki persyaratan data yang sangat beragam:
- Pertanian: suhu, kelembaban, amonia, level pakan
- Smart Home: intensitas cahaya, deteksi gerak, konsumsi energi
- Industri: getaran, tekanan, RPM, gradien suhu
Pendekatan saat ini: Perubahan skema hard-code untuk setiap domain

**Masalah 3: Keterikatan Hardware yang Kuat**
Logika spesifik tertanam dalam kode server:
- Integrasi firmware ESP32
- Format perintah/respons spesifik
- Manajemen state kustom
Pendekatan saat ini: Modifikasi kode untuk tipe hardware baru

**Masalah 4: Kompleksitas Operasional**
Men-deploy dan mengelola beberapa instance IoT sulit:
- Tidak ada pendekatan konfigurasi yang seragam
- Migrasi database manual
- Monitoring dan logging tidak konsisten
- Troubleshooting sulit di seluruh deployment

**Masalah 5: Skalabilitas Terbatas**
Desain saat ini tidak menangani:
- Skenario multi-tenant
- Registrasi perangkat dinamis
- Hot-reloading konfigurasi
- Scaling horizontal tanpa downtime

### Dampak

- **Biaya Pengembangan**: $20.000-$50.000 per proyek untuk pengembangan gateway kustom
- **Time to Market**: 4-8 minggu untuk integrasi MQTT dasar
- **Beban Pemeliharaan**: Biaya ongoing $5.000-$15.000/tahun per deployment
- **Masalah Kualitas**: Solusi kustom kurang teruji, lebih banyak bug
- **Keterikatan Vendor**: Sulit beralih penyedia atau migrasi sistem

---

## 5. Solusi yang Diusulkan

### Visi: Universal MQTT Gateway Server

Gateway server MQTT yang **digerakkan konfigurasi, domain-agnostik** yang dapat di-deploy dalam skenario IoT apa pun melalui file konfigurasi JSON sederhana, tanpa modifikasi kode.

### Inovasi Utama

**1. Arsitektur Berbasis JSON**
- Semua logika bisnis didefinisikan dalam file konfigurasi JSON
- Tidak perlu rekompilasi untuk use case berbeda
- Hot-reload perubahan konfigurasi tanpa restart

**2. Sistem Skema Fleksibel**
- Dukungan multiple backend database (PostgreSQL, MongoDB, InfluxDB, TimescaleDB)
- Definisi skema via JSON (tidak perlu ALTER TABLE)
- Migrasi dan validasi skema otomatis

**3. Pemetaan Topik Dinamis**
- Pola topik MQTT didefinisikan dalam konfigurasi
- Transformasi dan enrichments diterapkan via rules engine
- Dukungan wildcard, regex, dan parsing topik kustom

**4. Persistensi Data Yang Dapat Dipasang (Pluggable)**
- Opsi storage time-series yang dioptimalkan
- Storage berbasis dokumen untuk skema fleksibel
- Storage relasional untuk integritas transaksional
- Penambahan adapter storage baru dengan mudah

**5. Manajemen State Hardware yang Ditingkatkan**
- Sinkronisasi state terpadu di semua konsumen
- Resolusi konflik otomatis
- Riwayat dan audit trail state
- Notifikasi push real-time

### Ikhtisar Arsitektur

```
┌─────────────────────────────────────────────────────────────┐
│                    Server UM-Gateway                        │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ MQTT         │  │ HTTP/WebSocket│  │ gRPC (opsional)│   │
│  │ Subscriber   │  │ API Server   │  │ API Server   │      │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘      │
│         │                 │                 │               │
│  ┌──────▼─────────────────▼─────────────────▼───────┐      │
│  │           Mesin Pemrosesan Pesan                │      │
│  │  - Routing Topik (dikonfigurasi JSON)          │      │
│  │  - Validasi Skema (dikonfigurasi JSON)         │      │
│  │  - Transformasi Data (Rules Engine)            │      │
│  │  - Manajemen State (Ditingkatkan)              │      │
│  └──────┬────────────────────────────────────────────┘      │
│         │                                                  │
│  ┌──────▼────────────────────────────────────────────┐     │
│  │       Layer Storage (Dapat Dipasang)              │     │
│  │  ┌──────────────┐  ┌──────────────┐              │     │
│  │  │ Time-Series  │  │ Document     │              │     │
│  │  │ (InfluxDB/   │  │ (MongoDB/    │              │     │
│  │  │  TimescaleDB)│  │  PostgreSQL) │              │     │
│  │  └──────────────┘  └──────────────┘              │     │
│  │  ┌──────────────┐  ┌──────────────┐              │     │
│  │  │ Relasional   │  │ Cache/State  │              │     │
│  │  │ (PostgreSQL) │  │ (Redis)      │              │     │
│  │  └──────────────┘  └──────────────┘              │     │
│  └──────────────────────────────────────────────────┘     │
│                                                           │
│  ┌──────────────────────────────────────────────────┐    │
│  │        Manajer Konfigurasi (JSON)                │    │
│  │  - Definisi Skema                                 │    │
│  │  - Pemetaan Topik                                 │    │
│  │  - Aturan Transformasi                           │    │
│  │  - Konfigurasi Storage                            │    │
│  └──────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

---

## 6. Tujuan Proyek

### Tujuan Utama (Harus Ada)

**O1: Kembangkan Sistem Konfigurasi Berbasis JSON**
- Definisikan semua logika bisnis dalam file JSON
- Dukung hot-reload konfigurasi
- Validasi konfigurasi sebelum deployment
- Sediakan template konfigurasi untuk skenario umum

**O2: Implementasikan Dukungan Multi-Database**
- PostgreSQL (dengan ekstensi TimescaleDB)
- MongoDB (skema fleksibel)
- InfluxDB (time-series yang dioptimalkan)
- Penambahan adapter database baru dengan mudah

**O3: Buat Mesin Pemetaan Topik Dinamis**
- Routing berbasis pola
- Parser topik kustom
- Pipeline transformasi data
- Aturan enrichment dan validasi

**O4: Tingkatkan Manajemen State Hardware**
- Sinkronisasi state terpadu
- Strategi resolusi konflik
- Riwayat dan replay state
- Sistem notifikasi multi-konsumen

**O5: Bangun Layer API Komprehensif**
- RESTful API untuk semua operasi
- WebSocket untuk update real-time
- Autentikasi dan otorisasi
- Rate limiting dan kuota

### Tujuan Sekunder (Sebaiknya Ada)

**O6: Kembangkan Dashboard Admin**
- UI konfigurasi berbasis web
- Monitoring real-time
- Agregasi dan pencarian log
- Visualisasi metrik performa

**O7: Implementasikan Multi-Tenansi**
- Isolasi tenant
- Kuota sumber daya per tenant
- Konfigurasi spesifik tenant
- Logging audit per tenant

**O8: Buat Sistem Plugin**
- Prosesor data kustom
- Adapter storage kustom
- Provider autentikasi kustom
- Marketplace ekstensi

**O9: Tambahkan Integrasi Machine Learning**
- Deteksi anomali
- Pemeliharaan prediktif
- Skoring kualitas data
- Alerting otomatis

### Tujuan Tersier (Nice to Have)

**O10: Dukungan Edge Computing**
- Deploy pada perangkat edge
- Mode operasi offline
- Sinkronisasi data saat online
- Footprint ringan

**O11: Fitur Cloud-Native**
- Deployment Kubernetes
- Auto-scaling
- Integrasi service mesh
- Distributed tracing

---

## 7. Ruang Lingkup & Batasan Proyek

### Dalam Ruang Lingkup (Apa yang Akan Kami Bangun)

#### Fungsionalitas Inti
1. **Sistem Konfigurasi**
   - Validasi skema JSON
   - Mekanisme hot-reload
   - Versioning konfigurasi
   - Pustaka template

2. **Pemrosesan MQTT**
   - Subscribe ke multiple topik
   - Dukungan QoS (0, 1, 2)
   - Persistensi pesan
   - Last Will and Testament
   - Rekoneksi otomatis

3. **Persistensi Data**
   - Multiple backend database
   - Pembuatan skema otomatis
   - Insertion batch
   - Kebijakan retensi data
   - Backup dan restore

4. **Manajemen State Hardware**
   - Penyimpanan state (Redis)
   - Sinkronisasi state
   - Resolusi konflik
   - Riwayat state
   - Notifikasi real-time

5. **Layer API**
   - Operasi CRUD
   - Operasi query
   - Autentikasi/Otorisasi
   - Dukungan WebSocket
   - Rate limiting

6. **Monitoring & Logging**
   - Health checks
   - Pengumpulan metrik (Prometheus)
   - Logging terstruktur
   - Pelacakan error

#### Deployment & Operasi
- Container Docker
- Konfigurasi Docker Compose
- Manifest Kubernetes
- Dokumentasi (panduan pengguna, referensi API, tutorial)
- Tes integrasi

### Di Luar Ruang Lingkup (Apa yang Tidak Akan Kami Bangun - Awalnya)

1. **Hardware/Firmware**
   - Firmware perangkat (ESP32, Arduino, dll)
   - Provisioning hardware
   - Update OTA

2. **Fitur Spesifik Bisnis**
   - Analitik spesifik domain
   - Dashboard spesifik industri
   - Aturan logika bisnis (di luar transformasi dasar)

3. **Fitur Lanjutan (Fase Mendatang)**
   - Pipeline machine learning (fase 2)
   - Edge computing (fase 3)
   - Deployment multi-regional (fase 3)

4. **Integrasi Eksternal**
   - API pihak ketiga (kecuali sistem plugin diimplementasikan)
   - Spesifik penyedia cloud (di luar deployment dasar)
   - Konektor sistem legacy

### Batasan & Kendala

#### Kendala Teknis
- **Bahasa**: Go 1.21+ (untuk performa dan konkurensi)
- **Versi MQTT**: Dukungan 3.1.1 dan 5.0
- **Hardware Minimum**: 2 core CPU, 4GB RAM untuk deployment dasar
- **Jaringan**: Memerlukan konektivitas ke broker MQTT dan database yang dipilih

#### Kendala Waktu
- **Fase 1 (MVP)**: 6 bulan
- **Fase 2 (Ditingkatkan)**: +4 bulan
- **Fase 3 (Lanjutan)**: +6 bulan

#### Kendala Anggaran
- **Tim Pengembangan**: 3-5 pengembang full-time
- **Infrastruktur**: $2.000/bulan untuk pengembangan/tes
- **Layanan Pihak Ketiga**: Tier gratis awalnya (misalnya MongoDB Atlas free tier)

#### Kendala Sumber Daya
- **Tim Pengembangan**: Pengembang backend (Go), engineer DevOps, QA
- **Ahli Domain**: Spesialis IoT, konsultan domain pertanian/industri
- **Dukungan**: Technical writer part-time, desainer UI/UX (untuk dashboard)

---

## 8. Stack Teknologi

### Teknologi Inti

#### Framework Backend
| Komponen | Teknologi | Alasan |
|-----------|-----------|--------|
| Bahasa | Go 1.24+ | Performa, konkurensi, static typing |
| HTTP Router | Chi v5 | Ringan, idiomatic Go |
| Klien MQTT | Eclipse Paho MQTT v2 | Mature, didukung baik, dukungan MQTT 5.0 |
| WebSocket | gorilla/websocket | Standar industri untuk Go |

#### Opsi Database
| Database | Use Case | Kelebihan |
|----------|----------|-----------|
| **PostgreSQL + TimescaleDB** | Time-series dengan integritas relasional | SQL, ACID, ekosistem mature, kompresi |
| **MongoDB** | Storage dokumen fleksibel | Tidak perlu skema, scaling horizontal, query kaya |
| **InfluxDB** | Time-series yang dioptimalkan | purpose-built, throughput tinggi, downsampling |
| **Redis** | Manajemen state, caching | In-memory, cepat, pub/sub |

#### Konfigurasi & Validasi
| Komponen | Teknologi | Alasan |
|-----------|-----------|--------|
| Validasi Skema | JSON Schema + gojsonschema | Standar, terdokumentasi baik |
| Loading Konfigurasi | Viper | Fleksibel, multiple format, hot-reload |
| Mesin Template | text/template | Built-in Go, tanpa dependensi |

#### API & Dokumentasi
| Komponen | Teknologi | Alasan |
|-----------|-----------|--------|
| Dokumentasi API | OpenAPI 3.0 (Swagger) | Standar, dukungan tooling |
| Generasi Kode | oapi-codegen | Handler type-safe dari OpenAPI |
| Autentikasi | JWT (golang-jwt/jwt) | Stateless, scalable |

#### Monitoring & Observabilitas
| Komponen | Teknologi | Alasan |
|-----------|-----------|--------|
| Metrik | Prometheus + OpenTelemetry | Cloud-native, banyak diadopsi |
| Logging | Zap + Loki | Terstruktur, agregasi log cepat |
| Tracing | OpenTelemetry + Jaeger | Standar distributed tracing |
| Health Checks | Healthcheck library | Endpoint kesehatan terstandarisasi |

#### Deployment
| Komponen | Teknologi | Alasan |
|-----------|-----------|--------|
| Kontainerisasi | Docker | Standar industri |
| Orkestrasi | Kubernetes | Skalabilitas, self-healing |
| CI/CD | GitHub Actions | Terintegrasi dengan repository |
| Service Mesh | Istio (opsional) | Manajemen traffic, keamanan |

### Pola Arsitektur

**1. Arsitektur Plugin**
```go
type StorageAdapter interface {
    Connect(config Config) error
    Store(data Data) error
    Query(query Query) ([]Result, error)
    Close() error
}

type DatabasePlugin struct {
    Name      string
    Adapter   StorageAdapter
    ConfigSchema jsonschema.Schema
}
```

**2. Pemrosesan Berbasis Konfigurasi**
```json
{
  "version": "2.0",
  "mqtt": {
    "broker": "tcp://localhost:1883",
    "topics": [
      {
        "pattern": "sensors/+/telemetry",
        "qos": 1,
        "parser": "sensorTelemetryParser",
        "storage": "timeseries",
        "transformations": [...]
      }
    ]
  },
  "storage": {
    "timeseries": {
      "type": "influxdb",
      "config": {...}
    }
  }
}
```

**3. Registry Skema**
```json
{
  "schemas": {
    "sensorTelemetry": {
      "type": "object",
      "properties": {
        "temperature": {"type": "number"},
        "humidity": {"type": "number"},
        "timestamp": {"type": "string", "format": "date-time"}
      }
    }
  }
}
```

---

## 9. Fitur-fitur

### 9.1 Fitur Utama

#### F1: Sistem Konfigurasi Berbasis JSON

**Deskripsi**: Semua aspek perilaku sistem didefinisikan dalam file konfigurasi JSON

**Kemampuan**:
- Pengaturan koneksi MQTT (broker, topik, QoS)
- Definisi skema data
- Pemetaan topik-ke-storage
- Aturan transformasi
- Aturan validasi
- Definisi endpoint API

**Contoh Konfigurasi**:
```json
{
  "version": "2.0",
  "deployment": {
    "name": "smart-farm-prod",
    "environment": "production"
  },
  "mqtt": {
    "brokers": [
      {
        "url": "tcp://mqtt.example.com:1883",
        "username": "${MQTT_USERNAME}",
        "password": "${MQTT_PASSWORD}",
        "qos": 1,
        "clean_session": true
      }
    ],
    "subscriptions": [
      {
        "id": "farm-telemetry",
        "topic": "farms/+/sensors/#",
        "qos": 1,
        "parser": "farmSensorParser",
        "storage": "timeseries-db",
        "enrichments": [
          {
            "type": "add_timestamp",
            "field": "received_at",
            "timezone": "Asia/Jakarta"
          },
          {
            "type": "geo_lookup",
            "field": "location",
            "source": "farm_registry"
          }
        ]
      }
    ]
  },
  "parsers": {
    "farmSensorParser": {
      "schema": "farmSensorSchema",
      "transformations": [
        {
          "field": "temperature",
          "type": "convert_unit",
          "from": "fahrenheit",
          "to": "celsius"
        }
      ]
    }
  },
  "storage": {
    "timeseries-db": {
      "type": "timescaledb",
      "connection": "${TIMESCALEDB_URL}",
      "schema": {
        "table": "sensor_readings",
        "time_column": "timestamp",
        "columns": [
          {"name": "device_id", "type": "string", "tags": true},
          {"name": "sensor_type", "type": "string", "tags": true},
          {"name": "value", "type": "double"},
          {"name": "unit", "type": "string"},
          {"name": "location", "type": "geojson"}
        ]
      }
    }
  },
  "api": {
    "enabled": true,
    "port": 8080,
    "authentication": {
      "type": "jwt",
      "secret": "${JWT_SECRET}"
    },
    "endpoints": [
      {
        "path": "/api/v1/sensors/{device_id}/readings",
        "method": "GET",
        "storage": "timeseries-db",
        "query": {
          "time_range": "last_24h",
          "aggregation": "avg"
        }
      }
    ]
  }
}
```

**Manfaat**:
- Tidak perlu perubahan kode untuk deployment baru
- Version control untuk konfigurasi
- Rollback mudah ke konfigurasi sebelumnya
- Validasi konfigurasi sebelum deployment

---

#### F2: Dukungan Multi-Database dengan Adapter Pluggable

**Deskripsi**: Dukungan multiple backend database melalui interface adapter yang dapat dipasang

**Database yang Didukung**:

**A. PostgreSQL + TimescaleDB (Time-Series yang Dioptimalkan)**
- Use case: Ketika integritas relasional dibutuhkan dengan data time-series
- Kelebihan: Kepatuhan ACID, query SQL, ekosistem mature
- Pembuatan skema via JSON:
```json
{
  "type": "timescaledb",
  "connection": "postgresql://user:pass@localhost:5432/db",
  "hypertable": {
    "table": "sensor_data",
    "time_column": "timestamp",
    "chunk_interval": "1 day",
    "compression": true
  },
  "retention": {
    "policy": "drop_after",
    "interval": "90 days"
  },
  "columns": [
    {"name": "timestamp", "type": "timestamptz", "not_null": true},
    {"name": "device_id", "type": "varchar(50)", "tags": true},
    {"name": "sensor_type", "type": "varchar(50)", "tags": true},
    {"name": "value", "type": "double precision"},
    {"name": "metadata", "type": "jsonb"}
  ]
}
```

**B. MongoDB (Document Store Fleksibel)**
- Use case: Ketika skema bervariasi sering atau tidak diketahui sebelumnya
- Kelebihan: Tidak perlu migrasi skema, scaling horizontal, query kaya
- Pembuatan skema via JSON:
```json
{
  "type": "mongodb",
  "connection": "mongodb://localhost:27017",
  "database": "iot_data",
  "collection": "sensor_readings",
  "indexing": [
    {"keys": {"device_id": 1, "timestamp": -1}},
    {"keys": {"sensor_type": 1}},
    {"keys": {"timestamp": -1}, "expireAfterSeconds": 7776000}
  ],
  "sharding": {
    "enabled": true,
    "key": {"device_id": 1}
  }
}
```

**C. InfluxDB (Time-Series Purpose-Built)**
- Use case: Time-series volume tinggi dengan query sederhana
- Kelebihan: Kompresi yang dioptimalkan, fungsi khusus, throughput tinggi
- Pembuatan skema via JSON:
```json
{
  "type": "influxdb",
  "connection": "http://localhost:8086",
  "organization": "farm_sensors",
  "bucket": "telemetry",
  "retention": "30d",
  "schema": {
    "measurement": "sensor_reading",
    "tags": ["device_id", "sensor_type", "location"],
    "fields": ["value", "unit", "quality"]
  }
}
```

**Interface Adapter**:
```go
type StorageAdapter interface {
    // Inisialisasi koneksi
    Connect(config StorageConfig) error

    // Buat skema/tabel/koleksi berdasarkan definisi JSON
    CreateSchema(schema SchemaDefinition) error

    // Simpan data dengan batching otomatis
    Store(data []DataPoint) error

    // Query data dengan bahasa query fleksibel
    Query(query Query) (ResultSet, error)

    // Health check
    Ping() error

    // Tutup koneksi
    Close() error
}
```

---

#### F3: Pemetaan Topik Dinamis & Mesin Transformasi

**Deskripsi**: Routing topik fleksibel dengan pipeline transformasi data

**Fitur**:
1. **Routing Berbasis Pola**
   - Dukungan wildcard: `sensors/+/telemetry/#`
   - Pencocokan regex: `buildings/([^/]+)/floors/([^/]+)/sensors`
   - Ekstraksi variabel topik

2. **Pipeline Transformasi Data**
   - Konversi tipe (string → number)
   - Konversi unit (°F → °C)
   - Enrichment data (lookup nilai dari sumber eksternal)
   - Aturan validasi
   - Filtering

3. **Validasi Skema**
   - Validasi JSON Schema
   - Aturan validasi kustom
   - Skoring kualitas data

**Contoh Pipeline Transformasi**:
```json
{
  "topic": "sensors/+/telemetry",
  "transformations": [
    {
      "name": "Extract Device ID",
      "type": "extract_topic_variable",
      "source": "$topic",
      "variable": 1,
      "target": "device_id"
    },
    {
      "name": "Parse JSON Payload",
      "type": "parse_json",
      "source": "payload"
    },
    {
      "name": "Add Received Timestamp",
      "type": "add_field",
      "field": "received_at",
      "value": "$now",
      "format": "2006-01-02T15:04:05Z07:00"
    },
    {
      "name": "Convert Temperature to Celsius",
      "type": "convert_unit",
      "field": "temperature",
      "from": "kelvin",
      "to": "celsius"
    },
    {
      "name": "Enrich with Device Metadata",
      "type": "enrich_lookup",
      "field": "device_info",
      "lookup_table": "device_registry",
      "key": "$device_id",
      "cache_ttl": "1h"
    },
    {
      "name": "Validate Data",
      "type": "validate",
      "rules": [
        {
          "field": "temperature",
          "condition": "between",
          "min": -50,
          "max": 100,
          "on_fail": "drop"
        },
        {
          "field": "humidity",
          "condition": "between",
          "min": 0,
          "max": 100,
          "on_fail": "flag"
        }
      ]
    }
  ]
}
```

---

#### F4: Manajemen State Hardware yang Ditingkatkan

**Deskripsi**: Sinkronisasi state terpadu, andal di semua konsumen

**Masalah yang Diselesaikan**:
- Race condition ketika beberapa konsumen update state
- State tidak konsisten di Redis, database, dan MQTT
- Tidak ada riwayat perubahan state
- Sulit debugging masalah state

**Fitur**:

**1. Penyimpanan State dengan Versioning**
```json
{
  "state_management": {
    "enabled": true,
    "storage": "redis",
    "versioning": {
      "enabled": true,
      "history_ttl": "30d",
      "max_versions": 1000
    },
    "conflict_resolution": {
      "strategy": "last_write_wins",
      "compare_field": "timestamp"
    }
  }
}
```

**2. Pipeline Sinkronisasi State**
```
Request Konsumen
     ↓
[Ambil State Saat Ini + Versi]
     ↓
[Optimistic Lock dengan Versi]
     ↓
[Terapkan Update]
     ↓
[Persist ke Storage]
     ↓
[Publikasikan ke MQTT]
     ↓
[Notifikasikan Subscriber WebSocket]
     ↓
[Update Cache]
     ↓
Kembalikan Sukses
```

**3. Riwayat State & Audit Trail**
```json
{
  "state_change": {
    "id": "state_123",
    "entity_id": "farm_456",
    "entity_type": "coop",
    "field": "fan_1.state",
    "old_value": false,
    "new_value": true,
    "changed_by": "user_789",
    "changed_at": "2026-01-09T10:30:00Z",
    "version": 42,
    "reason": "manual_control"
  }
}
```

**4. Notifikasi Multi-Konsumen**
- Push WebSocket ke klien terhubung
- Broadcast MQTT ke topik perubahan state
- Callback webhook ke URL yang terdaftar
- Dukungan Server-Sent Events (SSE)

**5. Rekonsiliasi State**
- Verifikasi state berkala
- Deteksi konflik otomatis
- Penyembuhan state yang tidak konsisten
- Rollback ke state sebelumnya

**Endpoint API**:
```
GET  /api/v1/state/{entity_type}/{entity_id}
POST /api/v1/state/{entity_type}/{entity_id}
PUT  /api/v1/state/{entity_type}/{entity_id}/{field}
GET  /api/v1/state/{entity_type}/{entity_id}/history
GET  /api/v1/state/{entity_type}/{entity_id}/version/{version}
POST /api/v1/state/{entity_type}/{entity_id}/rollback
```

---

#### F5: REST & WebSocket API Komprehensif

**Deskripsi**: API berfitur lengkap untuk semua operasi

**Endpoint REST API**:

**Ingesti Data**:
```
POST /api/v2/ingest
POST /api/v2/ingest/batch
```

**Query Data**:
```
GET /api/v2/query
POST /api/v2/query/aggregate
GET /api/v2/query/time_range
```

**Manajemen State**:
```
GET  /api/v2/state/{entity}
POST /api/v2/state/{entity}
GET  /api/v2/state/{entity}/history
```

**Konfigurasi**:
```
GET  /api/v2/config
POST /api/v2/config/reload
GET  /api/v2/config/schema
```

**Monitoring**:
```
GET /api/v2/health
GET /api/v2/metrics
GET /api/v2/status
```

**API WebSocket**:
```
WS /api/v2/ws/subscribe
   - Subscribe ke data real-time
   - Subscribe ke perubahan state
   - Subscribe ke events sistem

WS /api/v2/ws/query
   - Hasil query live
   - Agregasi streaming
```

**Autentikasi**:
- Token JWT
- API keys
- OAuth 2.0 / OpenID Connect (opsional)

---

#### F6: Hot Configuration Reload

**Deskripsi**: Terapkan perubahan konfigurasi tanpa restart server

**Fitur**:
- Watch file konfigurasi untuk perubahan
- Validasi konfigurasi baru sebelum diterapkan
- Transisi graceful (tidak ada pesan yang terlewat)
- Rollback pada kegagalan
- Versioning konfigurasi

**Proses**:
1. Deteksi perubahan file konfigurasi
2. Validasi skema JSON
3. Tes koneksi database
4. Terapkan konfigurasi baru
5. Tutup koneksi lama
6. Mulai subscription baru
7. Pada error: rollback ke konfigurasi sebelumnya

---

### 9.2 Fitur Lanjutan

#### F7: Dukungan Multi-Tenansi

**Deskripsi**: Dukungan multiple deployment terisolasi dalam instance tunggal

**Fitur**:
- Isolasi tenant pada level data dan API
- Konfigurasi per-tenant
- Kuota sumber daya (pesan/detik, storage)
- Autentikasi spesifik tenant
- Logging audit per tenant

---

#### F8: Sistem Plugin

**Deskripsi**: Arsitektur yang dapat diperluas untuk fungsionalitas kustom

**Tipe Plugin**:
- Adapter storage kustom
- Parser kustom
- Transformasi kustom
- Provider autentikasi kustom
- Channel notifikasi kustom

**Interface Plugin**:
```go
type Plugin interface {
    Name() string
    Version() string
    Init(config map[string]interface{}) error
    Execute(ctx context.Context, input interface{}) (interface{}, error)
    Shutdown() error
}
```

---

#### F9: Integrasi Machine Learning (Fase 2)

**Deskripsi**: Kemampuan ML untuk pemrosesan data cerdas

**Fitur**:
- Deteksi anomali (isolation forest, autoencoders)
- Pemeliharaan prediktif
- Skoring kualitas data
- Alerting otomatis
- Forecasting

---

## 10. Roadmap Pengembangan

### Fase 1: Pondasi (Bulan 1-3)

**Sprint 1-2: Setup Proyek**
- [x] Struktur repository
- [ ] Setup environment pengembangan
- [ ] Pipeline CI/CD
- [ ] Template dokumentasi
- [ ] Setup pelacakan issue

**Sprint 3-4: Sistem Konfigurasi Inti**
- [ ] Definisi skema JSON
- [ ] Loader konfigurasi
- [ ] Validator konfigurasi
- [ ] Mekanisme hot-reload
- [ ] Pustaka template (3-5 template)

**Sprint 5-6: Mesin Pemrosesan MQTT**
- [ ] Refaktor klien MQTT
- [ ] Manajemen subscription dinamis
- [ ] Pencocokan pola topik
- [ ] Routing pesan
- [ ] Handling QoS

**Sprint 7-8: Interface Adapter Storage**
- [ ] Desain interface adapter
- [ ] Adapter PostgreSQL
- [ ] Adapter MongoDB
- [ ] Framework pengujian adapter

**Milestone**: Rilis Alpha - gateway MQTT yang dapat dikonfigurasi JSON dasar

---

### Fase 2: Fitur yang Ditingkatkan (Bulan 4-6)

**Sprint 9-10: Mesin Transformasi**
- [ ] Pipeline transformasi
- [ ] Transformasi built-in (20+)
- [ ] Dukungan transformasi kustom
- [ ] Enrichment data

**Sprint 11-12: Manajemen State yang Ditingkatkan**
- [ ] Versioning state
- [ ] Resolusi konflik
- [ ] Riwayat state
- [ ] Notifikasi multi-konsumen

**Sprint 13-14: Layer API**
- [ ] Implementasi REST API
- [ ] Dukungan WebSocket
- [ ] Autentikasi/Otorisasi
- [ ] Dokumentasi API (OpenAPI)

**Sprint 15-16: Adapter Storage Tambahan**
- [ ] Adapter InfluxDB
- [ ] Adapter TimescaleDB
- [ ] Optimasi performa
- [ ] Benchmarking

**Milestone**: Rilis Beta - UM-Gateway fitur lengkap

---

### Fase 3: Poles & Kesiapan Produksi (Bulan 7-8)

**Sprint 17-18: Pengujian & Kualitas**
- [ ] Suite tes integrasi
- [ ] Load testing (10K msg/detik)
- [ ] Audit keamanan
- [ ] Tuning performa

**Sprint 19-20: Operasi**
- [ ] Integrasi monitoring
- [ ] Peningkatan logging
- [ ] Health checks
- [ ] Prosedur backup/restore

**Sprint 21-22: Dokumentasi**
- [ ] Panduan pengguna
- [ ] Referensi API
- [ ] Panduan deployment
- [ ] Panduan troubleshooting
- [ ] Tutorial video

**Milestone**: v1.0 Ketersediaan Umum

---

### Fase 4: Fitur Lanjutan (Bulan 9-14)

**Sprint 23-26: Multi-Tenansi**
- [ ] Isolasi tenant
- [ ] Konfigurasi per-tenant
- [ ] Kuota sumber daya
- [ ] API manajemen tenant

**Sprint 27-30: Sistem Plugin**
- [ ] Framework plugin
- [ ] SDK plugin
- [ ] 5 plugin inti
- [ ] Marketplace plugin MVP

**Sprint 31-34: Dashboard Admin**
- [ ] UI konfigurasi
- [ ] Dashboard monitoring
- [ ] Viewer log
- [ ] Builder query

**Milestone**: v2.0 Enterprise Edition

---

### Fase 5: AI & Edge (Bulan 15-20)

**Sprint 35-38: Machine Learning**
- [ ] Deteksi anomali
- [ ] Pemeliharaan prediktif
- [ ] Skoring kualitas data
- [ ] Pipeline training model

**Sprint 39-42: Edge Computing**
- [ ] Versi ringan
- [ ] Mode offline
- [ ] Mekanisme sync
- [ ] Storage yang dioptimalkan untuk edge

**Milestone**: v3.0 dengan dukungan AI dan Edge

---

## 11. Metrik Keberhasilan

### Metrik Teknis

| Metrik | Target | Pengukuran |
|--------|--------|------------|
| Throughput Pesan | 100K msg/detik | Tes benchmark |
| Latensi (p99) | < 100ms | Monitoring performa |
| Uptime | 99.9% | Monitoring uptime |
| Waktu Reload Konfigurasi | < 5 detik | Tes otomatis |
| Waktu Respon API (p95) | < 200ms | Tool APM |

### Metrik Bisnis

| Metrik | Target | Timeline |
|--------|--------|----------|
| Instalasi Aktif | 100 | 6 bulan post-launch |
| Kontributor Komunitas | 20 | 12 bulan post-launch |
| Pelanggan Enterprise | 10 | 12 bulan post-launch |
| Bintang GitHub | 500 | 6 bulan post-launch |
| Pengguna Aktif Bulanan | 1.000 | 12 bulan post-launch |

### Metrik Kualitas

| Metrik | Target | Pengukuran |
|--------|--------|------------|
| Cakupan Tes | > 80% | Tool coverage kode |
| Bug Kritis | 0 di produksi | Pelacakan bug |
| Cakupan Dokumentasi | 100% API | Pengecekan otomatis |
| Waktu Deploy Instance Baru | < 30 menit | Survei pengguna |

---

## 12. Analisis Risiko

### Risiko Teknis

| Risiko | Dampak | Probabilitas | Mitigasi |
|--------|--------|--------------|----------|
| Performa tidak memenuhi target | Tinggi | Sedang | Benchmarking awal, profiling |
| Kompleksitas adapter database | Sedang | Tinggi | Batasi adapter awal, interface jelas |
| Hot-reload menyebabkan kehilangan data | Tinggi | Rendah | Pengujian ekstensif, deployment canary |
| Error konfigurasi | Sedang | Tinggi | Validasi, mode dry-run |

### Risiko Bisnis

| Risiko | Dampak | Probabilitas | Mitigasi |
|--------|--------|--------------|----------|
| Kompetitor merilis produk serupa | Tinggi | Sedang | Pengembangan cepat, fokus pada kemudahan penggunaan |
| Pasar tidak mengadopsi | Tinggi | Rendah | Pembangunan komunitas, tier gratis |
| Sumber daya terbatas | Sedang | Sedang | Pendekatan bertahap, prioritaskan fitur |

### Risiko Operasional

| Risiko | Dampak | Probabilitas | Mitigasi |
|--------|--------|--------------|----------|
| Kerentanan keamanan | Tinggi | Rendah | Audit keamanan, scanning dependensi |
| Kualitas dokumentasi | Sedang | Sedang | Technical writer, feedback pengguna |
| Beban dukungan | Sedang | Tinggi | Tool self-service, forum komunitas |

---

## 13. Kebutuhan Sumber Daya

### Struktur Tim

**Tim Inti (Fase 1-3)**
- 2x Pengembang Backend Senior (Go)
- 1x Engineer DevOps
- 1x Engineer QA
- 1x Technical Writer (part-time)

**Tim Diperluas (Fase 4-5)**
- +1x Pengembang Frontend (dashboard)
- +1x Engineer ML
- +1x Developer Advocate

**Pemangku Kepentingan**
- Manajer Produk
- Arsitek/Teknical Lead
- Ahli Domain (konsultan IoT)

### Estimasi Anggaran

| Kategori | Biaya (Bulanan) | Durasi |
|----------|------------------|---------|
| Gaji Tim Pengembangan | $25.000 | 14 bulan |
| Infrastruktur (Dev/Test) | $2.000 | 14 bulan |
| Tool & Layanan | $500 | Berkelanjutan |
| Dokumentasi & Desain | $2.000 | 8 bulan |
| Kontinjensi (20%) | $5.900 | - |
| **Total Fase 1-3** | **$35.400/bulan × 8** | **$283.200** |

---

## 14. Kriteria Go/No-Go

### Kriteria Keputusan Go

✅ **Lanjutkan jika**:
- Setidaknya 3 pelanggan potensial diwawancara menyatakan minat
- Kelayakan teknis dikonfirmasi melalui proof-of-concept
- Ketersediaan tim diamankan untuk minimum 8 bulan
- Anggaran disetujui untuk Fase 1-3

### Kriteria No-Go

❌ **Hentikan jika**:
- Riset pasar menunjukkan permintaan terbatas (< 10 pengguna potensial)
- Blocker teknis diidentifikasi tanpa solusi jelas
- Kompetitor memiliki produk mature dengan kemampuan serupa
- Kendala sumber daya mencegah pengembangan yang adekuat

### Go Bersyarat

⚠️ **Lanjutkan dengan modifikasi jika**:
- Minat pasar ada tetapi butuh set fitur berbeda
- Tantangan teknis memerlukan reduksi scope
- Kendala anggaran memerlukan pendekatan bertahap

---

## 15. Langkah Selanjutnya

### Tindakan Segera (Minggu 1-2)

1. **Review Pemangku Kepentingan**
   - Presentasikan brief proyek ke leadership
   - Amankan persetujuan anggaran
   - Dapatkan komitmen sumber daya

2. **Validasi Pasar**
   - Wawancara 5-10 pelanggan potensial
   - Survei komunitas developer IoT
   - Analisis lanskap kompetitor

3. **Proof-of-Concept Teknis**
   - Implementasikan prototype sistem konfigurasi
   - Tes interface adapter database
   - Validasikan mekanisme hot-reload

4. **Formasi Tim**
   - Rekrut anggota tim inti
   - Definisikan peran dan tanggung jawab
   - Siapkan saluran komunikasi

### Tindakan Jangka Pendek (Bulan 1)

5. **Environment Pengembangan**
   - Siapkan repository (Git)
   - Konfigurasikan pipeline CI/CD
   - Buat wiki proyek

6. **Perencanaan**
   - Perencanaan sprint detail
   - Review arsitektur
   - Asesmen keamanan

7. **Dokumentasi**
   - Buat pedoman kontribusi
   - Siapkan situs dokumentasi
   - Tulis README dan panduan memulai

---

## 16. Kesimpulan

Universal MQTT Gateway Server mewakili peluang signifikan untuk mengatasi titik nyeri yang meluas di ekosistem IoT. Dengan menciptakan solusi yang dapat digunakan kembali, dapat dikonfigurasi, dan dapat dipelihara, kita dapat:

- **Mempercepat adopsi IoT** di seluruh industri
- **Mengurangi biaya pengembangan** untuk integrator sistem
- **Membangun Rafflesia Agro** sebagai pemimpin teknologi
- **Membangun produk berkelanjutan** dengan multiple stream pendapatan

Kunci keberhasilan terletak pada menjaga sistem **sederhana namun powerful**, **terdokumentasi dengan baik**, dan **digerakkan komunitas**. Pendekatan konfigurasi berbasis JSON, adapter storage yang dapat dipasang, dan manajemen state yang ditingkatkan akan membedakan UM-Gateway dari solusi yang ada.

Dengan eksekusi yang fokus selama 8 bulan ke depan, kita dapat mengirimkan v1.0 produksi yang mengatasi kebutuhan pasar nyata dan membangun pondasi untuk pertumbuhan dan inovasi jangka panjang.

---

**Versi Dokumen**: 1.0
**Terakhir Diperbarui**: 9 Januari 2026
**Status**: Draft - Menunggu Review
**Review Berikutnya**: 16 Januari 2026

---

## Lampiran

### A. Glossarium

- **MQTT**: Message Queuing Telemetry Transport
- **Telemetry**: Pengumpulan data otomatis dari sensor jarak jauh
- **Database Time-Series**: Database yang dioptimalkan untuk data dengan timestamp
- **Skema**: Definisi struktur untuk data
- **Hot-Reload**: Memperbarui konfigurasi tanpa restart layanan
- **Multi-Tenansi**: Instance tunggal melayani beberapa pelanggan yang terisolasi

### B. Referensi

1. Codebase Go-MQTT saat ini (v1.0)
2. Spesifikasi MQTT 3.1.1 dan 5.0
3. Dokumentasi TimescaleDB
4. Koleksi Time-Series MongoDB
5. Dokumentasi InfluxDB
6. Laporan Riset Pasar IoT

### C. Dokumen Terkait

- `PROJECT_DESCRIPTION.md` - Dokumentasi sistem saat ini
- `PROJECT_DESCRIPTION_ID.md` - Sistem saat ini (Bahasa Indonesia)
- `openapi.yaml` - Spesifikasi API saat ini
- Diagram arsitektur (akan dibuat)

---

**Disiapkan oleh**: Tim Pengembangan
**Disetujui oleh**: [Pending]
**Tanggal Review**: [Pending]
