# Rencana Proyek: Universal MQTT Gateway Server (UM-Gateway) v4.0

## Ringkasan Eksekutif

**Nama Proyek**: Universal MQTT Gateway Server (UM-Gateway)
**Versi**: 4.0 (Internal Database untuk Konfigurasi Aplikasi)
**Status**: Proposal/Perencanaan
**Target Rilis**: Q2 2027
**Organisasi**: Divisi Teknologi Rafflesia Agro

---

## Riwayat Versi

| Versi | Tanggal | Perubahan | Penulis |
|-------|---------|-----------|--------|
| **v4.0** | 2026-01-09 | **Internal Database untuk Konfigurasi**: Menggantikan config.yaml dengan internal.db<br>• Menghapus config.yaml sepenuhnya<br>• Menggunakan embedded SQLite database (internal.db) untuk konfigurasi aplikasi<br>• Memisahkan internal.db (konfigurasi aplikasi) dan database.db (data sistem)<br>• Pendekatan seperti PocketBase untuk manajemen konfigurasi<br>• Setup mode dan runner mode tetap sama<br>• API untuk mengelola konfigurasi (bukan edit file)<br>• Migrasi otomatis dari versi sebelumnya<br>• Anggaran: $225,000 (10 bulan) | Tim Pengembangan |
| **v3.0** | 2026-01-09 | **Perubahan Technology Stack & Peningkatan Kejelasan**: UI berbasis React<br>• Mengganti SvelteKit dengan React + Vite<br>• Menambahkan TanStack Query untuk manajemen state server<br>• Menambahkan TanStack Router untuk routing<br>• Memperjelas Setup Mode = editor config.yaml<br>• Memperjelas Runner Mode = operasi produksi (config read-only)<br>• Pemisahan jelas: config.yaml (infrastruktur) vs database (data runtime)<br>• Anggaran: $215,500 (10 bulan) | Tim Pengembangan |

---

## 1. Latar Belakang

### Evolusi dari v3 ke v4

**Pendekatan v3**: Menggunakan config.yaml untuk konfigurasi infrastruktur
**Pendekatan v4**: Menggunakan embedded database (internal.db) untuk semua konfigurasi

### Keterbatasan v3 (config.yaml)

1. **Masalah Editing File**: Pengguna harus mengedit file YAML secara manual atau via UI yang kompleks
2. **Masalah Validasi**: Validasi YAML bisa rumit dan error messages tidak jelas
3. **Masalah Migrasi**: Migrasi antar versi memerlukan manipulasi file
4. **Masalah Versioning**: Tidak ada riwayat perubahan konfigurasi
5. **Masalah Backup**: Backup terpisah untuk config dan database
6. **Masalah Akses**: Permission file bisa menjadi masalah di beberapa OS
7. **Masalah Konsistensi**: Tidak ada ACID untuk memastikan konfigurasi konsisten

### Peningkatan Utama di v4

**1. Menggantikan config.yaml dengan internal.db**

**internal.db** (Konfigurasi Aplikasi):
- **Tujuan**: Menyimpan semua pengaturan aplikasi dan infrastruktur
- **Lokasi**: `./data/internal.db`
- **Format**: Embedded SQLite database
- **Akses**: Hanya via API (aman dan terstruktur)
- **Backup**: Termasuk dalam backup database rutin
- **Migrasi**: Otomatis dengan skema database
- **Transaksional**: ACID compliant

**database.db** (Data Sistem):
- **Tujuan**: Menyimpan data runtime MQTT (pesan, device, topic, dll)
- **Lokasi**: `./data/database.db` (default, dapat dikonfigurasi)
- **Format**: Embedded SQLite database dengan mode WAL
- **Adapter**: Dapat dikonfigurasi untuk database lain (PostgreSQL, MongoDB, InfluxDB)
- **Akses**: Via API untuk operasi CRUD
- **Backup**: Otomatis/periodik

**2. Pendekatan seperti PocketBase**

