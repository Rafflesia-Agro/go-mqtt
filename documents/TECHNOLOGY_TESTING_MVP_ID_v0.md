# Rencana Pengujian Teknologi: UM-Gateway MVP v0.1

## Ringkasan Eksekutif

**Nama Proyek**: Pengujian Teknologi UM-Gateway MVP
**Tujuan**: Memverifikasi semua teknologi/library sebelum pengembangan penuh
**Timeline**: 2 minggu (sebelum pengembangan utama dimulai)
**Target**: Menyelesaikan semua pengujian proof-of-concept (PoC)

---

## Riwayat Versi

| Versi | Tanggal | Perubahan | Penulis |
|-------|---------|-----------|--------|
| **v1.0** | 2026-01-09 | Rencana pengujian awal<br>• PoC broker mochi-mqtt<br>• Pengujian mode WAL SQLite<br>• Integrasi go-cache<br>• React + Vite dalam Go<br>• Konfigurasi YAML<br>• Verifikasi semua teknologi kritis | Tim Pengembangan |

---

## 1. Tujuan

### Kenapa Pengujian Teknologi?

Sebelum memulai pengembangan MVP 3 bulan, kita perlu **memverifikasi bahwa semua teknologi yang dipilih dapat bekerja bersama**. Ini mencegah:
- ❌ Library yang tidak kompatibel
- ❌ Blokir teknis yang tidak terduga
- ❌ Masalah performa
- ❌ Redesign arsitektur di tengah pengembangan

### Tujuan Pengujian

1. **Verifikasi Kelayakan**: Membuktikan setiap teknologi dapat melakukan apa yang kita butuhkan
2. **Uji Integrasi**: Memastikan teknologi bekerja bersama
3. **Ukur Performa**: Menetapkan metrik baseline
4. **Identifikasi Masalah**: Menemukan masalah sejak awal (ketika murah untuk diperbaiki)
5. **Bangun Kepercayaan**: Tahu tech stack bekerja sebelum berkomitmen

### Kriteria Sukses

- [ ] Semua PoC berjalan berhasil
- [ ] Semua teknologi terintegrasi tanpa konflik
- [ ] Performa memenuhi persyaratan minimum
- [ ] Tidak ada blokir kritis yang ditemukan
- [ ] Dokumentasi untuk setiap PoC dibuat

---

## 2. Ikhtisar Pengujian

### Timeline Pengujian

**Minggu 1: Teknologi Backend Inti**
- Hari 1-2: broker mochi-mqtt embedded
- Hari 3-4: mode WAL SQLite
- Hari 5-7: integrasi go-cache

**Minggu 2: Frontend & Integrasi**
- Hari 1-2: setup React + Vite
- Hari 3-4: Embed React dalam binary Go
- Hari 5-6: konfigurasi YAML
- Hari 7: Pengujian integrasi & dokumentasi

### Lingkungan Pengujian

**Hardware**:
- Mesin pengembangan (OS apa pun)
- Minimum: 2 core CPU, 4GB RAM

**Software**:
- Go 1.24+
- Node.js 20+
- Git
- Browser modern (Chrome/Firefox)

**Tools**:
- IDE (VS Code, GoLand, dll)
- Postman/cURL (untuk pengujian API)
- MQTT client (MQTTX, mosquitto_pub/sub)

---

## 3. Detail Rencana Pengujian

### Uji 1: mochi-mqtt Embedded Broker

**Tujuan**: Memverifikasi mochi-mqtt dapat bekerja sebagai embedded broker

**Apa yang Diuji**:
- [ ] Membuat broker MQTT dasar
- [ ] TCP listener pada port 1883
- [ ] Autentikasi klien (username/password)
- [ ] Fungsionalitas Publish/Subscribe
- [ ] Persistensi pesan (opsional)
- [ ] Koneksi konkuren (10+ klien)
- [ ] Graceful shutdown

**Kriteria Sukses**:
- ✅ Broker mulai dan listen pada port 1883
- ✅ Klien MQTT dapat terhubung dengan username/password
- ✅ Klien dapat publish ke topic
- ✅ Klien dapat subscribe ke topic
- ✅ Subscriber menerima pesan yang dipublish
- ✅ 10+ klien konkuren bekerja tanpa masalah
- ✅ Broker shutdown dengan graceful

**Struktur Kode yang Diharapkan**:

