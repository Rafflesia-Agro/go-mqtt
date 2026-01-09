# Rencana Proyek: Universal MQTT Gateway Server (UM-Gateway) v3.0

## Ringkasan Eksekutif

**Nama Proyek**: Universal MQTT Gateway Server (UM-Gateway)
**Versi**: 3.0 (UI Berbasis React dengan Kejelasan Konfigurasi yang Ditingkatkan)
**Status**: Proposal/Perencanaan
**Target Rilis**: Q1 2027
**Organisasi**: Divisi Teknologi Rafflesia Agro

---

## Riwayat Versi

| Versi | Tanggal | Perubahan | Penulis |
|-------|---------|-----------|--------|
| **v3.0** | 2026-01-09 | **Perubahan Technology Stack & Peningkatan Kejelasan**: UI berbasis React<br>• Mengganti SvelteKit dengan React + Vite<br>• Menambahkan TanStack Query untuk manajemen state server<br>• Menambahkan TanStack Router untuk routing<br>• Memperjelas Setup Mode = editor config.yaml<br>• Memperjelas Runner Mode = operasi produksi (config read-only)<br>• Pemisahan jelas: config.yaml (infrastruktur) vs database (data runtime)<br>• Dokumentasi yang ditingkatkan tentang apa ditaruh di mana<br>• Anggaran: $215,500 (10 bulan) | Tim Pengembangan |
| **v2.0** | 2026-01-09 | **Perubahan Arsitektur Utama**: Implementasi embedded MQTT broker<br>• Menambahkan embedded MQTT broker menggunakan library mochi-mqtt<br>• Memperkenalkan operasi dual-mode (setup/runner)<br>• Mengubah konfigurasi dari database ke YAML config.yaml<br>• Menggantikan cache eksternal dengan in-memory cache (go-cache)<br>• Mode WAL SQLite sebagai default (dioptimalkan untuk konkurensi)<br>• Arsitektur single executable dengan koordinasi proses<br>• Contoh kode baru untuk embedded broker, mode WAL, adapter cache<br>• Anggaran: $238,150 (10 bulan) | Tim Pengembangan |
| **v1.0** | 2026-01-09 | **Dashboard Admin Web**: Pendekatan konfigurasi visual<br>• Menambahkan dashboard admin SvelteKit untuk konfigurasi<br>• Database SQLite dengan wizard onboarding<br>• Menggantikan konfigurasi JSON dengan form UI web<br>• Model koneksi broker MQTT eksternal<br>• Anggaran: $178,650 (berkurang dari v0) | Tim Pengembangan |
| **v0.0** | 2026-01-09 | **Proposal Asli**: Konfigurasi berbasis JSON<br>• Konsep awal dengan definisi file JSON<br>• Koneksi broker MQTT eksternal diperlukan<br>• PostgreSQL/MongoDB dapat dikonfigurasi via JSON<br>• Fokus pada kasus penggunaan farm telemetry<br>• Anggaran: $283,200 | Tim Pengembangan |

---

## 1. Latar Belakang

### Evolusi dari v2 ke v3

**Pendekatan v2**: UI berbasis SvelteKit dengan pemisahan setup dan runner mode yang tidak jelas
**Pendekatan v3**: UI berbasis React dengan ekosistem TanStack, definisi mode yang sangat jelas dan pemisahan data

### Keterbatasan v2

1. **Tujuan Mode Tidak Jelas**: Tujuan setup mode vagu - apa sebenarnya yang bisa dikonfigurasi?
2. **Kebingungan Tentang config.yaml**: Tidak jelas data apa yang masuk ke config.yaml vs database
3. **Pilihan Framework**: SvelteKit kurang familiar bagi sebagian besar developer dibanding React
4. **Manajemen State**: Tidak ada strategi manajemen state server yang jelas
5. **Data Infrastruktur vs Runtime**: Tidak ada panduan jelas tentang apa ditaruh di mana

### Peningkatan Utama di v3

**1. Perubahan Technology Stack**
- **React + Vite**: Build lebih cepat, ekosistem lebih baik, pool talent lebih besar
- **TanStack Query**: Manajemen state server yang kuat, caching, sinkronisasi
- **TanStack Router**: Routing dengan type safety dan pengalaman developer yang excellent

**2. Definisi Mode yang Sangat Jelas**

**Setup Mode** (Mode Konfigurasi):
- **Tujuan**: Edit pengaturan infrastruktur `config.yaml`
- **Kapan**: First run, atau saat mengubah konfigurasi infrastruktur
- **Apa yang Bisa Anda Ubah**:
  - Pengaturan koneksi database (tipe, path, kredensial)
  - Pengaturan broker MQTT (embedded/eksternal, port, batas)
  - Konfigurasi cache (tipe, TTL, ukuran)
  - Pengaturan web server (port, autentikasi)
  - Konfigurasi logging (level, format, output)
  - Kebijakan retensi data
  - Path file (data, log, PID)
- **Apa yang TIDAK Bisa Anda Ubah**: Data runtime (topic, subscription, pesan, status device)
- **Bagaimana**: Web UI dengan editor YAML atau antarmuka berbasis form
- **Keamanan**: Perubahan divalidasi sebelum disimpan, backup dibuat otomatis

**Runner Mode** (Mode Produksi):
- **Tujuan**: Operasi produksi normal dengan embedded MQTT broker aktif
- **Kapan**: Operasi sehari-hari setelah konfigurasi awal
- **Apa yang Terjadi**:
  - `config.yaml` bersifat **read-only** (tidak bisa diedit)
  - Embedded MQTT broker berjalan dan menerima koneksi
  - Web UI menampilkan dashboard monitoring (config view read-only)
  - Semua operasi runtime menggunakan database, bukan config.yaml
  - Tidak bisa beralih ke setup mode tanpa menghentikan proses
- **Apa yang Bisa Anda Lakukan**: Monitoring sistem, melihat log, mengecek metrik, restart layanan
- **Apa yang TIDAK Bisa Anda Lakukan**: Edit konfigurasi infrastruktur