Mengadopsi konsep dari PocketBase:
- Single file executable dengan embedded database
- Admin UI untuk mengelola konfigurasi
- API REST untuk semua operasi
- Migrasi skema otomatis
- Collections untuk data terstruktur

**3. API untuk Mengelola Konfigurasi**

Daripada mengedit file YAML:
- `GET /api/config` - Ambil semua konfigurasi
- `PUT /api/config/:key` - Update nilai konfigurasi
- `POST /api/config/validate` - Validasi konfigurasi
- `POST /api/config/reset` - Reset ke default
- `GET /api/config/history` - Riwayat perubahan
- `POST /api/config/rollback` - Kembalikan ke versi sebelumnya

**4. Setup Mode dan Runner Mode**

**Setup Mode**:
- Dapat mengubah konfigurasi di internal.db
- UI form untuk pengaturan (bukan editor YAML)
- Validasi real-time
- Test koneksi
- Simpan & Restart

**Runner Mode**:
- Konfigurasi di internal.db menjadi read-only
- Tidak dapat mengubah pengaturan aplikasi
- Operasi produksi normal
- Dashboard monitoring

### Peluang Pasar

Lanskap IoT membutuhkan **MQTT gateway dengan manajemen konfigurasi modern**:
- **Tanpa Editing File**: Semua konfigurasi via API/Web UI
- **Transaksional**: Konfigurasi yang konsisten dan andal
- **Versioning**: Riwayat perubahan dan rollback
- **Mudah Dimigrasi**: Upgrade yang mulus
- **Pendekatan Modern**: Seperti PocketBase, Strapi, dll

---

## 2. Model Bisnis

### Proposition Nilai

**Untuk Edge Deployments**:
- Tidak perlu mengedit file konfigurasi
- Semua pengaturan via Web UI modern
- Konfigurasi yang aman dan transaksional
- Riwayat perubahan dan rollback

**Untuk Lingkungan Air-Gapped**:
- Operasi offline sepenuhnya
- Embedded database (tidak perlu database eksternal)
- Backup dan restore sederhana
- Konfigurasi yang mudah dikloning

**Untuk Pengembangan/Pengujian**:
- Reset konfigurasi dengan satu klik
- Snapshot dan restore state
- Lingkungan yang dapat diulang
- Migrasi otomatis antar versi

### Aliran Pendapatan

1. **Free Self-Hosted**: Single binary dengan embedded databases (Lisensi MIT)
2. **Lisensi Pro** ($299 sekali): Adapter database eksternal, fitur lanjutan
3. **Lisensi Enterprise** ($1,499/tahun): Dukungan prioritas, build kustom

### Strategi Harga

| Edisi | Target Pasar | Model Harga | Fitur |
|---------|--------------|---------------|----------|
| Komunitas | Hobiwan, edge computing | Gratis (MIT) | Internal.db, database.db embedded, setup mode, UI React |
| Profesional | UKM, produksi | $299 sekali | Adapter DB eksternal, akses API, dukungan prioritas |
| Enterprise | Organisasi besar | $1,499/tahun | Multiple gateway, clustering, dukungan prioritas |

---

## 3. Target Audiens

### Pengguna Utama

**1. Insinyur Edge Computing**
- Deploy IoT gateway pada perangkat edge
- Tidak ingin mengedit file YAML
- Butuh konfigurasi yang transaksional
- Titik nyeri: Error parsing YAML

**2. Operator Lingkungan Air-Gapped**
- IoT industri dalam jaringan terisolasi
- Butuh manajemen konfigurasi yang mudah
- Tidak dapat mengedit file di remote
- Titik nyeri: Akses file yang sulit

**3. Tim Pengembangan**
- Butuh lingkungan pengujian lokal
- Ingin reset konfigurasi cepat
- Perlu snapshot state
- Titik nyeri: Konfigurasi yang sulit diulang

**4. Integrator Sistem**
- Deploy solusi di lokasi pelanggan
- Ingin konfigurasi via UI, bukan file
- Butuh riwayat perubahan
- Titik nyeri: Tidak ada audit trail

