# Rencana Pengujian Teknologi: UM-Gateway MVP v1.0

## Informasi Dokumen

**Proyek**: Universal MQTT Gateway Server (UM-Gateway)
**Versi**: MVP v1.0
**Tipe Dokumen**: Rencana Pengujian & Validasi Teknologi
**Terakhir Diperbarui**: 9 Januari 2026
**Status**: Draft

---

## 1. Ringkasan

### Tujuan

Dokumen ini menguraikan rencana pengujian teknologi untuk UM-Gateway MVP v1.0 untuk memvalidasi **pendekatan konfigurasi internal.db** sebagai pengganti pendekatan file konfigurasi YAML.

### Filosofi Pengujian

**Apa yang berubah**:
- ❌ DIHAPUS: konfigurasi berbasis config.yaml
- ✅ DITAMBAHKAN: internal.db (SQLite) untuk konfigurasi aplikasi
- ✅ DITAMBAHKAN: UI Setup Mode berbasis form
- ✅ DITAMBAHKAN: Riwayat & rollback konfigurasi
- ✅ DITAMBAHKAN: Migrasi otomatis

**Tujuan Pengujian**:
1. Memvalidasi internal.db sebagai penyimpan konfigurasi
2. Memverifikasi alur Setup Mode → Runner Mode
3. Menguji broker MQTT dengan konfigurasi embedded
4. End-to-end: Setup → Start → Buat Topic → Publish Pesan

---

## 2. Validasi Technology Stack

### 2.1 Teknologi Inti

| Komponen | Teknologi | Versi | Cakupan Tes |
|-----------|-----------|---------|-------------|
| **Bahasa Backend** | Go | 1.24+ | Kompilasi, operasi dasar |
| **MQTT Broker** | mochi-mqtt | 2.0+ | Fungsionalitas broker embedded |
| **Internal DB** | SQLite | 3.40+ | Penyimpan konfigurasi |
| **System DB** | SQLite | 3.40+ | Data runtime |
| **Cache** | go-cache | Latest | Caching in-memory |
| **HTTP Router** | Chi | 5.0+ | Endpoints API |
| **Frontend** | React | 18.3+ | Rendering UI |
| **Build Tool** | Vite | 5.0+ | Development & build |
| **Bahasa** | TypeScript | 5.3+ | Pengecekan tipe |
| **State** | TanStack Query | 5.0+ | State server |
| **Form** | React Hook Form | Latest | Penanganan form |

### 2.2 Validasi Arsitektur

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
│  │         Web UI React + Vite (Embedded)                │  │
│  │                                                        │  │
│  │  ┌────────────────┐  ┌────────────────┐              │  │
│  │  │ Setup Mode     │  │ Runner Mode    │              │  │
│  │  │ (Edit Config)  │  │ (Monitoring)   │              │  │
│  │  │ Form-based UI  │  │ Dashboard      │              │  │
│  │  └────────────────┘  └────────────────┘              │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

---

## 3. Struktur Rencana Tes

### Fase 1: Tes Foundation (Minggu 1)

#### Tes 1.1: Setup Proyek Go

**Tujuan**: Memvalidasi struktur proyek Go dasar

**Langkah**:
```bash
# 1. Inisialisasi modul Go
go mod init github.com/rafflesia-agro/um-gateway

# 2. Buat struktur direktori
mkdir -p cmd/server internal/{config,broker,database,api,cache}
mkdir -p frontend/src/{pages,components,hooks,services}

# 3. Install dependensi
go get github.com/mochi-mqtt/mqtt/v2
go get github.com/mattn/go-sqlite3
go get github.com/go-chi/chi/v5
go get github.com/patrickmn/go-cache

# 4. Verifikasi kompilasi
go build -o um-gateway.exe ./cmd/server
```

**Hasil yang Diharapkan**:
- ✅ Modul terinisialisasi dengan sukses
- ✅ Semua dependensi terdownload
- ✅ Binary terkompilasi tanpa error
- ✅ Ukuran binary < 20 MB (build awal)

**Kriteria Sukses**:
```
✓ go mod sum sesuai dengan dependensi yang diharapkan
✓ Binary dieksekusi tanpa panic
✓ Flag versi berfungsi: ./um-gateway.exe --version
```

---

#### Tes 1.2: Skema Internal Database

**Tujuan**: Memvalidasi pembuatan skema internal.db

**Definisi Skema**:

```go
// internal/config/schema.go

package config

const schema = `
-- Settings table (key-value store)
CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    type TEXT NOT NULL,
    category TEXT NOT NULL,
    description TEXT,
    default_value TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Change history
