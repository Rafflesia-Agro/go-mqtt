# Rencana Proyek: Universal MQTT Gateway Server (UM-Gateway) MVP v0.1

## Ringkasan Eksekutif

**Nama Proyek**: Universal MQTT Gateway Server (UM-Gateway)
**Versi**: MVP v0.1 (Minimum Viable Product)
**Status**: Perencanaan
**Target Rilis**: Q2 2026 (3 bulan pengembangan)
**Organisasi**: Divisi Teknologi Rafflesia Agro

---

## Riwayat Versi

| Versi | Tanggal | Perubahan | Penulis |
|-------|---------|-----------|--------|
| **MVP v0.1** | 2026-01-09 | **Minimum Viable Product**: Single executable dengan embedded MQTT<br>• Single binary deployment (mode embedded saja)<br>• Dashboard admin dasar (React + Vite)<br>• Setup mode & Runner mode<br>• Wizard onboarding<br>• Cache lokal (in-memory)<br>• Database lokal (SQLite WAL)<br>• Broker MQTT embedded lokal (mochi-mqtt)<br>• Konfigurasi dasar (config.yaml)<br>• Fitur esensial saja (tanpa fitur SaaS)<br>• Waktu pengembangan: 3 bulan<br>• Anggaran: $65,000 | Tim Pengembangan |

---

## 1. Latar Belakang

### Apa itu MVP?

**MVP (Minimum Viable Product)**: Produk dengan fitur yang cukup untuk dapat digunakan oleh pelanggan awal dan memberikan umpan balik untuk pengembangan selanjutnya.

**Ini BUKAN**:
- ❌ Platform SaaS lengkap
- ❌ Sistem multi-tenant
- ❌ Deployment cloud
- ❌ Fitur enterprise
- ❌ Monitoring lanjutan

**Ini ADALAH**:
- ✅ Single executable yang dapat digunakan kembali
- ✅ Deployment mandiri
- ✅ Konfigurasi dan manajemen dasar
- ✅ Fungsionalitas MQTT gateway esensial
- ✅ Pengembangan cepat (3 bulan)

### Kenapa MVP?

1. **Kecepatan Pasar**: Produk kerja dalam 3 bulan
2. **Fokus**: Hanya fitur inti, tanpa gangguan
3. **Kesederhanaan**: Mudah dipahami, di-deploy, dan dirawat
4. **Validasi**: Uji asumsi inti sebelum membangun fitur lanjutan
5. **Risiko Lebih Rendah**: Investasi lebih sedikit, umpan balik lebih cepat

### Kebutuhan Pasar

Banyak deployment IoT kecil membutuhkan **MQTT gateway yang sederhana dan mandiri**:
- Penggemar smart home
- Peternakan kecil (seperti Rafflesia Agro)
- Setup industri kecil
- Lingkungan pengembangan/pengujian
- Prototipe edge computing

Mereka membutuhkan sesuatu yang **langsung bekerja** tanpa setup yang rumit.

---

## 2. Model Bisnis

### Proposition Nilai

**Untuk Deployment Kecil**:
- Single binary, tanpa setup kompleks
- Berjalan di komputer apa pun
- Tanpa ketergantungan eksternal
- Gratis dan open source

**Untuk Pengembangan**:
- Pengujian lokal cepat
- Konfigurasi mudah
- Broker MQTT bawaan
- Manajemen berbasis web

### Model Pendapatan

**Fase MVP (6 bulan pertama)**:
- **100% Gratis**: Open source (Lisensi MIT)
- Fokus pada adopsi dan umpan balik
- Belum ada fitur berbayar

**Fase Masa Depan** (Setelah MVP terbukti bernilai):
- Lisensi Pro untuk fitur lanjutan
- Kontrak dukungan
- Pengembangan kustom

### Target Skala

**Tujuan MVP**:
- 100 instalasi aktif dalam 3 bulan setelah rilis
- Umpan balik positif dari early adopters
- Pemahaman jelas tentang fitur apa yang akan ditambahkan selanjutnya

---

## 3. Target Audiens

### Pengguna Utama

**1. Proyek IoT Kecil**
- Otomatisasi smart home (1-50 perangkat)
- Monitoring peternakan kecil (seperti peternakan unggas)
- Proyek hobi
- Proyek mahasiswa