---

## 4. Pernyataan Masalah

### Masalah Inti

**Masalah 1: Editing File Konfigurasi**
- v3 menggunakan config.yaml yang harus diedit
- Pengguna harus memahami sintaks YAML
- Error indentation dan format yang sulit di-debug
- Pendekatan saat ini: Editor YAML dengan validasi
- Solusi: API/Web UI untuk semua konfigurasi (tanpa file editing)

**Masalah 2: Tidak Ada Riwayat Konfigurasi**
- Tidak ada jejak perubahan konfigurasi
- Sulit untuk rollback ke pengaturan sebelumnya
- Pendekatan saat ini: Tidak ada riwayat
- Solusi: Tabel history di internal.db

**Masalah 3: Migrasi Konfigurasi Sulit**
- Upgrade versi memerlukan edit manual config.yaml
- Tidak otomatis
- Pendekatan saat ini: Panduan migrasi manual
- Solusi: Migrasi otomatis dengan skema database

**Masalah 4: Backup Terpisah**
- config.yaml dan database.db harus dibackup terpisah
- Tidak konsisten
- Pendekatan saat ini: Backup script manual
- Solusi: Satu backup untuk kedua database

**Masalah 5: Validasi Konfigurasi Rumit**
- Validasi YAML bisa kompleks
- Error messages tidak jelas
- Pendekatan saat ini: Validasi YAML dasar
- Solusi: Validasi di level API dengan pesan error yang jelas

### Dampak

- **Error Konfigurasi**: Pengguna membuat kesalahan sintaks YAML
- **Upgrade yang Sulit**: Pengguna takut upgrade karena migrasi manual
- **Tidak Ada Audit**: Tidak bisa melihat siapa mengubah apa dan kapan
- **Backup Tidak Konsisten**: config dan database tidak sinkron
- **DX yang Buruk**: Debugging konfigurasi yang sulit

---

## 5. Solusi yang Diusulkan

### Visi: Konfigurasi Database-First

**Single executable binary** dengan:
- **internal.db**: Embedded SQLite untuk konfigurasi aplikasi
- **database.db**: Embedded SQLite untuk data sistem (dengan adapter untuk DB lain)
- **API/Web UI**: Semua konfigurasi via API (tanpa file editing)
- **Setup Mode**: Form-based configuration UI
- **Runner Mode**: Konfigurasi read-only
- **Riwayat & Rollback**: Version control untuk konfigurasi
- **Migrasi Otomatis**: Skema database untuk upgrade

### Inovasi Utama

