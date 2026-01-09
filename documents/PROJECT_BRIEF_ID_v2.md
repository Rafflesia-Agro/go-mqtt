# Rencana Proyek: Universal MQTT Gateway Server (UM-Gateway) v2.0

## Ringkasan Eksekutif

**Nama Proyek**: Universal MQTT Gateway Server (UM-Gateway)
**Versi**: 2.0 (Single Executable dengan Embedded MQTT)
**Status**: Proposal/Perencanaan
**Target Rilis**: Q4 2026
**Organisasi**: Rafflesia Agro Technology Division

---

## Riwayat Versi

| Versi | Tanggal | Perubahan | Penulis |
|-------|---------|-----------|--------|
| **v2.0** | 2026-01-09 | **Perubahan Arsitektur Utama**: Implementasi embedded MQTT broker<br>• Menambahkan embedded MQTT broker menggunakan library mochi-mqtt<br>• Memperkenalkan operasi dual-mode (setup/runner)<br>• Mengubah konfigurasi dari database ke YAML config.yaml<br>• Menggantikan cache eksternal dengan in-memory cache (go-cache)<br>• Mode WAL SQLite sebagai default (dioptimalkan untuk konkurensi)<br>• Arsitektur single executable dengan koordinasi proses<br>• Contoh kode baru untuk embedded broker, mode WAL, adapter cache<br>• Anggaran: $238,150 (10 bulan) | Tim Pengembangan |
| **v1.0** | 2026-01-09 | **Dashboard Admin Web**: Pendekatan konfigurasi visual<br>• Menambahkan dashboard admin SvelteKit untuk konfigurasi<br>• Database SQLite dengan wizard onboarding<br>• Menggantikan konfigurasi JSON dengan form UI web<br>• Model koneksi broker MQTT eksternal<br>• Anggaran: $178,650 (berkurang dari v0) | Tim Pengembangan |
| **v0.0** | 2026-01-09 | **Proposal Asli**: Konfigurasi berbasis JSON<br>• Konsep awal dengan definisi file JSON<br>• Koneksi broker MQTT eksternal diperlukan<br>• PostgreSQL/MongoDB dapat dikonfigurasi via JSON<br>• Fokus pada kasus penggunaan farm telemetry<br>• Anggaran: $283,200 | Tim Pengembangan |

---

## 1. Latar Belakang

### Evolusi dari v1 ke v2

**Pendekatan v1**: Dashboard admin web dengan SQLite, koneksi broker MQTT terpisah
**Pendekatan v2**: Single executable yang sepenuhnya mandiri dengan embedded MQTT broker, operasi dual-mode

### Keterbatasan v1

1. **Ketergantungan External MQTT Broker**: Masih memerlukan broker MQTT eksternal (Mosquitto, HiveMQ, dll)
2. **Tidak Ada Operasi Offline yang Sebenarnya**: Tidak dapat berfungsi tanpa layanan eksternal
3. **Kebingungan Setup vs Runtime**: Perubahan konfigurasi memerlukan menjalankan ulang setup
4. **Skenario Deploy Terbatas**: Tidak dapat berjalan di lingkungan air-gapped dengan mudah
5. **Ketergantungan Cache/State**: Redis diperlukan untuk penggunaan produksi

### Peluang Pasar

Lanskap IoT membutuhkan **MQTT gateway yang sepenuhnya mandiri** yang:
- **Edge Deployment**: Berjalan pada perangkat edge tanpa ketergantungan eksternal
- **Lingkungan Air-Gapped**: Berfungsi dalam jaringan terisolasi
- **Operasi yang Disederhanakan**: Single executable, tanpa infrastruktur terpisah
- **Kasus Penggunaan Embedded**: Deploy pada perangkat dengan sumber daya terbatas

---

## 2. Model Bisnis

### Proposition Nilai

**Untuk Edge Deployments**:
- Single binary mencakup semua yang dibutuhkan
- Tidak perlu mengelola broker MQTT terpisah
- Kompleksitas infrastruktur berkurang
- Biaya operasional lebih rendah

**Untuk Lingkungan Air-Gapped**:
- Operasi offline sepenuhnya
- Tidak ada ketergantungan eksternal
- Broker MQTT bawaan
- Kontrol penuh atas data

**Untuk Pengembangan/Pengujian**:
- Setup lokal cepat tanpa layanan eksternal
- Mudah di-reset dan rekonfigurasi
- Setup mode untuk konfigurasi awal
- Runner mode untuk produksi

### Aliran Pendapatan

1. **Free Self-Hosted**: Single binary dengan embedded MQTT (Lisensi MIT)
2. **Lisensi Pro** ($299 sekali): Adapter database eksternal, fitur lanjutan
3. **Lisensi Enterprise** ($1,499/tahun): Dukungan prioritas, build kustom
4. **Edge appliances**: Bundel perangkat keras + perangkat lunak yang sudah dikonfigurasi

### Strategi Harga

| Edisi | Target Pasar | Model Harga | Fitur |
|---------|--------------|---------------|----------|
| Komunitas | Hobiwan, edge computing | Gratis (MIT) | Embedded MQTT, SQLite WAL, cache in-memory, setup mode |
| Profesional | UKM, produksi | $299 sekali | Adapter DB eksternal, cache persisten, akses API |
| Enterprise | Organisasi besar | $1,499/tahun | Multiple gateway, clustering, dukungan prioritas |

---

## 3. Target Audiens

### Pengguna Utama

**1. Insinyur Edge Computing**
- Deploy IoT gateway pada perangkat edge
- Butuh footprint minimal
- Memerlukan operasi offline
- Titik nyeri: Ketergantungan eksternal meningkatkan kompleksitas

**2. Operator Lingkungan Air-Gapped**
- IoT industri dalam jaringan terisolasi
- Tidak dapat mengakses layanan eksternal
- Butuh kontrol penuh
- Titik nyeri: Mengelola multiple layanan terpisah

**3. Tim Pengembangan**
- Butuh lingkungan pengujian lokal
- Ingin setup/teardown cepat
- Memerlukan fleksibilitas konfigurasi
- Titik nyeri: Setup multiple layanan untuk pengembangan

**4. Integrator Sistem**
- Deploy solusi di lokasi pelanggan
- Ingin operasi yang disederhanakan
- Butuh rekonfigurasi mudah
- Titik nyeri: Mengelola broker MQTT eksternal per pelanggan

---

## 4. Pernyataan Masalah

