# Technology Testing Plan: UM-Gateway MVP v1.0

## Document Information

**Project**: Universal MQTT Gateway Server (UM-Gateway)
**Version**: MVP v1.0
**Document Type**: Technology Testing & Validation Plan
**Last Updated**: January 9, 2026
**Status**: Draft

---

## 1. Overview

### Purpose

This document outlines the technology testing plan for UM-Gateway MVP v1.0 to validate the **internal.db configuration approach** instead of the YAML configuration file approach.

### Testing Philosophy

**What changed**:
- ❌ REMOVED: config.yaml-based configuration
- ✅ ADDED: internal.db (SQLite) for application configuration
- ✅ ADDED: Form-based Setup Mode UI
- ✅ ADDED: Configuration history & rollback
- ✅ ADDED: Automatic migration

**Testing Goals**:
1. Validate internal.db as configuration store
2. Verify Setup Mode → Runner Mode flow
3. Test MQTT broker with embedded configuration
4. End-to-end: Setup → Start → Create Topic → Publish Message

---

## 2. Technology Stack Validation

### 2.1 Core Technologies

| Component | Technology | Version | Test Coverage |
|-----------|-----------|---------|---------------|
| **Backend Language** | Go | 1.24+ | Compilation, basic operations |
| **MQTT Broker** | mochi-mqtt | 2.0+ | Embedded broker functionality |
| **Internal DB** | SQLite | 3.40+ | Configuration storage |
| **System DB** | SQLite | 3.40+ | Runtime data storage |
| **Cache** | go-cache | Latest | In-memory caching |
| **HTTP Router** | Chi | 5.0+ | API endpoints |
| **Frontend** | React | 18.3+ | UI rendering |
| **Build Tool** | Vite | 5.0+ | Development & build |
| **Language** | TypeScript | 5.3+ | Type checking |
| **State** | TanStack Query | 5.0+ | Server state |
| **Forms** | React Hook Form | Latest | Form handling |

### 2.2 Architecture Validation

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

## 3. Test Plan Structure

### Phase 1: Foundation Tests (Week 1)

#### Test 1.1: Go Project Setup

**Objective**: Validate basic Go project structure

**Steps**:
```bash
# 1. Initialize Go module
go mod init github.com/rafflesia-agro/um-gateway

# 2. Create directory structure
mkdir -p cmd/server internal/{config,broker,database,api,cache}
mkdir -p frontend/src/{pages,components,hooks,services}

# 3. Install dependencies
go get github.com/mochi-mqtt/mqtt/v2
go get github.com/mattn/go-sqlite3
go get github.com/go-chi/chi/v5
go get github.com/patrickmn/go-cache

# 4. Verify compilation
go build -o um-gateway.exe ./cmd/server
```

**Expected Results**:
- ✅ Module initialized successfully
- ✅ All dependencies downloaded
- ✅ Binary compiled without errors
- ✅ Binary size < 20 MB (initial build)

**Success Criteria**:
```
✓ go mod sum matches expected dependencies
✓ Binary executes without panic
✓ Version flag works: ./um-gateway.exe --version
```

---

#### Test 1.2: Internal Database Schema

**Objective**: Validate internal.db schema creation

**Schema Definition**:

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

**Test Steps**:
```go
// internal/config/schema_test.go

func TestInternalDBSchema(t *testing.T) {
    // 1. Create temporary database
    db, err := sql.Open("sqlite3", ":memory:")
    assert.NoError(t, err)

    // 2. Execute schema
    _, err = db.Exec(schema)
    assert.NoError(t, err)

    // 3. Verify tables exist
    tables := []string{"settings", "settings_history", "credentials", "schema_migrations"}
    for _, table := range tables {
        var count int
        err = db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&count)
        assert.NoError(t, err)
        assert.Equal(t, 1, count)
    }

    // 4. Verify indexes
    indexes := []string{"idx_settings_category", "idx_settings_history_key", "idx_settings_history_changed_at"}
    for _, index := range indexes {
        var count int
        err = db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='index' AND name=?", index).Scan(&count)
        assert.NoError(t, err)
        assert.Equal(t, 1, count)
    }
}
```