**1. Arsitektur Dual Database**

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
│  │  ┌────────────────────────────────────────────────┐ │  │
│  │  │         Internal Database (internal.db)         │ │  │
│  │  │         Application Configuration               │ │  │
│  │  │                                                  │ │  │
│  │  │  • Pengaturan aplikasi                           │ │  │
│  │  │  • Konfigurasi MQTT broker                       │ │  │
│  │  │  • Pengaturan database                           │ │  │
│  │  │  • Pengaturan web server                         │ │  │
│  │  │  • Kredensial autentikasi                        │ │  │
│  │  │  • Konfigurasi logging                           │ │  │
│  │  │  • Riwayat perubahan                             │ │  │
│  │  │  • Migrasi skema                                 │ │  │
│  │  └────────────────────────────────────────────────┘ │  │
│  │                                                        │  │
│  │  ┌────────────────────────────────────────────────┐ │  │
│  │  │      System Database (database.db - Default)    │ │  │
│  │  │      Runtime Data (with Adapter Pattern)        │ │  │
│  │  │                                                  │ │  │
│  │  │  • Registry device                               │ │  │
│  │  │  • Topics dan subscriptions                      │ │  │
│  │  │  • Riwayat pesan MQTT                            │ │  │
│  │  │  • Status device                                 │ │  │
│  │  │  • Metrik dan statistik                          │ │  │
│  │  │  • Akun pengguna (jika ada)                      │ │  │
│  │  └────────────────────────────────────────────────┘ │  │
│  │                                                        │  │
│  │  ┌────────────┐  ┌────────────┐                    │  │
│  │  │ Mode       │  │ Config     │                    │  │
│  │  │ Manager    │  │ Service    │                    │  │
│  │  └────────────┘  └────────────┘                    │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                               │
│  ┌───────────────────────────────────────────────────────┐  │
│  │         React + Vite Web UI (Embedded SPA)            │  │
│  │                                                        │  │
│  │  ┌────────────────────────────────────────────────┐  │  │
│  │  │ Setup Mode UI (/setup)                          │  │  │
│  │  │ - Form-based configuration                      │  │  │
│  │  │ - Real-time validation                          │  │  │
│  │  │ - Test connection                               │  │  │
│  │  │ - Save & Restart                                │  │  │
│  │  └────────────────────────────────────────────────┘  │  │
│  │                                                        │  │
│  │  ┌────────────────────────────────────────────────┐  │  │
│  │  │ Runner Mode UI (/)                              │  │  │
│  │  │ - Dashboard monitoring                          │  │  │
│  │  │ - Config viewer (read-only)                     │  │  │
│  │  │ - Metrics real-time                             │  │  │
│  │  │ - Log viewer                                    │  │  │
│  │  └────────────────────────────────────────────────┘  │  │
│  │                                                        │  │
│  │  ┌────────────────────────────────────────────────┐  │  │
│  │  │ TanStack Router & Query                         │  │  │
│  │  └────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

**2. Skema Internal Database (internal.db)**

```sql
-- Tabel pengaturan aplikasi (key-value store)
CREATE TABLE settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    type TEXT NOT NULL,  -- string, int, bool, json
    category TEXT NOT NULL,  -- mqtt, database, web, logging, etc.
    description TEXT,
    default_value TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Tabel riwayat perubahan
CREATE TABLE settings_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    setting_key TEXT NOT NULL,
    old_value TEXT,
    new_value TEXT NOT NULL,
    changed_by TEXT,  -- 'system', 'admin', or user_id
    changed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    reason TEXT,
    FOREIGN KEY (setting_key) REFERENCES settings(key)
);

-- Tabel kredensial (terenkripsi)
CREATE TABLE credentials (
    id TEXT PRIMARY KEY,
    service TEXT NOT NULL,  -- 'mqtt', 'database', etc.
    username TEXT NOT NULL,
    password_hash TEXT NOT NULL,  -- Hashed
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Tabel migrasi skema
CREATE TABLE schema_migrations (
    version INTEGER PRIMARY KEY,
    applied_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    description TEXT
);

-- Indexes
CREATE INDEX idx_settings_category ON settings(category);
CREATE INDEX idx_settings_history_key ON settings_history(setting_key);
CREATE INDEX idx_settings_history_changed_at ON settings_history(changed_at);
```

**3. API untuk Konfigurasi**

```go
// API endpoints untuk konfigurasi

// GET /api/config - Ambil semua konfigurasi
func GetConfig(c *gin.Context) {
    config := configService.GetAll()
    c.JSON(200, config)
}

// GET /api/config/:category - Ambil konfigurasi per kategori
func GetConfigByCategory(c *gin.Context) {
    category := c.Param("category")
    config := configService.GetByCategory(category)
    c.JSON(200, config)
}

// PUT /api/config/:key - Update nilai konfigurasi
func UpdateConfig(c *gin.Context) {
    var req UpdateConfigRequest
    if err := c.BindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // Validasi
    if err := configService.Validate(req.Key, req.Value); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // Update dengan riwayat
    if err := configService.Set(req.Key, req.Value, req.Reason); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"status": "ok"})
}

// POST /api/config/validate - Validasi konfigurasi
func ValidateConfig(c *gin.Context) {
    var req map[string]interface{}
    if err := c.BindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    errors := configService.ValidateAll(req)
    if len(errors) > 0 {
        c.JSON(400, gin.H{"errors": errors})
        return
    }

    c.JSON(200, gin.H{"status": "valid"})
}

// GET /api/config/history - Riwayat perubahan
func GetConfigHistory(c *gin.Context) {
    history := configService.GetHistory()
    c.JSON(200, history)
}

// POST /api/config/rollback/:version - Rollback ke versi sebelumnya
func RollbackConfig(c *gin.Context) {
    version := c.Param("version")
    if err := configService.Rollback(version); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"status": "rolled_back"})
}

// POST /api/config/reset - Reset ke default
func ResetConfig(c *gin.Context) {
    if err := configService.ResetToDefaults(); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"status": "reset"})
}
```