### Masalah Inti

**Masalah 1: Ketergantungan External MQTT Broker**
- Harus deploy dan mengelola broker MQTT terpisah (Mosquitto, HiveMQ, EMQX)
- Kompleksitas operasional tambahan
- Titik kegagalan lebih banyak
- Pendekatan saat ini: Koneksi broker eksternal diperlukan

**Masalah 2: Tidak Ada Operasi Offline yang Sebenarnya**
- v1 masih memerlukan broker MQTT eksternal
- Tidak dapat berfungsi di jaringan air-gapped tanpa setup
- Ketergantungan eksternal mengurangi keandalan
- Pendekatan saat ini: Mode klien saja, tanpa kemampuan broker

**Masalah 3: Kebingungan Konfigurasi vs Runtime**
- Tidak ada pemisahan yang jelas antara setup dan runtime
- Perubahan konfigurasi selama runtime menyebabkan ketidakstabilan
- Tidak ada setup mode yang didedikasikan
- Pendekatan saat ini: Web admin selalu tersedia, dapat menyebabkan masalah

**Masalah 4: Kompleksitas Manajemen Cache/State**
- Redis diperlukan untuk caching tingkat produksi
- Ketergantungan eksternal untuk data sementara
- Sulit untuk deploy di lingkungan dengan sumber daya terbatas
- Pendekatan saat ini: Cache in-memory saja atau Redis eksternal

**Masalah 5: Koordinasi Multi-Proses**
- Web server, API server, klien MQTT berjalan dalam proses yang sama
- Potensi kontensi sumber daya
- Kompleksitas graceful shutdown
- Pendekatan saat ini: Single proses dengan goroutines, tidak teroptimasi dengan baik

### Dampak

- **Kompleksitas Deploy**: Multiple layanan untuk deploy vs single binary
- **Biaya Infrastruktur**: Server/broker tambahan yang dibutuhkan
- **Titik Kegagalan**: Ketergantungan eksternal meningkatkan skenario kegagalan
- **Beban Operasional**: Mengelola multiple komponen
- **Keterbatasan Edge**: Tidak dapat deploy di lingkungan yang benar-benar terbatas sumber dayanya

---

## 5. Solusi yang Diusulkan

### Visi: MQTT Gateway yang Sepenuhnya Mandiri

**Single executable binary** yang mencakup:
- **Embedded MQTT Broker** (menggunakan library mochi-mqtt)
- **Operasi Dual-Mode**: Setup mode dan Runner mode
- **SQLite dalam Mode WAL**: Dioptimalkan untuk concurrent reads/writes
- **In-Memory Cache**: Cache bawaan dengan adapter eksternal opsional
- **Konfigurasi YAML**: Pemisahan jelas antara konfigurasi setup dan runtime
- **Koordinasi Proses**: Arsitektur multi-server single-proses yang dioptimalkan

### Inovasi Utama

**1. Embedded MQTT Broker**
- Menggunakan library `mochi-mqtt` Go
- Dukungan broker MQTT 3.1.1 dan 5.0 penuh
- Tidak ada ketergantungan broker eksternal
- Masih dapat terhubung ke broker eksternal jika diperlukan

**2. Operasi Dual-Mode**
```
[Setup Mode]  ←→  Konfigurasi file YAML, database, pengaturan MQTT
     ↓
[Runner Mode] ←→  Operasi produksi dengan embedded broker
```

**3. SQLite dalam Mode WAL**
- Write-Ahead Logging untuk konkurensi lebih baik
- Performa lebih baik daripada SQLite standar
- Cukup untuk beban kerja sedang
- Upgrade opsional ke database eksternal

**4. In-Memory Cache dengan Adapter**
- Cache bawaan menggunakan `github.com/patrickmn/go-cache`
- Adapter cache yang dapat ditukar (Redis, Memcached tersedia)
- Ekspirasi berbasis TTL
- Ukuran cache yang dapat dikonfigurasi

**5. Konfigurasi config.yaml**
```yaml
mode: runner  # setup | runner
model: embedded  # embedded | external

database:
  type: sqlite  # sqlite | postgres | mongodb | influxdb
  connection:
    path: ./data/gateway.db
    wal_mode: true
    max_open_conns: 25
    cache:
      type: memory  # memory | redis | memcached
      ttl: 3600
      max_size: 10000

mqtt:
  type: embedded  # embedded | external
  embedded:
    listen_address: :1883
    max_connections: 1000
    max_message_size: 256KB
  external:
    url: tcp://localhost:1883
    client_id: um-gateway
    username: ""
    password: ""

paths:
  data: ./data
  logs: ./logs
  pid: ./um-gateway.pid

retention:
  enabled: true
  default_ttl: 2592000  # 30 hari
  cleanup_interval: 3600

web:
  address: :8080
  admin_api:
    enabled: true
  public_api:
    enabled: true

logging:
  level: info
  format: json
```

### Ikhtisar Arsitektur