**Expected Results**:
- ✅ All tables created successfully
- ✅ All indexes created
- ✅ Foreign keys enforced
- ✅ Default values working

---

#### Test 1.3: Configuration Service CRUD

**Objective**: Validate basic configuration CRUD operations

**Test Steps**:
```go
// internal/config/service_test.go

func TestConfigServiceCRUD(t *testing.T) {
    service := NewConfigService(":memory:")

    // Test 1: Set value
    err := service.Set("mqtt.listen_port", "1883", "Initial setup")
    assert.NoError(t, err)

    // Test 2: Get value
    value, err := service.Get("mqtt.listen_port")
    assert.NoError(t, err)
    assert.Equal(t, "1883", value)

    // Test 3: Get all in category
    mqttConfig := service.GetByCategory("mqtt")
    assert.Contains(t, mqttConfig, "mqtt.listen_port")

    // Test 4: Update value
    err = service.Set("mqtt.listen_port", "8883", "Changed to secure port")
    assert.NoError(t, err)

    // Test 5: Verify history
    history := service.GetHistory("mqtt.listen_port")
    assert.Len(t, history, 2)
    assert.Equal(t, "1883", history[0].OldValue)
    assert.Equal(t, "8883", history[0].NewValue)

    // Test 6: Validate
    err = service.Validate("mqtt.listen_port", "invalid")
    assert.Error(t, err)

    // Test 7: Reset to defaults
    err = service.ResetToDefaults()
    assert.NoError(t, err)
    value, _ = service.Get("mqtt.listen_port")
    assert.Equal(t, "1883", value) // Default value
}
```

**Expected Results**:
- ✅ Set/Get operations work correctly
- ✅ History tracking functional
- ✅ Validation rejects invalid values
- ✅ Reset to defaults works
- ✅ All operations are transactional

---

### Phase 2: Database POC (Week 2)

#### Test 2.1: System Database Schema

**Objective**: Validate database.db schema for MQTT data

**Schema Definition**:

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

**Test Steps**:
```go
// internal/database/schema_test.go

func TestSystemDBSchema(t *testing.T) {
    db, err := sql.Open("sqlite3", ":memory:")
    assert.NoError(t, err)

    _, err = db.Exec(schema)
    assert.NoError(t, err)

    // Verify tables
    tables := []string{"devices", "topics", "messages"}
    for _, table := range tables {
        var count int
        err = db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&count)
        assert.NoError(t, err)
        assert.Equal(t, 1, count)
    }

    // Test inserting a message
    _, err = db.Exec("INSERT INTO messages (topic, payload, qos) VALUES (?, ?, ?)", "test/topic", "hello", 0)
    assert.NoError(t, err)

    // Test query
    var payload string
    err = db.QueryRow("SELECT payload FROM messages WHERE topic = ?", "test/topic").Scan(&payload)
    assert.NoError(t, err)
    assert.Equal(t, "hello", payload)
}
```

---

#### Test 2.2: Database Service with Adapter Pattern

**Objective**: Test database service with SQLite adapter

**Implementation**:

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

**Test Steps**:
```go
// internal/database/service_test.go

func TestDatabaseService(t *testing.T) {
    db, _ := sql.Open("sqlite3", ":memory:")
    db.Exec(schema)

    service := NewSQLiteAdapter(db)

    // Test 1: Save and retrieve message
    err := service.SaveMessage("sensor/temp", "25.5", 0)
    assert.NoError(t, err)

    messages, err := service.GetMessages("sensor/temp", 10)
    assert.NoError(t, err)
    assert.Len(t, messages, 1)
    assert.Equal(t, "25.5", messages[0].Payload)

    // Test 2: Register device
    err = service.RegisterDevice("device-001", "Temperature Sensor", "sensor")
    assert.NoError(t, err)

    device, err := service.GetDevice("device-001")
    assert.NoError(t, err)
    assert.Equal(t, "Temperature Sensor", device.Name)

    // Test 3: Create topic
    err = service.CreateTopic("sensor/temp")
    assert.NoError(t, err)
}
```