**4. Setup Mode UI (Form-Based)**

```typescript
// frontend/src/pages/SetupMode.tsx
import { useMutation, useQuery } from '@tanstack/react-query';

export function SetupMode() {
  // Fetch konfigurasi saat ini
  const { data: config, isLoading } = useQuery({
    queryKey: ['config'],
    queryFn: () => fetch('/api/config').then(r => r.json()),
  });

  // Simpan mutation konfigurasi
  const saveConfig = useMutation({
    mutationFn: (newConfig: Record<string, any>) =>
      fetch('/api/config', {
        method: 'POST',
        body: JSON.stringify(newConfig),
      }),
    onSuccess: () => {
      alert('Konfigurasi tersimpan! Me-restart gateway...');
      setTimeout(() => window.location.reload(), 2000);
    },
  });

  if (isLoading) return <div>Loading...</div>;

  return (
    <div className="setup-mode">
      <h1>Setup Mode - Konfigurasi Aplikasi</h1>

      <form onSubmit={(e) => {
        e.preventDefault();
        saveConfig.mutate(config);
      }}>
        {/* MQTT Configuration */}
        <Section title="MQTT Broker">
          <FormField
            label="Tipe Broker"
            type="select"
            value={config.mqtt.type}
            onChange={(v) => setConfig('mqtt.type', v)}
            options={[
              { value: 'embedded', label: 'Embedded' },
              { value: 'external', label: 'External' },
            ]}
          />
          <FormField
            label="Port Listen"
            type="number"
            value={config.mqtt.listen_port}
            onChange={(v) => setConfig('mqtt.listen_port', v)}
          />
        </Section>

        {/* Database Configuration */}
        <Section title="Database">
          <FormField
            label="Tipe Database"
            type="select"
            value={config.database.type}
            onChange={(v) => setConfig('database.type', v)}
            options={[
              { value: 'sqlite', label: 'SQLite (Embedded)' },
              { value: 'postgres', label: 'PostgreSQL' },
              { value: 'mongodb', label: 'MongoDB' },
            ]}
          />
        </Section>

        <button type="submit">Simpan & Restart</button>
      </form>
    </div>
  );
}
```

**5. Migrasi Otomatis**

```go
// Migrasi dari v3 ke v4
func Migration_v3_to_v4() error {
    // Baca config.yaml lama jika ada
    if oldConfigExists() {
        yamlConfig := readYAMLConfig()

        // Konversi ke database
        for key, value := range yamlConfig {
            configService.Set(key, value, "Migration from v3")
        }

        // Backup dan hapus config.yaml
        backupFile("config.yaml")
        removeFile("config.yaml")
    }

    return nil
}
```

---

## 6. Tujuan Proyek

### Tujuan Utama (Must Have)

**O1: Implementasi Internal Database (internal.db)**
- Gantikan config.yaml dengan embedded SQLite
- Skema untuk settings, credentials, history
- API CRUD untuk konfigurasi
- Migrasi otomatis dari v3

**O2: Setup Mode Form-Based**
- UI form untuk konfigurasi (bukan editor YAML)
- Validasi real-time
- Kategorisasi pengaturan
- Test koneksi

**O3: Riwayat Konfigurasi**
- Track semua perubahan
- Rollback ke versi sebelumnya
- Export/import konfigurasi
- Audit log