```
┌─────────────────────────────────────────────────────────────┐
│              Single Executable Binary (um-gateway)          │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌───────────────────────────────────────────────────────┐  │
│  │               Main Process (Go)                        │  │
│  │                                                        │  │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐  │  │
│  │  │ Embedded  │  │ HTTP/Web  │  │ In-Memory │  │  │
│  │  │ MQTT      │  │ Server     │  │ Cache      │  │  │
│  │  │ Broker     │  │            │  │            │  │  │
│  │  │ (mochi-mqtt)│  │            │  │            │  │  │
│  │  └────────────┘  └────────────┘  └────────────┘  │  │
│  │                                                        │  │
│  │  ┌────────────┐  ┌────────────┐                    │  │
│  │  │ YAML       │  │ SQLite     │                    │  │
│  │  │ Config     │  │ (WAL Mode) │                    │  │
│  │  │ Loader     │  │            │                    │  │
│  │  └────────────┘  └────────────┘                    │  │
│  │                                                        │  │
│  │  ┌────────────┐  ┌────────────┐                    │  │
│  │  │ Setup      │  │ Mode       │                    │  │
│  │  │ Mode       │  │ Switcher   │                    │  │
│  │  │ Manager    │  │            │                    │  │
│  │  └────────────┘  └────────────┘                    │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                               │
│  ┌───────────────────────────────────────────────────────┐  │
│  │              Storage & Cache Layer                   │  │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────────┐      │  │
│  │  │ SQLite   │  │ Memory   │  │ External     │      │  │
│  │  │ (WAL)    │  │ Cache    │  │ Adapters    │      │  │
│  │  │          │  │ (go-cache)│  │ (Optional)  │      │  │
│  │  └──────────┘  └──────────┘  └──────────────┘      │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                               │
│  ┌───────────────────────────────────────────────────────┐  │
│  │         Web Admin Dashboard (Embedded SPA)            │  │
│  │  - Setup Mode UI (first run atau explicit)            │  │
│  │  - Manajemen konfigurasi (YAML editing)               │  │
│  │  - Dashboard monitoring                               │  │
│  │  - Viewer log                                         │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

---

## 6. Tujuan Proyek

### Tujuan Utama (Must Have)

**O1: Implementasi Embedded MQTT Broker**
- Integrasi library `mochi-mqtt` untuk fungsionalitas broker penuh
- Dukungan protokol MQTT 3.1.1 dan 5.0
- Alamat listener dan batas yang dapat dikonfigurasi
- Koneksi broker eksternal opsional (bridging)

**O2: Implementasi Operasi Dual-Mode**
- Setup mode: Konfigurasi awal dan migrasi
- Runner mode: Operasi produksi
- Perpindahan mode melalui config.yaml atau flag command-line
- Perlindungan terhadap perpindahan mode dalam produksi

**O3: Optimasi Mode WAL SQLite**
- Mengaktifkan Write-Ahead Logging untuk konkurensi lebih baik
- Konfigurasi connection pooling
- Optimalisasi untuk kasus penggunaan embedded
- Adapter database eksternal opsional

**O4: In-Memory Cache dengan Adapter yang Dapat Ditukar**
- Cache bawaan menggunakan `go-cache`
- Ekspirasi berbasis TTL
- Adapter cache yang dapat ditukar (Redis, Memcached)
- Ukuran dan kebijakan cache yang dapat dikonfigurasi

**O5: Koordinasi Proses**
- Optimalisasi penggunaan goroutine untuk server web, API, MQTT
- Penanganan graceful shutdown
- Manajemen dan batas sumber daya
- Pemeriksaan kesehatan untuk semua komponen

### Tujuan Sekunder (Should Have)

**O6: Manajemen Konfigurasi YAML**
- File config.yaml yang jelas dan mudah dibaca
- Dukungan hot-reload dalam setup mode
- Validasi sebelum menerapkan perubahan
- Helper migrasi konfigurasi

**O7: Web UI Setup Mode**
- Editor konfigurasi dengan validasi YAML
- Wizard migrasi database
- Konfirmasi perpindahan mode
- Backup/pemulihan konfigurasi

**O8: Adapter Database Eksternal**
- Adapter PostgreSQL
- Adapter MongoDB
- Adapter InfluxDB
- Migrasi satu-klik dari SQLite

**O9: Adapter Cache Eksternal**
- Adapter Redis
- Adapter Memcached
- Strategi pemanasan cache
- Statistik dan monitoring cache

### Tujuan Tersier (Nice to Have)

**O10: Mode High Availability**
- Clustering multi-gateway
- Replikasi data
- Mekanisme failover

**O11: Optimasi Performa**
- Connection pooling
- Optimasi batch
- Profiling memori dan tuning

---

## 7. Ruang Lingkup & Batasan Proyek

### Dalam Ruang Lingkup (Apa yang Akan Kami Bangun)

#### Fungsionalitas Inti

1. **Embedded MQTT Broker**
   - Implementasi broker penuh menggunakan mochi-mqtt
   - Dukungan MQTT 3.1.1 dan 5.0
   - Listener yang dapat dikonfigurasi
   - Autentikasi klien
   - Persistensi pesan (opsional)

2. **Sistem Dual-Mode**
   - Setup mode: First-run atau aktivasi eksplisit
   - Runner mode: Operasi produksi
   - Deteksi dan validasi mode
   - Perpindahan mode yang aman

3. **Database Mode WAL SQLite**
   - Write-Ahead Logging diaktifkan
   - Connection pooling
   - Migrasi otomatis
   - Utilitas backup/pemulihan

4. **In-Memory Cache**
   - Implementasi go-cache
   - Ekspirasi TTL
   - Eviksi berbasis ukuran
   - API statistik

5. **Konfigurasi YAML**
   - Struktur config.yaml
   - Validasi saat load
   - Hot-reload dalam setup mode
   - Versioning konfigurasi

6. **Manajemen Proses**
   - Single proses dengan goroutines yang dioptimalkan
   - Graceful shutdown
   - Manajemen file PID
   - Penanganan sinyal (SIGTERM, SIGINT)

#### Fitur Web UI

**UI Setup Mode**:
- Editor konfigurasi YAML
- Wizard migrasi database
- Konfirmasi perpindahan mode
- Pemeriksaan status sistem

**UI Runner Mode**:
- Dashboard monitoring
- Viewer log
- Visualisasi metrik
- Tampilan konfigurasi read-only

### Di Luar Ruang Lingkup (Apa yang Tidak Akan Kami Bangun - Awalnya)

1. **Perangkat Keras/Firmware**
   - Firmware perangkat
   - Update OTA
   - Provisioning perangkat keras

2. **Fitur MQTT Lanjutan** (Fase 2)
   - Clustered MQTT brokers
   - Konfigurasi bridge
   - Keamanan lanjutan (setup SSL/TLS)

3. **High Availability** (Fase 3)
   - Clustering multi-gateway
   - Replikasi data
   - Failover otomatis

### Batasan & Kendala

#### Kendala Teknis
- **Bahasa**: Go 1.24+
- **Library MQTT**: mochi-mqtt (embedded) atau Eclipse Paho (external)
- **Perangkat Keras Minimal**: 1 core CPU, 1GB RAM (mode WAL membutuhkan lebih banyak)
- **Jaringan**: konektivitas opsional untuk broker/database eksternal

#### Kendala Waktu
- **Fase 1 (MVP)**: 10 bulan
- **Fase 2 (Adapter Eksternal)**: +3 bulan
- **Fase 3 (Fitur HA)**: +6 bulan

#### Kendala Anggaran
- **Tim Pengembangan**: 2-3 pengembang full-time
- **Infrastruktur**: $1,500/bulan untuk pengembangan/pengujian
- **Layanan Pihak Ketiga**: Tier gratis awalnya

#### Kendala Sumber Daya
- **Tim Pengembangan**: 1-2 Full-stack developers (Go + SvelteKit), 1 engineer DevOps
- **Ahli Domain**: Spesialis MQTT, ahli cache
- **Dukungan**: Penulis teknis part-time

---

## 8. Teknologi Stack

### Teknologi Inti

#### Backend (Go)
| Komponen | Teknologi | Rasional |
|-----------|-----------|-----------|
| Bahasa | Go 1.24+ | Performa, konkurensi, single binary |
| MQTT Broker | mochi-mqtt | Embedded Go MQTT broker |
| MQTT Klien (eksternal) | Eclipse Paho | Untuk koneksi broker eksternal |
| HTTP Router | Chi v5 | Ringan, idiomatic Go |
| YAML | gopkg.in/yaml.v3 | Parsing dan serialisasi YAML |
| SQLite | mattn/go-sqlite3 | Dukungan mode WAL |
| Cache | patrickmn/go-cache | Caching in-memory dengan TTL |
| WebSocket | gorilla/websocket | Update UI real-time |

#### Konfigurasi Embedded MQTT Broker

```go
import mqtt "github.com/mochi-mqtt/mqtt/v2"