**3. Pemisahan Data yang Jelas**

**config.yaml** (Konfigurasi Infrastruktur - Statis):
- **Tujuan**: Mendefinisikan BAGAIMANA gateway berjalan (pengaturan infrastruktur)
- **Lokasi**: Sistem file (direktori yang sama dengan executable)
- **Format**: YAML (file teks yang mudah dibaca)
- **Editing**: Hanya mode setup
- **Perubahan Memerlukan**: Restart proses
- **Berisi**:
  ```yaml
  # Mode operasi
  mode: runner  # setup | runner

  # Konfigurasi database
  database:
    type: sqlite
    connection:
      path: ./data/gateway.db
      wal_mode: true
      max_open_conns: 25
    cache:
      type: memory
      ttl: 3600

  # Konfigurasi broker MQTT
  mqtt:
    type: embedded
    embedded:
      listen_address: :1883
      max_connections: 1000

  # Konfigurasi web server
  web:
    address: :8080
    admin_api:
      enabled: true

  # Konfigurasi logging
  logging:
    level: info
    format: json
  ```

**Database** (Data Runtime - Dinamis):
- **Tujuan**: Menyimpan APA yang diproses gateway (data sebenarnya)
- **Lokasi**: Sistem file (SQLite: ./data/gateway.db)
- **Format**: File database biner
- **Akses**: Mode setup dan runner (via API)
- **Perubahan**: Real-time, tidak perlu restart
- **Berisi**:
  - Registry device (device yang terhubung, metadata mereka)
  - Definisi topic (topic apa yang ada)
  - Aturan subscription (siapa subscribe ke apa)
  - Riwayat pesan (pesan MQTT sebenarnya)
  - Status device (status saat ini dari setiap device)
  - Metrik dan statistik (data performa)
  - Akun pengguna (jika autentikasi diaktifkan)
  - API key dan token (jika berlaku)

**Perbedaan Utama**:

| Aspek | config.yaml | Database |
|--------|-------------|----------|
| **Tujuan** | Pengaturan infrastruktur | Data runtime |
| **Contoh** | "Listen pada port 1883" | "Device X mengirim pesan Y" |
| **Editing** | Hanya mode setup | Kedua mode (via API) |
| **Perubahan Memerlukan** | Restart proses | Tidak perlu restart |
| **Format** | Teks YAML | Database biner |
| **Lokasi** | ./config.yaml | ./data/gateway.db |
| **Backup** | Manual atau mode setup | Otomatis/periodik |
| **Version Control** | Ya (Git-friendly) | Tidak (data biner) |

### Peluang Pasar

Lanskap IoT membutuhkan **MQTT gateway yang sepenuhnya mandiri** dengan:
- **Konfigurasi yang Jelas**: Pemisahan yang sangat jelas antara infrastruktur vs data runtime
- **Teknologi yang Familiar**: UI berbasis React, lebih mudah mencari developer
- **Manajemen State yang Lebih Baik**: TanStack Query untuk sinkronisasi state server yang andal
- **Pengalaman Developer**: Routing type-safe, tooling yang excellent

---

## 2. Model Bisnis

### Proposition Nilai

**Untuk Edge Deployments**:
- Single binary mencakup semua yang dibutuhkan
- Pemisahan jelas: config.yaml (infrastruktur) vs database (data)
- Tidak ada kebingungan tentang apa yang diedit dan kapan
- UI berbasis React = pool talent lebih besar

**Untuk Lingkungan Air-Gapped**:
- Operasi offline sepenuhnya
- Tidak ada ketergantungan eksternal
- Broker MQTT bawaan
- Kontrol penuh atas data

**Untuk Pengembangan/Pengujian**:
- Setup lokal cepat tanpa layanan eksternal
- Mudah di-reset dan rekonfigurasi
- Setup mode untuk perubahan infrastruktur
- Runner mode untuk pengujian produksi

### Aliran Pendapatan

1. **Free Self-Hosted**: Single binary dengan embedded MQTT (Lisensi MIT)
2. **Lisensi Pro** ($299 sekali): Adapter database eksternal, fitur lanjutan
3. **Lisensi Enterprise** ($1,499/tahun): Dukungan prioritas, build kustom
4. **Edge appliances**: Bundel perangkat keras + perangkat lunak yang sudah dikonfigurasi

### Strategi Harga

| Edisi | Target Pasar | Model Harga | Fitur |
|---------|--------------|---------------|----------|
| Komunitas | Hobiwan, edge computing | Gratis (MIT) | Embedded MQTT, SQLite WAL, cache in-memory, setup mode, UI React |
| Profesional | UKM, produksi | $299 sekali | Adapter DB eksternal, cache persisten, akses API, dukungan prioritas |
| Enterprise | Organisasi besar | $1,499/tahun | Multiple gateway, clustering, dukungan prioritas, build kustom |

---

## 3. Target Audiens

### Pengguna Utama

**1. Insinyur Edge Computing**
- Deploy IoT gateway pada perangkat edge
- Butuh footprint minimal
- Memerlukan operasi offline
- Titik nyeri: Tidak jelas apa yang dikonfigurasi di mana

**2. Operator Lingkungan Air-Gapped**
- IoT industri dalam jaringan terisolasi
- Tidak dapat mengakses layanan eksternal
- Butuh kontrol penuh
- Titik nyeri: Kebingungan antara config dan data

**3. Tim Pengembangan**
- Butuh lingkungan pengujian lokal
- Ingin setup/teardown cepat
- Lebih prefer ekosistem React
- Titik nyeri: Batas konfigurasi yang tidak jelas

**4. Integrator Sistem**
- Deploy solusi di lokasi pelanggan
- Ingin operasi yang disederhanakan
- Butuh rekonfigurasi mudah
- Titik nyeri: Tidak tahu apa yang memerlukan restart

---