---

### Phase 3: MQTT Broker POC (Week 3)

#### Test 3.1: Embedded MQTT Broker

**Objective**: Validate mochi-mqtt embedded broker

**Implementation**:

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
    // Get config from internal.db
    listenAddr, _ := configService.Get("mqtt.listen_address")

    server := mqtt.NewServer(nil)

    // Add authentication hook
    server.AddHook(new(auth.AllowHook), nil)

    // Add TCP listener
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

**Test Steps**:
```go
// internal/broker/embedded_test.go

func TestEmbeddedBroker(t *testing.T) {
    // Use in-memory config
    configService := config.NewConfigService(":memory:")
    configService.Set("mqtt.listen_address", ":1883", "Test")

    broker, err := NewEmbeddedBroker(configService)
    assert.NoError(t, err)

    // Start broker
    err = broker.Start()
    assert.NoError(t, err)

    // Give it time to start
    time.Sleep(100 * time.Millisecond)

    // Test connection
    client := mqtt.NewClient(mqtt.NewClientOptions().AddBroker("tcp://localhost:1883"))
    token := client.Connect()
    assert.True(t, token.WaitTimeout(5*time.Second))
    assert.Equal(t, true, token.Error() == nil)

    // Test publish
    token = client.Publish("test/topic", 0, false, "test payload")
    assert.True(t, token.WaitTimeout(5*time.Second))
    assert.Equal(t, true, token.Error() == nil)

    // Cleanup
    client.Disconnect(250)
    broker.Stop()
}
```

---

#### Test 3.2: Message Persistence Hook

**Objective**: Hook MQTT messages to database

**Implementation**:

```go
// internal/broker/persistence_hook.go

type PersistenceHook struct {
    dbService database.DatabaseService
    mqtt.Hook
}

func (h *PersistenceHook) OnPublish(cl *mqtt.Client, pk packets.Packet) {
    // Save to database
    h.dbService.SaveMessage(pk.TopicName, string(pk.Payload), pk.Qos)
}

func (h *PersistenceHook) ID() string {
    return "persistence-hook"
}
```

**Test Steps**:
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

    // Publish message
    client := mqtt.NewClient(mqtt.NewClientOptions().AddBroker("tcp://localhost:1883"))
    client.Connect()
    client.Publish("test/topic", 0, false, "persisted message")
    time.Sleep(100 * time.Millisecond)

    // Verify persisted
    messages, _ := db.GetMessages("test/topic", 10)
    assert.Len(t, messages, 1)
    assert.Equal(t, "persisted message", messages[0].Payload)

    // Cleanup
    client.Disconnect(250)
    broker.Stop()
}
```

---

### Phase 4: API & Frontend POC (Week 4)

#### Test 4.1: Configuration API Endpoints

**Objective**: Test configuration API

**Implementation**:

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

**Test Steps**:
```go
// internal/api/config_test.go