```go
package main

import (
    "log"
    mqtt "github.com/mochi-mqtt/mqtt/v2"
    "github.com/mochi-mqtt/mqtt/hooks/auth"
    "github.com/mochi-mqtt/mqtt/listeners"
)

func main() {
    // Buat server MQTT
    server := mqtt.NewServer(nil)

    // Tambahkan hook autentikasi
    authHook := auth.NewHook("username", "password")
    server.AddHook(authHook)

    // Buat TCP listener
    tcp := listeners.NewTCPListener("localhost:1883", nil)
    server.AddListener(tcp)

    // Mulai server
    log.Println("Memulai MQTT broker pada :1883")
    go server.Serve()

    // Tetap berjalan
    select {}
}
```

**Prosedur Pengujian**:
1. Buat direktori baru: `tests/mochi-mqtt-poc`
2. Inisialisasi Go module
3. Install mochi-mqtt
4. Salin kode yang diharapkan
5. Jalankan: `go run main.go`
6. Gunakan MQTTX untuk terhubung ke `localhost:1883`
7. Publish pesan ke topic `test/topic`
8. Subscribe ke `test/topic`
9. Verifikasi pesan diterima
10. Uji dengan 10+ koneksi konkuren
11. Stop server (Ctrl+C) dan verifikasi graceful shutdown

**Perkiraan Waktu**: 4-6 jam

**Deliverables**:
- Kode PoC yang berfungsi
- Dokumentasi hasil pengujian
- Catatan performa (penggunaan memori, CPU)
- Masalah atau keterbatasan yang diketahui

---

### Uji 2: Mode WAL SQLite

**Tujuan**: Memverifikasi SQLite dengan mode WAL bekerja untuk akses konkuren

**Apa yang Diuji**:
- [ ] Membuat database SQLite
- [ ] Mengaktifkan mode WAL
- [ ] Read konkuren
- [ ] Write konkuren
- [ ] Connection pooling
- [ ] Migrasi skema
- [ ] Backup/pemulihan

**Kriteria Sukses**:
- ✅ Database berhasil dibuat
- ✅ Mode WAL berhasil diaktifkan
- ✅ Multiple read konkuren bekerja
- ✅ Write konkuren tidak memblokir read
- ✅ Connection pool bekerja
- ✅ Migrasi berjalan tanpa error
- ✅ Backup dapat dibuat

**Struktur Kode yang Diharapkan**:

```go
package main

import (
    "database/sql"
    "log"
    "sync"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    // Buka database
    db, err := sql.Open("sqlite3", "./test.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Aktifkan mode WAL
    _, err = db.Exec("PRAGMA journal_mode=WAL")
    if err != nil {
        log.Fatal(err)
    }

    // Buat tabel
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS messages (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            topic TEXT,
            payload TEXT,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    `)
    if err != nil {
        log.Fatal(err)
    }

    // Uji write konkuren
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            _, err := db.Exec(
                "INSERT INTO messages (topic, payload) VALUES (?, ?)",
                "test/topic",
                "message-"+string(rune(n)),
            )
            if err != nil {
                log.Printf("Write error: %v", err)
            }
        }(i)
    }
    wg.Wait()

    // Verifikasi data
    rows, _ := db.Query("SELECT COUNT(*) FROM messages")
    defer rows.Close()
    var count int
    rows.Next()
    rows.Scan(&count)
    log.Printf("Total pesan: %d", count)
}
```

**Prosedur Pengujian**:
1. Buat direktori: `tests/sqlite-wal-poc`
2. Inisialisasi Go module
3. Install driver sqlite3
4. Salin kode yang diharapkan
5. Jalankan: `go run main.go`
6. Verifikasi mode WAL diaktifkan (cek file `.db-wal` dan `.db-shm`)
7. Uji read konkuren (multiple goroutines membaca)
8. Uji write konkuren (multiple goroutines menulis)
9. Verifikasi connection pooling bekerja
10. Uji pembuatan backup
11. Dokumentasikan performa

**Perkiraan Waktu**: 3-4 jam

**Deliverables**:
- Kode PoC yang berfungsi
- Hasil verifikasi mode WAL
- Hasil uji konkurensi
- Catatan performa
- Prosedur backup/pemulihan

---

### Uji 3: Integrasi go-cache

**Tujuan**: Memverifikasi go-cache bekerja untuk caching in-memory

**Apa yang Diuji**:
- [ ] Membuat instance cache
- [ ] Set nilai dengan TTL
- [ ] Get nilai
- [ ] Ekspirasi otomatis
- [ ] Delete nilai
- [ ] Flush semua
- [ ] Thread safety

**Kriteria Sukses**:
- ✅ Cache berhasil dibuat
- ✅ Nilai dapat di-set dengan TTL
- ✅ Nilai dapat di-retrieve
- ✅ Item kadaluarsa setelah TTL
- ✅ Nilai dapat di-delete
- ✅ Cache dapat di-flush
- ✅ Akses konkuren aman

**Struktur Kode yang Diharapkan**:

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/patrickmn/go-cache"
)

func main() {
    // Buat cache dengan TTL default 5 menit
    // Cleanup interval 10 menit
    c := cache.New(5*time.Minute, 10*time.Minute)

    // Set nilai
    c.Set("foo", "bar", cache.DefaultExpiration)

    // Get nilai
    if x, found := c.Get("foo"); found {
        fmt.Println("Menemukan foo:", x.(string))
    }

    // Set dengan TTL pendek untuk pengujian
    c.Set("temp", "kadaluarsa segera", 2*time.Second)

    // Tunggu ekspirasi
    time.Sleep(3 * time.Second)

    // Cek jika kadaluarsa
    if _, found := c.Get("temp"); !found {
        fmt.Println("temp kadaluarsa sesuai harapan")
    }

    // Uji akses konkuren
    done := make(chan bool)
    for i := 0; i < 10; i++ {
        go func(n int) {
            key := fmt.Sprintf("key%d", n)
            c.Set(key, n, cache.DefaultExpiration)
            if x, found := c.Get(key); found {
                log.Printf("Goroutine %d: nilai = %v", n, x)
            }
            done <- true
        }(i)
    }

    // Tunggu semua goroutine
    for i := 0; i < 10; i++ {
        <-done
    }

    // Jumlah item
    fmt.Printf("Item cache: %d\n", c.ItemCount())
}
```

**Prosedur Pengujian**:
1. Buat direktori: `tests/go-cache-poc`
2. Inisialisasi Go module
3. Install go-cache
4. Salin kode yang diharapkan
5. Jalankan: `go run main.go`
6. Verifikasi set/get dasar bekerja
7. Verifikasi ekspirasi TTL bekerja
8. Uji akses konkuren (100+ goroutines)
9. Verifikasi thread safety
10. Dokumentasikan penggunaan memori

**Perkiraan Waktu**: 2-3 jam

**Deliverables**:
- Kode PoC yang berfungsi
- Hasil verifikasi TTL
- Hasil uji konkurensi
- Catatan penggunaan memori

---

### Uji 4: Frontend React + Vite

**Tujuan**: Memverifikasi React + Vite bekerja untuk UI admin

**Apa yang Diuji**:
- [ ] Membuat proyek React + Vite
- [ ] Konfigurasi TypeScript
- [ ] Routing dasar
- [ ] Panggilan API (fetch)
- [ ] Build untuk produksi
- [ ] Integrasi TanStack Query
- [ ] Desain responsif

**Kriteria Sukses**:
- ✅ Proyek berhasil dibuat
- ✅ TypeScript bekerja
- ✅ Routing dasar bekerja
- ✅ Dapat memanggil endpoint API
- ✅ Build produksi berhasil
- ✅ TanStack Query mengambil data
- ✅ UI responsif

**Struktur Kode yang Diharapkan**:

```bash
# Buat proyek
npm create vite@latest mqtt-gateway-poc -- --template react-ts
cd mqtt-gateway-poc

# Install dependensi
npm install

# Install paket tambahan
npm install @tanstack/react-query
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p
```

**src/App.tsx**:
```typescript
import { useQuery } from '@tanstack/react-query';

function App() {
  const { data, isLoading, error } = useQuery({
    queryKey: ['test'],
    queryFn: async () => {
      const response = await fetch('http://localhost:8080/api/test');
      if (!response.ok) throw new Error('Network error');
      return response.json();
    },
  });

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div className="p-4">
      <h1>MQTT Gateway PoC</h1>
      <pre>{JSON.stringify(data, null, 2)}</pre>
    </div>
  );
}

export default App;
```