## 4. Pernyataan Masalah

### Masalah Inti

**Masalah 1: Tujuan Mode Tidak Jelas**
- v2 tidak menjelaskan dengan jelas apa yang dilakukan Setup Mode
- Pengguna bingung: "Bisakah saya mengubah semua di setup mode?"
- Pendekatan saat ini: Deskripsi "perubahan konfigurasi" yang vagu
- Solusi: Definisi jelas: Setup Mode = editor config.yaml SAJA

**Masalah 2: Kebingungan Tentang Penempatan Data**
- Tidak jelas apa yang masuk ke config.yaml vs database
- Pengguna bertanya: "Di mana saya menyimpan metadata device?"
- Pendekatan saat ini: Pemisahan data yang ambigu
- Solusi: Aturan yang sangat jelas: Infrastruktur (config.yaml) vs Data Runtime (database)

**Masalah 3: Familiaritas Framework**
- SvelteKit kurang familiar daripada React
- Pool talent lebih kecil
- Pendekatan saat ini: SvelteKit untuk UI web
- Solusi: React + Vite dengan ekosistem TanStack

**Masalah 4: Manajemen State**
- Tidak ada manajemen state server yang jelas
- Sinkronisasi manual antara UI dan backend
- Pendekatan saat ini: Panggilan API dasar
- Solusi: TanStack Query untuk caching otomatis, sinkronisasi, retry

**Masalah 5: Kompleksitas Routing**
- Masalah keamanan tipe dengan routing
- Penanganan parameter manual
- Pendekatan saat ini: Routing dasar
- Solusi: TanStack Router untuk routing type-safe

### Dampak

- **Error Konfigurasi**: Pengguna mengedit file/tempat yang salah
- **Restart yang Tidak Perlu**: Mengubah data database dengan berpikir perlu restart
- **Ketersediaan Developer**: Lebih sulit mencari developer SvelteKit
- **Inkonsistensi State**: UI menampilkan data yang basi
- **Pengalaman Developer yang Buruk**: Pengecekan tipe manual, bug parameter route

---

## 5. Solusi yang Diusulkan

### Visi: Konfigurasi yang Sangat Jelas dengan Stack React Modern

**Single executable binary** dengan:
- **React + Vite UI**: Stack teknologi modern dan familiar
- **TanStack Query**: Manajemen state server otomatis, caching, sinkronisasi
- **TanStack Router**: Routing type-safe dengan DX yang excellent
- **Definisi Mode yang Sangat Jelas**: Setup Mode = editor config.yaml, Runner Mode = produksi
- **Pemisahan Data yang Jelas**: Infrastruktur (config.yaml) vs Data Runtime (database)

### Inovasi Utama

**1. Stack React + Vite + TanStack**

```bash
# Frontend stack
- React 18+ (library UI)
- Vite 5+ (build tool)
- TypeScript (keamanan tipe)
- TanStack Query (manajemen state server)
- TanStack Router (routing)
- TailwindCSS (styling)

# Backend stack
- Go 1.24+
- mochi-mqtt (embedded MQTT broker)
- SQLite WAL mode (database default)
- go-cache (in-memory cache)
```

**2. Definisi Mode yang Sangat Jelas**

**Alur Setup Mode**:
```
1. Pengguna memulai gateway dengan --mode=setup ATAU first run terdeteksi
2. Web UI diluncurkan dalam setup mode di http://localhost:8080/setup
3. Pengguna melihat editor config.yaml dengan validasi
4. Pengguna mengubah pengaturan infrastruktur (database, MQTT, cache, dll)
5. Pengguna klik "Save & Restart"
6. Sistem memvalidasi config.yaml
7. Sistem membuat backup dari config.yaml lama
8. Sistem menulis config.yaml baru
9. Sistem me-restart dan beralih ke runner mode
```

**Alur Runner Mode**:
```
1. Pengguna memulai gateway normal (atau setelah setup)
2. config.yaml dimuat dan menjadi READ-ONLY
3. Embedded MQTT broker mulai pada port yang dikonfigurasi
4. Web UI diluncurkan dalam runner mode di http://localhost:8080
5. Pengguna melihat dashboard monitoring
6. Pengguna dapat melihat config.yaml (read-only)
7. Pengguna TIDAK dapat mengedit config.yaml
8. Pengguna berinteraksi dengan data runtime via database (panggilan API)
```

**3. Aturan Pemisahan Data yang Jelas**

**Aturan 1: Pengaturan Infrastruktur → config.yaml**

Tanya: "Apakah mengubah ini memerlukan restart gateway?"
- YA → Masuk ke config.yaml
- TIDAK → Masuk ke database

Contoh:
- ✅ "Database mana yang digunakan?" → config.yaml (memerlukan restart)
- ✅ "Port mana untuk didengar?" → config.yaml (memerlukan restart)
- ✅ "Berapa banyak item cache?" → config.yaml (memerlukan restart)
- ❌ "Topic apa yang ada?" → database (tidak perlu restart)
- ❌ "Device mana yang terhubung?" → database (tidak perlu restart)

**Aturan 2: Data Runtime → Database**

Tanya: "Apakah ini data yang berubah selama operasi normal?"
- YA → Masuk ke database
- TIDAK → Masuk ke config.yaml