**O4: Runner Mode Read-Only**
- Konfigurasi tidak dapat diubah
- Dashboard monitoring
- Viewer konfigurasi

**O5: Adapter Pattern untuk Database**
- database.db sebagai default (SQLite embedded)
- Adapter untuk PostgreSQL, MongoDB, InfluxDB
- Konfigurasi di internal.db
- Data runtime di database.db (atau eksternal)

### Tujuan Sekunder (Should Have)

**O6: Validasi Konfigurasi yang Ditingkatkan**
- Validasi per-field
- Dependency checking
- Conflict detection
- Pesan error yang jelas

**O7: Snapshot & Restore**
- Snapshot konfigurasi lengkap
- Restore dari snapshot
- Kloning konfigurasi

**O8: UI Pengaturan Lanjutan**
- Search dan filter pengaturan
- Reset per-kategori
- Import/export konfigurasi

---

## 7. Ruang Lingkup & Batasan Proyek

### Dalam Ruang Lingkup (Apa yang Akan Kami Bangun)

#### Fungsionalitas Inti

1. **Internal Database (internal.db)**
   - Skema untuk settings, credentials, history
   - API CRUD
   - Migrasi otomatis
   - Backup/restore

2. **System Database (database.db)**
   - Skema untuk device, topic, message, dll
   - Adapter pattern untuk DB eksternal
   - Default SQLite embedded

3. **Setup Mode**
   - UI form-based
   - Validasi real-time
   - Test koneksi
   - Simpan & Restart

4. **Runner Mode**
   - Konfigurasi read-only
   - Dashboard monitoring
   - Viewer konfigurasi

5. **Riwayat Konfigurasi**
   - Tracking perubahan
   - Rollback
   - Audit log

### Di Luar Ruang Lingkup (Apa yang Tidak Akan Kami Bangun - Awalnya)

1. **Fitur Lanjutan** (Fase 2)
   - Multi-user management
   - Role-based access control
   - Advanced authentication

2. **High Availability** (Fase 3)
   - Clustering
   - Replikasi data
   - Failover otomatis

### Batasan & Kendala

#### Kendala Teknis
- **Frontend**: React 18+, Vite 5+, TypeScript 5+
- **Backend**: Go 1.24+
- **Library MQTT**: mochi-mqtt
- **Database**: SQLite untuk internal.db dan database.db (default)

#### Kendala Waktu
- **Fase 1 (v4.0)**: 10 bulan
- **Fase 2 (Fitur Lanjutan)**: +3 bulan

#### Kendala Anggaran
- **Tim Pengembangan**: 2-3 developer full-time
- **Infrastruktur**: $1,500/bulan

---

## 8. Teknologi Stack

### Stack Frontend

| Komponen | Teknologi | Versi | Rasional |
|-----------|-----------|-------|-----------|
| Library UI | React | 18.3+ | Ekosistem terbesar |
| Build Tool | Vite | 5.0+ | HMR cepat, build teroptimasi |
| Bahasa | TypeScript | 5.3+ | Keamanan tipe |
| Manajemen State | TanStack Query | 5.0+ | Manajemen state server |
| Routing | TanStack Router | 1.0+ | Routing type-safe |
| Form | React Hook Form | Latest | Form handling yang baik |

### Stack Backend

| Komponen | Teknologi | Versi | Rasional |
|-----------|-----------|-------|-----------|
| Bahasa | Go | 1.24+ | Performa, single binary |
| MQTT Broker | mochi-mqtt | 2.0+ | Embedded Go MQTT broker |
| Internal DB | SQLite | 3.40+ | Embedded, single file |
| System DB | SQLite | 3.40+ | Default, dengan adapter |
| Cache | go-cache | Latest | In-memory caching |
| HTTP Router | Chi | 5.0+ | Ringan, idiomatic |

---

## 9. Peta Jalan Pengembangan

### Fase 1: Pondasi (Bulan 1-4)