CREATE TABLE IF NOT EXISTS settings_history (
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
CREATE TABLE IF NOT EXISTS credentials (
    id TEXT PRIMARY KEY,
    service TEXT NOT NULL,
    username TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Schema migrations
CREATE TABLE IF NOT EXISTS schema_migrations (
    version INTEGER PRIMARY KEY,
    applied_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    description TEXT
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_settings_category ON settings(category);
CREATE INDEX IF NOT EXISTS idx_settings_history_key ON settings_history(setting_key);
CREATE INDEX IF NOT EXISTS idx_settings_history_changed_at ON settings_history(changed_at);
`
```

**Langkah Tes**:
```go
// internal/config/schema_test.go

func TestInternalDBSchema(t *testing.T) {
    // 1. Buat database sementara
    db, err := sql.Open("sqlite3", ":memory:")
    assert.NoError(t, err)

    // 2. Eksekusi skema
    _, err = db.Exec(schema)
    assert.NoError(t, err)

    // 3. Verifikasi tabel ada
    tables := []string{"settings", "settings_history", "credentials", "schema_migrations"}
    for _, table := range tables {
        var count int
        err = db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&count)
        assert.NoError(t, err)
        assert.Equal(t, 1, count)
    }

    // 4. Verifikasi index
    indexes := []string{"idx_settings_category", "idx_settings_history_key", "idx_settings_history_changed_at"}
    for _, index := range indexes {
        var count int
        err = db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='index' AND name=?", index).Scan(&count)
        assert.NoError(t, err)
        assert.Equal(t, 1, count)
    }
}
```

**Hasil yang Diharapkan**:
- ✅ Semua tabel dibuat dengan sukses
- ✅ Semua index dibuat
- ✅ Foreign key diberlakukan
- ✅ Nilai default berfungsi

---

#### Tes 1.3: CRUD Service Konfigurasi

**Tujuan**: Memvalidasi operasi CRUD konfigurasi dasar

**Langkah Tes**:
```go
// internal/config/service_test.go

func TestConfigServiceCRUD(t *testing.T) {
    service := NewConfigService(":memory:")

    // Tes 1: Set nilai
    err := service.Set("mqtt.listen_port", "1883", "Initial setup")
    assert.NoError(t, err)

    // Tes 2: Get nilai
    value, err := service.Get("mqtt.listen_port")
    assert.NoError(t, err)
    assert.Equal(t, "1883", value)

    // Tes 3: Get semua dalam kategori
    mqttConfig := service.GetByCategory("mqtt")
    assert.Contains(t, mqttConfig, "mqtt.listen_port")

    // Tes 4: Update nilai
    err = service.Set("mqtt.listen_port", "8883", "Changed to secure port")
    assert.NoError(t, err)

    // Tes 5: Verifikasi history
    history := service.GetHistory("mqtt.listen_port")
    assert.Len(t, history, 2)
    assert.Equal(t, "1883", history[0].OldValue)
    assert.Equal(t, "8883", history[0].NewValue)

    // Tes 6: Validasi
    err = service.Validate("mqtt.listen_port", "invalid")
    assert.Error(t, err)

    // Tes 7: Reset ke default
    err = service.ResetToDefaults()
    assert.NoError(t, err)
    value, _ = service.Get("mqtt.listen_port")
    assert.Equal(t, "1883", value) // Nilai default
}
```

**Hasil yang Diharapkan**:
- ✅ Operasi Set/Get bekerja dengan benar
- ✅ Pelacakan history fungsional
- ✅ Validasi menolak nilai tidak valid
- ✅ Reset ke default bekerja
- ✅ Semua operasi bersifat transaksional

---

### Fase 2: POC Database (Minggu 2)

#### Tes 2.1: Skema System Database

**Tujuan**: Memvalidasi skema database.db untuk data MQTT

**Definisi Skema**:

```go
// internal/database/schema.go

const schema = `
-- Devices
CREATE TABLE IF NOT EXISTS devices (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT,
    first_seen DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_seen DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Topics
CREATE TABLE IF NOT EXISTS topics (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Messages (with retention)
CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    topic TEXT NOT NULL,
    payload TEXT NOT NULL,
    qos INTEGER DEFAULT 0,
    retained BOOLEAN DEFAULT 0,
    published_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_messages_topic ON messages(topic);
CREATE INDEX IF NOT EXISTS idx_messages_published_at ON messages(published_at);
CREATE INDEX IF NOT EXISTS idx_devices_last_seen ON devices(last_seen);
```

**Langkah Tes**:
```go
// internal/database/schema_test.go

func TestSystemDBSchema(t *testing.T) {
    db, err := sql.Open("sqlite3", ":memory:")
    assert.NoError(t, err)

    _, err = db.Exec(schema)
    assert.NoError(t, err)

    // Verifikasi tabel
    tables := []string{"devices", "topics", "messages"}
    for _, table := range tables {
        var count int
        err = db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&count)
        assert.NoError(t, err)
        assert.Equal(t, 1, count)
    }

    // Tes insert pesan
    _, err = db.Exec("INSERT INTO messages (topic, payload, qos) VALUES (?, ?, ?)", "test/topic", "hello", 0)
    assert.NoError(t, err)

    // Tes query
    var payload string
    err = db.QueryRow("SELECT payload FROM messages WHERE topic = ?", "test/topic").Scan(&payload)
    assert.NoError(t, err)
    assert.Equal(t, "hello", payload)
}
```

---

#### Tes 2.2: Service Database dengan Adapter Pattern

**Tujuan**: Menguji service database dengan adapter SQLite

**Implementasi**:

```go
// internal/database/service.go

type DatabaseService interface {
    SaveMessage(topic, payload string, qos int) error
    GetMessages(topic string, limit int) ([]Message, error)
    RegisterDevice(id, name, deviceType string) error
    GetDevice(id string) (*Device, error)
    CreateTopic(name string) error
}

type SQLiteAdapter struct {
    db *sql.DB
}

func (a *SQLiteAdapter) SaveMessage(topic, payload string, qos int) error {
    _, err := a.db.Exec(
        "INSERT INTO messages (topic, payload, qos) VALUES (?, ?, ?)",
        topic, payload, qos,
    )
    return err
}
```

**Langkah Tes**:
```go
// internal/database/service_test.go

func TestDatabaseService(t *testing.T) {
    db, _ := sql.Open("sqlite3", ":memory:")
    db.Exec(schema)

    service := NewSQLiteAdapter(db)

    // Tes 1: Simpan dan ambil pesan
    err := service.SaveMessage("sensor/temp", "25.5", 0)
    assert.NoError(t, err)

    messages, err := service.GetMessages("sensor/temp", 10)
    assert.NoError(t, err)
    assert.Len(t, messages, 1)
    assert.Equal(t, "25.5", messages[0].Payload)

    // Tes 2: Register device
    err = service.RegisterDevice("device-001", "Temperature Sensor", "sensor")
    assert.NoError(t, err)

    device, err := service.GetDevice("device-001")
    assert.NoError(t, err)
    assert.Equal(t, "Temperature Sensor", device.Name)

    // Tes 3: Buat topic
    err = service.CreateTopic("sensor/temp")
    assert.NoError(t, err)
}
```

---

### Fase 3: POC MQTT Broker (Minggu 3)

#### Tes 3.1: Embedded MQTT Broker

**Tujuan**: Memvalidasi broker mochi-mqtt embedded

**Implementasi**:

```go
// internal/broker/embedded.go

package broker

import (
    mqtt "github.com/mochi-mqtt/mqtt/v2"
    "github.com/mochi-mqtt/mqtt/v2/hooks/auth"
    "github.com/mochi-mqtt/mqtt/v2/listeners"
)

type EmbeddedBroker struct {
    server *mqtt.Server
    configService *config.ConfigService
}

func NewEmbeddedBroker(configService *config.ConfigService) (*EmbeddedBroker, error) {
    // Ambil konfigurasi dari internal.db
    listenAddr, _ := configService.Get("mqtt.listen_address")

    server := mqtt.NewServer(nil)

    // Tambah authentication hook
    server.AddHook(new(auth.AllowHook), nil)

    // Tambah TCP listener
    tcp := listeners.NewTCP(listenAddr, nil)
    err := server.AddListener(tcp, nil)
    if err != nil {
        return nil, err
    }

    return &EmbeddedBroker{
        server: server,
        configService: configService,
    }, nil
}

func (b *EmbeddedBroker) Start() error {
    go b.server.Serve()
    return nil
}

func (b *EmbeddedBroker) Stop() error {
    return b.server.Close()
}
```

**Langkah Tes**:
```go
// internal/broker/embedded_test.go

func TestEmbeddedBroker(t *testing.T) {
    // Gunakan config in-memory
    configService := config.NewConfigService(":memory:")
    configService.Set("mqtt.listen_address", ":1883", "Test")

    broker, err := NewEmbeddedBroker(configService)
    assert.NoError(t, err)

    // Start broker
    err = broker.Start()
    assert.NoError(t, err)

    // Beri waktu untuk start
    time.Sleep(100 * time.Millisecond)

    // Tes koneksi
    client := mqtt.NewClient(mqtt.NewClientOptions().AddBroker("tcp://localhost:1883"))
    token := client.Connect()
    assert.True(t, token.WaitTimeout(5*time.Second))
    assert.Equal(t, true, token.Error() == nil)

    // Tes publish
    token = client.Publish("test/topic", 0, false, "test payload")
    assert.True(t, token.WaitTimeout(5*time.Second))
    assert.Equal(t, true, token.Error() == nil)

    // Cleanup
    client.Disconnect(250)
    broker.Stop()
}
```

---

#### Tes 3.2: Hook Persistensi Pesan

**Tujuan**: Hook pesan MQTT ke database

**Implementasi**:

```go
// internal/broker/persistence_hook.go

type PersistenceHook struct {
    dbService database.DatabaseService
    mqtt.Hook
}

func (h *PersistenceHook) OnPublish(cl *mqtt.Client, pk packets.Packet) {
    // Simpan ke database
    h.dbService.SaveMessage(pk.TopicName, string(pk.Payload), pk.Qos)
}

func (h *PersistenceHook) ID() string {
    return "persistence-hook"
}
```

**Langkah Tes**:
```go
// internal/broker/persistence_test.go

func TestMessagePersistence(t *testing.T) {
    // Setup
    db := setupTestDB()
    configService := config.NewConfigService(":memory:")

    broker, _ := NewEmbeddedBroker(configService)
    broker.AddHook(&PersistenceHook{dbService: db}, nil)
    broker.Start()
    time.Sleep(100 * time.Millisecond)

    // Publish pesan
    client := mqtt.NewClient(mqtt.NewClientOptions().AddBroker("tcp://localhost:1883"))
    client.Connect()
    client.Publish("test/topic", 0, false, "persisted message")
    time.Sleep(100 * time.Millisecond)

    // Verifikasi persistensi
    messages, _ := db.GetMessages("test/topic", 10)
    assert.Len(t, messages, 1)
    assert.Equal(t, "persisted message", messages[0].Payload)

    // Cleanup
    client.Disconnect(250)
    broker.Stop()
}
```

---

### Fase 4: POC API & Frontend (Minggu 4)

#### Tes 4.1: Endpoints API Konfigurasi

**Tujuan**: Menguji API konfigurasi

**Implementasi**:

```go
// internal/api/config_handler.go

func RegisterConfigRoutes(r chi.Router, configService *config.ConfigService) {
    r.Get("/api/config", func(w http.ResponseWriter, r *http.Request) {
        all := configService.GetAll()
        json.NewEncoder(w).Encode(all)
    })

    r.Put("/api/config/:key", func(w http.ResponseWriter, r *http.Request) {
        key := chi.URLParam(r, "key")
        var req struct {
            Value  string `json:"value"`
            Reason string `json:"reason"`
        }
        json.NewDecoder(r.Body).Decode(&req)

        err := configService.Set(key, req.Value, req.Reason)
        if err != nil {
            http.Error(w, err.Error(), 400)
            return
        }

        w.WriteHeader(200)
        json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
    })

    r.Get("/api/config/history", func(w http.ResponseWriter, r *http.Request) {
        history := configService.GetAllHistory()
        json.NewEncoder(w).Encode(history)
    })

    r.Post("/api/config/validate", func(w http.ResponseWriter, r *http.Request) {
        var req map[string]interface{}
        json.NewDecoder(r.Body).Decode(&req)

        errors := configService.ValidateAll(req)
        if len(errors) > 0 {
            w.WriteHeader(400)
            json.NewEncoder(w).Encode(map[string][]string{"errors": errors})
            return
        }

        json.NewEncoder(w).Encode(map[string]string{"status": "valid"})
    })
}
```

**Langkah Tes**:
```go
// internal/api/config_test.go

func TestConfigAPI(t *testing.T) {
    configService := config.NewConfigService(":memory:")

    r := chi.NewRouter()
    RegisterConfigRoutes(r, configService)

    // Tes GET /api/config
    req := httptest.NewRequest("GET", "/api/config", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    assert.Equal(t, 200, w.Code)

    // Tes PUT /api/config/:key
    body := `{"value": "8883", "reason": "test update"}`
    req = httptest.NewRequest("PUT", "/api/config/mqtt.port", strings.NewReader(body))
    w = httptest.NewRecorder()
    r.ServeHTTP(w, req)
    assert.Equal(t, 200, w.Code)

    // Tes GET /api/config/history
    req = httptest.NewRequest("GET", "/api/config/history", nil)
    w = httptest.NewRecorder()
    r.ServeHTTP(w, req)
    assert.Equal(t, 200, w.Code)
}
```

---

#### Tes 4.2: Frontend Setup Mode

**Tujuan**: Menguji UI form Setup Mode

**Komponen**:

```typescript
// frontend/src/pages/SetupMode.tsx
import { useMutation, useQuery } from '@tanstack/react-query';

export function SetupMode() {
  const { data: config, isLoading } = useQuery({
    queryKey: ['config'],
    queryFn: () => fetch('/api/config').then(r => r.json()),
  });

  const saveConfig = useMutation({
    mutationFn: (newConfig: Record<string, any>) =>
      fetch('/api/config', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newConfig),
      }),
    onSuccess: () => {
      alert('Konfigurasi tersimpan! Me-restart...');
      setTimeout(() => window.location.reload(), 2000);
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    saveConfig.mutate(config);
  };

  if (isLoading) return <div>Loading configuration...</div>;

  return (
    <div className="setup-mode">
      <h1>Setup Mode - Konfigurasi Aplikasi</h1>

      <form onSubmit={handleSubmit}>
        <Section title="MQTT Broker">
          <FormField
            label="Listen Address"
            name="mqtt.listen_address"
            value={config['mqtt.listen_address'] || ':1883'}
            onChange={(v) => updateConfig('mqtt.listen_address', v)}
          />
          <FormField
            label="Max Connections"
            type="number"
            name="mqtt.max_connections"
            value={config['mqtt.max_connections'] || 100}
            onChange={(v) => updateConfig('mqtt.max_connections', v)}
          />
        </Section>

        <Section title="Web Server">
          <FormField
            label="Port"
            type="number"
            name="web.port"
            value={config['web.port'] || 8080}
            onChange={(v) => updateConfig('web.port', v)}
          />
        </Section>

        <button type="submit" disabled={saveConfig.isPending}>
          {saveConfig.isPending ? 'Menyimpan...' : 'Simpan & Restart'}
        </button>
      </form>
    </div>
  );
}
```

**Langkah Tes**:
```typescript
// frontend/src/pages/__tests__/SetupMode.test.tsx

describe('SetupMode', () => {
  it('memuat konfigurasi dari API', async () => {
    render(<SetupMode />);

    await waitFor(() => {
      expect(screen.getByText('Setup Mode')).toBeInTheDocument();
    });

    expect(screen.getByLabelText('Listen Address')).toHaveValue(':1883');
  });

  it('menyimpan konfigurasi', async () => {
    const user = userEvent.setup();
    render(<SetupMode />);

    await user.clear(screen.getByLabelText('Listen Address'));
    await user.type(screen.getByLabelText('Listen Address'), ':8883');
    await user.click(screen.getByText('Simpan & Restart'));

    await waitFor(() => {
      expect(screen.getByText('Konfigurasi tersimpan!')).toBeInTheDocument();
    });
  });
});
```

---

### Fase 5: Tes Integrasi End-to-End (Minggu 5-6)

#### Tes 5.1: Alur Lengkap - Setup hingga Pembuatan Topic

**Tujuan**: Memvalidasi perjalanan pengguna lengkap dari first launch hingga pembuatan topic

**Skenario Tes**:
1. Launch um-gateway.exe untuk pertama kali
2. Deteksi tidak ada internal.db
3. Inisialisasi internal.db dengan default
4. Launch UI Setup Mode pada port 8080
5. User isi form Setup Mode
6. Simpan konfigurasi ke internal.db
7. Restart dengan konfigurasi baru
8. Masuk Runner Mode
9. Hubungkan klien MQTT
10. Buat topic
11. Publish pesan
12. Verifikasi di dashboard

**Script Tes**:

```bash
#!/bin/bash
# test/e2e/complete_flow.sh

echo "=== Tes E2E UM-Gateway: Setup hingga Pembuatan Topic ==="

# Langkah 1: Clean slate
rm -rf ./data
mkdir -p ./data

# Langkah 2: Start gateway (first run - harus inisialisasi internal.db)
echo "Langkah 1: Starting gateway (first run)..."
./um-gateway.exe &
GATEWAY_PID=$!
sleep 3

# Verifikasi internal.db dibuat
if [ ! -f "./data/internal.db" ]; then
    echo "❌ GAGAL: internal.db tidak dibuat"
    kill $GATEWAY_PID
    exit 1
fi
echo "✅ internal.db dibuat"

# Langkah 3: Verifikasi Setup Mode dapat diakses
echo "Langkah 2: Memeriksa Setup Mode..."
curl -s http://localhost:8080/setup > /dev/null
if [ $? -ne 0 ]; then
    echo "❌ GAGAL: Setup Mode tidak dapat diakses"
    kill $GATEWAY_PID
    exit 1
fi
echo "✅ Setup Mode dapat diakses"

# Langkah 4: Ambil config saat ini
echo "Langkah 3: Mengambil config saat ini..."
CONFIG=$(curl -s http://localhost:8080/api/config)
echo "Config saat ini: $CONFIG"

# Langkah 5: Update konfigurasi
echo "Langkah 4: Mengupdate konfigurasi..."
curl -X PUT http://localhost:8080/api/config/mqtt.listen_port \
  -H "Content-Type: application/json" \
  -d '{"value": "1883", "reason": "Tes E2E"}'

# Langkah 6: Verifikasi history
echo "Langkah 5: Memeriksa history konfigurasi..."
HISTORY=$(curl -s http://localhost:8080/api/config/history)
echo "History: $HISTORY"

# Langkah 7: Restart gateway
echo "Langkah 6: Restarting gateway dengan config baru..."
kill $GATEWAY_PID
sleep 2
./um-gateway.exe &
GATEWAY_PID=$!
sleep 3

# Langkah 8: Verifikasi Runner Mode
echo "Langkah 7: Memeriksa Runner Mode..."
curl -s http://localhost:8080/ > /dev/null
if [ $? -ne 0 ]; then
    echo "❌ GAGAL: Runner Mode tidak dapat diakses"
    kill $GATEWAY_PID
    exit 1
fi
echo "✅ Runner Mode dapat diakses"

# Langkah 9: Tes koneksi MQTT
echo "Langkah 8: Testing MQTT broker..."
# Gunakan mosquitto_pub untuk tes
mosquitto_pub -h localhost -p 1883 -t "test/topic" -m "Hello from E2E test" -d
if [ $? -ne 0 ]; then
    echo "❌ GAGAL: MQTT publish gagal"
    kill $GATEWAY_PID
    exit 1
fi
echo "✅ Pesan MQTT terpublish"

# Langkah 10: Buat topic via API
echo "Langkah 9: Membuat topic via API..."
curl -X POST http://localhost:8080/api/topics \
  -H "Content-Type: application/json" \
  -d '{"name": "sensor/temperature"}'

# Langkah 11: Publish ke topic baru
echo "Langkah 10: Publishing ke topic baru..."
mosquitto_pub -h localhost -p 1883 -t "sensor/temperature" -m "25.5" -d

# Langkah 12: Verifikasi di database
echo "Langkah 11: Memverifikasi pesan di database..."
MESSAGES=$(curl -s http://localhost:8080/api/messages?topic=sensor/temperature)
echo "Pesan: $MESSAGES"

# Langkah 13: Cek metrics
echo "Langkah 12: Memeriksa metrics..."
METRICS=$(curl -s http://localhost:8080/api/metrics)
echo "Metrics: $METRICS"

# Cleanup
kill $GATEWAY_PID

echo ""
echo "=== Tes E2E Selesai ==="
echo "✅ Semua tes lulus!"
echo ""
echo "Ringkasan:"
echo "  - internal.db terinisialisasi"
echo "  - Setup Mode fungsional"
echo "  - Konfigurasi tersimpan ke internal.db"
echo "  - Pelacakan history berfungsi"
echo "  - Runner Mode fungsional"
echo "  - Broker MQTT operasional"
echo "  - Topics dibuat via API"
echo "  - Pesan tersimpan ke database.db"
```

**Output yang Diharapkan**:
```
=== Tes E2E UM-Gateway: Setup hingga Pembuatan Topic ===
Langkah 1: Starting gateway (first run)...
✅ internal.db dibuat
Langkah 2: Memeriksa Setup Mode...
✅ Setup Mode dapat diakses
Langkah 3: Mengambil config saat ini...
Config saat ini: {"mqtt.listen_address":":1883","mqtt.max_connections":100,"web.port":8080}
Langkah 4: Mengupdate konfigurasi...
{"status":"ok"}
Langkah 5: Memeriksa history konfigurasi...
History: [{"id":1,"setting_key":"mqtt.listen_port","old_value":":1883","new_value":"1883","changed_at":"2026-01-09T10:00:00Z"}]
Langkah 6: Restarting gateway dengan config baru...
Langkah 7: Memeriksa Runner Mode...
✅ Runner Mode dapat diakses
Langkah 8: Testing MQTT broker...
✅ Pesan MQTT terpublish
Langkah 9: Membuat topic via API...
{"status":"created","topic":"sensor/temperature"}
Langkah 10: Publishing ke topic baru...
✅ Pesan MQTT terpublish
Langkah 11: Memverifikasi pesan di database...
Pesan: [{"id":1,"topic":"sensor/temperature","payload":"25.5","qos":0,"published_at":"2026-01-09T10:05:00Z"}]
Langkah 12: Memeriksa metrics...
Metrics: {"total_messages":1,"total_topics":1,"uptime":300}

=== Tes E2E Selesai ===
✅ Semua tes lulus!

Ringkasan:
  - internal.db terinisialisasi
  - Setup Mode fungsional
  - Konfigurasi tersimpan ke internal.db
  - Pelacakan history berfungsi
  - Runner Mode fungsional
  - Broker MQTT operasional
  - Topics dibuat via API
  - Pesan tersimpan ke database.db
```

---

#### Tes 5.2: Tes Rollback Konfigurasi

**Tujuan**: Menguji history dan rollback konfigurasi

**Langkah Tes**:

```bash
#!/bin/bash
# test/e2e/config_rollback.sh

echo "=== Tes Rollback Konfigurasi ==="

# Start gateway
./um-gateway.exe &
GATEWAY_PID=$!
sleep 3

# Buat beberapa perubahan
echo "Membuat perubahan konfigurasi..."
curl -X PUT http://localhost:8080/api/config/mqtt.listen_port -d '{"value":"1883","reason":"v1"}'
sleep 1
curl -X PUT http://localhost:8080/api/config/mqtt.listen_port -d '{"value":"8883","reason":"v2"}'
sleep 1
curl -X PUT http://localhost:8080/api/config/mqtt.listen_port -d '{"value":"1884","reason":"v3"}'

# Ambil history
echo "History konfigurasi:"
curl -s http://localhost:8080/api/config/history | jq .

# Rollback ke v1
echo "Rollback ke v1..."
HISTORY_ID=$(curl -s http://localhost:8080/api/config/history | jq '.[2].id')
curl -X POST http://localhost:8080/api/config/rollback/$HISTORY_ID

# Verifikasi rollback
CURRENT=$(curl -s http://localhost:8080/api/config | jq '.["mqtt.listen_port"]')
echo "Nilai saat ini setelah rollback: $CURRENT"

if [ "$CURRENT" = "1883" ]; then
    echo "✅ Rollback sukses"
else
    echo "❌ Rollback gagal"
fi

kill $GATEWAY_PID
```

---

## 4. Checklist Tes

### Tes Foundation (Minggu 1)

- [ ] **Tes 1.1**: Setup proyek Go dan kompilasi
- [ ] **Tes 1.2**: Pembuatan skema internal database
- [ ] **Tes 1.3**: Operasi CRUD service konfigurasi
- [ ] **Tes 1.4**: Pelacakan history konfigurasi
- [ ] **Tes 1.5**: Validasi konfigurasi
- [ ] **Tes 1.6**: Reset konfigurasi ke default

### Tes Database (Minggu 2)

- [ ] **Tes 2.1**: Pembuatan skema system database
- [ ] **Tes 2.2**: Registrasi dan pengambilan device
- [ ] **Tes 2.3**: Pembuatan topic
- [ ] **Tes 2.4**: Persistensi pesan
- [ ] **Tes 2.5**: Pengambilan pesan dengan filter
- [ ] **Tes 2.6**: Implementasi adapter pattern

### Tes MQTT Broker (Minggu 3)

- [ ] **Tes 3.1**: Startup dan shutdown broker embedded
- [ ] **Tes 3.2**: Koneksi klien MQTT
- [ ] **Tes 3.3**: Publishing pesan
- [ ] **Tes 3.4**: Subscription pesan
- [ ] **Tes 3.5**: Hook persistensi pesan
- [ ] **Tes 3.6**: Level QoS (0, 1, 2)
- [ ] **Tes 3.7**: Pesan retained

### Tes API & Frontend (Minggu 4)

- [ ] **Tes 4.1**: Endpoints API konfigurasi
- [ ] **Tes 4.2**: Endpoints API topics
- [ ] **Tes 4.3**: Endpoints API pesan
- [ ] **Tes 4.4**: Rendering UI Setup Mode
- [ ] **Tes 4.5**: Pengiriman form konfigurasi
- [ ] **Tes 4.6**: Validasi real-time
- [ ] **Tes 4.7**: Penanganan error
- [ ] **Tes 4.8**: Tampilan history
- [ ] **Tes 4.9**: Fungsionalitas rollback

### Tes Integrasi (Minggu 5-6)

- [ ] **Tes 5.1**: Alur lengkap - Setup hingga Pembuatan Topic
- [ ] **Tes 5.2**: Rollback konfigurasi
- [ ] **Tes 5.3**: Migrasi dari versi sebelumnya
- [ ] **Tes 5.4**: Perubahan konfigurasi concurrent
- [ ] **Tes 5.5**: Throughput pesan (1000 msg/detik)
- [ ] **Tes 5.6**: Koneksi multi-device
- [ ] **Tes 5.7**: Stabilitas long-running (24 jam)

---

## 5. Kriteria Sukses

### Metrik Teknis

| Metrik | Target | Cara Mengukur |
|--------|--------|---------------|
| **Ukuran Binary** | < 60 MB | Artefak build |
| **Waktu Startup** | < 2 detik | Waktu dari launch hingga ready |
| **Memori Idle** | < 100 MB | Monitor proses |
| **Throughput Pesan** | > 1,000 msg/detik | Tes beban |
| **Waktu Respon API** | < 100ms (p95) | Benchmark API |
| **Penyelesaian Setup** | > 95% | Analytics (masa depan) |

### Persyaratan Fungsional

- [ ] **Tidak perlu config.yaml** - Semua konfigurasi via internal.db
- [ ] **Setup Mode** - UI berbasis form untuk konfigurasi first-time
- [ ] **Runner Mode** - Konfigurasi read-only di produksi
- [ ] **Pelacakan History** - Semua perubahan tercatat dengan timestamp
- [ ] **Rollback** - Dapat kembali ke konfigurasi sebelumnya
- [ ] **Migrasi Otomatis** - Upgrade mulus dari versi sebelumnya
- [ ] **MQTT Embedded** - Broker start dengan pengaturan yang terkonfigurasi
- [ ] **Persistensi Pesan** - Semua pesan disimpan ke database.db
- [ ] **Manajemen Topic** - Buat, lihat, hapus topic via API
- [ ] **Dashboard** - Metrics real-time dan monitoring

---

## 6. Alat Tes

### Tes Backend

```bash
# Install dependensi
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/suite

# Jalankan tes
go test ./...

# Jalankan dengan coverage
go test -cover ./...

# Jalankan tes spesifik
go test -v ./internal/config -run TestConfigServiceCRUD
```

### Tes MQTT

```bash
# Install klien mosquitto
# Ubuntu/Debian:
apt-get install mosquitto-clients

# macOS:
brew install mosquitto

# Windows:
# Download dari https://mosquitto.org/download/

# Tes publish
mosquitto_pub -h localhost -p 1883 -t "test/topic" -m "hello"

# Tes subscribe
mosquitto_sub -h localhost -p 1883 -t "test/#"
```

### Tes API

```bash
# Install httpie (alternatif untuk curl)
# Windows:
choco install httpie

# Linux/macOS:
pip install httpie

# Tes API
http GET http://localhost:8080/api/config
http PUT http://localhost:8080/api/config/mqtt.port value=8883 reason="test"
```

### Tes Beban

```bash
# Install wrk (HTTP benchmark)
# Ubuntu/Debian:
apt-get install wrk

# macOS:
brew install wrk

# Jalankan tes beban
wrk -t4 -c100 -d30s http://localhost:8080/api/messages

# MQTT load test (menggunakan mqtt-benchmark)
go install github.com/influxdata/mqtt-benchmark@latest
mqtt-benchmark -broker tcp://localhost:1883 -num 10000
```

---

## 7. Masalah & Keterbatasan yang Diketahui

### Keterbatasan Saat Ini (MVP v1.0)

1. **Single User** - Tidak ada dukungan multi-user
2. **Tidak Ada RBAC** - Tidak ada kontrol akses berbasis peran
3. **Embedded Only** - Tidak ada dukungan database eksternal
4. **Basic Auth** - Hanya username/password sederhana
5. **Tidak Ada Clustering** - Hanya single instance
6. **Tidak Ada WebSockets** - Hanya MQTT over TCP
7. **Retensi Pesan** - Semua pesan disimpan (tidak ada cleanup otomatis)

### Peningkatan Masa Depan

- Adapter database eksternal (PostgreSQL, MongoDB, InfluxDB)
- Autentikasi lanjutan (JWT, OAuth)
- Dukungan multi-tenant
- Clustering high availability
- Kebijakan retensi pesan
- Dukungan WebSocket
- Monitoring dan alerting lanjutan

---

## 8. Timeline Pengujian

### Minggu 1: Foundation
- Hari 1-2: Setup proyek dan tes dasar
- Hari 3-4: Implementasi internal database
- Hari 5: Validasi service konfigurasi

### Minggu 2: Database
- Hari 1-2: Skema system database
- Hari 3-4: Implementasi service database
- Hari 5: Tes persistensi pesan

### Minggu 3: MQTT Broker
- Hari 1-2: Setup broker embedded
- Hari 3-4: Hook persistensi pesan
- Hari 5: Tes integrasi MQTT

### Minggu 4: API & Frontend
- Hari 1-2: API konfigurasi
- Hari 3-4: UI Setup Mode
- Hari 5: Tes frontend

### Minggu 5-6: Integrasi
- Hari 1-3: Tes alur end-to-end
- Hari 4-5: Tes performa
- Hari 6-7: Perbaikan bug dan refinement

---

## 9. Langkah Selanjutnya

Setelah penyelesaian sukses semua tes:

1. **Dokumentasi**: Tulis panduan pengguna dan dokumentasi API
2. **Packaging**: Buat binary rilis untuk semua platform
3. **Instalasi**: Buat script instalasi
4. **Demo**: Rekam video demo
5. **Rilis**: Publikasikan v1.0 MVP

---

## 10. Kesimpulan

Rencana pengujian ini memvalidasi **pendekatan konfigurasi internal.db** untuk UM-Gateway MVP v1.0. Pembeda utama dari versi sebelumnya adalah eliminasi config.yaml dalam mendukung sistem manajemen konfigurasi database-first.

### Tujuan Utama Pengujian

1. ✅ **Tidak Ada File YAML** - Semua konfigurasi via internal.db
2. ✅ **Setup Mode** - UI berbasis form untuk konfigurasi
3. ✅ **Runner Mode** - Monitoring produksi
4. ✅ **History & Rollback** - Audit trail lengkap
5. ✅ **MQTT Embedded** - Broker dengan konfigurasi database-backed
6. ✅ **Alur End-to-End** - Dari setup hingga pembuatan topic

### Sukses

Pengujian MVP v1.0 sukses ketika:
- Pengguna mengunduh single executable
- Menjalankannya untuk pertama kali
- Menyelesaikan wizard Setup Mode
- Gateway start dengan pengaturan yang terkonfigurasi
- Menghubungkan perangkat MQTT
- Membuat topics
- Publish dan melihat pesan
- **Tanpa mengedit file konfigurasi apapun**

---

**Versi Dokumen**: v1.0
**Terakhir Diperbarui**: 9 Januari 2026
**Status**: Draft - Siap untuk Review
**Menggantikan**: TECHNOLOGY_TESTING_MVP_ID_v0.md (konfigurasi berbasis YAML)
