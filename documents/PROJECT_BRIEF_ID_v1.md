# Rencana Proyek: Universal MQTT Gateway Server (UM-Gateway) v1.0

## Ringkasan Eksekutif

**Nama Proyek**: Universal MQTT Gateway Server (UM-Gateway)
**Versi**: 1.0 (Self-Hosted dengan Web Admin)
**Status**: Proposal/Perencanaan
**Rilis Target**: Q3 2026
**Organisasi**: Divisi Teknologi Rafflesia Agro

---

## Riwayat Versi

| Versi | Tanggal | Perubahan | Penulis |
|-------|---------|-----------|--------|
| **v1.0** | 2026-01-09 | **Dashboard Admin Web**: Pendekatan konfigurasi visual<br>• Menambahkan dashboard admin SvelteKit untuk konfigurasi<br>• Database SQLite dengan wizard onboarding<br>• Menggantikan konfigurasi JSON dengan form UI web<br>• Model koneksi broker MQTT eksternal<br>• Anggaran: $178,650 (berkurang dari v0)<br>• Deploy mandiri seperti PocketBase<br>• Setup zero-config dengan inisialisasi database otomatis | Tim Pengembangan |
| **v0.0** | 2026-01-09 | **Proposal Asli**: Konfigurasi berbasis JSON<br>• Konsep awal dengan definisi file JSON<br>• Koneksi broker MQTT eksternal diperlukan<br>• PostgreSQL/MongoDB dapat dikonfigurasi via JSON<br>• Fokus pada kasus penggunaan farm telemetry<br>• Anggaran: $283,200<br>• Konfigurasi dan deployment manual | Tim Pengembangan |

---

## 1. Latar Belakang

### Kondisi Saat Ini
Server Go-MQTT yang ada (v1.0) dikembangkan sebagai solusi khusus untuk operasi peternakan ayam cerdas Rafflesia Agro. Meskipun sukses di domainnya, implementasi saat ini memiliki beberapa keterbatasan:

1. **Keterikatan Kuat**: Logika bisnis di-hard-code untuk kasus telemetry pertanian
2. **Kompleksitas Konfigurasi JSON**: Memerlukan pengetahuan teknis untuk mengedit file JSON
3. **Tidak Ada Antarmuka Pengguna**: Semua konfigurasi dilakukan melalui file teks
4. **Deployment Manual**: Setiap perubahan memerlukan restart server
5. **Aksesibilitas Terbatas**: Tidak ada antarmuka manajemen berbasis web

### Evolusi dari v0 ke v1
**Pendekatan v0**: File konfigurasi JSON yang memerlukan pengeditan manual dan restart server
**Pendekatan v1**: Dashboard admin web tertanam (seperti PocketBase) dengan database SQLite tertanam untuk setup zero-config

### Peluang Pasar
Lanskap IoT membutuhkan **gateway MQTT self-hosted yang ramah pengguna** yang dapat di-deploy dan dikelola oleh pengguna non-teknis:

- **Tim Kecil**: Ingin menghindari biaya cloud SaaS tapi butuh manajemen mudah
- **Tim IT**: Lebih memilih solusi self-hosted dengan UI web daripada config berbasis file
- **Startup**: Perlu deployment cepat tanpa belajar sintaks konfigurasi kompleks
- **Deployment Edge**: Memerlukan manajemen lokal tanpa ketergantungan eksternal

---

## 2. Model Bisnis

### Proposition Nilai

**Untuk Pengguna Non-Teknis**:
- Deployment zero-config - cukup jalankan dan akses interface web
- Setup visual daripada mengedit file JSON
- Panduan onboarding memandu setup pertama kali
- Tidak perlu belajar sinteks konfigurasi

**Untuk Integrator Sistem**:
- Deployment lebih cepat (menit vs jam)
- Demonstrasi klien yang mudah
- Manajemen remote melalui UI web
- Self-service klien mengurangi beban dukungan

**Untuk Rafflesia Agro**:
- Barrier masuk lebih rendah meningkatkan adopsi
- Membedakan diri dari kompetitor config JSON
- Penampilan profesional dengan dashboard admin
- Potensi versi hosted dengan UX yang sama

### Stream Pendapatan

1. **Free Self-Hosted**: Versi dasar dengan SQLite (100% gratis, Lisensi MIT)
2. **Pro License** ($199 sekali bayar): Fitur lanjutan, dukungan multi-database
3. **Enterprise License** ($999/tahun): Dukungan prioritas, opsi white-label
4. **Cloud Service** ($29/bulan): Versi terkelola dengan UI/UX yang sama
5. **Support Contracts**: Integrasi kustom dan konsultasi

### Strategi Harga

| Edisi | Pasar Target | Model Harga | Fitur |
|-------|--------------|-------------|-------|
| Komunitas | Hobiis, mahasiswa, proyek kecil | Gratis (MIT) | SQLite, web admin, gateway MQTT dasar |
| Pro | UMKM, penggunaan profesional | $199 sekali bayar | PostgreSQL/MongoDB, fitur lanjutan, dukungan email |
| Enterprise | Organisasi besar | $999/tahun | Multi-user, RBAC, dukungan prioritas, SLA |
| Cloud | Semua pasar | $29/bulan | Terkelola penuh, auto-backup, SLA 99.9% |

---

## 3. Target Pengguna

### Pengguna Utama

