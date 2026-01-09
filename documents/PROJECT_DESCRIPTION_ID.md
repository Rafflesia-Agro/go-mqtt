# Dokumentasi Proyek Go-MQTT (Bahasa Indonesia)

## Daftar Isi
1. [Ikhtisar Proyek](#ikhtisar-proyek)
2. [Latar Belakang](#latar-belakang)
3. [Tujuan & Sasaran](#tujuan--sasaran)
4. [Target Pengguna](#target-pengguna)
5. [Konteks Bisnis](#konteks-bisnis)
6. [Arsitektur Teknis](#arsitektur-teknis)
7. [Struktur Database](#struktur-database)
8. [Alur Bisnis](#alur-bisnis)
9. [Alur Pengguna](#alur-pengguna)
10. [Fitur-fitur](#fitur-fitur)
11. [Struktur File](#struktur-file)
12. [Spesifikasi Teknis](#spesifikasi-teknis)
13. [Struktur Router](#struktur-router)
14. [Controller & Logika Bisnis](#controller--logika-bisnis)
15. [Referensi Kode Penting](#referensi-kode-penting)
16. [Area Perbaikan](#area-perbaikan)
17. [Instalasi & Pengaturan](#instalasi--pengaturan)
18. [Lisensi](#lisensi)
19. [Penulis & Tim](#penulis--tim)

---

## Ikhtisar Proyek

**Go-MQTT** adalah layanan backend berkinerja tinggi berbasis Go yang menjembatani perangkat keras IoT dengan aplikasi web melalui protokol MQTT. Sistem ini berfungsi sebagai pusat kendali untuk manajemen kandang poultry (ayam), memungkinkan pemantauan dan kontrol perangkat keras (kipas, lampu, feeder, sistem air) secara realtime sekaligus mengumpulkan data sensor untuk analisis.

**Karakteristik Utama:**
- Komunikasi dua arah secara realtime melalui MQTT
- Pemrosesan data sensor berkapasitas tinggi dengan batch insertion
- Manajemen state hardware berbasis Redis untuk kontrol instan
- Autentikasi berbasis JWT dengan role-based access control
- RESTful API untuk integrasi aplikasi web/mobile
- Dioptimalkan untuk operasi 24/7 di lingkungan IoT pertanian

---

## Latar Belakang

Proyek ini dikembangkan untuk **Rafflesia Agro**, perusahaan teknologi pertanian yang berfokus pada solusi peternakan ayam cerdas di Indonesia. Sistem ini dibuat untuk menjawab kebutuhan:

1. **Manajemen Perangkat Keras Jarak Jauh**: Peternak perlu mengontrol perangkat kandang (lampu, kipas, feeder) tanpa harus hadir secara fisik
2. **Pemantauan Realtime**: Pengumpulan data sensor terus-menerus (suhu, kelembaban, dll) untuk kontrol lingkungan
3. **Skalabilitas**: Dukungan untuk banyak kandang dengan transmisi data sensor frekuensi tinggi
4. **Keandalan**: Operasi 24/7 dengan auto-reconnect dan penanganan kegagalan yang graceful

Proyek ini berkembang dari penanganan pesan MQTT sederhana menjadi platform IoT komprehensif dengan persistensi database, caching, dan kemampuan web API.

---

## Tujuan & Sasaran

### Tujuan Utama
Menyediakan lapisan middleware yang andal dan skalabel yang:
- Menerima telemetri sensor dari perangkat ESP32/IoT melalui MQTT
- Menyimpan data sensor berkapasitas tinggi secara efisien di PostgreSQL
- Mengekspos endpoint kontrol hardware melalui HTTP REST API
- Menjaga state hardware secara realtime di Redis
- Mengautentikasi dan mengotorisasi pengguna berdasarkan peran (role)

### Tujuan Teknis
- **Performa**: Menangani 10.000+ pembacaan sensor per menit dengan batch insertion
- **Keandalan**: Uptime 99.9% dengan auto-reconnect ke MQTT/Database/Redis
- **Skalabilitas**: Mendukung kandang tak terbatas dengan kemampuan scaling horizontal
- **Realtime**: Latensi sub-detik untuk update state hardware
- **Keamanan**: Autentikasi JWT dengan role-based access control (RBAC)

### Tujuan Bisnis
- Memungkinkan peternak memantau banyak kandang dari satu dashboard
- Mengurangi biaya tenaga kerja melalui kontrol hardware otomatis
- Meningkatkan kesehatan unggas melalui pemantauan lingkungan yang presisi
- Menyediakan wawasan berbasis data untuk optimasi operasional

---

## Target Pengguna

### Pengguna Utama
1. **Peternak** (Pemilik Kandang - Farmer)
   - Memiliki dan mengelola banyak kandang
   - Membutuhkan kontrol dan pemantauan jarak jauh
   - Memerlukan data historis untuk analisis

2. **ABK** (Ahli Budidaya Kandang - Operator)
   - Ditugaskan ke kandang tertentu
   - Memantau operasi harian dan pembacaan sensor
   - Mengeksekusi kontrol hardware sesuai kebutuhan

3. **Admin**
   - Administrator sistem
   - Akses penuh ke semua kandang untuk pemeliharaan dan troubleshooting
   - Mengelola peran pengguna dan izin

### Pihak Berkepentingan Lainnya
- **Manajemen**: Akses data agregat untuk wawasan bisnis
- **Tim Teknis**: Memantau kesehatan dan performa sistem
- **Integrator Hardware**: Pengembang ESP32 yang menghubungkan perangkat ke sistem

---

## Konteks Bisnis

### Domain: Peternakan Ayam Cerdas (Rafflesia Agro)

Sistem ini beroperasi di domain IoT pertanian, khususnya untuk peternakan ayam di Indonesia. Persyaratan bisnis utama:

1. **Manajemen Multi-Kandang**
   - Satu peternak dapat memiliki banyak kandang di lokasi berbeda
   - Setiap kandang dapat memiliki ABK (operator) yang ditugaskan
   - Kontrol akses hierarkis berdasarkan kepemilikan

2. **Jenis Hardware**
   - **Kontrol Iklim**: Kipas, heater, lampu (untuk regulasi suhu)
   - **Sistem Pakan**: Feeder otomatis dengan penjadwalan
   - **Sistem Air**: Dispenser air otomatis
   - **Sensor**: Suhu, kelembaban, cahaya, level amonia

3. **Mode Operasional**
   - **Mode Manual**: Kontrol on/off langsung oleh peternak/ABK
   - **Mode Otomatis**: Hardware merespons threshold sensor (misal, kipas menyala jika suhu > 30°C)
   - **Mode Terjadwal**: Operasi berbasis waktu (misal, lampu nyala jam 6 pagi, mati jam 6 sore)

4. **Kebutuhan Data**
   - Data sensor disimpan setiap ~10-15 detik per kandang
   - Perubahan state hardware harus instan (< 1 detik latency)
   - Data historis untuk analisis tren (minggu/bulan)

---

## Arsitektur Teknis

### Ikhtisar Arsitektur Sistem

```
┌─────────────┐         MQTT          ┌──────────────┐
│   ESP32     │ ◄────────────────────► │  Go-MQTT     │
│  (Hardware) │   tcp://mqtt:1883      │   Service    │
└─────────────┘                        └──────┬───────┘
                                               │
                    ┌──────────────────────────┼───────────────────┐
                    │                          │                   │
                    ▼                          ▼                   ▼
            ┌───────────┐           ┌─────────────┐       ┌─────────────┐
            │  Redis    │           │ PostgreSQL  │       │ Web/Mobile  │
            │  (State)  │           │  (Telemetry)│       │   Client    │
            └───────────┘           └─────────────┘       └─────────────┘
```

### Pola Arsitektur
**Microservice Event-Driven dengan Pemisahan Seperti CQRS**

- **Write Path**: Perintah state hardware → Redis → MQTT publish → Perangkat ESP32
- **Read Path**: Telemetri sensor → MQTT subscribe → Worker queue → Batch insert PostgreSQL
- **API Layer**: Endpoint REST untuk aplikasi klien dengan autentikasi JWT

### Komponen Utama

1. **Klien MQTT Broker** (`src/broker/mqtt.go`)
   - Subscribe ke: `$share/backend-workers/coops/+/telemetry`
   - Publish ke: `coops/{coop_id}/state`
   - Menangani kehilangan koneksi dan auto-reconnect

2. **Server HTTP** (`main.go:608-637`)
   - RESTful API menggunakan router Chi
   - Middleware autentikasi JWT
   - Middleware otorisasi berbasis peran

3. **Pool Worker Database** (`src/database/database.go:173-211`)
   - 10 worker konkuren memproses data sensor
   - Batch insertion (maksimal 1000 baris per batch)
   - Timeout 1 detik untuk flushing batch

4. **Layer Cache** (`src/cache/redis.go`)
   - Menyimpan state hardware untuk pengambilan instan
   - Melacak timestamp sensor terakhir untuk status online/offline
   - Pola: `hardware_state_{coop_id}_{hardware_name}`

---

## Struktur Database

### Skema PostgreSQL

**Tabel Utama:**

1. **`rec_sensor_coops`** (Pembacaan Sensor)
   ```sql
   - id: bigserial (PK)
   - value: double precision
   - sensor_type_id: integer (FK ke sensor_types)
   - coop_id: bigint (FK ke coops)
   - created_at: timestamp
   - updated_at: timestamp
   ```
   - **Indexing**: Index komposit pada `(coop_id, created_at)` untuk query time-series
   - **Insert Rate**: ~10-15 detik per kandang, di-insert batch via protokol COPY

2. **`sensor_types`** (Definisi Sensor)
   ```sql
   - id: serial (PK)
   - name: varchar (unique) - misal: "temperature", "humidity"
   ```
   - **Di-cache**: Dimuat ke memori saat startup (map[string]int32)

3. **`coops`** (Informasi Kandang)
   ```sql
   - id: bigint (PK)
   - farm_id: integer (FK ke farms)
   - abk_id: bigint (FK ke users, nullable)
   ```

4. **`farms`** (Informasi Peternakan)
   ```sql
   - id: integer (PK)
   - farmer_id: bigint (FK ke users)
   ```

5. **`users`** (Akun Pengguna)
   ```sql
   - id: bigint (PK)
   - current_role_id: integer (FK ke roles)
   ```

6. **`roles`** (Peran Pengguna)
   ```sql
   - id: integer (PK)
   - name: varchar - "admin", "farmer", "abk"
   ```

### Koneksi Database
- **Driver**: `pgx/v5` (driver PostgreSQL murni Go)
- **Connection Pool**: `pgxpool` dengan ukuran yang dapat dikonfigurasi
- **Timezone**: Asia/Jakarta (di-set pada level session)
- **Batch Insert**: Menggunakan protokol `COPY FROM` untuk bulk insert berkinerja tinggi

---

## Alur Bisnis

### 1. Alur Pengumpulan Data Sensor

```
Perangkat ESP32
    │
    │ Mempublikasikan telemetri ke: coops/{coop_id}/telemetry
    ▼
MQTT Broker
    │
    │ Shared subscription: $share/backend-workers/coops/+/telemetry
    ▼
Layanan Go-MQTT (Handler MQTT)
    │
    │ - Mengekstrak coop_id dari topik
    │ - Mengurai payload JSON
    │ - Menambahkan timestamp (Asia/Jakarta)
    ▼
Metode Save() → Channel Job (Buffer: 20.000)
    │
    │ 10 Worker DB Konkuren
    ▼
Batch Insert (maks 1000 baris atau timeout 1 detik)
    │
    │ Protokol PostgreSQL COPY FROM
    ▼
Tabel rec_sensor_coops
    │
    └── Update Redis key: last_sensor_stored_{coop_id}
```

### 2. Alur Kontrol Hardware

```
Klien Web/Mobile
    │
    │ POST /coops/{id}/state
    │ Headers: Authorization: Bearer <JWT>
    │ Body: { "hardware": "fan1", "state": true, "meta": {...} }
    ▼
Verifikasi JWT
    │
    │ - Memvalidasi signature token
    │ - Mengekstrak user_id dan role
    ▼
CoopAccessMiddleware
    │
    │ - Mengecek apakah user memiliki kandang (role farmer)
    │ - Mengecek apakah user adalah ABK yang ditugaskan (role abk)
    │ - Bypass untuk role admin
    ▼
Logika Handler
    │
    │ - Memvalidasi request (update parsial atau penuh)
    │ - Menggabungkan dengan state Redis yang ada
    ▼
Redis SET: hardware_state_{coop_id}_{hardware_name}
    │
    │ Publish ke MQTT: coops/{coop_id}/state
    ▼
Perangkat ESP32
    │
    │ Subscribe ke: coops/{coop_id}/state
    │ Mengeksekusi perintah hardware
    ▼
State Hardware Diperbarui
```

### 3. Alur Pengambilan State

```
Klien
    │
    │ GET /coops/{id}/state
    ▼
Auth JWT + Pengecekan Akses Coop
    │
    ▼
Redis SCAN: hardware_state_{coop_id}_*
    │
    │ Mengambil semua state hardware untuk kandang
    ▼
Mengembalikan array JSON state hardware
```

### 4. Alur Status Online/Offline

```
Klien
    │
    │ GET /coops/{id}/status
    ▼
Redis GET: last_sensor_stored_{coop_id}
    │
    │ Jika ada:
    │   - Mengurai timestamp
    │   - Menghitung waktu yang berlalu
    │   - Jika elapsed ≤ 45 detik → is_online: true
    │ Jika tidak ada:
    │   - is_online: false
    ▼
Mengembalikan status dengan last_seen dan time_elapsed
```

---

## Alur Pengguna

### Alur Pengguna Peternak (Farmer)

1. **Login ke Dashboard Web**
   - Mengautentikasi dengan username/password
   - Menerima token JWT dengan role: "farmer"

2. **Melihat Daftar Kandang**
   - Panggilan API untuk mengambil semua kandang yang dimiliki farmer
   - Menampilkan status online/offline untuk setiap kandang

3. **Memantau Kandang Tertentu**
   - Navigasi ke halaman detail kandang
   - **Data Sensor Realtime**: Mengambil pembacaan sensor terbaru dari database
   - **Status Hardware**: Mengambil state saat ini via `GET /coops/{id}/state`
   - **Status Online**: Mengecek apakah ESP32 aktif via `GET /coops/{id}/status`

4. **Mengontrol Hardware**
   - Klik tombol toggle untuk kipas/lampu/feeder
   - Klien mengirim `POST /coops/{id}/state` dengan state baru
   - State hardware diperbarui di Redis dalam hitungan milidetik
   - ESP32 menerima perintah MQTT dan mengeksekusinya

5. **Mengonfigurasi Mode Otomatis**
   - Mengatur threshold suhu (min/max)
   - Mengonfigurasi jadwal waktu (misal, lampu nyala jam 6 pagi)
   - Klien mengirim update parsial dengan data meta

6. **Melihat Data Historis**
   - Memilih rentang tanggal
   - Mengambil pembacaan sensor dari PostgreSQL
   - Menampilkan grafik untuk tren suhu, kelembaban

### Alur Pengguna ABK

1. **Login dengan Kredensial ABK**
   - Menerima token JWT dengan role: "abk"

2. **Melihat Hanya Kandang yang Ditugaskan**
   - API mengembalikan hanya kandang di mana abk_id cocok dengan ID user

3. **Memantau dan Mengontrol**
   - Kemampuan pemantauan yang sama dengan farmer
   - Dapat mengontrol hardware dalam kandang yang ditugaskan

### Alur Pengguna Admin

1. **Melewati Pengecekan Kepemilikan**
   - Mengakses kandang apa pun tanpa verifikasi kepemilikan
   - Digunakan untuk troubleshooting dan pemeliharaan

2. **Pemantauan Sistem**
   - Mengecek kesehatan layanan via `GET /health`
   - Memantau status koneksi MQTT

---

## Fitur-fitur

### Fitur Utama

#### 1. Pengumpulan Data Sensor Real-time
- **MQTT Subscription**: Mendengarkan topik `coops/+/telemetry`
- **Shared Subscription**: Menggunakan `$share/backend-workers/` untuk load balancing
- **Throughput Tinggi**: Memproses ribuan pembacaan per menit
- **Batch Insertion**: Mengakumulasi hingga 1000 pembacaan sebelum insert DB
- **Auto-timestamp**: Menambahkan timestamp Asia/Jakarta saat penerimaan

#### 2. Manajemen State Hardware
- **Penyimpanan Berbasis Redis**: Pengambilan state instan (< 10ms)
- **Update Parsial**: Memperbarui hanya state, mode, atau jadwal secara independen
- **MQTT Publishing**: Perubahan state disiarkan ke perangkat ESP32
- **Pola Wildcard**: `hardware_state_{coop_id}_{hardware_name}`

#### 3. Autentikasi & Otorisasi
- **Autentikasi JWT**: Algoritma HS256 dengan kunci secret
- **Role-Based Access Control**:
  - `admin`: Akses penuh ke semua kandang
  - `farmer`: Akses hanya ke kandang yang dimiliki
  - `abk`: Akses hanya ke kandang yang ditugaskan
- **Verifikasi Role Dinamis**: Role diambil dari database pada setiap request

#### 4. Deteksi Online/Offline
- **Pelacakan Last Seen**: Redis key `last_sensor_stored_{coop_id}`
- **Threshold 45 Detik**: Kandang dianggap online jika data diterima dalam 45 detik
- **Kalkulasi Waktu Berlalu**: Mengembalikan durasi sejak data sensor terakhir

#### 5. RESTful API
- **Spesifikasi OpenAPI**: Dokumentasi API lengkap di `openapi.yaml`
- **CORS Diaktifkan**: Request cross-origin diizinkan
- **Kode HTTP Standar**: Respons error yang tepat 400, 401, 403, 404, 500

#### 6. Graceful Shutdown
- **Penanganan Sinyal**: Menangkap SIGINT/SIGTERM
- **Pembersihan Koneksi**: Menutup koneksi DB, Redis, MQTT dengan benar
- **Worker Drain**: Memproses job yang tersisa sebelum shutdown
- **Timeout 10 Detik**: Memaksa shutdown setelah timeout

#### 7. Endpoint Testing
- **`POST /coops/{id}/test-sensor`**: Bypass MQTT untuk testing
- **Logika Database yang Sama**: Menggunakan mekanisme penyimpanan yang identik
- **Berguna Untuk**: Pengembangan, debugging, load testing

### Fitur Lanjutan

#### Batch Processing Berkinerja Tinggi
- **Worker Pool**: 10 goroutine konkuren memproses data sensor
- **Buffered Channel**: Kapasitas 20.000 job mencegah kehilangan data
- **Smart Batching**: Insert pada ukuran batch maksimum ATAU timeout (mana yang lebih dulu)
- **Protokol COPY**: Menggunakan metode bulk insert tercepat PostgreSQL

#### Auto-Reconnection
- **Klien MQTT**: Otomatis reconnect pada kehilangan koneksi
- **Retry Logic**: Built-in di paho.mqtt.golang client
- **Status Koneksi**: Dilacak dan diekspos via endpoint health

#### Metadata Hardware yang Fleksibel
- **Pemilihan Mode**: Operasi "auto" atau "manual"
- **Threshold Suhu**: Min/max triggers untuk mode auto
- **Jadwal Waktu**: Multiple jadwal ON/OFF per hari
- **Deteksi Tabrakan**: Memvalidasi overlap jadwal

---

## Struktur File

```
go-mqtt/
├── main.go                      # Entry point, server HTTP, router, handler
├── go.mod                       # Definisi module Go
├── go.sum                       # Checksum dependensi
├── .env                         # Variabel lingkungan (tidak di git)
├── .env.example                 # Template variabel lingkungan
├── .gitignore                   # Aturan ignore git
├── Dockerfile                   # Definisi container Docker
├── build.sh                     # Script build untuk produksi
├── openapi.yaml                 # Spesifikasi OpenAPI 3.0
├── README.md                    # Panduan quick start
├── PROJECT_DESCRIPTION.md       # Dokumentasi komprehensif bahasa Inggris
├── PROJECT_DESCRIPTION_ID.md    # Dokumentasi ini - bahasa Indonesia
├── esp.ino                      # Referensi firmware Arduino/ESP32
├── server                       # Binary yang dikompilasi (tidak di git)
├── server_openapi.yml           # Spesifikasi OpenAPI spesifik server
├── logs/                        # Direktori log aplikasi
│
└── src/                         # Paket source code
    ├── broker/
    │   └── mqtt.go              # Setup klien MQTT, handling pesan
    ├── cache/
    │   └── redis.go             # Konfigurasi klien Redis
    ├── config/
    │   └── config.go            # Variabel lingkungan, setup logger
    └── database/
        └── database.go          # PostgreSQL, worker pool, batch insert
```

### Deskripsi File

#### `main.go` (638 baris)
- **Tujuan**: Entry point aplikasi dan server HTTP
- **Fungsi Utama**:
  - `main()`: Menginisialisasi semua layanan, memulai server
  - `SetupRouter()`: Mengkonfigurasi router Chi dengan middleware
  - `CoopAccessMiddleware()`: Logika otorisasi
  - `StartServer()`: Server HTTP dengan graceful shutdown
  - `HTTPError()`: Respons error yang distandarisasi
  - `CORSMiddleware()`: Header CORS
  - `mergePartialUpdates()`: Menggabungkan update hardware parsial
  - `validateSchedules()`: Memvalidasi jadwal waktu

#### `src/broker/mqtt.go` (123 baris)
- **Tujuan**: Klien MQTT dan handling pesan
- **Fungsi Utama**:
  - `SetupMQTTClient()`: Membuat dan mengkonfigurasi klien MQTT
  - `MqttPublish()`: Mempublikasikan pesan ke topik MQTT
- **Struct Utama**:
  - `TelemetryMessage`: Struktur data sensor masuk
  - `MQTTConfig`: Struct konfigurasi

#### `src/database/database.go` (259 baris)
- **Tujuan**: Operasi database dan worker pool
- **Fungsi Utama**:
  - `NewPostgresStore()`: Menginisialisasi connection pool
  - `Save()`: Mengantrikan data sensor untuk batch insert
  - `DBWorker()`: Worker goroutine untuk pemrosesan batch
  - `batchInsert()`: Melakukan bulk insert COPY FROM
  - `LoadSensorTypes()`: Meng-cache mapping sensor types
  - `GetCoopOwnerAndABK()`: Mengambil kepemilikan kandang
  - `GetUserRole()`: Mengambil role user dari DB
- **Konstanta**:
  - `JOB_BUFFER_SIZE = 20000`
  - `MAX_BATCH_SIZE = 1000`
  - `BATCH_TIMEOUT = 1s`
  - `NUM_DB_WORKERS = 10`

#### `src/cache/redis.go` (49 baris)
- **Tujuan**: Setup klien Redis
- **Fungsi Utama**:
  - `SetupRedisClient()`: Membuat koneksi Redis
- **Struct Utama**:
  - `RedisConfig`: Struct konfigurasi

#### `src/config/config.go` (77 baris)
- **Tujuan**: Manajemen konfigurasi
- **Fungsi Utama**:
  - `LoadConfig()`: Memuat semua variabel lingkungan
  - `SetupLogger()`: Mengkonfigurasi logging terstruktur
  - `GetRequiredEnv()`: Memvalidasi env vars yang diperlukan
- **Struct Utama**:
  - `Config`: Menyimpan semua nilai konfigurasi

---

## Spesifikasi Teknis

### Stack Teknologi

#### Framework Backend
- **Bahasa**: Go 1.24.5
- **HTTP Router**: Chi v5.2.3 (ringan, idiomatic)
- **Klien MQTT**: Eclipse Paho MQTT v1.5.1

#### Database & Cache
- **Driver PostgreSQL**: pgx v5.7.6 (kinerja tinggi)
- **Connection Pooling**: pgxpool (manajemen koneksi built-in)
- **Klien Redis**: go-redis v9.14.0

#### Autentikasi
- **Library JWT**: go-chi/jwtauth v5.3.3
- **Algoritma**: HS256 (HMAC-SHA256)
- **Token Claims**: `sub` (ID user), `role` (nama role)

#### Utilitas
- **Environment**: godotenv v1.5.1 (memuat file .env)
- **Logging**: slog (logging terstruktur Go 1.21+)

### Spesifikasi Performa

#### Throughput
- **Data Sensor**: 10.000+ pembacaan/menit per worker
- **Request HTTP**: 1.000+ request/detik (tergantung latensi DB/Redis)
- **Pesan MQTT**: Subscribe ke topik shared untuk scaling horizontal

#### Latensi
- **Update State Hardware**: < 100ms (Redis write + MQTT publish)
- **Penyimpanan Data Sensor**: 1-2 detik (batch timeout)
- **Respon API HTTP**: 50-200ms rata-rata

#### Penggunaan Sumber Daya
- **Memori**: ~50-100 MB (bervariasi dengan ukuran worker pool)
- **CPU**: 5-15% (single core) di bawah beban normal
- **Koneksi Database**: 10-20 (ukuran pool yang dapat dikonfigurasi)

### Persyaratan Konfigurasi

#### Variabel Lingkungan
```bash
# Server
PORT=21999

# PostgreSQL
DB_HOST=103.197.190.23
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=<password>
DB_DATABASE=rafflesiaagro.migration

# Redis
REDIS_HOST=103.197.190.23
REDIS_PORT=6379
REDIS_PASSWORD=<password>
REDIS_DB=2

# MQTT
MQTT_URL=tcp://mqttserver.rafflesiaagro.com:1883
MQTT_USERNAME=<username>
MQQT_PASSWORD=<password>

# Keamanan
JWT_SECRET=<secret-key>

# Logging
LOG_LEVEL=info|debug
```

### Arsitektur Deployment

#### Deployment Produksi
- **Container**: Image Docker berbasis Alpine Linux
- **Port**: 21999 (internal), diekspos via reverse proxy
- **Replicas**: Instance tunggal (scaling horizontal dimungkinkan via shared subscription MQTT)
- **Reverse Proxy**: Nginx/Traefik untuk terminasi SSL
- **Monitoring**: Log terstruktur ke stdout untuk agregator log

#### Konfigurasi Docker
```dockerfile
FROM alpine:latest
RUN apk add --no-cache tzdata
ENV TZ=Asia/Jakarta
COPY . .
EXPOSE 21999
CMD ["/server"]
```

---

## Struktur Router

### Konfigurasi Router Chi

**File**: `main.go:317-605`

```go
func SetupRouter(client mqtt.Client, redisClient *redis.Client,
                 tokenAuth *jwtauth.JWTAuth, store *database.PostgresStore) http.Handler
```

### Definisi Route

#### Route Publik
```
GET  /health
     └── Mengembalikan: { "status": "ok", "mqtt_online": true }
     └── Middleware: Tidak ada
```

#### Route Terproteksi (JWT Diperlukan)

Semua route di bawah `/coops/{id}/*` memerlukan:
1. `jwtauth.Verifier` - Mengekstrak token dari header Authorization
2. `jwtauth.Authenticator` - Memvalidasi signature token
3. `CoopAccessMiddleware` - Mengecek perizinan pengguna

```
GET  /coops/{id}/state
     └── Mengembalikan: Semua state hardware untuk sebuah kandang
     └── Middleware: JWT Verifier, Authenticator, CoopAccess
     └── Operasi Redis: SCAN hardware_state_{coop_id}_*

POST /coops/{id}/state
     └── Menerima: HardwareState atau PartialHardwareState
     └── Mengembalikan: Konfirmasi sukses
     └── Middleware: JWT Verifier, Authenticator, CoopAccess
     └── Operasi Redis: SET hardware_state_{coop_id}_{hardware}
     └── MQTT Publish: coops/{id}/state

GET  /coops/{id}/status
     └── Mengembalikan: Status online, last seen, waktu berlalu
     └── Middleware: JWT Verifier, Authenticator, CoopAccess
     └── Operasi Redis: GET last_sensor_stored_{coop_id}

POST /coops/{id}/test-sensor
     └── Menerima: Data sensor JSON
     └── Mengembalikan: Konfirmasi tersimpan
     └── Middleware: JWT Verifier, Authenticator, CoopAccess
     └── Tujuan: Testing tanpa MQTT
```

### Rantai Middleware

```
Request → CORSMiddleware → RequestID → RealIP → Logger → Recoverer
       ↓
   JWT Verifier (ekstrak token)
       ↓
   JWT Authenticator (validasi token)
       ↓
   CoopAccessMiddleware (cek perizinan)
       ↓
   Route Handler
```

---

## Controller & Logika Bisnis

### Arsitektur Controller

**Pola**: Controller fungsional (tidak ada controller berbasis struct)
Handler didefinisikan sebagai fungsi anonim dalam `SetupRouter()`

### Controller Utama

#### 1. Controller Health Check
**Lokasi**: `main.go:323-326`

```go
r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]any{
        "status": "ok",
        "mqtt_online": client.IsConnected()
    })
})
```

**Logika**:
- Mengecek status koneksi klien MQTT
- Mengembalikan JSON dengan status layanan

#### 2. Controller Get Hardware States
**Lokasi**: `main.go:336-383`

**Logika Bisnis**:
1. Ekstrak `coop_id` dari parameter URL
2. Bangun pola Redis: `hardware_state_{coop_id}_*`
3. Gunakan `SCAN` untuk menemukan semua key yang cocok (iteratif, non-blocking)
4. Untuk setiap key:
   - `GET` value dari Redis
   - Unmarshal JSON ke `HardwareState`
   - Append ke array hasil
5. Kembalikan JSON dengan count

**Penanganan Error**:
- Error scan Redis → 500 Internal Server Error
- Error unmarshal → Log warning, lewati entry yang korup
- Hasil kosong → Kembalikan array kosong (bukan error)

#### 3. Controller Update Hardware State
**Lokasi**: `main.go:433-545`

**Logika Bisnis**:
1. Parse request body (limit 1MB)
2. Deteksi tipe update:
   - **Update Parsial**: `PartialHardwareState` (hanya field yang disediakan)
   - **Update Penuh**: `HardwareState` (penggantian lengkap)
3. Untuk update parsial:
   - Ambil state yang ada dari Redis
   - Merge dengan nilai baru menggunakan `mergePartialUpdates()`
4. Validasi jadwal (jika disediakan dalam meta)
5. Marshal ke JSON
6. `SET` ke Redis: `hardware_state_{coop_id}_{hardware}`
7. `PUBLISH` ke MQTT: `coops/{coop_id}/state`
8. Kembalikan 202 Accepted

**Penanganan Error**:
- JSON tidak valid → 400 Bad Request
- Jadwal tidak valid → 400 dengan error validasi
- Kegagalan Redis → 500 Internal Server Error
- Timeout publish MQTT → 500 Internal Server Error

#### 4. Controller Get Coop Status
**Lokasi**: `main.go:385-431`

**Logika Bisnis**:
1. Bangun Redis key: `last_sensor_stored_{coop_id}`
2. `GET` value dari Redis
3. Jika key ada:
   - Parse timestamp RFC3339
   - Hitung waktu berlalu
   - Cek jika elapsed ≤ 45 detik
   - Set `is_online: true` atau `false`
4. Jika key tidak ada:
   - Set `is_online: false`
5. Kembalikan JSON dengan detail status

**Penanganan Error**:
- Error parse timestamp → 500 Internal Server Error
- Error Redis → 500 Internal Server Error

#### 5. Controller Test Sensor
**Lokasi**: `main.go:548-601`

**Logika Bisnis**:
1. Parse `coop_id` dari URL
2. Parse body JSON (limit 1MB)
3. Buat `TelemetryMessage` dengan waktu Jakarta saat ini
4. Reload sensor types dari database
5. Panggil `store.Save()` untuk antrian batch insert
6. Kembalikan 201 Created

**Tujuan**: Endpoint pengembangan/testing yang melewati MQTT

### Logika Middleware

#### CoopAccessMiddleware
**Lokasi**: `main.go:248-315`

**Alur Otorisasi**:
1. Ekstrak `coop_id` dari URL
2. Ekstrak `user_id` dari klaim JWT (`sub`)
3. Query database untuk role user saat ini
4. Jika role adalah "admin":
   - Lewati semua pengecekan, izinkan akses
5. Jika role adalah "farmer" atau "abk":
   - Query `GetCoopOwnerAndABK(coop_id)`
   - Pengecekan Farmer: `farmer_id == user_id AND farmer_id != 0`
   - Pengecekan ABK: `abk_id == user_id AND abk_id != 0`
6. Jika tidak ada kondisi yang terpenuhi:
   - Kembalikan 403 Forbidden

**Query SQL**:
```sql
-- Get user role
SELECT r.name FROM users u
JOIN roles r ON u.current_role_id = r.id
WHERE u.id = $1

-- Get kepemilikan kandang
SELECT f.farmer_id, c.abk_id
FROM coops c
LEFT JOIN farms f ON c.farm_id = f.id
WHERE c.id = $1
```

### Fungsi Helper

#### mergePartialUpdates()
**Lokasi**: `main.go:83-124`

**Logika**:
- Mulai dengan state hardware yang ada
- Update field hanya jika disediakan dalam update parsial
- Menangani update meta bersarang (mode, temperature, schedules)
- Mengembalikan state yang sudah digabung

#### validateSchedules()
**Lokasi**: `main.go:126-160`

**Validasi**:
1. Nomor urut harus unik
2. Format waktu harus HH:MM
3. `time_off` harus setelah `time_on`
4. Tidak ada jadwal yang tumpang tindih (waktu mulai berikutnya harus setelah waktu selesai sebelumnya)
5. Mengembalikan error deskriptif untuk kegagalan validasi

---

## Referensi Kode Penting

### Bagian Kode Kritis

#### 1. Handler Pesan MQTT
**File**: `src/broker/mqtt.go:59-95`

**Mengapa Penting**: Logika inti pemrosesan telemetri

```go
if token := c.Subscribe(telemetryTopic, 1, func(client mqtt.Client, msg mqtt.Message) {
    parts := strings.Split(msg.Topic(), "/")
    if len(parts) != 3 {
        slog.Warn("Received message on unexpected topic format", ...)
        return
    }
    coopID := parts[1] // coops/[0] coopID/[1] telemetry/[2]

    var sensorData map[string]any
    if err := json.Unmarshal(msg.Payload(), &sensorData); err != nil {
        slog.Error("Failed to unmarshal telemetry", ...)
        return
    }

    telemetry := TelemetryMessage{
        CoopID:    coopID,
        Data:      sensorData,
        Timestamp: GetJakartaTime(),
    }

    if err := store.Save(telemetry, sensorTypes); err != nil {
        slog.Error("Failed to queue sensor data", ...)
    }
}); token.Wait() && token.Error() != nil {
    slog.Error("Failed to subscribe to telemetry topic", ...)
}
```

#### 2. Batch Insert dengan Protokol COPY
**File**: `src/database/database.go:214-238`

**Mengapa Penting**: Metode insert database berkinerja tertinggi

```go
func batchInsert(pool *pgxpool.Pool, readings []SensorReading) error {
    if len(readings) == 0 {
        return nil
    }

    rows := make([][]interface{}, len(readings))
    for i, r := range readings {
        rows[i] = []interface{}{r.Value, r.SensorTypeID, r.CoopID, r.Timestamp, r.Timestamp}
    }

    _, err := pool.CopyFrom(
        context.Background(),
        pgx.Identifier{"rec_sensor_coops"},
        []string{"value", "sensor_type_id", "coop_id", "created_at", "updated_at"},
        pgx.CopyFromRows(rows),
    )

    if err != nil {
        slog.Error("COPY From failed", "error", err)
        return err
    }

    slog.Info("Successfully inserted batch", "rows", len(readings))
    return nil
}
```

#### 3. Inisialisasi Worker Pool
**File**: `main.go:204-216`

**Mengapa Penting**: Arsitektur pemrosesan data berkapasitas tinggi

```go
sensorTypes, err := database.LoadSensorTypes(store.Pool)
if err != nil {
    slog.Error("Failed to load sensor types", "error", err)
    os.Exit(1)
}

var wg sync.WaitGroup
for i := 0; i < database.NUM_DB_WORKERS; i++ {
    wg.Add(1)
    go database.DBWorker(i, &wg, store.Pool, store.JobChan, sensorTypes)
}
slog.Info("🚀 Database worker pool started", "workers", database.NUM_DB_WORKERS)
```

#### 4. Penanganan Update Parsial
**File**: `main.go:440-508`

**Mengapa Penting**: Desain API yang fleksibel untuk UX yang lebih baik

```go
// Pertama coba decode sebagai update parsial
var partialUpdate PartialHardwareState
partialErr := decoder.Decode(&partialUpdate)

var finalState HardwareState
var hardwareName string

if partialErr == nil && partialUpdate.Hardware != "" {
    // Ini adalah update parsial - ambil state yang ada dan merge
    hardwareName = partialUpdate.Hardware
    redisKey := fmt.Sprintf("hardware_state_%s_%s", coopID, hardwareName)

    // Ambil state yang ada dari Redis
    existingVal, err := redisClient.Get(r.Context(), redisKey).Result()
    // ... unmarshal state yang ada

    // Merge update parsial
    finalState = mergePartialUpdates(finalState, partialUpdate)
} else {
    // Ini adalah update state penuh (backward compatibility)
    // ... decode sebagai HardwareState
}
```

#### 5. Graceful Shutdown
**File**: `main.go:228-245`

**Mengapa Penting**: Manajemen sumber daya yang bersih untuk operasi 24/7

```go
stop := make(chan os.Signal, 1)
signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
srv := StartServer(router, cfg)
<-stop

slog.Info("⏳ Shutting down services...")
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
if err := srv.Shutdown(ctx); err != nil {
    slog.Error("Server shutdown failed", "err", err)
}
mqttClient.Disconnect(250)
slog.Info("MQTT client disconnected")
close(store.JobChan)
wg.Wait()
slog.Info("All database workers have finished.")
slog.Info("✅ Server gracefully stopped.")
```

#### 6. Pengecekan Otorisasi
**File**: `main.go:277-309`

**Mengapa Penting**: Keamanan dan multi-tenancy

```go
// BARU: Ambil peran pengguna langsung dari database.
role, err := store.GetUserRole(userID)
if err != nil {
    HTTPError(w, http.StatusForbidden, err)
    return
}

// Bypass otorisasi untuk admin
if role == "admin" {
    slog.Info("Admin access granted, bypassing ownership checks.")
    next.ServeHTTP(w, r)
    return
}

// Query database untuk mendapatkan pemilik & ABK kandang
farmerID, abkID, err := store.GetCoopOwnerAndABK(coopID)
if err != nil {
    if err == pgx.ErrNoRows {
        HTTPError(w, http.StatusNotFound, errors.New("coop not found"))
        return
    }
    HTTPError(w, http.StatusInternalServerError, err)
    return
}

// Terapkan logika otorisasi untuk non-admin
isFarmOwner := (role == "farmer" && farmerID == userID && farmerID != 0)
isAssignedAbk := (role == "abk" && abk_id == userID && abk_id != 0)

if !isFarmOwner && !isAssignedAbk {
    HTTPError(w, http.StatusForbidden, fmt.Errorf("you are not authorized to access coop '%s'", coopIDStr))
    return
}
```

---

## Area Perbaikan

### Keterbatasan Saat Ini & Peningkatan Potensial

#### 1. Peningkatan Keamanan

**Masalah**: Secret JWT disimpan dalam variabel lingkungan
- **Peningkatan**: Gunakan layanan manajemen kunci (HashiCorp Vault, AWS Secrets Manager)

**Masalah**: Tidak ada rate limiting pada endpoint API
- **Peningkatan**: Implementasikan middleware rate limiting (misalnya tollbooth, slowcache)

**Masalah**: CORS mengizinkan semua origin (`*`)
- **Peningkatan**: Whitelist domain tertentu di produksi

**Masalah**: Tidak ada penandatanganan request untuk MQTT
- **Peningkatan**: Implementasikan MQTT TLS + sertifikat klien

#### 2. Optimasi Performa

**Masalah**: Sensor types dimuat hanya saat startup
- **Peningkatan**: Implementasikan refresh berkala atau invalidasi cache

**Masalah**: Konfigurasi connection pool tidak diekspos
- **Peningkatan**: Jadikan ukuran pool, koneksi max dapat dikonfigurasi via env vars

**Masalah**: Pengecekan otorisasi berurutan (2 query DB per request)
- **Peningkatan**: Cache role user di Redis dengan TTL (5 menit)

**Masalah**: Tidak ada pagination untuk pengambilan state hardware
- **Peningkatan**: Tambahkan pagination jika kandang memiliki banyak hardware

#### 3. Peningkatan Keandalan

**Masalah**: Tidak ada logika retry untuk operasi Redis yang gagal
- **Peningkatan**: Implementasikan mekanisme retry dengan exponential backoff

**Masalah**: Timeout publish MQTT (5 detik) mungkin terlalu lama
- **Peningkatan**: Gunakan publish background dengan antrian acknowledgment

**Masalah**: Tidak ada dead letter queue untuk data sensor yang gagal
- **Peningkatan**: Implementasikan DLQ untuk batch insert yang gagal

**Masalah**: Tidak ada health check untuk koneksi database
- **Peningkatan**: Tambah endpoint `/health/db` dengan ping check

#### 4. Monitoring & Observabilitas

**Masalah**: Logging terstruktur tapi tanpa logging terpusat
- **Peningkatan**: Integrasikan ELK Stack, Loki, atau CloudWatch

**Masalah**: Tidak ada pengumpulan metrik
- **Peningkatan**: Tambah metrik Prometheus (request rate, batch insert time, queue depth)

**Masalah**: Tidak ada distributed tracing
- **Peningkatan**: Tambah OpenTelemetry untuk tracing request

**Masalah**: Tidak ada alerting pada kegagalan
- **Peningkatan**: Integrasikan PagerDuty, webhook Slack pada error kritis

#### 5. Peningkatan Fitur

**Masalah**: Tidak ada riwayat state hardware
- **Peningkatan**: Simpan perubahan state di tabel audit log

**Masalah**: Tidak ada kontrol hardware massal
- **Peningkatan**: Tambah endpoint untuk update banyak hardware sekaligus

**Masalah**: Tidak ada agregasi data sensor
- **Peningkatan**: Tambah endpoint untuk min/max/avg selama periode waktu

**Masalah**: Tidak ada konfigurasi alert
- **Peningkatan**: Izinkan user mengatur threshold untuk alert (misal, temp > 35°C)

**Masalah**: Tidak ada notifikasi webhook
- **Peningkatan**: Kirim webhook pada perubahan state hardware atau pelanggaran threshold

#### 6. Kualitas Kode

**Masalah**: Fungsi handler besar di main.go
- **Peningkatan**: Ekstrak handler ke paket terpisah (`handlers/`)

**Masalah**: Kepedulian campuran dalam handler (auth, logika bisnis, akses data)
- **Peningkatan**: Implementasikan arsitektur bersih dengan layer service

**Masalah**: Konstanta hard-coded (buffer size, batch size)
- **Peningkatan**: Pindahkan ke file konfigurasi

**Masalah**: Tes integrasi terbatas
- **Peningkatan**: Tambah suite tes dengan MQTT/DB/Redis yang di-mock

#### 7. Optimasi Database

**Masalah**: Tidak ada kebijakan retensi data
- **Peningkatan**: Implementasikan partitioning per bulan, auto-drop data lama

**Masalah**: Tidak ada materialized view untuk data agregat
- **Peningkatan**: Buat ringkasan per jam/hari untuk query dashboard cepat

**Masalah**: Strategi pengindeksan tidak didokumentasikan
- **Peningkatan**: Dokumentasikan dan implementasikan indeks pada pola query umum

#### 8. DevOps & Deployment

**Masalah**: Deployment binary tunggal
- **Peningkatan**: Implementasikan strategi deployment blue-green

**Masalah**: Tidak ada mekanisme rollback
- **Peningkatan**: Penandaan image container dan script rollback

**Masalah**: Konfigurasi environment manual
- **Peningkatan**: Infrastructure as Code (Terraform/Ansible)

**Masalah**: Tidak ada strategi backup yang didokumentasikan
- **Peningkatan**: Prosedur backup otomatis untuk Redis dan PostgreSQL

---

## Instalasi & Pengaturan

### Prasyarat

- **Go**: 1.24.5 atau lebih tinggi
- **PostgreSQL**: 12+ dengan skema database
- **Redis**: 6+ (edisi apa pun)
- **MQTT Broker**: Mosquitto, HiveMQ, atau AWS IoT Core
- **Sistem Operasi**: Linux (direkomendasikan), macOS, Windows

### Pengaturan Pengembangan Lokal

#### 1. Clone Repository
```bash
git clone https://github.com/Anjasfedo/go-mqtt.git
cd go-mqtt
```

#### 2. Instal Dependensi
```bash
go mod download
```

#### 3. Konfigurasi Environment
```bash
cp .env.example .env
# Edit .env dengan konfigurasi Anda
nano .env
```

#### 4. Setup Database
```sql
-- Buat database
CREATE DATABASE "rafflesiaagro.migration";

-- Koneksi ke database
\c "rafflesiaagro.migration"

-- Buat tabel (skema harus cocok dengan ekspektasi kode)
-- Tabel yang dibutuhkan: rec_sensor_coops, sensor_types, coops, farms, users, roles

-- Insert sensor types
INSERT INTO sensor_types (name) VALUES
('temperature'), ('humidity'), ('light'), ('ammonia');

-- Insert roles
INSERT INTO roles (name) VALUES ('admin'), ('farmer'), ('abk');
```

#### 5. Jalankan Aplikasi
```bash
# Mode development
go run main.go

# Atau build dan jalankan
go build -o server
./server
```

#### 6. Verifikasi Instalasi
```bash
# Cek endpoint health
curl http://localhost:21999/health

# Respons yang diharapkan:
# {"status":"ok","mqtt_online":true}
```

### Deployment Docker

#### 1. Build Docker Image
```bash
docker build -t go-mqtt:latest .
```

#### 2. Jalankan Container
```bash
docker run -d \
  --name go-mqtt \
  --env-file .env \
  -p 21999:21999 \
  go-mqtt:latest
```

#### 3. Docker Compose (Direkomendasikan)
```yaml
version: '3.8'
services:
  go-mqtt:
    build: .
    ports:
      - "21999:21999"
    env_file:
      - .env
    depends_on:
      - postgres
      - redis
      - mqtt
    restart: unless-stopped

  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: rafflesiaagro.migration
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: your_password
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass your_redis_password
    volumes:
      - redis_data:/data

  mqtt:
    image: eclipse-mosquitto:2
    ports:
      - "1883:1883"
    volumes:
      - ./mosquitto.conf:/mosquitto/config/mosquitto.conf

volumes:
  postgres_data:
  redis_data:
```

### Deployment Produksi

#### 1. Build Binary Produksi
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o server

# Dengan optimasi
go build -ldflags="-s -w" -o server
```

#### 2. Menggunakan Build Script
```bash
chmod +x build.sh
./build.sh
```

#### 3. Deploy ke Server
```bash
# Copy ke server
scp server user@server:/path/to/deploy/

# SSH ke server
ssh user@server

# Jalankan sebagai service
sudo systemctl start go-mqtt
```

#### 4. Konfigurasi Service Systemd
```ini
[Unit]
Description=Go-MQTT Service
After=network.target

[Service]
Type=simple
User=go-mqtt
WorkingDirectory=/opt/go-mqtt
ExecStart=/opt/go-mqtt/server
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

### Panduan Variabel Lingkungan

#### Variabel Wajib (Tidak Ada Default)
- `PORT`: Port server HTTP (misal, 21999)
- `DB_HOST`: Host PostgreSQL
- `DB_PORT`: Port PostgreSQL (biasanya 5432)
- `DB_USERNAME`: User database
- `DB_PASSWORD`: Password database
- `DB_DATABASE`: Nama database
- `REDIS_HOST`: Host Redis
- `REDIS_PORT`: Port Redis (biasanya 6379)
- `REDIS_PASSWORD`: Password Redis
- `REDIS_DB`: Nomor database Redis (0-15)
- `MQTT_URL`: URL broker MQTT (misal, tcp://localhost:1883)
- `MQTT_USERNAME`: Username MQTT
- `MQTT_PASSWORD`: Password MQTT
- `JWT_SECRET`: Kunci secret untuk penandatanganan JWT

#### Variabel Opsional
- `LOG_LEVEL`: Level logging (info, debug) - Default: info

### Testing

#### 1. Health Check
```bash
curl http://localhost:21999/health
```

#### 2. Test Data Sensor (Tanpa MQTT)
```bash
curl -X POST http://localhost:21999/coops/123/test-sensor \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "temperature": 28.5,
    "humidity": 75.2,
    "light": 500
  }'
```

#### 3. Get Hardware States
```bash
curl http://localhost:21999/coops/123/state \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

#### 4. Update Hardware State
```bash
curl -X POST http://localhost:21999/coops/123/state \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "hardware": "fan1",
    "state": true,
    "meta": {
      "mode": "auto",
      "temperature_min": 25.0,
      "temperature_max": 30.0
    }
  }'
```

---

## Lisensi

Proyek ini tampaknya adalah proyek pribadi/proprietary untuk Rafflesia Agro. Tidak ada file lisensi eksplisit yang ada di repository.

**Untuk Distribusi Open Source** (jika berlaku):
Pertimbangkan untuk menambahkan file `LICENSE`. Opsi umum:
- **MIT License**: Permisif, sederhana
- **Apache 2.0**: Perlindungan paten, banyak digunakan
- **GPL v3**: Copyleft, mengharuskan karya turunan open source

---

## Penulis & Tim

### Informasi Proyek
- **Perusahaan**: Rafflesia Agro
- **Domain**: IoT Pertanian / Peternakan Ayam Cerdas
- **Lokasi**: Indonesia (Timezone Asia/Jakarta)

### Tim Pengembangan
- **Pengembang Utama**: Anjasfedo (GitHub: @Anjasfedo)
- **Status Proyek**: Aktif (per commit terakhir: 10 Oktober 2024)

### Kontak & Dukungan
- **Repository**: https://github.com/Anjasfedo/go-mqtt
- **Pelacakan Issue**: GitHub Issues
- **Dokumentasi**: File ini (PROJECT_DESCRIPTION_ID.md) dan README.md

### Proyek Terkait
- **Firmware ESP32**: `esp.ino` (kode Arduino untuk perangkat hardware)
- **Spesifikasi API**: `openapi.yaml` (dokumentasi OpenAPI 3.0)
- **Klien Web/Mobile**: (Repository terpisah, tidak termasuk)

### Penghargaan
- **Library MQTT**: Eclipse Paho MQTT Go Client
- **HTTP Router**: Chi - router HTTP yang ringan dan idiomatic
- **Driver Database**: pgx - driver PostgreSQL berkinerja tinggi
- **Library JWT**: go-chi/jwtauth

---

## Sumber Daya Tambahan

### Dokumentasi
- **README.md**: Panduan quick start dan contoh API
- **openapi.yaml**: Spesifikasi OpenAPI 3.0 lengkap
- **server_openapi.yml**: Dokumentasi API spesifik server

### Integrasi ESP32
- **esp.ino**: Firmware Arduino untuk perangkat ESP32
  - Logika koneksi MQTT
  - Format publikasi data sensor
  - Handling subscription state hardware

### Logging & Monitoring
- **Direktori Logs**: `/logs` (log aplikasi)
- **Format Log**: JSON terstruktur via slog
- **Level Log**: info (default), debug (verbose)

### Checklist Deployment Produksi
- [ ] Variabel lingkungan dikonfigurasi
- [ ] Skema database dibuat
- [ ] Koneksi Redis diuji
- [ ] Broker MQTT dapat diakses
- [ ] Secret JWT dibuat (32+ karakter)
- [ ] Aturan firewall dikonfigurasi (port 21999)
- [ ] SSL/TLS dikonfigurasi (reverse proxy)
- [ ] Agregasi log dikonfigurasi
- [ ] Setup monitoring (Prometheus/DataDog)
- [ ] Strategi backup diimplementasikan
- [ ] Prosedur rollback didokumentasikan

---

## Referensi Cepat

### Perintah Penting

```bash
# Build
go build -o server

# Run
./server

# Docker
docker build -t go-mqtt .
docker run -p 21999:21999 --env-file .env go-mqtt

# Dependensi
go mod tidy
go mod download

# Testing
curl http://localhost:21999/health
```

### Angka Konfigurasi Kunci
- **Port HTTP**: 21999
- **Worker DB**: 10
- **Buffer Job**: 20.000
- **Ukuran Batch**: 1.000
- **Timeout Batch**: 1 detik
- **Threshold Online**: 45 detik

### Pola Key Redis
- **State Hardware**: `hardware_state_{coop_id}_{hardware_name}`
- **Sensor Terakhir**: `last_sensor_stored_{coop_id}`

### Topik MQTT
- **Subscribe**: `$share/backend-workers/coops/+/telemetry`
- **Publish**: `coops/{coop_id}/state`

---

**Versi Dokumen**: 1.0
**Terakhir Diperbarui**: 9 Januari 2026
**Dibuat Untuk**: Proyek Go-MQTT oleh Rafflesia Agro
