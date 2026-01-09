# Rencana Proyek: Universal MQTT Gateway Server (UM-Gateway) MVP v1.0

## Ringkasan Eksekutif

**Nama Proyek**: Universal MQTT Gateway Server (UM-Gateway)
**Versi**: MVP v1.0 (Minimum Viable Product dengan Internal Database)
**Status**: Perencanaan
**Target Rilis**: Q2 2026 (3 bulan pengembangan)
**Organisasi**: Divisi Teknologi Rafflesia Agro

---

## Riwayat Versi

| Versi | Tanggal | Perubahan | Penulis |
|-------|---------|-----------|--------|
| **MVP v1.0** | 2026-01-09 | **MVP dengan Internal Database**: Menggantikan config.yaml dengan internal.db<br>• Single executable dengan embedded MQTT<br>• Dashboard admin dasar (React + Vite)<br>• Setup mode & Runner mode<br>• Wizard onboarding berbasis form<br>• Internal.db untuk konfigurasi aplikasi (seperti PocketBase)<br>• Database.db untuk data sistem MQTT (default embedded SQLite)<br>• Adapter pattern untuk database eksternal (post-MVP)<br>• Cache lokal (in-memory)<br>• Broker MQTT embedded lokal (mochi-mqtt)<br>• Tanpa file konfigurasi YAML<br>• Riwayat konfigurasi & rollback<br>• Migrasi otomatis<br>• Fitur esensial saja (tanpa fitur SaaS)<br>• Waktu pengembangan: 3 bulan<br>• Anggaran: $68,000 | Tim Pengembangan |

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
- ✅ Tanpa file konfigurasi manual

### Perubahan Utama dari MVP v0.1

**v0.1**: Menggunakan config.yaml untuk konfigurasi
**v1.0**: Menggunakan internal.db (embedded SQLite) untuk konfigurasi

### Kenapa Internal Database?

**Masalah dengan config.yaml**:
1. Pengguna harus mengedit file YAML secara manual
2. Validasi YAML bisa rumit
3. Tidak ada riwayat perubahan
4. Migrasi antar versi sulit
5. Backup terpisah dari database

**Solusi dengan internal.db**:
1. Semua konfigurasi via Web UI / API
2. Validasi di level aplikasi
3. Riwayat lengkap dengan rollback
4. Migrasi otomatis dengan skema database
5. Backup terpadu

---

## 2. Model Bisnis

### Proposition Nilai

**Untuk Deployment Kecil**:
- Single binary, tanpa setup kompleks
- Tidak perlu mengedit file konfigurasi
- Berjalan di komputer apa pun
- Gratis dan open source

**Untuk Pengembangan**:
- Pengujian lokal cepat
- Konfigurasi mudah (form-based UI)
- Broker MQTT bawaan
- Manajemen berbasis web

### Model Pendapatan

**Fase MVP (6 bulan pertama)**:
- **100% Gratis**: Open source (Lisensi MIT)
- Fokus pada adopsi dan umpan balik
- Belum ada fitur berbayar

**Fase Masa Depan** (Setelah MVP terbukti bernilai):
- Lisensi Pro untuk fitur lanjutan
- Dukungan database eksternal
- Kontrak dukungan

---

## 3. Target Audiens

### Pengguna Utama

**1. Proyek IoT Kecil**
- Otomatisasi smart home (1-50 perangkat)
- Monitoring peternakan kecil
- Proyek hobi
- Proyek mahasiswa

**2. Tim Pengembangan**
- Pengujian aplikasi IoT
- Prototyping solusi
- Lingkungan pengembangan lokal

**3. Integrator Sistem**
- Deployment pelanggan kecil
- Proyek proof of concept
- Kebutuhan gateway sederhana

---

## 4. Pernyataan Masalah

### Masalah Inti

Deployment IoT kecil membutuhkan **MQTT gateway yang sederhana dan mandiri** tetapi solusi yang ada:
1. **Terlalu Kompleks**: Memerlukan broker, database, cache terpisah
2. **Sulit Dikonfigurasi**: Memerlukan keahlian teknis / editing file
3. **Mahal**: Solusi komersial terlalu mahal
4. **Over-engineered**: Platform SaaS lengkap untuk kebutuhan sederhana