**1. Pemilik Bisnis Non-Teknis**
- Butuh pemantauan IoT tapi tidak memiliki keterampilan coding
- Ingin solusi "tinggal pakai"
- Nyaman dengan form web dan dashboard
- Titik nyeri: Tidak dapat mengedit file konfigurasi JSON

**2. Tim IT Kecil**
- Mengelola proyek IoT internal
- Lebih memilih UI web daripada file config
- Butuh mendelegasi ke staf non-teknis
- Titik nyeri: Anggota tim memiliki keterampilan teknis yang bervariasi

**3. Founder Startup**
- Membangun MVP IoT dengan cepat
- Butuh iterasi pada konfigurasi
- Tidak ingin memelihara infrastruktur kompleks
- Titik nyeri: Tekanan time-to-market, sumber daya terbatas

**4. MSP (Managed Service Providers)**
- Mengelola deployment untuk banyak klien
- Butuh interface manajemen yang konsisten
- Ingin mendelegasi konfigurasi ke klien
- Titik nyeri: Mengelola deployment klien yang beragam secara efisien

### Pengguna Sekunder

**5. Mahasiswa & Pendidik**
- Belajar konsep IoT
- Butuh feedback visual
- Ingin memahami tanpa pengetahuan teknis yang mendalam

---

## 4. Pernyataan Masalah

### Masalah Utama

**Masalah 1: Konfigurasi JSON Terlalu Teknis**
- Error sintaks JSON merusak deployment
- Tidak ada validasi sampai server restart
- Sulit untuk pengguna non-teknis
- Tidak ada feedback visual pada struktur konfigurasi

**Masalah 2: Tidak Ada UI untuk Manajemen**
- Semua konfigurasi memerlukan akses SSH/file
- Tidak dapat mengkonfigurasi secara remote tanpa akses server
- Tidak ada tinjauan visual pada subscription
- Sulit troubleshooting tanpa keterampilan teknis

**Masalah 3: Konfigurasi Berbasis File Rentan Error**
- Pengeditan file manual menyebabkan error sintaks
- Tidak ada versioning konfigurasi
- Sulit rollback perubahan
- Tidak ada audit trail pada perubahan konfigurasi

**Masalah 4: Barrier Masuk Tinggi**
- Memerlukan pengetahuan tentang MQTT, database, JSON
- Tidak ada proses setup terpandu
- Kurva pembelajaran yang berat dokumentasi
- Tidak ada feedback visual selama konfigurasi

**Masalah 5: Manajemen Remote Sulit**
- Tidak dapat mengkonfigurasi tanpa akses server
- Tidak ada kolaborasi pada konfigurasi
- Tidak ada manajemen multi-user
- Sulit mengelola beberapa deployment

### Dampak

- **Barrier Adopsi**: Kompleksitas teknis mencegah penggunaan luas
- **Beban Dukungan**: Error konfigurasi sederhana memerlukan dukungan teknis
- **Waktu Deployment**: Jam/hari vs menit untuk solusi berbasis web
- **Frustrasi Pengguna**: Pengguna teknis frustrasi dengan keterbatasan UI, pengguna non-teknis tidak dapat menggunakannya
- **Peluang Hilang**: Pasar lebih memilih solusi terkelola web (dibuktikan dengan kesuksesan Home Assistant, Node-RED)

---

## 5. Solusi yang Diusulkan

### Visi: PocketBase untuk MQTT Gateway

**Gateway MQTT self-hosted dengan dashboard admin web tertanam** yang menyediakan:
- **Deployment zero-konfigurasi** - single binary, tidak perlu setup
- **Administrasi berbasis web** - semua dikelola melalui browser
- **Database SQLite** - tertanam, tidak perlu database eksternal untuk penggunaan dasar
- **Konfigurasi visual** - form, wizard, interface drag-and-drop
- **Pengalaman Onboarding** - setup terpandu untuk pertama kali

### Inovasi Utama

**1. Dashboard Admin Web Tertanam**
- Dibangun dengan SvelteKit (seperti PocketBase)
- Dilayani langsung oleh binary Go
- Single-page application untuk feedback instan
- Desain responsif untuk akses mobile/tablet

**2. Pengalaman Onboarding First-Run**
- Setup wizard saat akses pertama
- Pembuatan akun admin
- Konfigurasi awal broker MQTT
- Setup subscription topik pertama
- Pemilihan database (SQLite default, PostgreSQL/MongoDB opsional)

**3. SQLite untuk Konfigurasi & Data**
- Database tertanam (nol ketergantungan eksternal)
- Menyimpan semua konfigurasi, user, subscription
- Opsi upgrade ke PostgreSQL/MongoDB untuk scale
- Migrasi otomatis pada update versi

**4. Visual Topic Builder**
- Pembuatan subscription topik berbasis form
- Validasi real-time
- Preview konfigurasi subscription
- Tes subscription segera

**5. Arsitektur Self-Contained**
- Single binary mencakup aset UI web
- Tidak perlu server web terpisah
- Tidak ada ketergantungan eksternal untuk penggunaan dasar
- Deployment mudah - cukup jalankan executable

### Ikhtisar Arsitektur