// Setup embedded broker
func NewEmbeddedBroker(config MQTTConfig) (*mqtt.Server, error) {
    server := mqtt.NewServer(nil)

    // Konfigurasi TCP listener
    tcp := mqtt.NewTCPListener(config.ListenAddress, nil)
    server.AddListener(tcp)

    return server, nil
}
```

#### Implementasi Cache

```go
import "github.com/patrickmn/go-cache"

// In-memory cache
type Cache struct {
    store *cache.Cache
    ttl   time.Duration
}

func NewCache(defaultTTL time.Duration, maxSize int) *Cache {
    return &Cache{
        store: cache.New(defaultTTL, 10*time.Minute),
        ttl:   defaultTTL,
    }
}

// Interface adapter cache yang dapat ditukar
type CacheAdapter interface {
    Get(key string) (interface{}, error)
    Set(key string, value interface{}, ttl time.Duration) error
    Delete(key string) error
    Clear() error
}
```

#### Mode WAL SQLite

```go
import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

// Aktifkan mode WAL
func enableWALMode(db *sql.DB) error {
    _, err := db.Exec("PRAGMA journal_mode=WAL")
    if err != nil {
        return err
    }

    _, err = db.Exec("PRAGMA synchronous=NORMAL")
    if err != nil {
        return err
    }

    _, err = db.Exec("PRAGMA cache_size=-10000") // Cache 10MB
    return err
}
```

### Struktur Konfigurasi

#### Skema config.yaml

```yaml
# Mode Operasi
mode: runner  # setup | runner

# Tipe Model
model: embedded  # embedded | external

# Konfigurasi Database
database:
  # Tipe: sqlite, postgres, mongodb, influxdb
  type: sqlite
  connection:
    # SQLite spesifik
    path: ./data/gateway.db
    wal_mode: true
    max_open_conns: 25
    max_idle_conns: 5
    conn_max_lifetime: 300s

    # Konfigurasi cache
    cache:
      # Tipe: memory, redis, memcached
      type: memory
      ttl: 3600
      max_size: 10000

      # Redis (jika tipe adalah redis)
      redis:
        network: tcp
        address: localhost:6379
        password: ""
        db: 0

      # Memcached (jika tipe adalah memcached)
      memcached:
        address: localhost:11211
        timeout: 1000

# Konfigurasi MQTT
mqtt:
  # Tipe: embedded, external
  type: embedded

  # Pengaturan broker embedded
  embedded:
    listen_address: :1883
    max_connections: 1000
    max_message_size: 256KB
    keepalive: 60s

  # Broker eksternal (jika tipe adalah external)
  external:
    url: tcp://localhost:1883
    client_id: um-gateway
    username: ""
    password: ""
    clean_session: true
    auto_reconnect: true
    keepalive: 60s

# Konfigurasi Path
paths:
  data: ./data       # Direktori data
  logs: ./logs       # Direktori log
  pid: ./um-gateway.pid  # Lokasi file PID

# Retensi Data
retention:
  enabled: true
  default_ttl: 2592000  # 30 hari dalam detik
  cleanup_interval: 3600   # Interval cleanup dalam detik

# Konfigurasi Web Server
web:
  address: :8080
  admin_api:
    enabled: true
    authentication:
      type: basic  # basic, jwt, none
  public_api:
    enabled: true
    rate_limit:
      enabled: true
      requests_per_minute: 1000

# Konfigurasi Logging
logging:
  level: info  # debug, info, warn, error
  format: json  # json, text
  output: stdout  # stdout, file
  file:
    path: ./logs/gateway.log
    max_size: 100MB
    max_backups: 10
    max_age: 30
```

---

## 9. Fitur

### 9.1 Fitur Inti

#### F1: Operasi Dual-Mode

**Setup Mode**:
- Diaktifkan pada first run atau flag eksplisit (`--mode=setup`)
- Web UI memungkinkan perubahan konfigurasi
- Editing konfigurasi YAML
- Migrasi database
- Perpindahan mode ke runner mode

**Runner Mode**:
- Operasi produksi normal
- Konfigurasi bersifat read-only
- Embedded MQTT broker aktif
- Tidak dapat beralih ke setup mode tanpa menghentikan proses

**Perpindahan Mode**:

```go
// Deteksi mode
func DetermineMode(config *Config) string {
    // Periksa flag command line
    if *setupFlag {
        return "setup"
    }

    // Periksa file konfigurasi
    if config.Mode == "" {
        // First run - default ke setup
        if isFirstRun() {
            return "setup"
        }
        return "runner"
    }

    return config.Mode
}

// Penegakan mode
func RunInMode(config *Config) error {
    mode := DetermineMode(config)

    switch mode {
    case "setup":
        return runSetupMode(config)
    case "runner":
        return runRunnerMode(config)
    default:
        return fmt.Errorf("unknown mode: %s", mode)
    }
}
```

---

#### F2: Embedded MQTT Broker

**Implementasi menggunakan mochi-mqtt**:

```go
package broker

import (
    mqtt "github.com/mochi-mqtt/mqtt/v2"
    "github.com/mochi-mqtt/mqtt/hooks"
)

type EmbeddedBroker struct {
    server   *mqtt.Server
    config   *MQTTConfig
    hooks    *hooks.Hooks
}