**Prosedur Pengujian**:
1. Buat direktori: `tests/react-vite-poc`
2. Buat proyek Vite + React
3. Install dependensi
4. Buat komponen dasar
5. Tambahkan TanStack Query
6. Buat endpoint API mock (gunakan json-server atau Go server sederhana)
7. Uji panggilan API
8. Uji kompilasi TypeScript
9. Build untuk produksi: `npm run build`
10. Preview build: `npm run preview`
11. Uji responsivitas (browser devTools)

**Perkiraan Waktu**: 4-5 jam

**Deliverables**:
- Kode PoC yang berfungsi
- Artefak build
- Catatan performa (ukuran bundle, waktu load)
- Konfigurasi TypeScript

---

### Uji 5: Embed React dalam Binary Go

**Tujuan**: Memverifikasi app React dapat di-embed dalam binary Go

**Apa yang Diuji**:
- [ ] Build app React untuk produksi
- [ ] Embed static files dalam Go
- [ ] Serve app React dari Go
- [ ] Proxing API
- [ ] Dukungan routing SPA
- [ ] Hot reload dalam pengembangan

**Kriteria Sukses**:
- ✅ App React berhasil di-build
- ✅ Static files di-embed dalam binary Go
- ✅ Go serve app React
- ✅ Panggilan API bekerja
- ✅ Routing SPA bekerja
- ✅ Single binary executable

**Struktur Kode yang Diharapkan**:

**Go Backend**:
```go
package main

import (
    "embed"
    "io/fs"
    "log"
    "net/http"
)

//go:embed frontend/dist
var frontendFS embed.FS

func main() {
    // Dapatkan subdirectory untuk frontend
    distFS, err := fs.Sub(frontendFS, "frontend/dist")
    if err != nil {
        log.Fatal(err)
    }

    // Serve frontend
    http.Handle("/", http.FileServer(http.FS(distFS)))

    // Endpoint API
    http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"status":"ok"}`))
    })

    log.Println("Server mulai pada :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

**Prosedur Pengujian**:
1. Buat direktori: `tests/embed-react-poc`
2. Buat app React (gunakan PoC sebelumnya)
3. Build app React: `npm run build`
4. Buat backend Go dengan embed
5. Tempatkan build React di `frontend/dist`
6. Jalankan server Go: `go run main.go`
7. Buka browser ke `http://localhost:8080`
8. Verifikasi app React dimuat
9. Verifikasi panggilan API bekerja
10. Uji routing SPA (navigasi antar halaman)
11. Build binary Go: `go build -o gateway.exe`
12. Uji binary standalone
13. Ukur ukuran binary

**Perkiraan Waktu**: 4-6 jam

**Deliverables**:
- Kode PoC yang berfungsi
- Binary standalone
- Pengukuran ukuran binary
- Dokumentasi prosedur embed
- Workflow pengembangan (hot reload)

---

### Uji 6: Konfigurasi YAML

**Tujuan**: Memverifikasi loading konfigurasi YAML bekerja

**Apa yang Diuji**:
- [ ] Membuat file config YAML
- [ ] Parse YAML dalam Go
- [ ] Validasi konfigurasi
- [ ] Nilai default
- [ ] Variabel lingkungan
- [ ] Hot reload (opsional)
- [ ] Penanganan error

**Kriteria Sukses**:
- ✅ File YAML berhasil dibuat
- ✅ YAML berhasil di-parse
- ✅ Validasi bekerja
- ✅ Default diterapkan dengan benar
- ✅ Variabel lingkungan override config
- ✅ Error ditangani dengan graceful

**Struktur Kode yang Diharapkan**:

**config.yaml**:
```yaml
mqtt:
  embedded:
    enabled: true
    listen_address: :1883
    username: admin
    password: secret
    max_connections: 100

database:
  type: sqlite
  path: ./data/gateway.db
  wal_mode: true

web:
  address: :8080
```

**config.go**:
```go
package main

import (
    "fmt"
    "log"
    "os"

    "gopkg.in/yaml.v3"
)

type Config struct {
    MQTT struct {
        Embedded struct {
            Enabled         bool   `yaml:"enabled"`
            ListenAddress   string `yaml:"listen_address"`
            Username        string `yaml:"username"`
            Password        string `yaml:"password"`
            MaxConnections  int    `yaml:"max_connections"`
        } `yaml:"embedded"`
    } `yaml:"mqtt"`

    Database struct {
        Type    string `yaml:"type"`
        Path    string `yaml:"path"`
        WALMode bool   `yaml:"wal_mode"`
    } `yaml:"database"`

    Web struct {
        Address string `yaml:"address"`
    } `yaml:"web"`
}

func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    // Set default
    if config.MQTT.Embedded.ListenAddress == "" {
        config.MQTT.Embedded.ListenAddress = ":1883"
    }
    if config.Web.Address == "" {
        config.Web.Address = ":8080"
    }

    return &config, nil
}

func main() {
    config, err := LoadConfig("config.yaml")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("MQTT Listen: %s\n", config.MQTT.Embedded.ListenAddress)
    fmt.Printf("Web Address: %s\n", config.Web.Address)
    fmt.Printf("Database: %s\n", config.Database.Path)
}
```

**Prosedur Pengujian**:
1. Buat direktori: `tests/yaml-config-poc`
2. Buat config.yaml sample
3. Buat loader config Go
4. Uji parse YAML valid
5. Uji parse YAML invalid (verifikasi penanganan error)
6. Uji nilai default
7. Uji override variabel lingkungan
8. Uji file config tidak ada
9. Uji aturan validasi
10. Dokumentasikan skema konfigurasi

**Perkiraan Waktu**: 2-3 jam

**Deliverables**:
- Kode PoC yang berfungsi
- Skema konfigurasi
- Aturan validasi
- Contoh penanganan error

---

### Uji 7: Uji Integrasi (Semua Bersama)

**Tujuan**: Memverifikasi semua teknologi bekerja bersama

**Apa yang Diuji**:
- [ ] Mulai embedded MQTT broker
- [ ] Mulai HTTP server dengan React embedded
- [ ] Hubungkan klien MQTT
- [ ] Publish/Subscribe pesan
- [ ] Simpan pesan di SQLite
- [ ] Cache pesan terkini
- [ ] Tampilkan di UI React
- [ ] Update config via YAML

**Kriteria Sukses**:
- ✅ Semua layanan mulai tanpa error
- ✅ Broker MQTT menerima koneksi
- ✅ Web UI berhasil dimuat
- ✅ Pesan mengalir melalui sistem
- ✅ Database menyimpan pesan
- ✅ Cache meningkatkan performa
- ✅ UI update real-time
- ✅ Perubahan config bekerja

**Arsitektur yang Diharapkan**:

```
┌─────────────────────────────────────────────┐
│      Proses Go (Single Executable)          │
│                                             │
│  ┌──────────────┐  ┌──────────────┐       │
│  │ mochi-mqtt   │  │ HTTP Server  │       │
│  │ Broker       │  │              │       │
│  │              │  │ Serve:       │       │
│  │ Port: 1883   │  │ - React UI   │       │
│  └──────────────┘  │ - API        │       │
│        ↓           │              │       │
│  ┌──────────────┐  │              │       │
│  │ Hooks        │  │              │       │
│  │ (on message) │  │              │       │
│  └──────────────┘  │              │       │
│        ↓           └──────────────┘       │
│  ┌──────────────┐                          │
│  │ SQLite       │                          │
│  │ (mode WAL)   │                          │
│  └──────────────┘                          │
│        ↓                                   │
│  ┌──────────────┐                          │
│  │ go-cache     │                          │
│  │ (in-memory)  │                          │
│  └──────────────┘                          │
└─────────────────────────────────────────────┘
```

**Prosedur Pengujian**:
1. Buat direktori: `tests/integration-poc`
2. Salin semua PoC sebelumnya
3. Integrasikan ke dalam aplikasi tunggal
4. Mulai semua layanan
5. Hubungkan dengan klien MQTT
6. Publish pesan percobaan
7. Verifikasi penyimpanan database
8. Verifikasi cache hits
9. Buka web UI
10. Verifikasi update real-time
11. Uji reload config
12. Ukur performa keseluruhan
13. Uji graceful shutdown

**Perkiraan Waktu**: 8-12 jam

**Deliverables**:
- PoC integrasi lengkap
- Metrik performa
- Masalah yang diketahui
- Rekomendasi
- Keputusan arsitektur final

---

## 4. Organisasi Pengujian

### Struktur Direktori

```
tests/
├── README.md                           (ikhtisar file ini)
├── mochi-mqtt-poc/                     (Uji 1)
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   └── HASIL.md
├── sqlite-wal-poc/                     (Uji 2)
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   └── HASIL.md
├── go-cache-poc/                       (Uji 3)
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   └── HASIL.md
├── react-vite-poc/                     (Uji 4)
│   ├── src/
│   ├── package.json
│   ├── vite.config.ts
│   └── HASIL.md
├── embed-react-poc/                    (Uji 5)
│   ├── main.go
│   ├── frontend/
│   │   └── (output build React)
│   └── HASIL.md
├── yaml-config-poc/                    (Uji 6)
│   ├── main.go
│   ├── config.yaml
│   └── HASIL.md
└── integration-poc/                    (Uji 7)
    ├── main.go
    ├── embedded/
    │   └── (build React)
    └── HASIL.md
```

### Template Dokumentasi

Setiap pengujian harus memiliki file `HASIL.md` dengan:

```markdown
# Hasil Pengujian: [Nama Pengujian]

## Ringkasan
- **Tanggal**: [Tanggal]
- **Penguji**: [Nama]
- **Status**: ✅ LULUS / ❌ GAGAL

## Apa yang Diuji
[Daftar tujuan pengujian]

## Hasil
### Kriteria Sukses
- [ ] Kriteria 1 - ✅ LULUS
- [ ] Kriteria 2 - ✅ LULUS
- [ ] Kriteria 3 - ❌ GAGAL

### Performa
- Metrik 1: [Nilai]
- Metrik 2: [Nilai]

### Masalah yang Ditemukan
1. [Deskripsi masalah]
2. [Deskripsi masalah]

### Rekomendasi
1. [Rekomendasi]
2. [Rekomendasi]

## Kesimpulan
[Penilaian keseluruhan - keputusan lanjut/tidak]

## Langkah Selanjutnya
[Apa yang harus dilakukan berdasarkan hasil]
```

---

## 5. Kriteria Sukses (Keseluruhan)

### Keputusan Lanjut/Tidak

**LANJUT (Mulai Pengembangan) jika**:
- ✅ Semua 7 uji lulus
- ✅ Tidak ada blokir kritis
- ✅ Performa memenuhi persyaratan
- ✅ Integrasi bekerja
- ✅ Tim percaya diri

**TIDAK LANJUT (Re-evaluasi) jika**:
- ❌ Setiap uji gagal secara kritis
- ❌ Performa tidak dapat diterima
- ❌ Integrasi memiliki masalah besar
- ❌ Alternatif yang lebih baik ditemukan

### Persyaratan Minimum

**Persyaratan Fungsional**:
- Semua fitur inti bekerja sesuai harapan
- Tidak ada bug show-stopper
- Integrasi stabil

**Persyaratan Performa**:
- Ukuran binary < 50 MB
- Waktu startup < 2 detik
- Penggunaan memori < 100 MB (idle)
- Throughput MQTT > 1,000 msg/detik
- Load halaman web < 1 detik

**Persyaratan Developer Experience**:
- Kode jelas dan dapat dirawat
- Library memiliki dokumentasi yang baik
- Debugging mudah dilakukan
- Hot reload bekerja dalam pengembangan

---

## 6. Timeline & Sumber Daya

### Timeline

**Minggu 1**: Teknologi Backend (Uji 1-3)
**Minggu 2**: Frontend & Integrasi (Uji 4-7)

### Sumber Daya

**Orang**:
- 1 Developer Full-Stack (Go + React)
- 10-15 jam per minggu
- Total: 20-30 jam

**Tools**:
- Mesin pengembangan
- Go 1.24+
- Node.js 20+
- Repository Git
- MQTT client (MQTTX)

**Anggaran**:
- Waktu developer: ~$1,500
- Tools: $0 (semua gratis/open source)
- **Total**: ~$1,500

---

## 7. Manajemen Risiko

### Risiko Potensial

| Risiko | Dampak | Probabilitas | Mitigasi |
|------|--------|-------------|----------|
| mochi-mqtt tidak mendukung fitur yang dibutuhkan | TINGGI | RENDAH | Alternatif: Eclipse Paho |
| Masalah performa mode WAL SQLite | SEDANG | RENDAH | Dapat tune pengaturan PRAGMA |
| Embed React meningkatkan ukuran binary secara signifikan | RENDAH | SEDANG | Dapat optimasi build, gunakan kompresi |
| Masalah integrasi antar komponen | TINGGI | SEDANG | Uji integrasi akan mengungkap ini |
| Dokumentasi buruk | SEDANG | RENDAH | Buat dokumentasi kami sendiri |

### Rencana Kontinjensi

**Jika mochi-mqtt gagal**:
- Alternatif: Gunakan Eclipse Paho (mode klien saja, tanpa embedded broker)
- Dampak: Tanpa embedded broker, harus menggunakan Mosquitto eksternal
- Timeline: +1 minggu untuk implementasi alternatif

**Jika SQLite WAL gagal**:
- Alternatif: Gunakan SQLite standar (tanpa WAL)
- Dampak: Konkurensi lebih rendah
- Timeline: Perubahan minimal

**Jika embed React gagal**:
- Alternatif: Serve React terpisah (tidak embedded)
- Dampak: Bukan single binary
- Timeline: +1 minggu untuk menyesuaikan arsitektur

---

## 8. Deliverables

### Untuk Setiap Pengujian

1. **Source Code**
   - PoC lengkap dan berfungsi
   - Berkomentar dengan baik
   - Mengikuti best practices

2. **Dokumentasi Hasil**
   - File HASIL.md
   - Screenshot (jika berlaku)
   - Metrik performa
   - Masalah dan rekomendasi

3. **Instruksi Setup**
   - Cara menjalankan PoC
   - Dependensi
   - Persyaratan lingkungan

### Deliverables Final

1. **Suite Pengujian Lengkap**
   - Semua 7 PoC
   - Uji integrasi
   - Dokumentasi

2. **Laporan Ringkas**
   - Ringkasan eksekutif
   - Rekomendasi Lanjut/Tidak
   - Pelajaran yang dipelajari
   - Langkah selanjutnya

3. **Keputusan Arsitektur**
   - Tech stack final dikonfirmasi
   - Perubahan apa pun dari rencana asli
   - Justifikasi untuk keputusan

---

## 9. Langkah Selanjutnya

### Setelah Pengujian Selesai

**Jika LANJUT**:
1. Buat repository proyek utama
2. Setup lingkungan pengembangan
3. Mulai pengembangan MVP (sprint 3 bulan)
4. Referensi PoC untuk implementasi

**Jika TIDAK LANJUT**:
1. Identifikasi blokir spesifik
2. riset alternatif
3. Uji kembali komponen bermasalah
4. Buat keputusan arsitektur final

---

## 10. Kesimpulan

Rencana pengujian ini memastikan kita **memvalidasi pilihan teknologi sebelum berkomitmen ke 3 bulan pengembangan**. Dengan menghabiskan 2 minggu dan ~$1,500 sekarang, kita mengurangi risiko redo biaya atau perubahan arsitektur nanti.

**Manfaat Utama**:
- ✅ Kepercayaan pada tech stack
- ✅ Penemuan masalah awal
- ✅ Baseline performa
- ✅ Kode yang dapat digunakan kembali untuk MVP
- ✅ Pembelajaran tim

**Kesuksesan Berarti**:
- Semua teknologi diverifikasi
- Integrasi bekerja
- Tidak ada blokir
- Siap mulai pengembangan MVP

Dengan pengujian yang tepat, kita dapat memulai pengembangan MVP dengan kepercayaan, mengetahui teknologi yang dipilih akan bekerja sesuai harapan.

---

**Versi Dokumen**: 1.0
**Terakhir Diperbarui**: 9 Januari 2026
**Status**: Siap untuk Pengujian
**Periode Pengujian**: 2 minggu
**Pemilik**: Tim Pengembangan

---

## Lampiran

### A. Referensi Cepat

**Repository**: [Link ke repository pengujian]
**Wiki**: [Link ke wiki internal]
**Slack**: #mqtt-gateway-testing

### B. Sumber Daya

- [Dokumentasi mochi-mqtt](https://github.com/mochi-mqtt/mqtt)
- [Mode WAL SQLite](https://www.sqlite.org/wal.html)
- [go-cache](https://github.com/patrickmn/go-cache)
- [React + Vite](https://vitejs.dev/guide/)
- [TanStack Query](https://tanstack.com/query/latest)

### C. Kontak

**Pertanyaan**: Hubungi [Tech Lead]
**Masalah**: Buat issue GitHub
**Keadaan Darurat**: Slack #dev-urgent

---

**Disiapkan oleh**: Tim Pengembangan
**Disetujui oleh**: [Pending]
**Tanggal Review**: [Pending]