**2. Tim Pengembangan**
- Pengujian aplikasi IoT
- Prototyping solusi
- Lingkungan pengembangan lokal
- Pengujian integrasi

**3. Integrator Sistem**
- Deployment pelanggan kecil
- Proyek proof of concept
- Kebutuhan gateway sederhana
- Deployment cepat

### Profil Pengguna

**Keahlian Teknis**:
- Pemahaman dasar MQTT
- Dapat menjalankan file executable
- Dapat mengkonfigurasi pengaturan sederhana
- Mungkin bukan developer profesional

**Lingkungan Deployment**:
- Single machine (PC, server, Raspberry Pi)
- Jaringan lokal
- Tidak memerlukan konektivitas cloud
- Tanpa ketergantungan eksternal

---

## 4. Pernyataan Masalah

### Masalah Inti

Deployment IoT kecil membutuhkan **MQTT gateway yang sederhana dan mandiri** tetapi solusi yang ada:
1. **Terlalu Kompleks**: Memerlukan broker, database, cache terpisah
2. **Sulit Dikonfigurasi**: Memerlukan keahlian teknis
3. **Mahal**: Solusi komersial terlalu mahal
4. **Over-engineered**: Platform SaaS lengkap untuk kebutuhan sederhana

### Solusi MVP

**Single executable yang mencakup**:
- ✅ Embedded MQTT broker (tanpa broker terpisah)
- ✅ Database lokal (tanpa database eksternal)
- ✅ In-memory cache (tanpa cache eksternal)
- ✅ Antarmuka admin web (tanpa editing file config)
- ✅ Setup wizard (konfigurasi first-time mudah)

**Apa yang TIDAK Kami Bangun (Saat Ini)**:
- ❌ Dukungan database eksternal (PostgreSQL, MongoDB)
- ❌ Dukungan cache eksternal (Redis, Memcached)
- ❌ Integrasi cloud
- ❌ Multi-tenancy
- ❌ Monitoring & alerting lanjutan
- ❌ High availability & clustering
- ❌ Autentikasi & otorisasi pengguna
- ❌ API untuk integrasi third-party

---

## 5. Solusi yang Diusulkan

### Visi: MQTT Gateway yang Sederhana dan Mandiri

**Satu file executable** yang menyediakan:
1. **Embedded MQTT Broker** - Menangani koneksi MQTT
2. **Database Lokal** - Menyimpan pesan dan konfigurasi
3. **In-Memory Cache** - Meningkatkan performa
4. **Antarmuka Admin Web** - Manajemen mudah
5. **Setup Wizard** - Setup pertama kali terpandu

### Arsitektur