func NewEmbeddedBroker(config *MQTTConfig) (*EmbeddedBroker, error) {
    // Buat server
    server := mqtt.NewServer(nil)

    // Setup hooks
    h := hooks.NewHooks()
    h.OnConnect = onClientConnect
    h.OnMessage = onMessageReceived
    h.OnDisconnect = onClientDisconnect

    server.AddHook(h)

    // Konfigurasi TCP listener
    tcp := mqtt.NewTCPListener(config.ListenAddress, nil)
    server.AddListener(tcp)

    // Konfigurasi WebSocket listener (opsional)
    ws := mqtt.NewWebSocketListener(config.WSAddress, nil)
    server.AddListener(ws)

    return &EmbeddedBroker{
        server: server,
        config: config,
        hooks:  h,
    }, nil
}

func (b *EmbeddedBroker) Start() error {
    go b.server.Serve()
    return nil
}

func (b *EmbeddedBroker) Stop() {
    b.server.Close()
}
```

**Penanganan Koneksi Klien**:

```go
func onClientConnect(cl mqtt.Client, pk packets.Packet) {
    // Logika autentikasi
    clientID := cl.ClientInfo().ClientID
    username := cl.ClientInfo().Username
    password := cl.ClientInfo().Password

    // Validasi kredensial
    if !authenticateClient(clientID, username, password) {
        cl.Disconnect(0x80) // Tidak otorisasi
        return
    }

    // Log koneksi
    log.Info("Klien terhubung", "client_id", clientID)
}

func onMessageReceived(cl mqtt.Client, pk packets.Packet) {
    topic := pk.TopicName
    payload := pk.Payload

    // Proses pesan
    handleMessage(topic, payload)
}

func onClientDisconnect(cl mqtt.Client, err error) {
    clientID := cl.ClientInfo().ClientID
    log.Info("Klien terputus", "client_id", clientID, "error", err)
}
```

---

#### F3: Mode WAL SQLite

**Optimasi untuk Akses Konkuren**:

```go
package database

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

type SQLiteDatabase struct {
    db *sql.DB
}

func NewSQLiteDatabase(config SQLiteConfig) (*SQLiteDatabase, error) {
    // Buka database
    db, err := sql.Open("sqlite3", config.Path)
    if err != nil {
        return nil, err
    }

    // Aktifkan mode WAL untuk konkurensi lebih baik
    if err := enableWALMode(db); err != nil {
        return nil, err
    }

    // Konfigurasi connection pool
    db.SetMaxOpenConns(config.MaxOpenConns)
    db.SetMaxIdleConns(config.MaxIdleConns)
    db.SetConnMaxLifetime(config.ConnMaxLifetime)

    return &SQLiteDatabase{db: db}, nil
}

func enableWALMode(db *sql.DB) error {
    settings := []string{
        "PRAGMA journal_mode=WAL",
        "PRAGMA synchronous=NORMAL",
        "PRAGMA cache_size=-10000",        // -10000 berarti gunakan maksimal yang tersedia
        "PRAGMA temp_store=memory",
        "PRAGMA mmap_size=268435456",      // 256MB
    }

    for _, setting := range settings {
        if _, err := db.Exec(setting); err != nil {
            return fmt.Errorf("gagal menjalankan %s: %w", setting, err)
        }
    }

    return nil
}
```

**Manfaat Performa**:
- Konkurensi lebih baik (multiple readers + single writer)
- I/O disk berkurang
- Waktu commit lebih cepat
- Lebih baik untuk kasus penggunaan embedded/edge

---

#### F4: In-Memory Cache dengan Adapter

**Implementasi Cache Dasar**:

```go
package cache

import (
    "github.com/patrickmn/go-cache"
    "time"
)

type MemoryCache struct {
    store *cache.Cache
    ttl   time.Duration
}

func NewMemoryCache(defaultTTL time.Duration, maxSize int) *MemoryCache {
    return &MemoryCache{
        store: cache.New(defaultTTL, 10*time.Minute),
        ttl:   defaultTTL,
    }
}

func (c *MemoryCache) Get(key string) (interface{}, error) {
    val, found := c.store.Get(key)
    if !found {
        return nil, ErrKeyNotFound
    }
    return val, nil
}

func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) error {
    if ttl == 0 {
        ttl = c.ttl
    }
    c.store.Set(key, value, ttl)
    return nil
}

func (c *MemoryCache) Delete(key string) error {
    c.store.Delete(key)
    return nil
}

func (c *MemoryCache) Clear() error {
    c.store.Flush()
    return nil
}

// Statistik
func (c *MemoryCache) Stats() CacheStats {
    stats := c.store.ItemCount()
    return CacheStats{
        ItemCount: stats,
        Size:       c.estimateSize(),
    }
}
```

**Adapter Redis**:

```go
package cache

import (
    "github.com/redis/go-redis/v9"
    "context"
    "time"
)

type RedisCache struct {
    client *redis.Client
    ttl    time.Duration
}

func NewRedisCache(addr, password string, db int, ttl time.Duration) (*RedisCache, error) {
    client := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })

    if err := client.Ping(context.Background()).Err(); err != nil {
        return nil, err
    }

    return &RedisCache{
        client: client,
        ttl:    ttl,
    }, nil
}

func (c *RedisCache) Get(key string) (interface{}, error) {
    val, err := c.client.Get(context.Background(), key).Result()
    if err == redis.Nil {
        return nil, ErrKeyNotFound
    }
    return val, err
}

func (c *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
    if ttl == 0 {
        ttl = c.ttl
    }
    return c.client.Set(context.Background(), key, value, ttl).Err()
}

func (c *RedisCache) Delete(key string) error {
    return c.client.Del(context.Background(), key).Err()
}

func (c *RedisCache) Clear() error {
    // Peringatan: Ini akan menghapus semua kunci dalam database!
    return c.client.FlushDB(context.Background()).Err()
}
```

---

#### F5: Koordinasi Proses

**Single Proses dengan Goroutines yang Dioptimalkan**:

```go
package main

import (
    "sync"
    "context"
    "os"
    "os/signal"
    "syscall"
    "time"
)

type GatewayServer struct {
    config      *Config
    mqttBroker  *broker.EmbeddedBroker
    webServer  *WebServer
    apiServer  *APIServer
    cache       cache.Cache
    database   *database.Database
    wg          sync.WaitGroup
    shutdownCtx context.Context
    shutdown    func()
}