Contoh:
- ✅ "Pesan MQTT" → database (berubah terus-menerus)
- ✅ "Status device" → database (berubah selama operasi)
- ✅ "Subscription topic" → database (berubah selama operasi)
- ❌ "Path file database" → config.yaml (pengaturan statis)
- ❌ "Level log" → config.yaml (pengaturan statis)

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
│  │  │ Embedded  │  │ HTTP      │  │ In-Memory │  │  │
│  │  │ MQTT      │  │ API        │  │ Cache      │  │  │
│  │  │ Broker     │  │ Server     │  │            │  │  │
│  │  │ (mochi-mqtt)│  │            │  │            │  │  │
│  │  └────────────┘  └────────────┘  └────────────┘  │  │
│  │                                                        │  │
│  │  ┌────────────┐  ┌────────────┐                    │  │
│  │  │ YAML       │  │ SQLite     │                    │  │
│  │  │ Config     │  │ Database   │                    │  │
│  │  │ Loader     │  │ (WAL Mode) │                    │  │
│  │  └────────────┘  └────────────┘                    │  │
│  │                                                        │  │
│  │  ┌────────────┐  ┌────────────┐                    │  │
│  │  │ Mode       │  │ Config     │                    │  │
│  │  │ Manager    │  │ Validator  │                    │  │
│  │  └────────────┘  └────────────┘                    │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                               │
│  ┌───────────────────────────────────────────────────────┐  │
│  │         React + Vite Web UI (Embedded SPA)            │  │
│  │                                                        │  │
│  │  ┌────────────────────────────────────────────────┐  │  │
│  │  │ Setup Mode UI (/setup)                          │  │  │
│  │  │ - YAML editor dengan validasi                  │  │  │
│  │  │ - Antarmuka konfigurasi berbasis form          │  │  │
│  │  │ - Tombol test koneksi                          │  │  │
│  │  │ - Tombol Save & Restart                        │  │  │
│  │  └────────────────────────────────────────────────┘  │  │
│  │                                                        │  │
│  │  ┌────────────────────────────────────────────────┐  │  │
│  │  │ Runner Mode UI (/)                              │  │  │
│  │  │ - Dashboard monitoring                          │  │  │
│  │  │ - Viewer config read-only                       │  │  │
│  │  │ - Metrik real-time (TanStack Query)             │  │  │
│  │  │ - Viewer log                                    │  │  │
│  │  └────────────────────────────────────────────────┘  │  │
│  │                                                        │  │
│  │  ┌────────────────────────────────────────────────┐  │  │
│  │  │ TanStack Router (Routing Type-Safe)             │  │  │
│  │  │ TanStack Query (Manajemen State Server)         │  │  │
│  │  └────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

---

## 6. Tujuan Proyek

### Tujuan Utama (Must Have)

**O1: Implementasi Frontend React + Vite**
- Mengganti SvelteKit dengan React + Vite
- TypeScript untuk keamanan tipe
- TailwindCSS untuk styling
- Single-page application yang tertanam dalam binary

**O2: Integrasi TanStack Query**
- Manajemen state server
- Caching dan revalidasi otomatis
- Refetch di latar belakang
- Optimistic updates
- Penanganan error dan retry

**O3: Integrasi TanStack Router**
- Routing type-safe
- Parameter route dengan inferensi tipe
- Nested route dan layout
- Code splitting
- Penanganan parameter search

**O4: Perjelas Tujuan Setup Mode**
- Setup Mode = editor config.yaml SAJA
- Pesan yang jelas: "Konfigurasikan pengaturan infrastruktur"
- Editor YAML dengan validasi
- Antarmuka alternatif berbasis form
- Test konfigurasi sebelum menyimpan

**O5: Perjelas Tujuan Runner Mode**
- Runner Mode = operasi produksi
- config.yaml adalah READ-ONLY
- Dashboard monitoring
- Viewer config read-only
- Tidak dapat mengedit pengaturan infrastruktur

**O6: Dokumentasikan Pemisahan Data dengan Jelas**
- Infrastruktur → config.yaml
- Data runtime → database
- Pohon keputusan: "Apakah ini memerlukan restart?"
- Contoh dan aturan
- Diagram visual

### Tujuan Sekunder (Should Have)

**O7: Validasi Konfigurasi yang Ditingkatkan**
- Validasi YAML real-time
- Testing koneksi
- Validasi skema
- Pesan error dengan saran
- Viewer diff konfigurasi

**O8: Penanganan Error yang Lebih Baik**
- Pesan error yang ramah pengguna
- Saran pemulihan error
- Detail error konfigurasi
- Integrasi file log
- Laporan error

**O9: Pengalaman Developer yang Lebih Baik**
- Hot reload dalam pengembangan
- Mode TypeScript strict
- ESLint + Prettier
- Storybook untuk komponen
- Dokumentasi komprehensif

### Tujuan Tersier (Nice to Have)

**O10: Import/Ekspor Konfigurasi**
- Ekspor config.yaml
- Impor dari file
- Template konfigurasi
- Konfigurasi preset

**O11: Riwayat Konfigurasi**
- Track perubahan config.yaml
- Rollback ke versi sebelumnya
- Viewer perubahan diff
- Log audit

---

## 7. Ruang Lingkup & Batasan Proyek

### Dalam Ruang Lingkup (Apa yang Akan Kami Bangun)

#### Fungsionalitas Inti

1. **Frontend React + Vite**
   - SPA yang tertanam dalam binary Go
   - TypeScript untuk keamanan tipe
   - TailwindCSS untuk styling
   - Desain responsif

2. **Integrasi TanStack Query**
   - Query hooks untuk semua panggilan API
   - Mutation hooks untuk perubahan data
   - Caching dan revalidasi otomatis
   - Status error dan loading

3. **Integrasi TanStack Router**
   - Route type-safe
   - Proteksi route (setup vs runner)
   - Layout nested
   - Parameter search

4. **UI Setup Mode**
   - Editor YAML dengan syntax highlighting
   - Antarmuka konfigurasi berbasis form
   - Validasi real-time
   - Tombol test konfigurasi
   - Tombol Save & Restart
   - Backup konfigurasi

5. **UI Runner Mode**
   - Dashboard monitoring
   - Viewer config read-only
   - Metrik real-time
   - Viewer log
   - Manajemen device

6. **Dokumentasi yang Jelas**
   - Aturan pemisahan data
   - Definisi mode
   - Contoh konfigurasi
   - Pohon keputusan

### Di Luar Ruang Lingkup (Apa yang Tidak Akan Kami Bangun - Awalnya)