```
┌─────────────────────────────────────────────────────────────┐
│           Single Executable: um-gateway.exe                 │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌───────────────────────────────────────────────────────┐  │
│  │              Backend Go (Single Process)               │  │
│  │                                                        │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌────────────┐ │  │
│  │  │ Embedded     │  │ HTTP API    │  │ In-Memory  │ │  │
│  │  │ MQTT Broker  │  │ Server      │  │ Cache      │ │  │
│  │  │ (mochi-mqtt) │  │             │  │ (go-cache) │ │  │
│  │  └──────────────┘  └──────────────┘  └────────────┘ │  │
│  │                                                        │  │
│  │  ┌──────────────┐  ┌──────────────┐                  │  │
│  │  │ SQLite       │  │ Config YAML  │                  │  │
│  │  │ (Mode WAL)   │  │ Loader       │                  │  │
│  │  └──────────────┘  └──────────────┘                  │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                               │
│  ┌───────────────────────────────────────────────────────┐  │
│  │         UI Web React + Vite (Embedded)                │  │
│  │                                                        │  │
│  │  ┌────────────────┐  ┌────────────────┐              │  │
│  │  │ Wizard         │  │ Setup Mode     │              │  │
│  │  │ Onboarding     │  │ (Edit Config)  │              │  │
│  │  └────────────────┘  └────────────────┘              │  │
│  │                                                        │  │
│  │  ┌────────────────┐  ┌────────────────┐              │  │
│  │  │ Runner Mode    │  │ Dashboard      │              │  │
│  │  │ (Monitoring)   │  │ (Read-only)    │              │  │
│  │  └────────────────┘  └────────────────┘              │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### Fitur Utama

#### F1: Embedded MQTT Broker (mochi-mqtt)
- Protokol MQTT 3.1.1
- Hingga 100 koneksi konkuren
- TCP listener pada port yang dapat dikonfigurasi (default 1883)
- Autentikasi dasar (username/password)
- Persistensi pesan (opsional)

#### F2: Database Lokal (SQLite WAL)
- Database single file
- Write-Ahead Logging untuk konkurensi lebih baik
- Migrasi skema otomatis
- Menyimpan: devices, topics, messages, metrics
- Utilitas backup/pemulihan

#### F3: In-Memory Cache (go-cache)
- TTL default: 1 jam
- Maksimum 10,000 item
- Meningkatkan performa query
- Ekspirasi otomatis

#### F4: Antarmuka Admin Web (React + Vite)
- UI modern dan responsif
- Tampilan metrik real-time
- Manajemen device
- Browser topic
- Viewer riwayat pesan
- Viewer konfigurasi (read-only di runner mode)

#### F5: Setup Mode
- Edit config.yaml via web UI
- Editor YAML dengan validasi
- Konfigurasi berbasis form
- Tombol test konfigurasi
- Simpan & restart
- Backup konfigurasi

#### F6: Runner Mode
- Operasi produksi
- Konfigurasi read-only
- Dashboard monitoring
- Metrik real-time
- Tidak dapat mengedit config tanpa restart

#### F7: Wizard Onboarding
- Deteksi first-run
- Konfigurasi langkah-demi-langkah
- Pengaturan dasar:
  - Username/password admin
  - Port broker MQTT
  - Lokasi database
  - Port UI web
- Test koneksi
- Mulai runner mode

---

## 6. Ruang Lingkup & Batasan Proyek

### Dalam Ruang Lingkup (Fitur MVP)

#### Must Have (Fungsionalitas Inti)

1. **Embedded MQTT Broker**
   - Fungsionalitas broker dasar
   - TCP listener
   - Autentikasi klien
   - Routing pesan
   - Persistensi opsional

2. **Database Lokal**
   - SQLite dengan mode WAL
   - Registry device
   - Definisi topic
   - Penyimpanan pesan
   - Metrik dasar

3. **In-Memory Cache**
   - Implementasi go-cache
   - Ekspirasi berbasis TTL
   - Batas ukuran

4. **UI Admin Web**
   - Dashboard dengan metrik
   - Daftar device
   - Browser topic
   - Viewer pesan
   - Viewer config (read-only di runner mode)

5. **Setup Mode**
   - Editor config YAML
   - Config berbasis form
   - Validasi
   - Simpan & restart

6. **Runner Mode**
   - Operasi produksi
   - Config read-only
   - Monitoring

7. **Wizard Onboarding**
   - Deteksi first-run
   - Konfigurasi dasar
   - Test koneksi
   - Mulai gateway

8. **Konfigurasi (config.yaml)**
   - Pengaturan database
   - Pengaturan broker MQTT
   - Pengaturan web server
   - Pengaturan logging
   - Kebijakan retensi dasar

### Di Luar Ruang Lingkup (Versi Masa Depan)

**TIDAK termasuk di MVP**:
- ❌ Dukungan database eksternal (PostgreSQL, MongoDB, InfluxDB)
- ❌ Dukungan cache eksternal (Redis, Memcached)
- ❌ Koneksi broker MQTT eksternal
- ❌ Manajemen pengguna (multi-user)
- ❌ Role-based access control
- ❌ API untuk integrasi third-party
- ❌ Monitoring & alerting lanjutan
- ❌ High availability & clustering
- ❌ Export/import data
- ❌ Template konfigurasi
- ❌ Keamanan lanjutan (SSL/TLS)
- ❌ Dukungan WebSocket untuk MQTT
- ❌ Fitur MQTT 5.0
- ❌ Konfigurasi bridge

### Kendala Teknis

**Kendala MVP**:
- **Mode embedded SAJA** (tanpa broker/database/cache eksternal)
- **Deployment single machine** (tanpa clustering)
- **Database SQLite SAJA** (tanpa PostgreSQL/MongoDB)
- **Cache in-memory SAJA** (tanpa Redis)
- **Autentikasi dasar SAJA** (username/password di config)
- **Single admin user** (tanpa dukungan multi-user)
- **HTTP SAJA** (tanpa HTTPS/TLS di MVP)

**Persyaratan Hardware**:
- Minimum: 1 core CPU, 1GB RAM, 100MB disk
- Direkomendasikan: 2 core CPU, 2GB RAM, 1GB disk

**Persyaratan Software**:
- Go 1.24+
- SQLite 3.40+
- Browser web modern (Chrome, Firefox, Edge, Safari)

---

## 7. Teknologi Stack

### Backend (Go)

| Komponen | Teknologi | Versi | Kenapa |
|-----------|-----------|-------|--------|
| Bahasa | Go | 1.24+ | Performa, single binary |
| MQTT Broker | mochi-mqtt | 2.0+ | Embedded, pure Go |
| Database | SQLite | 3.40+ | Single file, mode WAL |
| Cache | go-cache | Latest | In-memory, sederhana |
| HTTP Router | Chi | 5.0+ | Ringan, idiomatic |
| YAML | yaml.v3 | Latest | Parsing config |
| Embed | embed | std | Embed UI dalam binary |

### Frontend (React)

| Komponen | Teknologi | Versi | Kenapa |
|-----------|-----------|-------|--------|
| Framework | React | 18.3+ | Populer, familiar |
| Build Tool | Vite | 5.0+ | Cepat, sederhana |
| Bahasa | TypeScript | 5.3+ | Keamanan tipe |
| Library UI | - | - | CSS biasa (tanpa lib komponen) |
| HTTP | Fetch API | std | Tidak perlu Axios |
| Router | - | - | Routing hash sederhana |

### Alat Pengembangan

| Alat | Tujuan |
|------|---------|
| Git | Kontrol versi |
| Go modules | Manajemen dependensi |
| npm/yarn | Dependensi frontend |
| air | Hot reload untuk Go |
| Vite HMR | Hot reload untuk React |

---

## 8. Fitur (Ruang Lingkup MVP)

### F1: Embedded MQTT Broker

**Implementasi**:

```go
package broker