func NewGatewayServer(configPath string) (*GatewayServer, error) {
    // Load konfigurasi YAML
    config, err := LoadConfig(configPath)
    if err != nil {
        return nil, err
    }

    // Inisialisasi cache
    cache, err := initializeCache(config.Database.Cache)
    if err != nil {
        return nil, err
    }

    // Inisialisasi database
    db, err := initializeDatabase(config.Database)
    if err != nil {
        return nil, err
    }

    // Inisialisasi embedded MQTT broker
    broker, err := initializeMQTTBroker(config.MQTT, db, cache)
    if err != nil {
        return nil, err
    }

    // Buat shutdown context
    shutdownCtx, shutdown := context.WithCancel(context.Background())

    return &GatewayServer{
        config:     config,
        mqttBroker: broker,
        cache:      cache,
        database:   db,
        shutdownCtx: shutdownCtx,
        shutdown:   shutdown,
    }, nil
}

func (s *GatewayServer) Start() error {
    // Mulai embedded MQTT broker
    s.wg.Add(1)
    go func() {
        defer s.wg.Done()
        if err := s.mqttBroker.Start(); err != nil {
            log.Error("Error broker MQTT", "error", err)
        }
    }()

    // Mulai web server (dashboard admin)
    s.wg.Add(1)
    go func() {
        defer s.wg.Done()
        s.webServer.Start(s.shutdownCtx)
    }()

    // Mulai API server
    s.wg.Add(1)
    go func() {
        defer s.wg.Done()
        s.apiServer.Start(s.shutdownCtx)
    }()

    // Mulai background jobs
    s.wg.Add(1)
    go func() {
        defer s.wg.Done()
        s.runBackgroundJobs(s.shutdownCtx)
    }()

    return nil
}

func (s *GatewayServer) runBackgroundJobs(ctx context.Context) {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            log.Info("Background jobs shutting down")
            return
        case <-ticker.C:
            // Lakukan tugas cleanup
            s.performCleanup()

            // Dapatkan statistik cache
            stats := s.cache.Stats()
            log.Debug("Statistik cache", "stats", stats)
        }
    }
}

func (s *GatewayServer) Stop() {
    log.Info("Mematikan gateway...")

    // Sinyalkan shutdown
    s.shutdown()

    // Hentikan embedded MQTT broker
    s.mqttBroker.Stop()

    // Tunggu semua goroutine (dengan timeout)
    done := make(chan struct{})
    go func() {
        s.wg.Wait()
        close(done)
    }()

    select {
    case <-done:
        log.Info("Semua layanan berhenti dengan grace")
    case <-time.After(30 * time.Second):
        log.Warn("Shutdown timeout, memaksa keluar")
    }

    log.Info("Gateway berhenti")
}