1. **Fitur Lanjutan** (Fase 2)
   - Riwayat/versioning konfigurasi
   - Import/ekspor konfigurasi
   - Monitoring dan alerting lanjutan

2. **High Availability** (Fase 3)
   - Clustering multi-gateway
   - Replikasi data
   - Failover otomatis

### Batasan & Kendala

#### Kendala Teknis
- **Frontend**: React 18+, Vite 5+, TypeScript 5+
- **Backend**: Go 1.24+
- **Library MQTT**: mochi-mqtt
- **Hardware Minimal**: 1 core CPU, 1GB RAM

#### Kendala Waktu
- **Fase 1 (MVP)**: 10 bulan
- **Fase 2 (Fitur Lanjutan)**: +3 bulan
- **Fase 3 (Fitur HA)**: +6 bulan

#### Kendala Anggaran
- **Tim Pengembangan**: 2-3 developer full-time
- **Infrastruktur**: $1,500/bulan untuk pengembangan/pengujian
- **Layanan Pihak Ketiga**: Tier gratis awalnya

---

## 8. Teknologi Stack

### Stack Frontend

| Komponen | Teknologi | Versi | Rasional |
|-----------|-----------|-------|-----------|
| Library UI | React | 18.3+ | Ekosistem terbesar, familiar bagi sebagian besar developer |
| Build Tool | Vite | 5.0+ | HMR cepat, build teroptimasi, DX hebat |
| Bahasa | TypeScript | 5.3+ | Keamanan tipe, pengalaman developer lebih baik |
| Manajemen State | TanStack Query | 5.0+ | Manajemen state server yang kuat, caching |
| Routing | TanStack Router | 1.0+ | Routing type-safe, DX excellent |
| Styling | TailwindCSS | 3.4+ | Utility-first, pengembangan cepat |
| Kualitas Kode | ESLint, Prettier | Latest | Pemformatan kode konsisten, linting |

### Stack Backend

| Komponen | Teknologi | Versi | Rasional |
|-----------|-----------|-------|-----------|
| Bahasa | Go | 1.24+ | Performa, konkurensi, single binary |
| MQTT Broker | mochi-mqtt | 2.0+ | Embedded Go MQTT broker |
| Database | SQLite | 3.40+ | Dukungan mode WAL, single file |
| Cache | go-cache | Latest | In-memory caching dengan TTL |
| HTTP Router | Chi | 5.0+ | Ringan, idiomatic Go |
| YAML | yaml.v3 | Latest | Parsing dan serialisasi YAML |

### Integrasi TanStack Query

```typescript
// Contoh: Query hook untuk mengambil metrik
import { useQuery } from '@tanstack/react-query';

function useMetrics() {
  return useQuery({
    queryKey: ['metrics'],
    queryFn: async () => {
      const response = await fetch('/api/metrics');
      if (!response.ok) throw new Error('Failed to fetch metrics');
      return response.json();
    },
    refetchInterval: 5000, // Auto-refresh setiap 5 detik
  });
}

// Penggunaan dalam komponen
function MetricsDashboard() {
  const { data, error, isLoading } = useMetrics();

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div>
      <h2>Metrik</h2>
      <p>Klien Terhubung: {data.connectedClients}</p>
      <p>Pesan/detik: {data.messagesPerSec}</p>
    </div>
  );
}
```

### Integrasi TanStack Router

```typescript
// Contoh: Route type-safe
import { createRootRoute, createRoute, createRouter } from '@tanstack/react-router';

// Route setup mode (hanya dapat diakses dalam setup mode)
const setupRoute = createRoute({
  path: '/setup',
  component: SetupMode,
  beforeLoad: async ({ location }) => {
    // Periksa apakah dalam setup mode
    const mode = await fetch('/api/mode').then(r => r.json());
    if (mode !== 'setup') {
      throw redirect({ to: '/' });
    }
  },
});

// Route runner mode (hanya dapat diakses dalam runner mode)
const rootRoute = createRootRoute({
  component: RunnerMode,
});

// Konfigurasi router
const router = createRouter({
  routeTree: rootRoute.addChildren([setupRoute]),
});

// Navigasi type-safe
function NavigateToSetup() {
  const navigate = useNavigate();
  return (
    <button onClick={() => navigate({ to: '/setup' })}>
      Go to Setup
    </button>
  );
}
```

### Struktur Konfigurasi

#### config.yaml (Infrastruktur - Statis)

```yaml
# =============================================================================
# File Konfigurasi UM-Gateway
# =============================================================================
# File ini berisi PENGATURAN INFRASTRUKTUR SAJA
# Pengaturan ini mendefinisikan BAGAIMANA gateway berjalan
# Perubahan pada file ini memerlukan restart gateway
#
# JANGAN simpan data runtime di sini (gunakan database sebagai gantinya)
# =============================================================================

# Mode Operasi
# - setup: Mode konfigurasi (edit file ini via web UI)
# - runner: Mode produksi (file ini menjadi read-only)
mode: runner

# Konfigurasi Database
database:
  # Tipe database: sqlite, postgres, mongodb, influxdb
  type: sqlite

  connection:
    # Pengaturan spesifik SQLite
    path: ./data/gateway.db
    wal_mode: true
    max_open_conns: 25
    max_idle_conns: 5
    conn_max_lifetime: 300s

    # Konfigurasi cache
    cache:
      # Tipe cache: memory, redis, memcached
      type: memory
      ttl: 3600
      max_size: 10000

# Konfigurasi Broker MQTT
mqtt:
  # Tipe broker: embedded, external
  type: embedded

  embedded:
    listen_address: :1883
    max_connections: 1000
    max_message_size: 256KB
    keepalive: 60s

  external:
    url: tcp://localhost:1883
    client_id: um-gateway
    username: ""
    password: ""

# Konfigurasi Path
paths:
  data: ./data
  logs: ./logs
  pid: ./um-gateway.pid

# Konfigurasi Web Server
web:
  address: :8080
  admin_api:
    enabled: true
    authentication:
      type: basic

# Kebijakan Retensi Data
retention:
  enabled: true
  default_ttl: 2592000  # 30 hari
  cleanup_interval: 3600

# Konfigurasi Logging
logging:
  level: info
  format: json
  output: stdout
```