### Solusi MVP v1.0

**Single executable** dengan:
- ✅ Embedded MQTT broker (tanpa broker terpisah)
- ✅ Internal.db untuk konfigurasi (seperti PocketBase)
- ✅ Database.db untuk data sistem (embedded SQLite)
- ✅ In-memory cache (tanpa cache eksternal)
- ✅ Antarmuka admin web (tanpa editing file)
- ✅ Setup wizard (konfigurasi first-time mudah)

---

## 5. Solusi yang Diusulkan

### Visi: MQTT Gateway dengan Internal Database

**Satu file executable** dengan dua embedded database:

**internal.db** (Konfigurasi Aplikasi):
- Pengaturan aplikasi
- Konfigurasi MQTT broker
- Pengaturan database
- Pengaturan web server
- Kredensial autentikasi
- Riwayat perubahan
- Migrasi skema

**database.db** (Data Sistem):
- Registry device
- Definisi topic
- Riwayat pesan
- Status device
- Metrik dan statistik

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
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ Internal Database (internal.db)                 │ │  │
│  │  │ Application Configuration                        │ │  │
│  │  │ • Settings (key-value)                          │ │  │
│  │  │ • Credentials (encrypted)                       │ │  │
│  │  │ • Change History                                │ │  │
│  │  │ • Schema Migrations                             │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                        │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ System Database (database.db)                    │ │  │
│  │  │ Runtime Data (Embedded SQLite)                   │ │  │
│  │  │ • Devices                                        │ │  │
│  │  │ • Topics                                         │ │  │
│  │  │ • Messages                                       │ │  │
│  │  │ • Metrics                                        │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                               │
│  ┌───────────────────────────────────────────────────────┐  │
│  │         UI Web React + Vite (Embedded)                │  │
│  │                                                        │  │
│  │  ┌────────────────┐  ┌────────────────┐              │  │
│  │  │ Setup Mode     │  │ Runner Mode    │              │  │
│  │  │ (Edit Config)  │  │ (Monitoring)   │              │  │
│  │  │ Form-based UI  │  │ Dashboard      │              │  │
│  │  └────────────────┘  └────────────────┘              │  │
│  │                                                        │  │
│  │  ┌────────────────┐  ┌────────────────┐              │  │
│  │  │ Wizard         │  │ Config History │              │  │
│  │  │ Onboarding     │  │ & Rollback     │              │  │
│  │  └────────────────┘  └────────────────┘              │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

---

## 6. Skema Database

### Internal Database (internal.db)

**Tujuan**: Menyimpan konfigurasi aplikasi

```sql
-- Settings table (key-value store)
CREATE TABLE settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    type TEXT NOT NULL,  -- string, int, bool, json
    category TEXT NOT NULL,  -- mqtt, database, web, logging
    description TEXT,
    default_value TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Change history
CREATE TABLE settings_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    setting_key TEXT NOT NULL,
    old_value TEXT,
    new_value TEXT NOT NULL,
    changed_by TEXT,
    changed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    reason TEXT,
    FOREIGN KEY (setting_key) REFERENCES settings(key)
);

-- Credentials (encrypted)
CREATE TABLE credentials (
    id TEXT PRIMARY KEY,
    service TEXT NOT NULL,
    username TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Schema migrations
CREATE TABLE schema_migrations (
    version INTEGER PRIMARY KEY,
    applied_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    description TEXT
);
```

### System Database (database.db)

**Tujuan**: Menyimpan data MQTT runtime

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

-- Messages (with retention)
CREATE TABLE messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    topic TEXT NOT NULL,
    payload TEXT NOT NULL,
    qos INTEGER DEFAULT 0,
    retained BOOLEAN DEFAULT 0,
    published_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_messages_topic ON messages(topic);