import mqtt "github.com/mochi-mqtt/mqtt/v2"

type EmbeddedBroker struct {
    server *mqtt.Server
    config *BrokerConfig
}

func NewEmbeddedBroker(config *BrokerConfig) (*EmbeddedBroker, error) {
    server := mqtt.NewServer(nil)

    // Tambahkan hooks untuk logging dan auth
    hooks := &hooks.Auth{
        Authenticate: func(cl *mqtt.Client, user, pass string) bool {
            // Pengecekan username/password sederhana
            return user == config.Username && pass == config.Password
        },
    }
    server.AddHook(hooks)

    // TCP listener
    tcp := mqtt.NewTCPListener(config.ListenAddress, nil)
    server.AddListener(tcp)

    return &EmbeddedBroker{server: server, config: config}, nil
}

func (b *EmbeddedBroker) Start() error {
    go b.server.Serve()
    return nil
}

func (b *EmbeddedBroker) Stop() {
    b.server.Close()
}
```

**Config (config.yaml)**:

```yaml
mqtt:
  embedded:
    enabled: true
    listen_address: :1883
    username: admin
    password: changeme
    max_connections: 100
```

---

### F2: Database Lokal (SQLite)

**Skema**:

```sql
-- Devices
CREATE TABLE devices (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT,
    first_seen DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_seen DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Topics
CREATE TABLE topics (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Messages (dengan retensi)
CREATE TABLE messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    topic TEXT NOT NULL,
    payload TEXT NOT NULL,
    qos INTEGER DEFAULT 0,
    retained BOOLEAN DEFAULT 0,
    published_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Buat indexes
CREATE INDEX idx_messages_topic ON messages(topic);
CREATE INDEX idx_messages_published_at ON messages(published_at);
```

**Implementasi**:

```go
package database

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

type SQLiteDatabase struct {
    db *sql.DB
}

func NewSQLiteDatabase(path string) (*SQLiteDatabase, error) {
    db, err := sql.Open("sqlite3", path)
    if err != nil {
        return nil, err
    }

    // Aktifkan mode WAL
    _, err = db.Exec("PRAGMA journal_mode=WAL")
    if err != nil {
        return nil, err
    }

    // Jalankan migrasi
    if err := runMigrations(db); err != nil {
        return nil, err
    }

    return &SQLiteDatabase{db: db}, nil
}
```

---

### F3: In-Memory Cache

**Implementasi**:

```go
package cache

import (
    "github.com/patrickmn/go-cache"
    "time"
)

type MemoryCache struct {
    store *cache.Cache
}

func NewMemoryCache() *MemoryCache {
    // TTL default: 1 jam, cleanup setiap 10 menit
    return &MemoryCache{
        store: cache.New(1*time.Hour, 10*time.Minute),
    }
}

func (c *MemoryCache) Get(key string) (interface{}, bool) {
    return c.store.Get(key)
}

func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) {
    if ttl == 0 {
        ttl = cache.DefaultExpiration
    }
    c.store.Set(key, value, ttl)
}

func (c *MemoryCache) Delete(key string) {
    c.store.Delete(key)
}
```

---

### F4: UI Admin Web

**Halaman**:

1. **Dashboard** (`/`)
   - Metrik sistem (koneksi aktif, pesan/menit, uptime)
   - Statistik cepat (total device, total topic, total pesan)
   - Aktivitas terkini

2. **Devices** (`/devices`)
   - Daftar semua device
   - Detail device
   - Search/filter

3. **Topics** (`/topics`)
   - Daftar semua topic
   - Detail topic
   - Jumlah pesan per topic

4. **Messages** (`/messages`)
   - Riwayat pesan (paginasi)
   - Filter berdasarkan topic
   - Lihat detail pesan

5. **Setup** (`/setup`)
   - Editor config YAML
   - Config berbasis form
   - Simpan & restart

6. **Config** (`/config`)
   - Viewer config read-only (runner mode)
   - Download config

**Contoh Komponen**:

```typescript
// frontend/src/pages/Dashboard.tsx
import { useQuery } from '@tanstack/react-query';

export function Dashboard() {
  const { data: metrics, isLoading } = useQuery({
    queryKey: ['metrics'],
    queryFn: async () => {
      const res = await fetch('/api/metrics');
      return res.json();
    },
    refetchInterval: 5000, // Auto-refresh
  });

  if (isLoading) return <div>Loading...</div>;

  return (
    <div className="dashboard">
      <h1>Dashboard</h1>

      <div className="metrics-grid">
        <MetricCard title="Klien Terhubung" value={metrics.clients} />
        <MetricCard title="Pesan/menit" value={metrics.msgPerMin} />
        <MetricCard title="Uptime" value={metrics.uptime} />
      </div>

      <div className="stats-grid">
        <StatCard title="Total Device" value={metrics.totalDevices} />
        <StatCard title="Total Topic" value={metrics.totalTopics} />
        <StatCard title="Total Pesan" value={metrics.totalMessages} />
      </div>
    </div>
  );
}
```

---

### F5: Setup Mode

**Tujuan**: Edit config.yaml

**Akses**: `/setup` (hanya saat mode=setup di config.yaml)

**Fitur**:
- Editor YAML dengan syntax highlighting
- Konfigurasi berbasis form
- Validasi
- Tombol test
- Tombol simpan & restart
- Backup sebelum menyimpan

**Alur**:
```
1. Hentikan gateway
2. Mulai gateway dengan --mode=setup
3. Buka http://localhost:8080/setup
4. Edit konfigurasi
5. Klik "Simpan & Restart"
6. Gateway restart dalam mode runner
```

---

### F6: Runner Mode

**Tujuan**: Operasi produksi

**Akses**: `/` (saat mode=runner di config.yaml)

**Fitur**:
- Viewer config read-only
- Dashboard monitoring
- Metrik real-time
- Tidak dapat mengedit config

**Perilaku**:
```
1. Mulai gateway (atau setelah setup)
2. config.yaml menjadi read-only
3. Embedded MQTT broker mulai
4. Web UI menampilkan dashboard
5. Gateway menerima koneksi MQTT
```

---

### F7: Wizard Onboarding

**Tujuan**: Setup pertama kali

**Trigger**: Tidak ada config.yaml ditemukan

**Langkah**:

**Langkah 1: Selamat Datang**
- Jelaskan apa yang dilakukan gateway
- Tombol "Mulai"

**Langkah 2: Kredensial Admin**
- Username (default: admin)
- Password (wajib)
- Konfirmasi password

**Langkah 3: Konfigurasi MQTT**
- Port listen (default: 1883)
- Max koneksi (default: 100)

**Langkah 4: Konfigurasi UI Web**
- Port (default: 8080)

**Langkah 5: Konfigurasi Database**
- Path (default: ./data/gateway.db)

**Langkah 6: Review**
- Tampilkan semua pengaturan
- Tombol "Selesai & Mulai"

**Langkah 7: Selesai**
- Gateway mulai
- Tampilkan URL dashboard
- Tampilkan URL broker MQTT

---

## 9. Peta Jalan Pengembangan

### Rencana 3 Bulan

#### Bulan 1: Pondasi

**Minggu 1-2: Setup Proyek**
- [ ] Setup repository (Go + React)
- [ ] Struktur proyek dasar
- [ ] CI/CD (GitHub Actions)
- [ ] Setup lingkungan pengembangan
- [ ] Server Go dasar dengan UI embedded

**Minggu 3-4: Backend Inti**
- [ ] Embedded MQTT broker (mochi-mqtt)
- [ ] Setup database SQLite
- [ ] Endpoint API dasar
- [ ] In-memory cache
- [ ] Loading config (YAML)

**Milestone**: Backend dasar berfungsi

---

#### Bulan 2: Fitur

**Minggu 5-6: MQTT & Database**
- [ ] Hooks broker MQTT (auth, logging)
- [ ] Persistensi pesan
- [ ] Skema & migrasi database
- [ ] Tracking device
- [ ] Tracking topic
- [ ] Penyimpanan pesan

**Minggu 7-8: Web UI**
- [ ] Setup React + Vite
- [ ] Layout dasar
- [ ] Halaman dashboard
- [ ] Halaman devices
- [ ] Halaman topics
- [ ] Halaman messages

**Milestone**: Fitur inti berfungsi

---

#### Bulan 3: Poles & Rilis

**Minggu 9-10: Setup & Onboarding**
- [ ] UI setup mode
- [ ] Editor config YAML
- [ ] Wizard onboarding
- [ ] Validasi config
- [ ] Fungsionalitas simpan & restart

**Minggu 11-12: Pengujian & Rilis**
- [ ] Pengujian end-to-end
- [ ] Perbaikan bug
- [ ] Dokumentasi
- [ ] Build untuk platform (Linux, Windows, macOS)
- [ ] Rilis v0.1.0

**Milestone**: Rilis MVP

---

## 10. Metrik Kesuksesan

### Metrik Teknis

| Metrik | Target | Pengukuran |
|--------|--------|-------------|
| Ukuran Binary | < 50 MB | Artefak build |
| Waktu Startup | < 2 detik | Pengujian manual |
| Penggunaan Memori | < 100 MB idle | Monitoring sumber daya |
| Pesan MQTT/detik | > 1,000 | Pengujian beban |
| Klien Konkuren | 100 | Pengujian beban |
| Load Halaman Web | < 1 detik | Browser DevTools |

### Metrik Adopsi

| Metrik | Target | Timeline |
|--------|--------|----------|
| GitHub Stars | 100 | 1 bulan pasca-rilis |
| Downloads | 500 | 3 bulan pasca-rilis |
| Instalasi Aktif | 50 | 3 bulan pasca-rilis |
| Isu Dilaporkan | < 10 kritis | 3 bulan pasca-rilis |
| Kontribusi Komunitas | 5+ | 6 bulan pasca-rilis |

### Metrik Kualitas

| Metrik | Target | Pengukuran |
|--------|--------|-------------|
| Tingkat Penyelesaian Onboarding | > 90% | Analytics |
| Tingkat Sukses Setup | 100% | Pengujian manual |
| Cakupan Dokumentasi | 100% fitur | Review manual |
| Bug Kritis Diketahui | 0 | GitHub Issues |

---

## 11. Kebutuhan Sumber Daya

### Struktur Tim

**Tim Minimum untuk MVP 3 Bulan**:
- **1 Developer Full-Stack** (Go + React)
- **1 QA Part-Time** (20 jam/minggu)
- **1 Technical Writer Part-Time** (10 jam/minggu)

### Estimasi Anggaran

| Kategori | Biaya | Durasi |
|----------|------|---------|
| Developer Full-Stack | $15,000/bulan | 3 bulan |
| QA (part-time) | $2,000/bulan | 3 bulan |
| Technical Writer (part-time) | $1,000/bulan | 3 bulan |
| Infrastruktur (Dev/Test) | $500/bulan | 3 bulan |
| Alat & Layanan | $200/bulan | 3 bulan |
| **Total** | **$18,700/bulan × 3** | **$56,100** |

**Buffer (15%)**: $8,400
**Total dengan Buffer**: **$65,000**

---

## 12. Perbandingan: MVP vs Versi Lengkap

| Fitur | MVP v0.1 | Full v3.0 |
|---------|----------|-----------|
| **Waktu Pengembangan** | 3 bulan | 10 bulan |
| **Anggaran** | $65,000 | $215,500 |
| **MQTT Broker** | Embedded only | Embedded + External |
| **Database** | SQLite only | SQLite + PostgreSQL + MongoDB |
| **Cache** | In-memory only | In-memory + Redis + Memcached |
| **Autentikasi** | Dasar | Lanjutan (JWT, OAuth) |
| **Manajemen Pengguna** | Single user | Multi-user, RBAC |
| **Monitoring** | Dashboard dasar | Monitoring & alerting lanjutan |
| **API** | None | RESTful API |
| **High Availability** | No | Yes (clustering) |
| **Fitur SaaS** | No | Yes |
| **Target** | Deployment kecil | Enterprise + Kecil |

---

## 13. Langkah Selanjutnya

### Segera (Minggu 1)

1. **Setup Repository**
   - Buat repository GitHub
   - Setup struktur proyek
   - Konfigurasi CI/CD

2. **Lingkungan Pengembangan**
   - Setup lingkungan Go lokal
   - Setup React + Vite
   - Konfigurasi hot reload

3. **Commit Awal**
   - Server Go dasar
   - UI embedded
   - Halaman Hello World

### Jangka Pendek (Bulan 1)

4. **Backend Inti**
   - Integrasi broker MQTT
   - Setup database
   - Implementasi cache

5. **UI Dasar**
   - Skeleton dashboard
   - Daftar device
   - Daftar topic

---

## 14. Kesimpulan

### Filosofi MVP

**Tujuannya adalah membangun sesuatu yang KECIL yang BEKERJA, bukan sesuatu yang BESAR yang tidak bekerja.**

**Yang Penting**:
- ✅ Single executable
- ✅ Embedded MQTT broker
- ✅ Langsung bekerja
- ✅ Mudah dikonfigurasi
- ✅ Monitoring dasar

**Yang Tidak Penting (Saat Ini)**:
- ❌ Database eksternal
- ❌ Integrasi cloud
- ❌ Fitur enterprise
- ❌ Monitoring lanjutan
- ❌ Multi-tenancy

### Kriteria Sukses

**MVP berhasil jika**:
1. Pengguna dapat mengunduh single executable
2. Menjalankannya di mesin mereka
3. Menyelesaikan wizard onboarding
4. Menghubungkan perangkat MQTT
5. Melihat pesan di dashboard
6. Semua bekerja tanpa ketergantungan eksternal

Dengan eksekusi yang terfokus selama 3 bulan, kami dapat mengirimkan **MVP yang berfungsi** yang membuktikan konsep dan memberikan nilai bagi early adopters.

---

**Versi Dokumen**: MVP v0.1
**Terakhir Diperbarui**: 9 Januari 2026
**Status**: Draft - Menunggu Review
**Waktu Pengembangan**: 3 bulan
**Anggaran**: $65,000

---

## Lampiran

### A. Checklist MVP

**Must Have untuk MVP**:
- [ ] Single executable binary
- [ ] Embedded MQTT broker
- [ ] Database lokal (SQLite)
- [ ] In-memory cache
- [ ] Antarmuka admin web
- [ ] Setup mode
- [ ] Runner mode
- [ ] Wizard onboarding
- [ ] Konfigurasi dasar (config.yaml)
- [ ] Dashboard dengan metrik
- [ ] Manajemen device
- [ ] Browser topic
- [ ] Viewer pesan
- [ ] Dokumentasi

### B. Fitur Masa Depan (Post-MVP)

**Fase 2** (6 bulan setelah MVP):
- Dukungan database eksternal
- REST API
- Monitoring lanjutan
- Autentikasi pengguna

**Fase 3** (12 bulan setelah MVP):
- High availability
- Clustering
- Integrasi cloud
- Aplikasi mobile

### C. Referensi

1. Dokumentasi mochi-mqtt
2. Dokumentasi SQLite
3. Dokumentasi go-cache
4. Dokumentasi React
5. Dokumentasi Vite

---

**Disiapkan oleh**: Tim Pengembangan
**Disetujui oleh**: [Pending]
**Tanggal Review**: [Pending]
**Target Rilis**: Q2 2026