func (s *GatewayServer) performCleanup() {
    // Hapus data kadaluarsa
    s.database.CleanupExpiredData(s.config.Retention.DefaultTTL)

    // Flush cache jika needed
    // s.cache.Flush()
}
```

---

### 9.2 Fitur Lanjutan

#### F6: Web UI Setup Mode

**Editor Konfigurasi**:

```
┌─────────────────────────────────────────────────────────────┐
│  Setup Mode - Editor Konfigurasi                             │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  Mode: ● Setup Mode  ○ Runner Mode                           │
│                                                               │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  config.yaml                                        │   │
│  │  ┌───────────────────────────────────────────────┐  │   │
│  │  │ mode: runner                                     │  │   │
│  │  │                                                │  │   │
│  │  │ model: embedded                                 │  │   │
│  │  │                                                │  │   │
│  │  │ database:                                      │  │   │
│  │  │   type: sqlite                                 │  │   │
│  │  │   connection:                                  │  │   │
│  │  │     path: ./data/gateway.db                  │  │   │
│  │  │     wal_mode: true                             │  │   │
│  │  │     cache:                                     │  │   │
│  │  │       type: memory                             │  │   │
│  │  │       ttl: 3600                               │  │   │
│  │  │                                                │  │   │
│  │  │ mqtt:                                          │  │   │
│  │  │   type: embedded                               │  │   │
│  │  │   embedded:                                    │  │   │
│  │  │     listen_address: :1883                     │  │   │
│  │  │     max_connections: 1000                      │  │   │
│  │  │                                                │  │   │
│  │  │ retention:                                     │  │   │
│  │  │   enabled: true                                │  │   │
│  │  │   default_ttl: 2592000                         │  │   │
│  │  │                                                │  │   │
│  │  └───────────────────────────────────────────────┘  │   │
│  │                                                      │   │
│  │  [Validasi]  [Simpan]  [Batal]                       │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  [Beralih ke Runner Mode] [Export Konfig] [Backup Data]      │
└─────────────────────────────────────────────────────────────┘
```

---

#### F7: Wizard Migrasi Database

**SQLite → PostgreSQL**:

```
┌─────────────────────────────────────────────────────────────┐
│  Wizard Migrasi Database                                     │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  Database Saat Ini: SQLite (mode WAL)                        │
│  Database Target: ○ PostgreSQL  ● MongoDB  ● InfluxDB     │
│                                                               │
│  Data Migrasi:                                               │
│  ☑ Data konfigurasi                                        │
│  ☑ Data time-series (mungkin butuh jam)                     │
│  ☑ Data cache                                              │
│  ☐ Akun pengguna (fitur Pro)                                │
│                                                               │
│  Koneksi PostgreSQL:                                         │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Host: [localhost]                                   │   │
│  │ Port: [5432]                                        │   │
│  │ Database: [um_gateway_prod]                         │   │
│  │ Username: [um_gateway]                              │   │
│  │ Password: [••••••••••]                              │   │
│  │                                                      │   │
│  │ [Test Koneksi] ✓ Terhubung berhasil               │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  Opsi Migrasi:                                               │
│  ☑ Hentikan broker MQTT selama migrasi                       │
│  ☑ Buat backup sebelum migrasi                              │
│  ☑ Verifikasi data setelah migrasi                           │
│                                                               │
│  Perkiraan waktu: ~2 jam untuk 1.2M record                   │
│                                                               │
│  [Mulai Migrasi]  [Batal]                                   │
└─────────────────────────────────────────────────────────────┘
```

---

#### F8: Dashboard Runner Mode

**Monitoring Produksi**:

```
┌─────────────────────────────────────────────────────────────┐
│  Runner Mode - Dashboard Monitoring                          │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  Status Sistem                                               │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐ │
│  │  ✓ MQTT  │  │  ✓ Web    │  │  ✓ API    │  │  ✓ DB     │ │
│  │  Broker  │  │  Server   │  │  Server   │  │  (WAL)   │ │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘ │
│                                                               │
│  Statistik Broker                                            │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Klien Terhubung: 47                                   │   │
│  │ Pesan/menit: 1,234                                   │   │
│  │ Uptime: 45 hari, 6 jam                                │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  Statistik Cache                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Tipe: Memory (go-cache)                               │   │
│  │ Item: 5,234                                            │   │
│  │ Hit Rate: 87.3%                                        │   │
│  │ Memori: ~45 MB                                          │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  Statistik Database                                          │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Tipe: SQLite (mode WAL)                               │   │
│  │ Ukuran: 127 MB                                         │   │
│  │ Record: 1,234,567                                       │   │
│  │ Koneksi: 5/25                                          │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  [Lihat Log]  [Lihat Metrik]  [Masuk Setup Mode]             │
└─────────────────────────────────────────────────────────────┘
```

---

## 10. Peta Jalan Pengembangan

### Fase 1: Pondasi (Bulan 1-5)

**Sprint 1-2: Setup Proyek & Arsitektur**
- [ ] Struktur monorepo (Go + SvelteKit)
- [ ] Loading konfigurasi YAML
- [ ] Logika deteksi dan perpindahan mode
- [ ] Pipeline CI/CD

**Sprint 3-4: Embedded MQTT Broker**
- [ ] Integrasi mochi-mqtt
- [ ] Autentikasi klien
- [ ] Persistensi pesan (opsional)
- [ ] Dukungan WebSocket

**Sprint 5-6: Mode WAL SQLite**
- [ ] SQLite dengan mode WAL
- [ ] Connection pooling
- [ ] Migrasi skema
- [ ] Utilitas backup/pemulihan

**Sprint 7-8: In-Memory Cache**
- [ ] Integrasi go-cache
- [ ] Ekspirasi TTL
- [ ] API statistik cache
- [ ] Manajemen memori

**Sprint 9-10: Implementasi Dual-Mode**
- [ ] Logika setup mode
- [ ] Logika runner mode
- [ ] Perpindahan mode
- [ ] Mekanisme perlindungan

**Milestone**: Alpha - Single executable dengan embedded broker

---

### Fase 2: Web UI & Fitur (Bulan 6-8)

**Sprint 11-12: UI Setup Mode**
- [ ] Editor konfigurasi YAML
- [ ] Feedback validasi
- [ ] Konfirmasi perpindahan mode
- [ ] Backup/pemulihan konfigurasi

**Sprint 13-14: Dashboard Runner Mode**
- [ ] Dashboard monitoring
- [ ] Viewer log
- [ ] Visualisasi metrik
- [ ] Indikator status

**Sprint 15-16: Koordinasi Proses**
- [ ] Manajemen goroutine yang dioptimalkan
- [ ] Graceful shutdown
- [ ] Monitoring sumber daya
- [ ] Pemeriksaan kesehatan

**Milestone**: Beta - v2 lengkap fitur

---

### Fase 3: Poles & Produksi (Bulan 9-10)

**Sprint 17-18: Pengujian & Kualitas**
- [ ] Pengujian beban (embedded broker)
- [ ] Pengujian konkurensi (mode WAL)
- [ ] Pengujian kebocoran memori
- [ ] Audit keamanan

**Sprint 19-20: Dokumentasi & Deploy**
- [ ] Panduan pengguna
- [ ] Panduan instalasi
- [ ] Dokumentasi arsitektur
- [ ] Build cross-platform

**Milestone**: v2.0 Ketersediaan Umum

---

## 11. Metrik Kesuksesan

### Metrik Teknis

| Metrik | Target | Pengukuran |
|--------|--------|-------------|
| Ukuran Binary | < 60 MB | Artefak build |
| Waktu Startup | < 3 detik | Pengujian otomatis |
| Penggunaan Memori (Idle) | < 100 MB | Monitoring sumber daya |
| Pesan/Detik (Embedded) | 5K pesan/detik | Pengujian benchmark |
| Klien MQTT | 1,000 konkuren | Pengujian beban |
| Cache Hit Rate | > 85% | Statistik cache |
| Database Write (SQLite WAL) | 1K write/detik | Pengujian benchmark |

### Metrik Bisnis

| Metrik | Target | Timeline |
|--------|--------|----------|
| Instalasi Aktif | 1,000 | 6 bulan pasca-rilis |
| GitHub Stars | 2,000 | 6 bulan pasca-rilis |
| Deploy Edge | 200 | 12 bulan pasca-rilis |
| Penjualan Lisensi Pro | 100 | 12 bulan pasca-rilis |

### Metrik Kualitas

| Metrik | Target | Pengukuran |
|--------|--------|-------------|
| Tingkat Keberhasilan Switch Mode | 100% | Analytics |
| Tingkat Penyelesaian Setup | > 95% | Analytics |
| Mean Time Between Failures | 720 jam (30 hari) | Monitoring uptime |
| Cakupan Dokumentasi | 100% fitur | Review manual |

---

## 12. Analisis Risiko

### Risiko Teknis

| Risiko | Dampak | Probabilitas | Mitigasi |
|------|--------|-------------|------------|
| Keterbatasan library mochi-mqtt | Sedang | Sedang | Fallback ke Eclipse Paho, pengujian ekstensif |
| Mode WAL SQLite tidak cukup | Sedang | Rendah | Adapter database eksternal dari awal |
| Kontensi sumber daya single-proses | Tinggi | Sedang | Manajemen goroutine yang hati-hati, batas sumber daya |
| Kehabisan memori cache | Sedang | Sedang | Batas ukuran, eviksi LRU |

### Risiko Bisnis

| Risiko | Dampak | Probabilitas | Mitigasi |
|------|--------|-------------|------------|
| Pengguna lebih suka broker eksternal | Rendah | Sedang | Pertahankan opsi broker eksternal |
| Keterbatasan embedded broker | Sedang | Tinggi | Komunikasi jelas tentang kemampuan |
| Kekhawatiran deploy single binary | Rendah | Rendah | Scan keamanan, penandatanganan kode |

### Risiko Operasional

| Risiko | Dampak | Probabilitas | Mitigasi |
|------|--------|-------------|------------|
| Proses hang saat shutdown | Sedang | Rendah | Mekanisme timeout, opsi force kill |
| Perpindahan mode dalam produksi | Tinggi | Rendah | Konfirmasi eksplisit, pemeriksaan keselamatan |
| Kehilangan data saat perpindahan mode | Tinggi | Rendah | Backup otomatis, validasi |

---

## 13. Kebutuhan Sumber Daya

### Struktur Tim

**Tim Inti (Fase 1-3)**
- 2x Full-Stack Developers (Go + SvelteKit)
- 1x DevOps Engineer
- 1x QA Engineer (part-time)
- 1x Technical Writer (part-time)

**Estimasi Anggaran**

| Kategori | Biaya (Bulanan) | Durasi |
|----------|----------------|---------|
| Gaji Tim Pengembangan | $18,000 | 10 bulan |
| Infrastruktur (Dev/Test) | $1,500 | 10 bulan |
| Alat & Layanan | $400 | Berkelanjutan |
| Dokumentasi & Desain | $1,000 | 8 bulan |
| Kontinjensi (15%) | $2,915 | - |
| **Total Fase 1-3** | **$23,815/bulan × 10** | **$238,150** |

---

## 14. Perbandingan: v1 vs v2

| Aspek | v1 (Web Admin + SQLite) | v2 (Single Executable + Embedded MQTT) |
|--------|---------------------------|-----------------------------------------|
| **MQTT Broker** | Koneksi eksternal diperlukan | Embedded (mochi-mqtt) |
| **Deployment** | Memerlukan broker MQTT terpisah | Single binary |
| **Operasi Offline** | Memerlukan layanan eksternal | Benar-benar offline-capable |
| **Konfigurasi** | Database SQLite + Web UI | File config YAML + Web UI |
| **Cache** | Eksternal (Redis) direkomendasikan | In-memory (go-cache) |
| **Proses Setup** | Wizard onboarding | Setup mode + Runner mode |
| **Kasus Penggunaan** | Deploy IoT umum | Edge, air-gapped, embedded |
| **Kompleksitas** | Lebih rendah (layanan terpisah) | Lebih tinggi (koordinasi single proses) |
| **Penggunaan Sumber Daya** | Lebih tinggi (multiple proses) | Lebih rendah (single proses teroptimasi) |

---

## 15. Langkah Selanjutnya

### Tindakan Segera (Minggu 1-2)

1. **Riset Teknis**
   - Evaluasi kemampuan library mochi-mqtt
   - Uji performa mode WAL SQLite
   - Prototipe implementasi go-cache

2. **Desain Arsitektur**
   - Desain sistem dual-mode
   - Rencanakan koordinasi proses
   - Definisikan manajemen goroutine

3. **Formasi Tim**
   - Rekrut developer dengan keahlian Go
   - Definisikan proses kolaborasi

### Tindakan Jangka Pendek (Bulan 1)

4. **Pengembangan Prototipe**
   - Prototipe embedded MQTT broker
   - Perpindahan dual-mode
   - Loading konfigurasi YAML

5. **Lingkungan Pengembangan**
   - Setup monorepo
   - CI/CD untuk build cross-platform
   - Framework pengujian

---

## 16. Kesimpulan

Pendekatan v2 merepresentasikan **MQTT gateway yang sepenuhnya mandiri**:

### Keunggulan Utama atas v1

**Kemandirian Sejati**
- ✅ Embedded MQTT broker (tanpa ketergantungan eksternal)
- ✅ Caching in-memory (tidak perlu Redis)
- ✅ Deploy single executable
- ✅ Berfungsi di lingkungan air-gapped

**Operasi yang Disederhanakan**
- ✅ Tidak perlu mengelola broker MQTT terpisah
- ✅ Pemisahan setup/runner mode yang jelas
- ✅ Konfigurasi YAML (mudah dibaca)
- ✅ Perpindahan mode untuk pemeliharaan

**Performa Lebih Baik**
- ✅ Mode WAL SQLite untuk konkurensi lebih baik
- ✅ Cache in-memory untuk lookup lebih cepat
- ✅ Arsitektur single-proses yang dioptimasi
- ✅ Overhead komunikasi inter-proses berkurang

**Kasus Penggunaan yang Ditargetkan**
- ✅ Deploy edge computing
- ✅ Lingkungan air-gapped
- ✅ Perangkat dengan sumber daya terbatas
- ✅ Pengembangan dan pengujian

### Jalur Evolusi

**v0** → **v1** → **v2**:
- v0: Konfigurasi JSON, pengguna teknis saja
- v1: Web admin, SQLite, broker MQTT eksternal
- **v2: Single executable, embedded MQTT, config YAML, dual-mode**

Versi v2 mencapai tujuan utama: **single executable yang menyediakan fungsionalitas MQTT gateway lengkap** tanpa ketergantungan eksternal apa pun, sempurna untuk deploy edge, air-gapped, dan embedded.

Dengan eksekusi yang terfokus selama 10 bulan, kami dapat mengirimkan **v2.0 production-ready** yang melayani pasar yang belum tersentuh dari IoT gateway yang sepenuhnya mandiri.

---

**Versi Dokumen**: 2.0
**Terakhir Diperbarui**: 9 Januari 2026
**Status**: Draft - Menunggu Review
**Menggantikan**: PROJECT_BRIEF_EN_v1.md

---

## Lampiran

### A. Daftar Istilah

- **Mode WAL**: Write-Ahead Logging - Optimasi SQLite untuk konkurensi
- **Dual-Mode**: Setup mode (konfigurasi) dan Runner mode (produksi)
- **Embedded Broker**: Broker MQTT yang berjalan dalam proses yang sama
- **In-Memory Cache**: Cache yang disimpan dalam RAM menggunakan go-cache
- **Single Executable**: Semua komponen dibundel dalam satu file binary

### B. Referensi

1. Dokumentasi library mochi-mqtt
2. Dokumentasi go-cache
3. Dokumentasi mode WAL SQLite
4. PocketBase (inspirasi untuk single executable)
5. Praktik terbaik edge computing

### C. Dokumen Terkait

- `PROJECT_BRIEF_EN_v0.md` - Pendekatan JSON-based asli
- `PROJECT_BRIEF_EN_v1.md` - Pendekatan web admin
- `PROJECT_BRIEF_ID_v0.md` - Pendekatan asli (Bahasa Indonesia)
- `PROJECT_BRIEF_ID_v1.md` - Pendekatan web admin (Bahasa Indonesia)

---

**Disiapkan oleh**: Tim Pengembangan
**Disetujui oleh**: [Pending]
**Tanggal Review**: [Pending]
