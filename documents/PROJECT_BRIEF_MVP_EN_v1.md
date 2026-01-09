# Project Plan: Universal MQTT Gateway Server (UM-Gateway) MVP v1.0

## Executive Summary

**Project Name**: Universal MQTT Gateway Server (UM-Gateway)
**Version**: MVP v1.0 (Minimum Viable Product with Internal Database)
**Status**: Planning
**Target Release**: Q2 2026 (3 months development)
**Organization**: Rafflesia Agro Technology Division

---

## Version History

| Version | Date | Changes | Author |
|---------|------|---------|--------|
| **MVP v1.0** | 2026-01-09 | **MVP with Internal Database**: Replacing config.yaml with internal.db<br>• Single executable with embedded MQTT<br>• Basic admin dashboard (React + Vite)<br>• Setup mode & Runner mode<br>• Form-based onboarding wizard<br>• Internal.db for application configuration (like PocketBase)<br>• Database.db for MQTT system data (default embedded SQLite)<br>• Adapter pattern for external database (post-MVP)<br>• Local cache (in-memory)<br>• Local embedded MQTT broker (mochi-mqtt)<br>• No YAML configuration files<br>• Configuration history & rollback<br>• Automatic migration<br>• Essential features only (no SaaS features)<br>• Development time: 3 months<br>• Budget: $68,000 | Development Team |

---

## 1. Background

### What is an MVP?

**MVP (Minimum Viable Product)**: A product with sufficient features to be used by early customers and provide feedback for further development.

**This is NOT**:
- ❌ Complete SaaS platform
- ❌ Multi-tenant system
- ❌ Cloud deployment
- ❌ Enterprise features
- ❌ Advanced monitoring

**This IS**:
- ✅ Reusable single executable
- ✅ Standalone deployment
- ✅ Basic configuration and management
- ✅ Essential MQTT gateway functionality
- ✅ Rapid development (3 months)
- ✅ No manual configuration files

### Key Changes from MVP v0.1

**v0.1**: Using config.yaml for configuration
**v1.0**: Using internal.db (embedded SQLite) for configuration

### Why Internal Database?

**Problems with config.yaml**:
1. Users must manually edit YAML files
2. YAML validation can be complex
3. No change history
4. Migration between versions is difficult
5. Backup separate from database

**Solution with internal.db**:
1. All configuration via Web UI / API
2. Validation at application level
3. Complete history with rollback
4. Automatic migration with database schema
5. Integrated backup

---

## 2. Business Model

### Value Proposition

**For Small Deployments**:
- Single binary, no complex setup
- No need to edit configuration files
- Runs on any computer
- Free and open source

**For Development**:
- Quick local testing
- Easy configuration (form-based UI)
- Built-in MQTT broker
- Web-based management

### Revenue Model

**MVP Phase (first 6 months)**:
- **100% Free**: Open source (MIT License)
- Focus on adoption and feedback
- No paid features yet

**Future Phase** (After MVP proves valuable):
- Pro license for advanced features
- External database support
- Support contracts

---

## 3. Target Audience

### Primary Users

**1. Small IoT Projects**
- Smart home automation (1-50 devices)
- Small farm monitoring
- Hobby projects
- Student projects

**2. Development Teams**
- IoT application testing
- Solution prototyping
- Local development environment

**3. System Integrators**
- Small customer deployments
- Proof of concept projects
- Simple gateway needs

---

## 4. Problem Statement

### Core Problem

Small IoT deployments need a **simple and standalone MQTT gateway** but existing solutions:
1. **Too Complex**: Require separate broker, database, cache
2. **Hard to Configure**: Require technical expertise / file editing
3. **Expensive**: Commercial solutions too expensive
4. **Over-engineered**: Complete SaaS platforms for simple needs

### MVP v1.0 Solution

**Single executable** with:
- ✅ Embedded MQTT broker (no separate broker)
- ✅ Internal.db for configuration (like PocketBase)
- ✅ Database.db for system data (embedded SQLite)
- ✅ In-memory cache (no external cache)
- ✅ Web admin interface (no file editing)
- ✅ Setup wizard (easy first-time configuration)

---

## 5. Proposed Solution

### Vision: MQTT Gateway with Internal Database

**One executable file** with two embedded databases:

**internal.db** (Application Configuration):
- Application settings
- MQTT broker configuration
- Database settings
- Web server settings
- Authentication credentials
- Change history
- Schema migrations

**database.db** (System Data):
- Device registry
- Topic definitions
- Message history
- Device status
- Metrics and statistics