CREATE INDEX idx_messages_published_at ON messages(published_at);
```

---

## 7. Fitur (Ruang Lingkup MVP)

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

    // Load credentials from internal.db
    authHook := NewDBAuthHook(config.DB)
    server.AddHook(authHook)

    // TCP listener
    listenAddr := configService.Get("mqtt.listen_address")
    tcp := mqtt.NewTCPListener(listenAddr, nil)
    server.AddListener(tcp)

    return &EmbeddedBroker{server: server, config: config}, nil
}
```

### F2: Internal Database Service

**Implementasi**:

```go
package config

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

type ConfigService struct {
    db *sql.DB
}

func NewConfigService(path string) (*ConfigService, error) {
    db, err := sql.Open("sqlite3", path)
    if err != nil {
        return nil, err
    }

    // Run migrations
    if err := runMigrations(db); err != nil {
        return nil, err
    }

    return &ConfigService{db: db}, nil
}

func (s *ConfigService) Get(key string) (string, error) {
    var value string
    err := s.db.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&value)
    return value, err
}

func (s *ConfigService) Set(key, value, reason string) error {
    tx, _ := s.db.Begin()

    // Get old value
    var oldValue string
    tx.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&oldValue)

    // Update setting
    tx.Exec(`
        INSERT INTO settings (key, value, type, category, updated_at)
        VALUES (?, ?, 'string', 'general', CURRENT_TIMESTAMP)
        ON CONFLICT(key) DO UPDATE SET value = ?, updated_at = CURRENT_TIMESTAMP
    `, key, value, value)

    // Record history
    tx.Exec(`
        INSERT INTO settings_history (setting_key, old_value, new_value, changed_by, reason)
        VALUES (?, ?, ?, 'system', ?)
    `, key, oldValue, value, reason)

    return tx.Commit()
}
```

### F3: Setup Mode (Form-Based)

**UI Form**:

```typescript
// frontend/src/pages/SetupMode.tsx
export function SetupMode() {
  const { data: config } = useQuery({
    queryKey: ['config'],
    queryFn: () => fetch('/api/config').then(r => r.json()),
  });

  const saveConfig = useMutation({
    mutationFn: (newConfig) => fetch('/api/config', {
      method: 'POST',
      body: JSON.stringify(newConfig),
    }),
  });

  return (
    <form onSubmit={(e) => {
      e.preventDefault();
      saveConfig.mutate(config);
    }}>
      <Section title="MQTT Broker">
        <Input label="Listen Port" name="mqtt.listen_port" type="number" />
        <Input label="Max Connections" name="mqtt.max_connections" type="number" />
      </Section>

      <Section title="Web Server">
        <Input label="Port" name="web.port" type="number" />
      </Section>

      <button type="submit">Save & Restart</button>
    </form>
  );
}
```

### F4: Runner Mode

- Konfigurasi read-only
- Dashboard monitoring
- Metrics real-time

### F5: Wizard Onboarding

```typescript
export function OnboardingWizard() {
  const [step, setStep] = useState(1);

  return (
    <Wizard currentStep={step}>
      <Step 1: Welcome />
      <Step 2: Admin Credentials />
      <Step 3: MQTT Configuration />
      <Step 4: Web UI Configuration />
      <Step 5: Review & Start />
    </Wizard>
  );
}
```

### F6: Configuration History & Rollback

```typescript
export function ConfigHistory() {
  const { data: history } = useQuery({
    queryKey: ['config-history'],
    queryFn: () => fetch('/api/config/history').then(r => r.json()),
  });

  const rollback = useMutation({
    mutationFn: (version) => fetch(`/api/config/rollback/${version}`),
  });

  return (
    <div>
      <h1>Configuration History</h1>
      {history.map(change => (
        <div key={change.id}>
          <span>{change.changed_at}</span>
          <span>{change.setting_key}</span>
          <span>{change.old_value} → {change.new_value}</span>
          <button onClick={() => rollback.mutate(change.id)}>
            Rollback
          </button>
        </div>
      ))}
    </div>
  );
}
```

---

## 8. Teknologi Stack

### Backend (Go)