#### Skema Database (Data Runtime - Dinamis)

```sql
-- Tabel device (data runtime: device mana yang ada)
CREATE TABLE devices (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  type TEXT NOT NULL,
  metadata TEXT, -- JSON
  first_seen DATETIME DEFAULT CURRENT_TIMESTAMP,
  last_seen DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Tabel topics (data runtime: topic apa yang ada)
CREATE TABLE topics (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  description TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Tabel subscriptions (data runtime: siapa subscribe ke apa)
CREATE TABLE subscriptions (
  id TEXT PRIMARY KEY,
  device_id TEXT NOT NULL,
  topic_id TEXT NOT NULL,
  qos INTEGER DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (device_id) REFERENCES devices(id),
  FOREIGN KEY (topic_id) REFERENCES topics(id)
);

-- Tabel messages (data runtime: pesan MQTT sebenarnya)
CREATE TABLE messages (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  topic TEXT NOT NULL,
  payload TEXT NOT NULL,
  qos INTEGER DEFAULT 0,
  retained BOOLEAN DEFAULT FALSE,
  published_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Tabel device_states (data runtime: status saat ini)
CREATE TABLE device_states (
  device_id TEXT PRIMARY KEY,
  status TEXT NOT NULL,
  last_payload TEXT,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (device_id) REFERENCES devices(id)
);
```

---

## 9. Fitur

### 9.1 Fitur Inti

#### F1: Frontend React + Vite

**Struktur Proyek**:

```
frontend/
├── src/
│   ├── main.tsx           # Entry point
│   ├── App.tsx            # Root component
│   ├── routes/            # Definisi route
│   │   ├── __root.tsx     # Root layout
│   │   ├── index.tsx      # Route runner mode
│   │   └── setup.tsx      # Route setup mode
│   ├── pages/             # Komponen halaman
│   │   ├── SetupMode.tsx  # Halaman setup mode
│   │   ├── RunnerMode.tsx # Halaman runner mode
│   │   └── Dashboard.tsx  # Halaman dashboard
│   ├── components/        # Komponen yang dapat digunakan kembali
│   │   ├── YAMLEditor.tsx
│   │   ├── MetricsCard.tsx
│   │   └── LogViewer.tsx
│   ├── hooks/             # Custom hooks
│   │   ├── useMetrics.ts
│   │   ├── useConfig.ts
│   │   └── useMode.ts
│   ├── lib/               # Utilitas
│   │   ├── api.ts         # Klien API
│   │   └── query.ts       # Konfigurasi TanStack Query
│   └── styles/            # Styles
│       └── index.css      # Tailwind CSS
├── index.html
├── vite.config.ts
├── tsconfig.json
└── package.json
```

**Konfigurasi TanStack Query**:

```typescript
// frontend/src/lib/query.ts
import { QueryClient } from '@tanstack/react-query';

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5000,        // Data segar selama 5 detik
      gcTime: 1000 * 60 * 5, // Cache selama 5 menit
      retry: 3,              // Retry permintaan yang gagal 3 kali
      refetchOnWindowFocus: true,
    },
    mutations: {
      retry: 1,
    },
  },
});
```

---

#### F2: UI Setup Mode

**Tujuan**: Edit pengaturan infrastruktur config.yaml

**Komponen**:

```typescript
// frontend/src/pages/SetupMode.tsx
import { useMutation, useQuery } from '@tanstack/react-query';

export function SetupMode() {
  // Fetch config saat ini
  const { data: config, isLoading } = useQuery({
    queryKey: ['config'],
    queryFn: () => fetch('/api/config').then(r => r.json()),
  });

  // Simpan mutation config
  const saveConfig = useMutation({
    mutationFn: (newConfig: string) =>
      fetch('/api/config', {
        method: 'POST',
        body: newConfig,
        headers: { 'Content-Type': 'text/yaml' },
      }),
    onSuccess: () => {
      alert('Config tersimpan! Me-restart gateway...');
      setTimeout(() => window.location.reload(), 2000);
    },
  });

  if (isLoading) return <div>Loading configuration...</div>;

  return (
    <div className="setup-mode">
      <h1>Setup Mode - Konfigurasi Infrastruktur</h1>
      <p className="warning">
        ⚠️ Anda mengedit PENGATURAN INFRASTRUKTUR (config.yaml)
        <br />
        Perubahan akan memerlukan restart gateway
      </p>

      <YAMLEditor
        value={config.yaml}
        onChange={(newValue) => setConfig(newValue)}
      />

      <div className="actions">
        <button onClick={() => saveConfig.mutate(config)}>
          Simpan & Restart
        </button>
        <button onClick={() => testConfig.mutate(config)}>
          Test Konfigurasi
        </button>
      </div>
    </div>
  );
}
```

---

#### F3: UI Runner Mode

**Tujuan**: Monitoring operasi produksi

**Komponen**:

```typescript
// frontend/src/pages/RunnerMode.tsx
import { useQuery } from '@tanstack/react-query';

export function RunnerMode() {
  // Fetch metrik (auto-refresh setiap 5 detik)
  const { data: metrics } = useQuery({
    queryKey: ['metrics'],
    queryFn: () => fetch('/api/metrics').then(r => r.json()),
    refetchInterval: 5000,
  });

  // Fetch config (read-only)
  const { data: config } = useQuery({
    queryKey: ['config'],
    queryFn: () => fetch('/api/config').then(r => r.json()),
  });

  return (
    <div className="runner-mode">
      <h1>Runner Mode - Dashboard Produksi</h1>
      <p className="info">
        ℹ️ Gateway berjalan dalam mode PRODUKSI
        <br />
        Konfigurasi bersifat read-only. Untuk membuat perubahan, beralih ke setup mode.
      </p>

      <MetricsCard metrics={metrics} />
      <ConfigViewer config={config} />
      <LogViewer />
    </div>
  );
}
```