### Architecture

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
│  │  │ Application Configuration                        │  │
│  │  │ • Settings (key-value)                          │ │  │
│  │  │ • Credentials (encrypted)                       │ │  │
│  │  │ • Change History                                │ │  │
│  │  │ • Schema Migrations                             │ │  │
│  │  └─────────────────────────────────────────────────┘ │  │
│  │                                                        │  │
│  │  ┌─────────────────────────────────────────────────┐ │  │
│  │  │ System Database (database.db)                    │ │  │
│  │  │ Runtime Data (Embedded SQLite)                   │  │
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
│  │                                                        │  │
│  │  ┌────────────────┐  ┌────────────────┐              │  │
│  │  │ Wizard         │  │ Config History │              │  │
│  │  │ Onboarding     │  │ & Rollback     │              │  │
│  │  └────────────────┘  └────────────────┘              │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

---

## 6. Database Schema

### Internal Database (internal.db)

**Purpose**: Store application configuration

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

**Purpose**: Store MQTT runtime data

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

## 7. Features (MVP Scope)

### F1: Embedded MQTT Broker

**Implementation**:

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

**Implementation**:

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

- Read-only configuration
- Monitoring dashboard
- Real-time metrics

### F5: Onboarding Wizard

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

## 8. Technology Stack

### Backend (Go)

| Component | Technology | Version |
|-----------|-----------|---------|
| Language | Go | 1.24+ |
| MQTT Broker | mochi-mqtt | 2.0+ |
| Internal DB | SQLite | 3.40+ |
| System DB | SQLite | 3.40+ |
| Cache | go-cache | Latest |
| HTTP Router | Chi | 5.0+ |

### Frontend (React)

| Component | Technology | Version |
|-----------|-----------|---------|
| Framework | React | 18.3+ |
| Build Tool | Vite | 5.0+ |
| Language | TypeScript | 5.3+ |
| State | TanStack Query | 5.0+ |
| Forms | React Hook Form | Latest |

---

## 9. Development Roadmap (3 Months)

### Month 1: Foundation

**Week 1-2: Project Setup**
- [ ] Setup repository
- [ ] Project structure
- [ ] CI/CD

**Week 3-4: Core Backend**
- [ ] Internal.db service
- [ ] System.db service
- [ ] Embedded MQTT broker
- [ ] In-memory cache

### Month 2: Features

**Week 5-6: MQTT & Database**
- [ ] Broker hooks
- [ ] Message persistence
- [ ] Device/topic tracking

**Week 7-8: Web UI**
- [ ] Setup React + Vite
- [ ] Setup mode form
- [ ] Runner mode dashboard

### Month 3: Polish & Release

**Week 9-10: Setup & Onboarding**
- [ ] Onboarding wizard
- [ ] Config history UI
- [ ] Validation

**Week 11-12: Testing & Release**
- [ ] E2E testing
- [ ] Bug fixes
- [ ] Documentation
- [ ] Build & release

---

## 10. Success Metrics

### Technical

| Metric | Target |
|--------|--------|
| Binary Size | < 60 MB |
| Startup Time | < 2 seconds |
| Idle Memory | < 100 MB |
| Messages/second | > 1,000 |

### Adoption

| Metric | Target |
|--------|--------|
| GitHub Stars | 100 (1 month) |
| Downloads | 500 (3 months) |
| Active Installations | 50 (3 months) |

---

## 11. Resource Requirements

### Budget Estimation

| Category | Cost | Duration |
|----------|------|----------|
| Full-Stack Developer | $15,000/month | 3 months |
| QA (part-time) | $2,000/month | 3 months |
| Technical Writer | $1,000/month | 3 months |
| Infrastructure | $500/month | 3 months |
| Tools | $200/month | 3 months |
| **Subtotal** | **$18,700/month × 3** | **$56,100** |
| **Buffer (20%)** | **$11,220** | - |
| **Total** | **$67,320** | **~$68,000** |

---

## 12. Conclusion

### MVP v1.0 Philosophy

**No Manual Configuration Files**:
- ✅ Internal.db for all configuration
- ✅ Form-based UI (not YAML editor)
- ✅ Complete history with rollback
- ✅ Automatic migration

**Focus on Essentials**:
- ✅ Single executable
- ✅ Embedded MQTT broker
- ✅ Two databases: internal.db + database.db
- ✅ Easy setup
- ✅ Basic monitoring

### Success Criteria

**MVP succeeds if**:
1. Users can download single executable
2. Complete onboarding wizard
3. Connect MQTT devices
4. View messages in dashboard
5. No need to edit configuration files

---

**Document Version**: MVP v1.0
**Last Updated**: January 9, 2026
**Status**: Draft
**Development Time**: 3 months
**Budget**: $68,000

---

**Prepared by**: Development Team
**Approved by**: [Pending]
**Target Release**: Q2 2026