```
┌─────────────────────────────────────────────────────────────┐
│                    Server UM-Gateway                        │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ MQTT         │  │ HTTP/WebSocket│  │ Embedded     │      │
│  │ Subscriber   │  │ API Server   │  │ Web Server   │      │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘      │
│         │                 │                 │               │
│  ┌──────▼─────────────────▼─────────────────▼───────┐      │
│  │           Mesin Pemrosesan Pesan                │      │
│  │  - Routing Topik (Dikonfigurasi Database)      │      │
│  │  - Validasi Skema                                │      │
│  │  - Transformasi Data                             │      │
│  │  - Manajemen State                                │      │
│  └──────┬────────────────────────────────────────────┘      │
│         │                                                  │
│  ┌──────▼────────────────────────────────────────────┐     │
│  │       Layer Storage (Dapat Dipasang)              │     │
│  │  ┌──────────────┐  ┌──────────────┐              │     │
│  │  │ SQLite       │  │ PostgreSQL   │              │     │
│  │  │ (Tertanam)   │  │ (Opsional)   │              │     │
│  │  └──────────────┘  └──────────────┘              │     │
│  │  ┌──────────────┐  ┌──────────────┐              │     │
│  │  │ MongoDB      │  │ InfluxDB     │              │     │
│  │  │ (Opsional)   │  │ (Opsional)   │              │     │
│  │  └──────────────┘  └──────────────┘              │     │
│  └──────────────────────────────────────────────────┘     │
│                                                           │
│  ┌──────────────────────────────────────────────────┐    │
│  │      Dashboard Admin (SPA SvelteKit)              │    │
│  │  - Wizard Onboarding                               │    │
│  │  - UI Manajemen Topik                              │    │
│  │  - Builder Konfigurasi                             │    │
│  │  - Monitoring Real-time                             │    │
│  │  - Manajemen User                                   │    │
│  └──────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

---

## 6. Tujuan Proyek

### Tujuan Utama (Harus Ada)

**O1: Kembangkan Dashboard Admin Web Self-Contained**
- SPA berbasis SvelteKit tertanam dalam binary Go
- Desain responsif untuk semua ukuran layar
- Dukungan tema gelap/terang
- Dukungan multi-bahasa (EN, ID awalnya)

**O2: Implementasikan Pengalaman Onboarding First-Run**
- Setup wizard untuk konfigurasi awal
- Pembuatan akun admin dengan password aman
- Pengujian koneksi broker MQTT
- Pembuatan subscription topik pertama
- Konfirmasi sukses dengan langkah selanjutnya

**O3: Database SQLite untuk Konfigurasi**
- SQLite tertanam untuk penyimpanan config
- Migrasi skema otomatis
- Jalur upgrade opsional ke PostgreSQL/MongoDB
- Fungsionalitas backup/restore

**O4: Visual Builder Subscription Topik**
- Pembuatan subscription topik berbasis form
- Validasi topik MQTT real-time
- Builder aturan transformasi (visual)
- Selector tujuan penyimpanan
- Tes subscription sebelum mengaktifkan

**O5: Autentikasi & Otorisasi User**
- Pembuatan akun admin saat first-run
- Manajemen user tambahan (fitur Pro)
- Kontrol akses berbasis peran (Admin, Editor, Viewer)
- Manajemen sesi dengan timeout

### Tujuan Sekunder (Sebaiknya Ada)

**O6: Dashboard Monitoring Real-Time**
- Tampilan throughput pesan langsung
- Status subscription aktif
- Indikator kesehatan koneksi
- Monitoring kapasitas penyimpanan

**O7: Import/Ekspor Konfigurasi**
- Ekspor konfigurasi sebagai JSON
- Import konfigurasi dari file
- Template konfigurasi
- Berbagi konfigurasi antar deployment

**O8: Dukungan Multi-Database (Pro)**
- Upgrade dari SQLite ke PostgreSQL
- MongoDB untuk skema fleksibel
- InfluxDB untuk optimasi time-series
- Wizard migrasi satu-klik

**O9: API untuk Otomasi**
- RESTful API yang sesuai dengan fungsionalitas UI web
- Generasi API key
- Dokumentasi API dalam admin UI
- Contoh kode untuk operasi umum

### Tujuan Tersier (Nice to Have)

**O10: Mobile App**
- React Native atau Flutter mobile app
- Notifikasi push untuk alert
- Cek status cepat
- UI mobile yang disederhanakan

**O11: Sistem Plugin**
- Prosesor data kustom
- Adapter storage kustom
- Provider autentikasi kustom
- Marketplace plugin

---

## 7. Ruang Lingkup & Batasan Proyek

### Dalam Ruang Lingkup (Apa yang Akan Kami Bangun)

#### Fungsionalitas Inti
1. **Server Web Tertanam**
   - Layani SPA SvelteKit dari filesystem tertanam
   - Dukungan WebSocket untuk update real-time
   - HTTPS otomatis dengan Let's Encrypt (opsional)

2. **Dashboard Admin**
   - Setup wizard untuk konfigurasi pertama kali
   - Manajemen subscription topik
   - Manajemen koneksi broker MQTT
   - Konfigurasi tujuan penyimpanan
   - Manajemen user (fitur Pro)

3. **Database SQLite**
   - Penyimpanan konfigurasi
   - Data autentikasi user
   - Subscription topik
   - Log audit
   - Penyimpanan data time-series opsional

4. **Pemrosesan MQTT**
   - Subscribe ke topik yang dikonfigurasi melalui UI
   - Routing pesan berdasarkan config database
   - Transformasi data
   - Manajemen state

5. **Opsi Storage**
   - SQLite (tertanam, default)
   - PostgreSQL (fitur Pro)
   - MongoDB (fitur Pro)
   - InfluxDB (fitur Pro)

#### Fitur Web UI
- Wizard onboarding
- Dashboard home
- Manajemen topik
- Manajemen storage
- Manajemen user (Pro)
- Halaman settings
- Status sistem

#### Deployment & Operasi
- Deployment single binary
- Image Docker
- File service systemd
- Build cross-platform (Linux, macOS, Windows)
- Mekanisme auto-update (opsional)

### Di Luar Ruang Lingkup (Apa yang Tidak Akan Kami Bangun - Awalnya)

1. **Hardware/Firmware**
   - Firmware perangkat (ESP32, Arduino, dll)
   - Provisioning hardware
   - Update OTA

2. **Analitik Lanjutan**
   - Dashboard visualisasi data (dapat integrasi dengan Grafana)
   - Business intelligence
   - Analitik prediktif

3. **Multi-Tenansi** (Fase 2)
   - Beberapa organisasi dalam instance tunggal
   - Kuota sumber daya per tenant

4. **Edge Computing** (Fase 3)
   - Versi ringan
   - Mode operasi offline

### Batasan & Kendala

#### Kendala Teknis
- **Bahasa**: Go 1.24+ (backend), SvelteKit (frontend)
- **Versi MQTT**: Dukungan 3.1.1 dan 5.0
- **Hardware Minimum**: 1 core CPU, 512MB RAM untuk penggunaan dasar
- **Jaringan**: Konektivitas ke broker MQTT dan database eksternal opsional

#### Kendala Waktu
- **Fase 1 (MVP)**: 9 bulan
- **Fase 2 (Fitur Pro)**: +3 bulan
- **Fase 3 (Lanjutan)**: +6 bulan

#### Kendala Anggaran
- **Tim Pengembangan**: 2-3 pengembang full-time
- **Infrastruktur**: $1.000/bulan untuk pengembangan/tes
- **Layanan Pihak Ketiga**: Tier gratis awalnya

#### Kendala Sumber Daya
- **Tim Pengembangan**: 1-2 Full-stack developer (Go + SvelteKit), 1 desainer UI/UX
- **Ahli Domain**: Spesialis IoT untuk konsultasi
- **Dukungan**: Technical writer part-time

---

## 8. Stack Teknologi

### Teknologi Inti

#### Backend
| Komponen | Teknologi | Alasan |
|-----------|-----------|--------|
| Bahasa | Go 1.24+ | Performa, konkurensi, single binary |
| HTTP Router | Chi v5 | Ringan, idiomatic Go |
| Klien MQTT | Eclipse Paho MQTT v2 | Mature, dukungan MQTT 5.0 |
| Embedded Files | embed package | Bundle web UI dalam binary |
| SQLite | mattn/go-sqlite3 | Driver SQLite murni Go |
| WebSocket | gorilla/websocket | Update UI real-time |

#### Frontend (Dashboard Admin)
| Komponen | Teknologi | Alasan |
|-----------|-----------|--------|
| Framework | SvelteKit | Modern, cepat, UX bagus, bundle kecil |
| UI Components | Skeleton UI | Library komponen Svelte, default bagus |
| State Management | Store Svelte | Built-in, sederhana, reaktif |
| Forms | SvelteKit form actions | Handling form native |
| Charts | Chart.js / ApexCharts | Visualisasi |
| Icons | Lucide Svelte | Library ikon ringan |
| Styling | TailwindCSS | Pengembangan UI cepat |

#### Opsi Database
| Database | Use Case | Kelebihan |
|----------|----------|-----------|
| **SQLite** | Default, tertanam | Zero-config, single file, cukup untuk beban moderat |
| **PostgreSQL** | Scaling produksi | SQL, ACID, ekosistem mature |
| **MongoDB** | Skema fleksibel | Tidak perlu migrasi, scaling horizontal |
| **InfluxDB** | Optimasi time-series | Purpose-built, throughput tinggi |

#### Autentikasi & Keamanan
| Komponen | Teknologi | Alasan |
|-----------|-----------|--------|
| Password Hashing | bcrypt | Standar industri |
| Session Storage | SQLite (default) / Redis (opsional) | Tertanam atau scalable |
| JWT | golang-jwt/jwt | Auth stateless untuk API |
| CSRF Protection | Built-in SvelteKit | Best practice keamanan |

#### Deployment
| Komponen | Teknologi | Alasan |
|-----------|-----------|--------|
| Build | Go embed + SvelteKit static adapter | Single binary |
| Kontainerisasi | Docker | Standar industri |
| Cross-platform | Go cross-compilation | Single codebase, multiple platform |

### Pola Arsitektur

**1. UI Web Tertanam**
```go
//go:embed all:dist
var uiFS embed.FS