---

#### F4: Dokumentasi Pemisahan Data yang Jelas

**Pohon Keputusan**:

```
┌─────────────────────────────────────────────────────────────┐
│  Data ini harus ditaruh di mana?                            │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  Q1: Apakah mengubah ini memerlukan restart gateway?         │
│                                                              │
│     YA → config.yaml (Pengaturan Infrastruktur)              │
│     ┌─────────────────────────────────────────────────┐     │
│     │ Contoh:                                          │     │
│     │ • Tipe database dan pengaturan koneksi          │     │
│     │ • Pengaturan broker MQTT (port, batas)          │     │
│     │ • Konfigurasi cache                             │     │
│     │ • Pengaturan web server                         │     │
│     │ • Konfigurasi logging                           │     │
│     │ • Path file                                     │     │
│     └─────────────────────────────────────────────────┘     │
│                                                              │
│     TIDAK → Database (Data Runtime)                          │
│     ┌─────────────────────────────────────────────────┐     │
│     │ Contoh:                                          │     │
│     │ • Registry device                                │     │
│     │ • Topics dan subscriptions                       │     │
│     │ • Pesan MQTT                                     │     │
│     │ • Status device                                  │     │
│     │ • Metrik dan statistik                           │     │
│     │ • Akun pengguna                                  │     │
│     └─────────────────────────────────────────────────┘     │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## 10. Peta Jalan Pengembangan

### Fase 1: Pondasi (Bulan 1-5)

**Sprint 1-2: Setup Proyek & Arsitektur**
- [ ] Struktur monorepo (Go + React + Vite)
- [ ] Setup frontend React + Vite
- [ ] Konfigurasi TanStack Query
- [ ] Setup TanStack Router
- [ ] Konfigurasi TypeScript
- [ ] Pipeline CI/CD

**Sprint 3-4: Backend Inti**
- [ ] Loading konfigurasi YAML
- [ ] Deteksi dan perpindahan mode
- [ ] Integrasi embedded MQTT broker
- [ ] Setup mode WAL SQLite
- [ ] Endpoint API untuk config dan metrik

**Sprint 5-6: Frontend Inti**
- [ ] Skeleton UI setup mode
- [ ] Skeleton UI runner mode
- [ ] Hooks TanStack Query
- [ ] Route TanStack Router
- [ ] Layout dasar

**Sprint 7-8: Implementasi Setup Mode**
- [ ] Komponen editor YAML
- [ ] Validasi konfigurasi
- [ ] Tombol test konfigurasi
- [ ] Fungsionalitas simpan & restart
- [ ] Backup konfigurasi

**Sprint 9-10: Implementasi Runner Mode**
- [ ] Dashboard monitoring
- [ ] Viewer config read-only
- [ ] Metrik real-time
- [ ] Viewer log
- [ ] Auto-refresh dengan TanStack Query

**Milestone**: Alpha - Fungsionalitas dasar selesai

---

### Fase 2: Poles & Fitur (Bulan 6-8)

**Sprint 11-12: UI yang Ditingkatkan**
- [ ] Desain responsif
- [ ] Penanganan error
- [ ] Status loading
- [ ] Notifikasi toast
- [ ] Dialog modal

**Sprint 13-14: Dokumentasi**
- [ ] Panduan pemisahan data
- [ ] Definisi mode
- [ ] Contoh konfigurasi
- [ ] Pohon keputusan
- [ ] Dokumentasi API

**Sprint 15-16: Pengujian**
- [ ] Unit test (komponen React)
- [ ] Integration test (API)
- [ ] E2E test (Playwright)
- [ ] Load testing
- [ ] Audit keamanan

**Milestone**: Beta - v3 lengkap fitur

---

### Fase 3: Produksi (Bulan 9-10)

**Sprint 17-18: Hardening**
- [ ] Implementasi error boundary
- [ ] Degradasi gracefully
- [ ] Optimasi performa
- [ ] Optimasi ukuran bundle
- [ ] Peningkatan aksesibilitas

**Sprint 19-20: Deploy**
- [ ] Build produksi
- [ ] Kompilasi cross-platform
- [ ] Panduan instalasi
- [ ] Dokumentasi pengguna
- [ | Catatan rilis

**Milestone**: v3.0 Ketersediaan Umum

---

## 11. Metrik Kesuksesan

### Metrik Teknis

| Metrik | Target | Pengukuran |
|--------|--------|-------------|
| Ukuran Bundle Frontend | < 500 KB (gzipped) | Artefak build |
| Load Halaman Awal | < 2 detik | Lighthouse |
| Time to Interactive | < 3 detik | Lighthouse |
| Waktu Respon API | < 100ms | Pengujian otomatis |
| Tingkat Cache Hit TanStack Query | > 90% | Statistik cache |
| Cakupan TypeScript | 100% | TSC --noEmit |

### Metrik Bisnis

| Metrik | Target | Timeline |
|--------|--------|----------|
| Instalasi Aktif | 1,000 | 6 bulan pasca-rilis |
| GitHub Stars | 2,500 | 6 bulan pasca-rilis |
| Durasi Sesi Rata-rata | > 10 menit | Analytics |
| Tingkat Penyelesaian Setup | > 95% | Analytics |

---

## 12. Analisis Risiko

### Risiko Teknis

| Risiko | Dampak | Probabilitas | Mitigasi |
|------|--------|-------------|------------|
| Kompleksitas TanStack Query | Sedang | Rendah | Dokumentasi komprehensif, contoh |
| Kurva belajar TanStack Router | Sedang | Rendah | Pelatihan, contoh kode |
| Ukuran bundle React | Rendah | Sedang | Code splitting, lazy loading |
| Masalah konfigurasi TypeScript | Rendah | Rendah | Konfigurasi strict dari awal |

### Risiko Bisnis

| Risiko | Dampak | Probabilitas | Mitigasi |
|------|--------|-------------|------------|
| Pengguna lebih prefer SvelteKit | Rendah | Rendah | React memiliki ekosistem lebih besar |
| Masalah adopsi TanStack | Rendah | Rendah | Library yang sudah terbukti |
| Kelelahan framework | Sedang | Sedang | Teknologi yang stabil dan terbukti |

---

## 13. Kebutuhan Sumber Daya

### Struktur Tim

**Tim Inti (Fase 1-3)**
- 2x Full-Stack Developers (Go + React)
- 1x Spesialis Frontend (React + TanStack)
- 1x DevOps Engineer
- 1x QA Engineer (part-time)
- 1x Technical Writer (part-time)

**Estimasi Anggaran**

| Kategori | Biaya (Bulanan) | Durasi |
|----------|----------------|---------|
| Gaji Tim Pengembangan | $20,000 | 10 bulan |
| Infrastruktur (Dev/Test) | $1,500 | 10 bulan |
| Alat & Layanan | $400 | Berkelanjutan |
| Dokumentasi & Desain | $1,000 | 8 bulan |
| Kontinjensi (15%) | $2,865 | - |
| **Total Fase 1-3** | **$21,500/bulan × 10** | **$215,500** |

---

## 14. Perbandingan: v2 vs v3

| Aspek | v2 (SvelteKit) | v3 (React + TanStack) |
|--------|----------------|----------------------|
| **Framework Frontend** | SvelteKit | React + Vite |
| **Manajemen State** | Dasar | TanStack Query |
| **Routing** | router SvelteKit | TanStack Router |
| **Keamanan Tipe** | Baik | Excellent (TanStack Router) |
| **Kemudahan Mode** | Vagu | Sangat jelas |
| **Pemisahan Data** | Tidak jelas | Aturan terdokumentasi |
| **Pool Talent** | Lebih kecil | Lebih besar (React) |
| **Ukuran Bundle** | Lebih kecil | Sedikit lebih besar |
| **Kurva Belajar** | Sedang | Rendah (React populer) |
| **Anggaran** | $238,150 | $215,500 |

---

## 15. Langkah Selanjutnya

### Tindakan Segera (Minggu 1-2)

1. **Riset Teknologi**
   - Evaluasi fitur TanStack Query
   - Uji kemampuan TanStack Router
   - Prototipe integrasi React + Vite

2. **Desain Arsitektur**
   - Definisikan struktur route
   - Rencanakan query hooks
   - Desain hierarki komponen

3. **Formasi Tim**
   - Rekrut developer React
   - Pelatihan TanStack jika diperlukan

### Tindakan Jangka Pendek (Bulan 1)

4. **Pengembangan Prototipe**
   - Skeleton frontend React + Vite
   - Integrasi TanStack Query
   - Setup TanStack Router

5. **Dokumentasi**
   - Panduan pemisahan data
   - Definisi mode
   - Pohon keputusan

---

## 16. Kesimpulan

Versi v3 merepresentasikan **pendekatan yang paling ramah developer dan paling jelas**:

### Keunggulan Utama atas v2

**Technology Stack**
- ✅ React + Vite (ekosistem lebih besar)
- ✅ TanStack Query (manajemen state yang kuat)
- ✅ TanStack Router (routing type-safe)
- ✅ TypeScript (DX excellent)

**Kemudahan**
- ✅ Definisi mode yang sangat jelas
- ✅ Setup Mode = editor config.yaml
- ✅ Runner Mode = operasi produksi
- ✅ Aturan pemisahan data yang jelas
- ✅ Pohon keputusan dan contoh

**Pengalaman Developer**
- ✅ Ekosistem React yang familiar
- ✅ Routing type-safe
- ✅ Sinkronisasi state otomatis
- ✅ Penanganan error yang lebih baik
- ✅ Pool talent lebih besar

**Anggaran Lebih Rendah**
- ✅ $215,500 (vs $238,150 untuk v2)
- ✅ Lebih banyak developer React yang tersedia
- ✅ Biaya pelatihan lebih rendah

Dengan eksekusi yang terfokus selama 10 bulan, kami dapat mengirimkan **v3.0 production-ready** dengan kemudahan dan pengalaman developer yang excellent.

---

**Versi Dokumen**: 3.0
**Terakhir Diperbarui**: 9 Januari 2026
**Status**: Draft - Menunggu Review
**Menggantikan**: PROJECT_BRIEF_EN_v2.md

---

## Lampiran

### A. Contoh Pemisahan Data

**config.yaml berisi** (Infrastruktur):
```yaml
database:
  type: sqlite           # Database mana yang digunakan
  path: ./data/gateway.db  # Di mana menyimpannya

mqtt:
  type: embedded         # Gunakan broker embedded
  embedded:
    listen_address: :1883  # Port mana untuk didengar
```

**Database berisi** (Data Runtime):
```sql
-- Device "temp-sensor-1" ada
INSERT INTO devices (id, name, type) VALUES ('temp-sensor-1', 'Temperature Sensor', 'sensor');

-- Topic "sensors/temperature" ada
INSERT INTO topics (id, name) VALUES ('topic-1', 'sensors/temperature');

-- Device subscribe ke topic
INSERT INTO subscriptions (device_id, topic_id) VALUES ('temp-sensor-1', 'topic-1');

-- Pesan diterbitkan
INSERT INTO messages (topic, payload) VALUES ('sensors/temperature', '{"value": 25.5}');
```

### B. Referensi

1. Dokumentasi TanStack Query
2. Dokumentasi TanStack Router
3. Dokumentasi React
4. Dokumentasi Vite
5. Dokumentasi TypeScript

---

**Disiapkan oleh**: Tim Pengembangan
**Disetujui oleh**: [Pending]
**Tanggal Review**: [Pending]