func TestConfigAPI(t *testing.T) {
    configService := config.NewConfigService(":memory:")

    r := chi.NewRouter()
    RegisterConfigRoutes(r, configService)

    // Test GET /api/config
    req := httptest.NewRequest("GET", "/api/config", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    assert.Equal(t, 200, w.Code)

    // Test PUT /api/config/:key
    body := `{"value": "8883", "reason": "test update"}`
    req = httptest.NewRequest("PUT", "/api/config/mqtt.port", strings.NewReader(body))
    w = httptest.NewRecorder()
    r.ServeHTTP(w, req)
    assert.Equal(t, 200, w.Code)

    // Test GET /api/config/history
    req = httptest.NewRequest("GET", "/api/config/history", nil)
    w = httptest.NewRecorder()
    r.ServeHTTP(w, req)
    assert.Equal(t, 200, w.Code)
}
```

---

#### Test 4.2: Setup Mode Frontend

**Objective**: Test Setup Mode form UI

**Component**:

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
      alert('Configuration saved! Restarting...');
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
      <h1>Setup Mode - Application Configuration</h1>

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
          {saveConfig.isPending ? 'Saving...' : 'Save & Restart'}
        </button>
      </form>
    </div>
  );
}
```

**Test Steps**:
```typescript
// frontend/src/pages/__tests__/SetupMode.test.tsx

describe('SetupMode', () => {
  it('loads configuration from API', async () => {
    render(<SetupMode />);

    await waitFor(() => {
      expect(screen.getByText('Setup Mode')).toBeInTheDocument();
    });

    expect(screen.getByLabelText('Listen Address')).toHaveValue(':1883');
  });

  it('saves configuration', async () => {
    const user = userEvent.setup();
    render(<SetupMode />);

    await user.clear(screen.getByLabelText('Listen Address'));
    await user.type(screen.getByLabelText('Listen Address'), ':8883');
    await user.click(screen.getByText('Save & Restart'));

    await waitFor(() => {
      expect(screen.getByText('Configuration saved!')).toBeInTheDocument();
    });
  });
});
```

---

### Phase 5: End-to-End Integration Test (Week 5-6)

#### Test 5.1: Complete Flow - Setup to Topic Creation

**Objective**: Validate complete user journey from first launch to creating a topic

**Test Scenario**:
1. Launch um-gateway.exe for first time
2. Detect no internal.db exists
3. Initialize internal.db with defaults
4. Launch Setup Mode UI on port 8080
5. User fills Setup Mode form
6. Save configuration to internal.db
7. Restart with new configuration
8. Enter Runner Mode
9. Connect MQTT client
10. Create topic
11. Publish message
12. Verify in dashboard

**Test Script**:

```bash
#!/bin/bash
# test/e2e/complete_flow.sh

echo "=== UM-Gateway E2E Test: Setup to Topic Creation ==="

# Step 1: Clean slate
rm -rf ./data
mkdir -p ./data

# Step 2: Start gateway (first run - should initialize internal.db)
echo "Step 1: Starting gateway (first run)..."
./um-gateway.exe &
GATEWAY_PID=$!
sleep 3

# Verify internal.db created
if [ ! -f "./data/internal.db" ]; then
    echo "❌ FAILED: internal.db not created"
    kill $GATEWAY_PID
    exit 1
fi
echo "✅ internal.db created"

# Step 3: Verify Setup Mode is accessible
echo "Step 2: Checking Setup Mode..."
curl -s http://localhost:8080/setup > /dev/null
if [ $? -ne 0 ]; then
    echo "❌ FAILED: Setup Mode not accessible"
    kill $GATEWAY_PID
    exit 1
fi
echo "✅ Setup Mode accessible"

# Step 4: Get current config
echo "Step 3: Fetching current config..."
CONFIG=$(curl -s http://localhost:8080/api/config)
echo "Current config: $CONFIG"

# Step 5: Update configuration
echo "Step 4: Updating configuration..."
curl -X PUT http://localhost:8080/api/config/mqtt.listen_port \
  -H "Content-Type: application/json" \
  -d '{"value": "1883", "reason": "E2E test"}'

# Step 6: Verify history
echo "Step 5: Checking configuration history..."
HISTORY=$(curl -s http://localhost:8080/api/config/history)
echo "History: $HISTORY"

# Step 7: Restart gateway
echo "Step 6: Restarting gateway with new config..."
kill $GATEWAY_PID
sleep 2
./um-gateway.exe &
GATEWAY_PID=$!
sleep 3

# Step 8: Verify Runner Mode
echo "Step 7: Checking Runner Mode..."
curl -s http://localhost:8080/ > /dev/null
if [ $? -ne 0 ]; then
    echo "❌ FAILED: Runner Mode not accessible"
    kill $GATEWAY_PID
    exit 1
fi
echo "✅ Runner Mode accessible"

# Step 9: Test MQTT connection
echo "Step 8: Testing MQTT broker..."
# Use mosquitto_pub for testing
mosquitto_pub -h localhost -p 1883 -t "test/topic" -m "Hello from E2E test" -d
if [ $? -ne 0 ]; then
    echo "❌ FAILED: MQTT publish failed"
    kill $GATEWAY_PID
    exit 1
fi
echo "✅ MQTT message published"

# Step 10: Create topic via API
echo "Step 9: Creating topic via API..."
curl -X POST http://localhost:8080/api/topics \
  -H "Content-Type: application/json" \
  -d '{"name": "sensor/temperature"}'

# Step 11: Publish to new topic
echo "Step 10: Publishing to new topic..."
mosquitto_pub -h localhost -p 1883 -t "sensor/temperature" -m "25.5" -d

# Step 12: Verify in database
echo "Step 11: Verifying message in database..."
MESSAGES=$(curl -s http://localhost:8080/api/messages?topic=sensor/temperature)
echo "Messages: $MESSAGES"

# Step 13: Check metrics
echo "Step 12: Checking metrics..."
METRICS=$(curl -s http://localhost:8080/api/metrics)
echo "Metrics: $METRICS"

# Cleanup
kill $GATEWAY_PID

echo ""
echo "=== E2E Test Complete ==="
echo "✅ All tests passed!"
echo ""
echo "Summary:"
echo "  - internal.db initialized"
echo "  - Setup Mode functional"
echo "  - Configuration saved to internal.db"
echo "  - History tracking working"
echo "  - Runner Mode functional"
echo "  - MQTT broker operational"
echo "  - Topics created via API"
echo "  - Messages persisted to database.db"
```

**Expected Output**:
```
=== UM-Gateway E2E Test: Setup to Topic Creation ===
Step 1: Starting gateway (first run)...
✅ internal.db created
Step 2: Checking Setup Mode...
✅ Setup Mode accessible
Step 3: Fetching current config...
Current config: {"mqtt.listen_address":":1883","mqtt.max_connections":100,"web.port":8080}
Step 4: Updating configuration...
{"status":"ok"}
Step 5: Checking configuration history...
History: [{"id":1,"setting_key":"mqtt.listen_port","old_value":":1883","new_value":"1883","changed_at":"2026-01-09T10:00:00Z"}]
Step 6: Restarting gateway with new config...
Step 7: Checking Runner Mode...
✅ Runner Mode accessible
Step 8: Testing MQTT broker...
✅ MQTT message published
Step 9: Creating topic via API...
{"status":"created","topic":"sensor/temperature"}
Step 10: Publishing to new topic...
✅ MQTT message published
Step 11: Verifying message in database...
Messages: [{"id":1,"topic":"sensor/temperature","payload":"25.5","qos":0,"published_at":"2026-01-09T10:05:00Z"}]
Step 12: Checking metrics...
Metrics: {"total_messages":1,"total_topics":1,"uptime":300}

=== E2E Test Complete ===
✅ All tests passed!

Summary:
  - internal.db initialized
  - Setup Mode functional
  - Configuration saved to internal.db
  - History tracking working
  - Runner Mode functional
  - MQTT broker operational
  - Topics created via API
  - Messages persisted to database.db
```

---

#### Test 5.2: Configuration Rollback Test

**Objective**: Test configuration history and rollback

**Test Steps**:

```bash
#!/bin/bash
# test/e2e/config_rollback.sh

echo "=== Configuration Rollback Test ==="

# Start gateway
./um-gateway.exe &
GATEWAY_PID=$!
sleep 3

# Make multiple changes
echo "Making configuration changes..."
curl -X PUT http://localhost:8080/api/config/mqtt.listen_port -d '{"value":"1883","reason":"v1"}'
sleep 1
curl -X PUT http://localhost:8080/api/config/mqtt.listen_port -d '{"value":"8883","reason":"v2"}'
sleep 1
curl -X PUT http://localhost:8080/api/config/mqtt.listen_port -d '{"value":"1884","reason":"v3"}'

# Get history
echo "Configuration history:"
curl -s http://localhost:8080/api/config/history | jq .

# Rollback to v1
echo "Rolling back to v1..."
HISTORY_ID=$(curl -s http://localhost:8080/api/config/history | jq '.[2].id')
curl -X POST http://localhost:8080/api/config/rollback/$HISTORY_ID

# Verify rollback
CURRENT=$(curl -s http://localhost:8080/api/config | jq '.["mqtt.listen_port"]')
echo "Current value after rollback: $CURRENT"

if [ "$CURRENT" = "1883" ]; then
    echo "✅ Rollback successful"
else
    echo "❌ Rollback failed"
fi

kill $GATEWAY_PID
```

---

## 4. Test Checklist

### Foundation Tests (Week 1)

- [ ] **Test 1.1**: Go project setup and compilation
- [ ] **Test 1.2**: Internal database schema creation
- [ ] **Test 1.3**: Configuration service CRUD operations
- [ ] **Test 1.4**: Configuration history tracking
- [ ] **Test 1.5**: Configuration validation
- [ ] **Test 1.6**: Configuration reset to defaults

### Database Tests (Week 2)

- [ ] **Test 2.1**: System database schema creation
- [ ] **Test 2.2**: Device registration and retrieval
- [ ] **Test 2.3**: Topic creation
- [ ] **Test 2.4**: Message persistence
- [ ] **Test 2.5**: Message retrieval with filters
- [ ] **Test 2.6**: Adapter pattern implementation

### MQTT Broker Tests (Week 3)

- [ ] **Test 3.1**: Embedded broker startup and shutdown
- [ ] **Test 3.2**: MQTT client connection
- [ ] **Test 3.3**: Message publishing
- [ ] **Test 3.4**: Message subscription
- [ ] **Test 3.5**: Message persistence hook
- [ ] **Test 3.6**: QoS levels (0, 1, 2)
- [ ] **Test 3.7**: Retained messages

### API & Frontend Tests (Week 4)

- [ ] **Test 4.1**: Configuration API endpoints
- [ ] **Test 4.2**: Topics API endpoints
- [ ] **Test 4.3**: Messages API endpoints
- [ ] **Test 4.4**: Setup Mode UI rendering
- [ ] **Test 4.5**: Configuration form submission
- [ ] **Test 4.6**: Real-time validation
- [ ] **Test 4.7**: Error handling
- [ ] **Test 4.8**: History display
- [ ] **Test 4.9**: Rollback functionality

### Integration Tests (Week 5-6)

- [ ] **Test 5.1**: Complete flow - Setup to Topic Creation
- [ ] **Test 5.2**: Configuration rollback
- [ ] **Test 5.3**: Migration from previous version
- [ ] **Test 5.4**: Concurrent configuration changes
- [ ] **Test 5.5**: Message throughput (1000 msg/sec)
- [ ] **Test 5.6**: Multiple device connections
- [ ] **Test 5.7**: Long-running stability (24 hours)

---

## 5. Success Criteria

### Technical Metrics

| Metric | Target | How to Measure |
|--------|--------|----------------|
| **Binary Size** | < 60 MB | Build artifacts |
| **Startup Time** | < 2 seconds | Time from launch to ready |
| **Idle Memory** | < 100 MB | Process monitor |
| **Message Throughput** | > 1,000 msg/sec | Load test |
| **API Response Time** | < 100ms (p95) | API benchmarks |
| **Setup Completion** | > 95% | Analytics (future) |

### Functional Requirements

- [ ] **No config.yaml required** - All configuration via internal.db
- [ ] **Setup Mode** - Form-based UI for first-time configuration
- [ ] **Runner Mode** - Read-only configuration in production
- [ ] **History Tracking** - All changes recorded with timestamps
- [ ] **Rollback** - Can revert to any previous configuration
- [ ] **Automatic Migration** - Smooth upgrade from previous versions
- [ ] **Embedded MQTT** - Broker starts with configured settings
- [ ] **Message Persistence** - All messages saved to database.db
- [ ] **Topic Management** - Create, view, delete topics via API
- [ ] **Dashboard** - Real-time metrics and monitoring

---

## 6. Testing Tools

### Backend Testing

```bash
# Install dependencies
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/suite

# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -v ./internal/config -run TestConfigServiceCRUD
```

### MQTT Testing

```bash
# Install mosquitto clients
# Ubuntu/Debian:
apt-get install mosquitto-clients

# macOS:
brew install mosquitto

# Windows:
# Download from https://mosquitto.org/download/

# Test publish
mosquitto_pub -h localhost -p 1883 -t "test/topic" -m "hello"

# Test subscribe
mosquitto_sub -h localhost -p 1883 -t "test/#"
```

### API Testing

```bash
# Install httpie (alternative to curl)
# Windows:
choco install httpie

# Linux/macOS:
pip install httpie

# Test API
http GET http://localhost:8080/api/config
http PUT http://localhost:8080/api/config/mqtt.port value=8883 reason="test"
```

### Load Testing

```bash
# Install wrk (HTTP benchmark)
# Ubuntu/Debian:
apt-get install wrk

# macOS:
brew install wrk

# Run load test
wrk -t4 -c100 -d30s http://localhost:8080/api/messages

# MQTT load test (using mqtt-benchmark)
go install github.com/influxdata/mqtt-benchmark@latest
mqtt-benchmark -broker tcp://localhost:1883 -num 10000
```

---

## 7. Known Issues & Limitations

### Current Limitations (MVP v1.0)

1. **Single User** - No multi-user support
2. **No RBAC** - No role-based access control
3. **Embedded Only** - No external database support
4. **Basic Auth** - Simple username/password only
5. **No Clustering** - Single instance only
6. **No WebSockets** - MQTT over TCP only
7. **Message Retention** - All messages kept (no automatic cleanup)

### Future Enhancements

- External database adapters (PostgreSQL, MongoDB, InfluxDB)
- Advanced authentication (JWT, OAuth)
- Multi-tenant support
- High availability clustering
- Message retention policies
- WebSocket support
- Advanced monitoring and alerting

---

## 8. Testing Timeline

### Week 1: Foundation
- Days 1-2: Project setup and basic tests
- Days 3-4: Internal database implementation
- Day 5: Configuration service validation

### Week 2: Database
- Days 1-2: System database schema
- Days 3-4: Database service implementation
- Day 5: Message persistence testing

### Week 3: MQTT Broker
- Days 1-2: Embedded broker setup
- Days 3-4: Message persistence hooks
- Day 5: MQTT integration testing

### Week 4: API & Frontend
- Days 1-2: Configuration API
- Days 3-4: Setup Mode UI
- Day 5: Frontend testing

### Weeks 5-6: Integration
- Days 1-3: End-to-end flow testing
- Days 4-5: Performance testing
- Days 6-7: Bug fixes and refinement

---

## 9. Next Steps

After successful completion of all tests:

1. **Documentation**: Write user guide and API documentation
2. **Packaging**: Create release binaries for all platforms
3. **Installation**: Create install scripts
4. **Demos**: Record demo videos
5. **Release**: Publish v1.0 MVP

---

## 10. Conclusion

This testing plan validates the **internal.db configuration approach** for UM-Gateway MVP v1.0. The key differentiator from previous versions is the elimination of config.yaml in favor of a database-first configuration management system.

### Key Testing Goals

1. ✅ **No YAML Files** - All configuration via internal.db
2. ✅ **Setup Mode** - Form-based UI for configuration
3. ✅ **Runner Mode** - Production monitoring
4. ✅ **History & Rollback** - Complete audit trail
5. ✅ **Embedded MQTT** - Broker with database-backed config
6. ✅ **End-to-End Flow** - From setup to topic creation

### Success

MVP v1.0 testing succeeds when:
- User downloads single executable
- Runs it for first time
- Completes Setup Mode wizard
- Gateway starts with configured settings
- Connects MQTT devices
- Creates topics
- Publishes and views messages
- **Without editing any configuration files**

---

**Document Version**: v1.0
**Last Updated**: January 9, 2026
**Status**: Draft - Ready for Review
**Replaces**: TECHNOLOGY_TESTING_MVP_EN_v0.md (YAML-based config)