| Komponen | Teknologi | Versi |
|-----------|-----------|-------|
| Bahasa | Go | 1.24+ |
| MQTT Broker | mochi-mqtt | 2.0+ |
| Internal DB | SQLite | 3.40+ |
| System DB | SQLite | 3.40+ |
| Cache | go-cache | Latest |
| HTTP Router | Chi | 5.0+ |

### Frontend (React)

| Komponen | Teknologi | Versi |
|-----------|-----------|-------|
| Framework | React | 18.3+ |
| Build Tool | Vite | 5.0+ |
| Bahasa | TypeScript | 5.3+ |
| State | TanStack Query | 5.0+ |
| Forms | React Hook Form | Latest |

---

## 9. Peta Jalan Pengembangan (3 Bulan)

### Bulan 1: Pondasi

**Minggu 1-2: Setup Proyek**
- [ ] Setup repository
- [ ] Struktur proyek
- [ ] CI/CD

**Minggu 3-4: Backend Inti**
- [ ] Internal.db service
- [ ] System.db service
- [ ] Embedded MQTT broker
- [ ] In-memory cache

### Bulan 2: Fitur

**Minggu 5-6: MQTT & Database**
- [ ] Hooks broker
- [ ] Persistensi pesan
- [ ] Tracking device/topic

**Minggu 7-8: Web UI**
- [ ] Setup React + Vite
- [ ] Setup mode form
- [ ] Runner mode dashboard

### Bulan 3: Poles & Rilis

**Minggu 9-10: Setup & Onboarding**
- [ ] Wizard onboarding
- [ ] Config history UI
- [ ] Validasi

**Minggu 11-12: Pengujian & Rilis**
- [ ] E2E testing
- [ ] Bug fixes
- [ ] Documentation
- [ ] Build & release

---

## 10. Metrik Kesuksesan

### Teknis

| Metrik | Target |
|--------|--------|
| Ukuran Binary | < 60 MB |
| Waktu Startup | < 2 detik |
| Memori Idle | < 100 MB |
| Pesan/detik | > 1,000 |

### Adopsi

| Metrik | Target |
|--------|--------|
| GitHub Stars | 100 (1 bulan) |
| Downloads | 500 (3 bulan) |
| Instalasi Aktif | 50 (3 bulan) |

---

## 11. Kebutuhan Sumber Daya

### Estimasi Anggaran

| Kategori | Biaya | Durasi |
|----------|------|---------|
| Developer Full-Stack | $15,000/bulan | 3 bulan |
| QA (part-time) | $2,000/bulan | 3 bulan |
| Technical Writer | $1,000/bulan | 3 bulan |
| Infrastruktur | $500/bulan | 3 bulan |
| Tools | $200/bulan | 3 bulan |
| **Subtotal** | **$18,700/bulan × 3** | **$56,100** |
| **Buffer (20%)** | **$11,220** | - |
| **Total** | **$67,320** | **~$68,000** |

---

## 12. Kesimpulan

### Filosofi MVP v1.0

**Tanpa File Konfigurasi Manual**:
- ✅ Internal.db untuk semua konfigurasi
- ✅ Form-based UI (bukan YAML editor)
- ✅ Riwayat lengkap dengan rollback
- ✅ Migrasi otomatis

**Fokus pada Esensial**:
- ✅ Single executable
- ✅ Embedded MQTT broker
- ✅ Dua database: internal.db + database.db
- ✅ Setup mudah
- ✅ Monitoring dasar

### Kriteria Sukses

**MVP berhasil jika**:
1. Pengguna dapat mengunduh single executable
2. Menyelesaikan wizard onboarding
3. Menghubungkan perangkat MQTT
4. Melihat pesan di dashboard
5. Tidak perlu mengedit file konfigurasi

---

**Versi Dokumen**: MVP v1.0
**Terakhir Diperbarui**: 9 Januari 2026
**Status**: Draft
**Waktu Pengembangan**: 3 bulan
**Anggaran**: $68,000

---

**Disiapkan oleh**: Tim Pengembangan
**Disetujui oleh**: [Pending]
**Target Rilis**: Q2 2026