**Sprint 1-2: Setup Proyek & Arsitektur**
- [ ] Struktur monorepo
- [ ] Setup frontend React + Vite
- [ ] Desain skema internal.db
- [ ] Desain skema database.db

**Sprint 3-4: Internal Database**
- [ ] Implementasi internal.db
- [ ] API CRUD untuk settings
- [ ] Riwayat konfigurasi
- [ ] Migrasi dari v3

**Sprint 5-6: System Database**
- [ ] Implementasi database.db
- [ ] Adapter pattern
- [ ] Skema untuk device, topic, message

**Sprint 7-8: Frontend**
- [ ] Setup mode form-based
- [ ] Runner mode dashboard
- [ ] Integrasi TanStack Query

### Fase 2: Fitur (Bulan 5-8)

**Sprint 9-12: Validasi & Testing**
- [ ] Validasi konfigurasi
- [ ] Integration testing
- [ ] E2E testing
- [ ] Load testing

**Sprint 13-16: Polesan**
- [ ] UI improvements
- [ ] Error handling
- [ ] Performance optimization
- [ ] Documentation

### Fase 3: Produksi (Bulan 9-10)

**Sprint 17-20: Deploy**
- [ ] Build produksi
- [ ] Cross-platform compilation
- [ ] Installation guide
- [ ] Release notes

---

## 10. Metrik Kesuksesan

### Metrik Teknis

| Metrik | Target | Pengukuran |
|--------|--------|-------------|
| Ukuran Binary | < 60 MB | Artefak build |
| Waktu Startup | < 2 detik | Pengujian manual |
| Penggunaan Memori | < 100 MB idle | Monitoring |
| Setup Completion Rate | > 95% | Analytics |

---

## 11. Kebutuhan Sumber Daya

### Estimasi Anggaran

| Kategori | Biaya (Bulanan) | Durasi |
|----------|----------------|---------|
| Gaji Tim Pengembangan | $20,000 | 10 bulan |
| Infrastruktur | $1,500 | 10 bulan |
| Alat & Layanan | $400 | Berkelanjutan |
| Kontinjensi (15%) | $3,150 | - |
| **Total** | **$22,500/bulan × 10** | **$225,000** |

---

## 12. Perbandingan: v3 vs v4

| Aspek | v3 (config.yaml) | v4 (internal.db) |
|--------|----------------|----------------|
| **Konfigurasi** | File YAML | Database SQLite |
| **Editing** | Edit file manual | API/Web UI |
| **Riwayat** | Tidak ada | Ya, lengkap |
| **Migrasi** | Manual | Otomatis |
| **Backup** | Terpisah | Satu backup |
| **Validasi** | YAML parsing | API validation |
| **Rollback** | Manual | Otomatis |
| **Anggaran** | $215,500 | $225,000 |

---

## 13. Kesimpulan

Versi v4 merepresentasikan **pendekatan modern dan database-first**:

### Keunggulan Utama atas v3

**Manajemen Konfigurasi**
- ✅ Tanpa editing file
- ✅ API/Web UI untuk semua konfigurasi
- ✅ Riwayat lengkap
- ✅ Rollback otomatis
- ✅ Migrasi otomatis

**Keandalan**
- ✅ Transaksional (ACID)
- ✅ Konsisten
- ✅ Backup terpadu
- ✅ Validasi di level API

**Pengalaman Pengguna**
- ✅ Form-based UI (bukan editor YAML)
- ✅ Validasi real-time
- ✅ Pesan error yang jelas
- ✅ Mudah digunakan

Dengan eksekusi yang terfokus selama 10 bulan, kami dapat mengirimkan **v4.0 production-ready** dengan manajemen konfigurasi yang modern dan andal.

---

**Versi Dokumen**: 4.0
**Terakhir Diperbarui**: 9 Januari 2026
**Status**: Draft - Menunggu Review
**Menggantikan**: PROJECT_BRIEF_ID_v3.md

---

**Disiapkan oleh**: Tim Pengembangan
**Disetujui oleh**: [Pending]
**Tanggal Review**: [Pending]