func (s *Server) ServeDashboard(w http.ResponseWriter, r *http.Request) {
    http.FileServer(http.FS(uiFS)).ServeHTTP(w, r)
}
```

**2. Konfigurasi Berbasis Database**
```sql
-- Schema SQLite
CREATE TABLE mqtt_brokers (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    username TEXT,
    password TEXT,
    enabled BOOLEAN DEFAULT 1
);

CREATE TABLE topic_subscriptions (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    topic_pattern TEXT NOT NULL,
    qos INTEGER DEFAULT 1,
    storage_destination_id INTEGER,
    transformation_config TEXT, -- JSON
    enabled BOOLEAN DEFAULT 1
);
```

**3. Update Real-Time**
```go
// WebSocket endpoint untuk dashboard
func (s *Server) DashboardWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, _ := upgrader.Upgrade(w, r, nil)
    for {
        // Push updates: count pesan, status koneksi, dll
    }
}
```

---

## 9. Fitur-fitur

### 9.1 Fitur Utama

#### F1: Wizard Onboarding First-Run

**Deskripsi**: Pengalaman setup terpandu untuk pengguna pertama kali

**Layar**:

**Layar 1: Selamat Datang**
```
┌─────────────────────────────────────┐
│  Selamat Datang di UM-Gateway!     │
│                                     │
│  Mari kita atur dalam 3 langkah    │
│                                     │
│  [Mulai]                            │
└─────────────────────────────────────┘
```

**Layar 2: Buat Akun Admin**
```
┌─────────────────────────────────────┐
│  Buat Akun Admin Anda              │
│                                     │
│  Email: [________________]          │
│  Password: [________________]       │
│  Konfirmasi: [________________]    │
│                                     │
│  [Kembali]  [Lanjut]               │
└─────────────────────────────────────┘
```

**Layar 3: Hubungkan ke Broker MQTT**
```
┌─────────────────────────────────────┐
│  Hubungkan ke Broker MQTT           │
│                                     │
│  Nama Broker: [Broker Saya]        │
│  URL Broker: [mqtt://localhost...]│
│  Port: [1883]                       │
│  Username: [________________]      │
│  Password: [________________]      │
│                                     │
│  [Tes Koneksi]                      │
│  ✓ Berhasil terhubung!             │
│                                     │
│  [Kembali]  [Lanjut]               │
└─────────────────────────────────────┘
```

**Layar 4: Buat Subscription Topik Pertama**
```
┌─────────────────────────────────────┐
│  Subscribe ke Topik                 │
│                                     │
│  Nama Subscription: [Sensor Farm]  │
│  Pola Topik: [farms/+/sensors/#]   │
│  QoS: [0 ▼]                         │
│                                     │
│  □ Aktifkan transformasi data       │
│  □ Tambah metadata                  │
│                                     │
│  Storage: SQLite ▼                  │
│                                     │
│  [Tes Subscription]                 │
│  ✓ Menerima pesan...                │
│                                     │
│  [Kembali]  [Selesai]               │
└─────────────────────────────────────┘
```

**Layar 5: Sukses**
```
┌─────────────────────────────────────┐
│  Semua Siap! 🎉                     │
│                                     │
│  Gateway MQTT Anda sekarang berjalan│
│                                     │
│  Subscription Aktif: 1            │
│  Pesan Diterima: 127                │
│                                     │
│  [Ke Dashboard]                     │
│  [Lihat Dokumentasi]                 │
└─────────────────────────────────────┘
```

---

#### F2: Dashboard Admin Home

**Layout**:
```
┌─────────────────────────────────────────────────────────────┐
│  UM-Gateway  [☰]  [Admin ▼]  [●]              [Cari...]  │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   127       │  │     1       │  │   Berjalan  │         │
│  │  Pesan/dtk   │  │  Sub Aktif  │  │    Status   │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│                                                               │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │  Aktivitas Terbaru                                    │ │
│  │  • Subscription baru: farm/#/sensors  (2 menit lalu) │ │
│  │  • Konfigurasi diperbarui (1 jam lalu)               │ │
│  │  • User 'admin' login (3 jam lalu)                   │ │
│  └─────────────────────────────────────────────────────────┘ │
│                                                               │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │  Aksi Cepat                                           │ │
│  │  [+ Tambah Subscription]  [+ Tambah Broker]  [⚙ Set] │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

---

#### F3: Visual Builder Subscription Topik

**Interface Berbasis Form**:

```
┌─────────────────────────────────────────────────────────────┐
│  Subscription Topik Baru                       [×]        │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  Informasi Dasar                                            │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Nama *        [________________]                    │   │
│  │ Deskripsi     [________________]                    │   │
│  │ Aktif         ☑                                     │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  Konfigurasi MQTT                                          │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Pola Topik * [sensors/+/telemetry/#]               │   │
│  │                  ✓ Pola valid                       │   │
│  │                                                      │   │
│  │ Level QoS    ○ None  ● Sekali saja  ○ Persis     │   │
│  │                                                      │   │
│  │ [Tes Koneksi]  Menerima 23 pesan/menit             │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  Tujuan Penyimpanan                                        │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Tujuan       SQLite ▼                              │   │
│  │ Tabel        [sensor_readings]                      │   │
│  │                                                      │   │
│  │ [+ Buat Tabel Baru]                                 │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  Transformasi (Opsional)                                  │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ ☐ Aktifkan transformasi data                        │   │
│  │                                                      │   │
│  │ [+ Tambah Aturan Transformasi]                       │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  [Batal]  [Simpan & Aktifkan]                               │
└─────────────────────────────────────────────────────────────┘
```

---

#### F4: Manajemen Database SQLite

**Visual Schema Builder**:

```
┌─────────────────────────────────────────────────────────────┐
│  Penyimpanan Data                                [×]        │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  Database Saat Ini                                         │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ 🔵 main (SQLite)                   [Kelola] [Ekspor]│   │
│  │    Tabel: 5  Ukuran: 127 MB  Record: 1.2M         │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  [+ Tambah Database]                                         │
│                                                               │
│  Tabel dalam 'main'                                          │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ sensor_readings                    [Lihat] [Edit]   │   │
│  │   Kolom: id, timestamp, device_id, value, ...      │   │
│  │   Record: 1.234.567                                  │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ device_states                      [Lihat] [Edit]   │   │
│  │   Kolom: id, device_id, state, updated_at          │   │
│  │   Record: 45                                        │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  [+ Buat Tabel]  [Backup Database]                           │
└─────────────────────────────────────────────────────────────┘
```

---

#### F5: Dashboard Monitoring Real-Time

```
┌─────────────────────────────────────────────────────────────┐
│  Dashboard Monitoring                                         │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  Kesehatan Sistem                                            │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  │
│  │  ✓ MQTT  │  │  ✓ DB    │  │  ✓ API   │  │  ✓ Web   │  │
│  │ Terhubung│  │  Sehat   │  │  Berjalan│  │  Melayani│  │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘  │
│                                                               │
│  Throughput Pesan (24 Jam Terakhir)                          │
│  ┌─────────────────────────────────────────────────────┐   │
│  │     ╭────╮                                            │   │
│  │  120 │╱╲│  ╭─╮  ╭─╮                                 │   │
│  │   90 │ ╲ │╱ ╰──╯╭─╯╭─╯                             │   │
│  │   60 │  ╲╱      │  │                                │   │
│  │   30 │   ╱      │  │                                │   │
│  │    0 │  ╱       ╰──╯                                │   │
│  │      └────────────────────────────────>             │   │
│  │       00  06  12  18  24                            │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  Subscription Aktif                                         │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ farm/#/sensors/#     45 msg/dtk  ● Aktif            │   │
│  │ building/+/temp/#     12 msg/dtk  ● Aktif            │   │
│  │ devices/+/state       3 msg/dtk   ● Aktif            │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

---

#### F6: Manajemen User (Fitur Pro)

```
┌─────────────────────────────────────────────────────────────┐
│  Users                                               [+ Tambah]│
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  admin@example.com                    [Edit] [×]    │   │
│  │  Role: Admin  ●  Aktif terakhir: 2 menit lalu   │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  user@company.com                      [Edit] [×]    │   │
│  │  Role: Editor  ○  Aktif terakhir: 1 hari lalu   │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  Roles                                                        │
│  • Admin - Akses penuh                                      │
│  • Editor - Dapat kelola subscription, bukan user       │
│  • Viewer - Akses read-only                                 │
└─────────────────────────────────────────────────────────────┘
```

---

### 9.2 Fitur Lanjutan (Pro)

#### F7: Wizard Upgrade Database

**SQLite → PostgreSQL**:

```
┌─────────────────────────────────────────────────────────────┐
│  Upgrade Database                                             │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  Saat Ini: SQLite (main)                    [Backup Sekarang]│
│  Upgrade ke: ○ PostgreSQL  ● MongoDB  ○ InfluxDB        │
│                                                               │
│  Koneksi PostgreSQL                                         │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Host: [localhost]                                   │   │
│  │ Port: [5432]                                        │   │
│  │ Database: [um_gateway_prod]                         │   │
│  │ Username: [____________]                             │   │
│  │ Password: [____________]                             │   │
│  │                                                      │   │
│  │ [Tes Koneksi] ✓ Terhubung                          │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  Opsi Migrasi                                               │
│  ☑ Migrasi data konfigurasi                              │
│  ☑ Migrasi data time-series (mungkin butuh beberapa jam)│
│  ☐ Pertahankan SQLite sebagai backup                     │
│                                                               │
│  Perkiraan waktu: ~2 jam untuk 1.2M record                   │
│                                                               │
│  [Mulai Migrasi]  [Batal]                                   │
└─────────────────────────────────────────────────────────────┘
```

---

#### F8: Import/Ekspor Konfigurasi

```
┌─────────────────────────────────────────────────────────────┐
│  Import/Ekspor Konfigurasi                                   │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  Ekspor Konfigurasi                                         │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Sertakan:                                             │   │
│  │ ☑ Subscription topik                                │   │
│  │ ☑ MQTT broker                                        │   │
│  │ ☑ Tujuan penyimpanan                                  │   │
│  │ ☐ Users (jika mengekspor untuk backup)            │   │
│  │                                                      │   │
│  │ [Ekspor sebagai JSON]  [Ekspor sebagai File]        │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  Import Konfigurasi                                         │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Drag & drop file di sini atau klik untuk browse     │   │
│  │                                                      │   │
│  │ File konfigurasi: um-gateway-config.json           │   │
│  │                                                      │   │
│  │ [Preview]  [Import]                                  │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  Template                                                    │
│  • Template Starter Smart Farm                            │
│  • Template Smart Home Dasar                              │
│  • Template Monitoring Industri                          │
└─────────────────────────────────────────────────────────────┘
```

---

## 10. Roadmap Pengembangan

### Fase 1: Pondasi (Bulan 1-4)

**Sprint 1-2: Setup Proyek**
- [ ] Struktur repository (Go + SvelteKit monorepo)
- [ ] Setup environment pengembangan
- [ ] Pipeline CI/CD
- [ ] Setup sistem desain

**Sprint 3-4: Backend Inti**
- [ ] Schema database SQLite
- [ ] Model konfigurasi
- [ ] Refaktor klien MQTT
- [ ] Server web tertanam
- [ ] Dukungan WebSocket

**Sprint 5-6: Pengalaman Onboarding**
- [ ] Deteksi first-run
- [ ] Pembuatan akun admin
- [ ] Wizard koneksi broker MQTT
- [ ] Setup subscription pertama
- [ ] Flow konfirmasi sukses

**Sprint 7-8: Dashboard Dasar**
- [ ] Setup proyek SvelteKit
- [ ] Layout dan navigasi
- [ ] Dashboard home
- [ ] UI autentikasi
- [ ] Halaman settings

**Milestone**: Alpha - Gateway self-hosted dengan onboarding

---

### Fase 2: Fitur Inti (Bulan 5-7)

**Sprint 9-10: Manajemen Topik**
- [ ] UI builder subscription topik
- [ ] Validasi pola topik
- [ ] Selector tujuan penyimpanan
- [ ] Fitur tes subscription
- [ ] Tampilan daftar subscription

**Sprint 11-12: Manajemen Storage**
- [ ] Browser database SQLite
- [ ] UI pembuatan tabel
- [ ] Viewer data (paginasi)
- [ ] Fungsionalitas backup/restore
- [ ] Ekspor ke CSV/JSON

**Sprint 13-14: Monitoring Real-Time**
- [ ] Update real-time WebSocket
- [ ] Chart throughput
- [ ] Indikator status koneksi
- [ ] Tampilan subscription aktif
- [ ] Halaman kesehatan sistem

**Milestone**: Beta - Gateway self-hosted fitur lengkap

---

### Fase 3: Poles & Produksi (Bulan 8-9)

**Sprint 15-16: Pengujian & Kualitas**
- [ ] Pengujian end-to-end
- [ ] Load testing (1K msg/detik)
- [ ] Audit keamanan
- [ ] Optimasi performa
- [ ] Pengujian cross-platform

**Sprint 17-18: Dokumentasi**
- [ ] Panduan pengguna (dengan screenshot)
- [ ] Panduan instalasi
- [ ] Panduan troubleshooting
- [ ] Tutorial video (onboarding, tugas umum)
- [ ] Dokumentasi API

**Sprint 19-20: Deployment**
- [ ] Build single binary (Linux, macOS, Windows)
- [ ] Image Docker
- [ ] Script instalasi
- [ ] Mekanisme auto-update
- [ ] Otomasi proses rilis

**Milestone**: v1.0 Ketersediaan Umum

---

## 11. Metrik Keberhasilan

### Metrik Teknis

| Metrik | Target | Pengukuran |
|--------|--------|------------|
| Waktu ke Subscription Pertama | < 5 menit | Pengujian otomatis |
| Waktu Load Dashboard | < 2 detik | Monitoring performa |
| Ukuran Binary | < 50 MB | Artifact build |
| Penggunaan Memori | < 200 MB (idle) | Monitoring sumber daya |
| Pesan/Detik (SQLite) | 10K msg/detik | Tes benchmark |

### Metrik Bisnis

| Metrik | Target | Timeline |
|--------|--------|----------|
| Instalasi Aktif | 500 | 6 bulan post-launch |
| Bintang GitHub | 1.000 | 6 bulan post-launch |
| Penjualan Lisensi Pro | 50 | 6 bulan post-launch |
| Tingkat Pengembalian | < 5% | 6 bulan post-launch |
| Durasi Sesi Rata-rata | > 10 menit | Analytics |

### Metrik Kualitas

| Metrik | Target | Pengukuran |
|--------|--------|------------|
| Tingkat Penyelesaian Setup | > 90% | Analytics |
| Drop-off Onboarding | < 10% per langkah | Analytics |
| Request Dukungan/100 User | < 5 | Pelacakan dukungan |
| Kepuasan Pengguna | > 4.5/5 | Survei |
| Kelengkapan Dokumentasi | 100% fitur | Review manual |

---

## 12. Analisis Risiko

### Risiko Teknis

| Risiko | Dampak | Probabilitas | Mitigasi |
|--------|--------|--------------|----------|
| Batas performa SQLite | Sedang | Sedang | Dokumentasikan batas, tawarkan upgrade PostgreSQL |
| Ukuran binary UI tertanam meningkat | Rendah | Tinggi | Kompresi, lazy loading |
| Kurva pembelajaran SvelteKit | Sedang | Rendah | Pelatihan tim, prototype awal |
| Kompleksitas build cross-platform | Sedang | Sedang | Matriks build GitHub Actions |

### Risiko Bisnis

| Risiko | Dampak | Probabilitas | Mitigasi |
|--------|--------|--------------|----------|
| Kompetitor menambahkan UI web | Tinggi | Sedang | Pengembangan cepat, fokus pada UX |
| Pasar lebih memilih cloud daripada self-hosted | Sedang | Rendah | Tawarkan kedua opsi |
| Beban dukungan tier gratis | Tinggi | Sedang | Forum komunitas, dokumentasi |

### Risiko Operasional

| Risiko | Dampak | Probabilitas | Mitigasi |
|--------|--------|--------------|----------|
| Kehilangan data user | Tinggi | Rendah | Prompt auto-backup, fitur ekspor |
| Kerentanan keamanan | Tinggi | Rendah | Audit keamanan, scanning dependensi |
| Upgrade sulit | Sedang | Sedang | Auto-migrasi, pengujian ekstensif |

---

## 13. Kebutuhan Sumber Daya

### Struktur Tim

**Tim Inti (Fase 1-3)**
- 1x Developer Full-Stack (Go + SvelteKit)
- 1x Developer Frontend (fokus SvelteKit)
- 1x Desainer UI/UX (part-time)
- 1x Engineer QA (part-time)

**Pemangku Kepentingan**
- Manajer Produk
- Technical Lead
- Ahli Domain (konsultan IoT)

### Estimasi Anggaran

| Kategori | Biaya (Bulanan) | Durasi |
|----------|------------------|---------|
| Gaji Tim Pengembangan | $15.000 | 9 bulan |
| Infrastruktur (Dev/Test) | $1.000 | 9 bulan |
| Tool & Layanan | $300 | Berkelanjutan |
| Sumber Daya Desain | $1.000 | 6 bulan |
| Kontinjensi (15%) | $2.550 | - |
| **Total Fase 1-3** | **$19.850/bulan × 9** | **$178.650** |

---

## 14. Perbandingan: v0 vs v1

| Aspek | v0 (Config JSON) | v1 (Web Admin) |
|-------|------------------|-----------------|
| Konfigurasi | File JSON | UI Web + SQLite |
| Setup Pertama | Config manual | Wizard onboarding |
| Tipe User | Teknis | Teknis + Non-teknis |
| Waktu Deployment | 30-60 menit | 2-5 menit |
| Barrier Masuk | Tinggi | Rendah |
| Manajemen | SSH/Akses file | Browser web |
| Kolaborasi | Git/Version control | Multi-user dengan RBAC |
| Feedback Visual | Tidak ada (edit file) | Validasi real-time |
| Pencegahan Error | Runtime saja | Sebelum submit |
| Kurva Pembelajaran | Curam | Landai |

---

## 15. Langkah Selanjutnya

### Tindakan Segera (Minggu 1-2)

1. **Review Pemangku Kepentingan**
   - Presentasikan perbandingan v0 vs v1
   - Validasikan pendekatan admin web
   - Amankan persetujuan anggaran

2. **Proof-of-Concept Teknis**
   - SvelteKit tertanam dalam binary Go
   - Integrasi SQLite untuk konfigurasi
   - Flow onboarding dasar

3. **Setup Sistem Desain**
   - Pilih library komponen UI
   - Buat mockup desain
   - Definisikan flow onboarding

4. **Formasi Tim**
   - Rekrut developer SvelteKit
   - Definisikan proses kolaborasi

### Tindakan Jangka Pendek (Bulan 1)

5. **Environment Pengembangan**
   - Struktur monorepo (Go + SvelteKit)
   - Hot reload untuk backend dan frontend
   - CI/CD untuk build cross-platform

6. **Perencanaan Sprint**
   - Perincian sprint detail
   - Keputusan arsitektur
   - Komitmen teknologi

7. **Prototype**
   - Wizard onboarding yang berfungsi
   - UI subscription topik dasar
   - Integrasi SQLite

---

## 16. Kesimpulan

Pendekatan v1 dengan **dashboard admin web tertanam** secara dramatis meningkatkan pendekatan v0 konfigurasi JSON:

### Peningkatan Kunci

**Aksesibilitas**
- ✅ Pengguna non-teknis sekarang dapat menggunakan sistem
- ✅ Tidak perlu belajar sinteks JSON
- ✅ Feedback visual di seluruh setup

**Pengalaman Pengguna**
- ✅ Onboarding terpandu mengurangi error
- ✅ Validasi real-time mencegah miskonfigurasi
- ✅ Manajemen berbasis web dari mana saja

**Time to Value**
- ✅ 2-5 menit ke subscription pertama vs 30-60 menit
- ✅ Self-service mengurangi beban dukungan
- ✅ Konfigurasi visual lebih cepat daripada mengedit file

**Posisi Pasar**
- ✅ Membedakan diri dari kompetitor config JSON
- ✅ Bersaing dengan SaaS cloud pada UX
- ✅ Menarik audiens yang lebih luas

Pendekatan **SvelteKit + SQLite + Go Tertanam** mengikuti pola sukses **PocketBase**, menciptakan solusi self-hosted yang semudah digunakan seperti alternatif cloud sambil mempertahankan privasi dan kontrol data.

Dengan eksekusi yang fokus selama 9 bulan, kita dapat mengirimkan **v1.0 produksi** yang mengatasi kebutuhan pasar nyata dan membangun pondasi untuk pertumbuhan jangka panjang.

---

**Versi Dokumen**: 1.0
**Terakhir Diperbarui**: 9 Januari 2026
**Status**: Draft - Menunggu Review
**Menggantikan**: PROJECT_BRIEF_ID_v0.md

---

## Lampiran

### A. Glossarium

- **Self-Hosted**: Software dijalankan di infrastruktur sendiri pengguna, bukan cloud
- **Tertanam**: Dibundle dalam single executable/binary
- **Onboarding**: Proses setup terpandu untuk pengguna baru
- **Zero-Config**: Tidak memerlukan konfigurasi manual
- **SPA**: Single Page Application

### B. Referensi

1. PocketBase (inspirasi untuk UI web tertanam)
2. Home Assistant (konfigurasi IoT visual)
3. Node-RED (pemrograman visual untuk IoT)
4. Dokumentasi SvelteKit
5. Driver SQLite Go tertanam

### C. Dokumen Terkait

- `PROJECT_BRIEF_EN_v0.md` - Pendekatan JSON asli
- `PROJECT_BRIEF_ID_v0.md` - Pendekatan asli (Bahasa Indonesia)
- `PROJECT_DESCRIPTION.md` - Dokumentasi sistem saat ini

---

**Disiapkan oleh**: Tim Pengembangan
**Disetujui oleh**: [Pending]
**Tanggal Review**: [Pending]
