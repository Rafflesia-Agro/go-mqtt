# Project Brief: Universal MQTT Gateway Server (UM-Gateway) MVP v0.1

## Executive Summary

**Project Name**: Universal MQTT Gateway Server (UM-Gateway)
**Version**: MVP v0.1 (Minimum Viable Product)
**Status**: Planning
**Target Release**: Q2 2026 (3 months development)
**Organization**: Rafflesia Agro Technology Division

---

## Version History

| Version | Date | Changes | Author |
|---------|------|---------|--------|
| **MVP v0.1** | 2026-01-09 | **Minimum Viable Product**: Single executable with embedded MQTT<br>• Single binary deployment (embedded mode only)<br>• Basic admin dashboard (React + Vite)<br>• Setup mode & Runner mode<br>• Onboarding wizard<br>• Local cache (in-memory)<br>• Local database (SQLite WAL)<br>• Local embedded MQTT broker (mochi-mqtt)<br>• Basic configuration (config.yaml)<br>• Essential features only (no SaaS features)<br>• Development time: 3 months<br>• Budget: $65,000 | Development Team |

---

## 1. Background

### What is MVP?

**MVP (Minimum Viable Product)**: A product with just enough features to be usable by early customers and provide feedback for future development.

**This is NOT**:
- ❌ Full SaaS platform
- ❌ Multi-tenant system
- ❌ Cloud deployment
- ❌ Enterprise features
- ❌ Advanced monitoring

**This IS**:
- ✅ Single reusable executable
- ✅ Self-contained deployment
- ✅ Basic configuration and management
- ✅ Essential MQTT gateway functionality
- ✅ Quick to develop (3 months)

### Why MVP?

1. **Speed to Market**: Get working product in 3 months
2. **Focus**: Core features only, no distractions
3. **Simplicity**: Easy to understand, deploy, and maintain
4. **Validation**: Test core assumptions before building advanced features
5. **Lower Risk**: Less investment, faster feedback

### Market Need

Many small IoT deployments need a **simple, self-contained MQTT gateway**:
- Smart home enthusiasts
- Small farms (like Rafflesia Agro)
- Small industrial setups
- Development/testing environments
- Edge computing prototypes

They need something that **just works** without complex setup.

---

## 2. Business Model

### Value Proposition

**For Small Deployments**:
- Single binary, no complex setup
- Runs on any computer
- No external dependencies
- Free and open source

**For Development**:
- Quick local testing
- Easy configuration
- Built-in MQTT broker
- Web-based management

### Revenue Model

**MVP Phase (First 6 months)**:
- **100% Free**: Open source (MIT License)
- Focus on adoption and feedback
- No paid features yet

**Future Phase** (After MVP proves value):
- Pro license for advanced features
- Support contracts
- Custom development

### Target Scale

**MVP Goal**:
- 100 active installations within 3 months of release
- Positive feedback from early adopters
- Clear understanding of what features to add next

---

## 3. Target Audience

### Primary Users

**1. Small IoT Projects**
- Smart home automation (1-50 devices)
- Small farm monitoring (like poultry farms)
- Hobbyist projects
- Student projects

**2. Development Teams**
- Testing IoT applications
- Prototyping solutions
- Local development environment
- Integration testing

**3. System Integrators**
- Small customer deployments
- Proof of concept projects
- Simple gateway needs
- Quick deployments

### User Profile

**Technical Skills**:
- Basic understanding of MQTT
- Can run executable files
- Can configure simple settings
- May not be professional developers

**Deployment Environment**:
- Single machine (PC, server, Raspberry Pi)
- Local network
- No cloud connectivity required
- No external dependencies

---

## 4. Problem Statement

### Core Problem

Small IoT deployments need a **simple, self-contained MQTT gateway** but existing solutions are:
1. **Too Complex**: Require separate brokers, databases, caches
2. **Hard to Configure**: Need technical expertise
3. **Expensive**: Commercial solutions cost too much
4. **Over-engineered**: Full SaaS platforms for simple needs

### MVP Solution

**A single executable that includes**:
- ✅ Embedded MQTT broker (no separate broker needed)
- ✅ Local database (no external database)
- ✅ In-memory cache (no external cache)
- ✅ Web admin interface (no config file editing)
- ✅ Setup wizard (easy first-time configuration)

**What We're NOT Building (Yet)**:
- ❌ External database support (PostgreSQL, MongoDB)
- ❌ External cache support (Redis, Memcached)
- ❌ Cloud integration
- ❌ Multi-tenancy
- ❌ Advanced monitoring & alerting
- ❌ High availability & clustering
- ❌ User authentication & authorization
- ❌ API for third-party integration

---

## 5. Proposed Solution

### Vision: Simple, Self-Contained MQTT Gateway

**One executable file** that provides:
1. **Embedded MQTT Broker** - Handles MQTT connections
2. **Local Database** - Stores messages and configuration
3. **In-Memory Cache** - Improves performance
4. **Web Admin Interface** - Easy management
5. **Setup Wizard** - Guided first-time setup

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│           Single Executable: um-gateway.exe                 │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌───────────────────────────────────────────────────────┐  │
│  │              Go Backend (Single Process)               │  │
│  │                                                        │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌────────────┐ │  │
│  │  │ Embedded     │  │ HTTP API    │  │ In-Memory  │ │  │
│  │  │ MQTT Broker  │  │ Server      │  │ Cache      │ │  │
│  │  │ (mochi-mqtt) │  │             │  │ (go-cache) │ │  │
│  │  └──────────────┘  └──────────────┘  └────────────┘ │  │
│  │                                                        │  │
│  │  ┌──────────────┐  ┌──────────────┐                  │  │
│  │  │ SQLite       │  │ Config YAML  │                  │  │
│  │  │ (WAL Mode)   │  │ Loader       │                  │  │
│  │  └──────────────┘  └──────────────┘                  │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                               │
│  ┌───────────────────────────────────────────────────────┐  │
│  │         React + Vite Web UI (Embedded)                │  │
│  │                                                        │  │
│  │  ┌────────────────┐  ┌────────────────┐              │  │
│  │  │ Onboarding     │  │ Setup Mode     │              │  │
│  │  │ Wizard         │  │ (Config Edit)  │              │  │
│  │  └────────────────┘  └────────────────┘              │  │
│  │                                                        │  │
│  │  ┌────────────────┐  ┌────────────────┐              │  │
│  │  │ Runner Mode    │  │ Dashboard      │              │  │
│  │  │ (Monitoring)   │  │ (Read-only)    │              │  │
│  │  └────────────────┘  └────────────────┘              │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### Key Features

#### F1: Embedded MQTT Broker (mochi-mqtt)
- MQTT 3.1.1 protocol
- Up to 100 concurrent connections
- TCP listener on configurable port (default 1883)
- Basic authentication (username/password)
- Message persistence (optional)

#### F2: Local Database (SQLite WAL)
- Single file database
- Write-Ahead Logging for better concurrency
- Automatic schema migrations
- Stores: devices, topics, messages, metrics
- Backup/restore utilities

#### F3: In-Memory Cache (go-cache)
- Default TTL: 1 hour
- Maximum 10,000 items
- Improves query performance
- Automatic expiration

#### F4: Web Admin Interface (React + Vite)
- Modern, responsive UI
- Real-time metrics display
- Device management
- Topic browser
- Message history viewer
- Configuration viewer (read-only in runner mode)

#### F5: Setup Mode
- Edit config.yaml via web UI
- YAML editor with validation
- Form-based configuration
- Test configuration button
- Save & restart
- Configuration backup

#### F6: Runner Mode
- Production operation
- Read-only configuration
- Monitoring dashboard
- Real-time metrics
- Cannot edit config without restart

#### F7: Onboarding Wizard
- First-run detection
- Step-by-step configuration
- Basic settings:
  - Admin username/password
  - MQTT broker port
  - Database location
  - Web UI port
- Test connection
- Start runner mode

---

## 6. Project Scope & Boundaries

### In Scope (MVP Features)

#### Must Have (Core Functionality)

1. **Embedded MQTT Broker**
   - Basic broker functionality
   - TCP listener
   - Client authentication
   - Message routing
   - Optional persistence

2. **Local Database**
   - SQLite with WAL mode
   - Device registry
   - Topic definitions
   - Message storage
   - Basic metrics

3. **In-Memory Cache**
   - go-cache implementation
   - TTL-based expiration
   - Size limits

4. **Web Admin UI**
   - Dashboard with metrics
   - Device list
   - Topic browser
   - Message viewer
   - Config viewer (read-only in runner mode)

5. **Setup Mode**
   - YAML config editor
   - Form-based config
   - Validation
   - Save & restart

6. **Runner Mode**
   - Production operation
   - Read-only config
   - Monitoring

7. **Onboarding Wizard**
   - First-run detection
   - Basic configuration
   - Test connections
   - Start gateway

8. **Configuration (config.yaml)**
   - Database settings
   - MQTT broker settings
   - Web server settings
   - Logging settings
   - Basic retention policy

### Out of Scope (Future Versions)

**NOT in MVP**:
- ❌ External database support (PostgreSQL, MongoDB, InfluxDB)
- ❌ External cache support (Redis, Memcached)
- ❌ External MQTT broker connection
- ❌ User management (multi-user)
- ❌ Role-based access control
- ❌ API for third-party integration
- ❌ Advanced monitoring & alerting
- ❌ High availability & clustering
- ❌ Data export/import
- ❌ Configuration templates
- ❌ Advanced security (SSL/TLS)
- ❌ WebSocket support for MQTT
- ❌ MQTT 5.0 features
- ❌ Bridge configurations

### Technical Constraints

**MVP Constraints**:
- **Embedded mode ONLY** (no external broker/database/cache)
- **Single machine deployment** (no clustering)
- **SQLite database ONLY** (no PostgreSQL/MongoDB)
- **In-memory cache ONLY** (no Redis)
- **Basic authentication ONLY** (username/password in config)
- **Single admin user** (no multi-user support)
- **HTTP ONLY** (no HTTPS/TLS in MVP)

**Hardware Requirements**:
- Minimum: 1 CPU core, 1GB RAM, 100MB disk
- Recommended: 2 CPU cores, 2GB RAM, 1GB disk

**Software Requirements**:
- Go 1.24+
- SQLite 3.40+
- Modern web browser (Chrome, Firefox, Edge, Safari)

---

## 7. Technology Stack

### Backend (Go)

| Component | Technology | Version | Why |
|-----------|-----------|---------|-----|
| Language | Go | 1.24+ | Performance, single binary |
| MQTT Broker | mochi-mqtt | 2.0+ | Embedded, pure Go |
| Database | SQLite | 3.40+ | Single file, WAL mode |
| Cache | go-cache | Latest | In-memory, simple |
| HTTP Router | Chi | 5.0+ | Lightweight, idiomatic |
| YAML | yaml.v3 | Latest | Config parsing |
| Embed | embed | std | Embed UI in binary |

### Frontend (React)

| Component | Technology | Version | Why |
|-----------|-----------|---------|-----|
| Framework | React | 18.3+ | Popular, familiar |
| Build Tool | Vite | 5.0+ | Fast, simple |
| Language | TypeScript | 5.3+ | Type safety |
| UI Library | - | - | Plain CSS (no component lib) |
| HTTP | Fetch API | std | No Axios needed |
| Router | - | - | Simple hash routing |

### Dev Tools

| Tool | Purpose |
|------|---------|
| Git | Version control |
| Go modules | Dependency management |
| npm/yarn | Frontend dependencies |
| air | Hot reload for Go |
| Vite HMR | Hot reload for React |

---

## 8. Features (MVP Scope)

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

    // Add hooks for logging and auth
    hooks := &hooks.Auth{
        Authenticate: func(cl *mqtt.Client, user, pass string) bool {
            // Simple username/password check
            return user == config.Username && pass == config.Password
        },
    }
    server.AddHook(hooks)

    // TCP listener
    tcp := mqtt.NewTCPListener(config.ListenAddress, nil)
    server.AddListener(tcp)

    return &EmbeddedBroker{server: server, config: config}, nil
}

func (b *EmbeddedBroker) Start() error {
    go b.server.Serve()
    return nil
}

func (b *EmbeddedBroker) Stop() {
    b.server.Close()
}
```

**Config (config.yaml)**:

```yaml
mqtt:
  embedded:
    enabled: true
    listen_address: :1883
    username: admin
    password: changeme
    max_connections: 100
```

---

### F2: Local Database (SQLite)

**Schema**:

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

-- Create indexes
CREATE INDEX idx_messages_topic ON messages(topic);
CREATE INDEX idx_messages_published_at ON messages(published_at);
```

**Implementation**:

```go
package database

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

type SQLiteDatabase struct {
    db *sql.DB
}

func NewSQLiteDatabase(path string) (*SQLiteDatabase, error) {
    db, err := sql.Open("sqlite3", path)
    if err != nil {
        return nil, err
    }

    // Enable WAL mode
    _, err = db.Exec("PRAGMA journal_mode=WAL")
    if err != nil {
        return nil, err
    }

    // Run migrations
    if err := runMigrations(db); err != nil {
        return nil, err
    }

    return &SQLiteDatabase{db: db}, nil
}
```

---

### F3: In-Memory Cache

**Implementation**:

```go
package cache

import (
    "github.com/patrickmn/go-cache"
    "time"
)

type MemoryCache struct {
    store *cache.Cache
}

func NewMemoryCache() *MemoryCache {
    // Default TTL: 1 hour, cleanup every 10 minutes
    return &MemoryCache{
        store: cache.New(1*time.Hour, 10*time.Minute),
    }
}

func (c *MemoryCache) Get(key string) (interface{}, bool) {
    return c.store.Get(key)
}

func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) {
    if ttl == 0 {
        ttl = cache.DefaultExpiration
    }
    c.store.Set(key, value, ttl)
}

func (c *MemoryCache) Delete(key string) {
    c.store.Delete(key)
}
```

---

### F4: Web Admin UI

**Pages**:

1. **Dashboard** (`/`)
   - System metrics (connected clients, messages/min, uptime)
   - Quick stats (total devices, total topics, total messages)
   - Recent activity

2. **Devices** (`/devices`)
   - List of all devices
   - Device details
   - Search/filter

3. **Topics** (`/topics`)
   - List of all topics
   - Topic details
   - Message count per topic

4. **Messages** (`/messages`)
   - Message history (paginated)
   - Filter by topic
   - View message details

5. **Setup** (`/setup`)
   - YAML config editor
   - Form-based config
   - Save & restart

6. **Config** (`/config`)
   - Read-only config viewer (runner mode)
   - Download config

**Example Component**:

```typescript
// frontend/src/pages/Dashboard.tsx
import { useQuery } from '@tanstack/react-query';

export function Dashboard() {
  const { data: metrics, isLoading } = useQuery({
    queryKey: ['metrics'],
    queryFn: async () => {
      const res = await fetch('/api/metrics');
      return res.json();
    },
    refetchInterval: 5000, // Auto-refresh
  });

  if (isLoading) return <div>Loading...</div>;

  return (
    <div className="dashboard">
      <h1>Dashboard</h1>

      <div className="metrics-grid">
        <MetricCard title="Connected Clients" value={metrics.clients} />
        <MetricCard title="Messages/min" value={metrics.msgPerMin} />
        <MetricCard title="Uptime" value={metrics.uptime} />
      </div>

      <div className="stats-grid">
        <StatCard title="Total Devices" value={metrics.totalDevices} />
        <StatCard title="Total Topics" value={metrics.totalTopics} />
        <StatCard title="Total Messages" value={metrics.totalMessages} />
      </div>
    </div>
  );
}
```

---

### F5: Setup Mode

**Purpose**: Edit config.yaml

**Access**: `/setup` (only when mode=setup in config.yaml)

**Features**:
- YAML editor with syntax highlighting
- Form-based configuration
- Validation
- Test button
- Save & restart button
- Backup before saving

**Flow**:
```
1. Stop gateway
2. Start gateway with --mode=setup
3. Open http://localhost:8080/setup
4. Edit configuration
5. Click "Save & Restart"
6. Gateway restarts in runner mode
```

---

### F6: Runner Mode

**Purpose**: Production operation

**Access**: `/` (when mode=runner in config.yaml)

**Features**:
- Read-only config viewer
- Monitoring dashboard
- Real-time metrics
- Cannot edit config

**Behavior**:
```
1. Start gateway (or after setup)
2. config.yaml becomes read-only
3. Embedded MQTT broker starts
4. Web UI shows dashboard
5. Gateway accepts MQTT connections
```

---

### F7: Onboarding Wizard

**Purpose**: First-time setup

**Trigger**: No config.yaml found

**Steps**:

**Step 1: Welcome**
- Explain what the gateway does
- "Get Started" button

**Step 2: Admin Credentials**
- Username (default: admin)
- Password (required)
- Confirm password

**Step 3: MQTT Configuration**
- Listen port (default: 1883)
- Max connections (default: 100)

**Step 4: Web UI Configuration**
- Port (default: 8080)

**Step 5: Database Configuration**
- Path (default: ./data/gateway.db)

**Step 6: Review**
- Show all settings
- "Finish & Start" button

**Step 7: Complete**
- Gateway starts
- Show dashboard URL
- Show MQTT broker URL

---

## 9. Development Roadmap

### 3-Month Sprint Plan

#### Month 1: Foundation

**Week 1-2: Project Setup**
- [ ] Repository setup (Go + React)
- [ ] Basic project structure
- [ ] CI/CD (GitHub Actions)
- [ ] Development environment setup
- [ ] Basic Go server with embedded UI

**Week 3-4: Backend Core**
- [ ] Embedded MQTT broker (mochi-mqtt)
- [ ] SQLite database setup
- [ ] Basic API endpoints
- [ ] In-memory cache
- [ ] Config loading (YAML)

**Milestone**: Backend basics working

---

#### Month 2: Features

**Week 5-6: MQTT & Database**
- [ ] MQTT broker hooks (auth, logging)
- [ ] Message persistence
- [ ] Database schema & migrations
- [ ] Device tracking
- [ ] Topic tracking
- [ ] Message storage

**Week 7-8: Web UI**
- [ ] React + Vite setup
- [ ] Basic layout
- [ ] Dashboard page
- [ ] Devices page
- [ ] Topics page
- [ ] Messages page

**Milestone**: Core features working

---

#### Month 3: Polish & Release

**Week 9-10: Setup & Onboarding**
- [ ] Setup mode UI
- [ ] YAML config editor
- [ ] Onboarding wizard
- [ ] Config validation
- [ ] Save & restart functionality

**Week 11-12: Testing & Release**
- [ ] End-to-end testing
- [ ] Bug fixes
- [ ] Documentation
- [ ] Build for platforms (Linux, Windows, macOS)
- [ ] Release v0.1.0

**Milestone**: MVP release

---

## 10. Success Metrics

### Technical Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| Binary Size | < 50 MB | Build artifacts |
| Startup Time | < 2 seconds | Manual testing |
| Memory Usage | < 100 MB idle | Resource monitoring |
| MQTT Messages/sec | > 1,000 | Load testing |
| Concurrent Clients | 100 | Load testing |
| Web Page Load | < 1 second | Browser DevTools |

### Adoption Metrics

| Metric | Target | Timeline |
|--------|--------|----------|
| GitHub Stars | 100 | 1 month post-release |
| Downloads | 500 | 3 months post-release |
| Active Installations | 50 | 3 months post-release |
| Issues Reported | < 10 critical | 3 months post-release |
| Community Contributions | 5+ | 6 months post-release |

### Quality Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| Onboarding Completion Rate | > 90% | Analytics |
| Setup Success Rate | 100% | Manual testing |
| Documentation Coverage | 100% features | Manual review |
| Known Critical Bugs | 0 | GitHub Issues |

---

## 11. Resource Requirements

### Team Structure

**Minimum Team for 3-Month MVP**:
- **1 Full-Stack Developer** (Go + React)
- **1 Part-Time QA** (20 hours/week)
- **1 Part-Time Technical Writer** (10 hours/week)

### Budget Estimate

| Category | Cost | Duration |
|----------|------|---------|
| Full-Stack Developer | $15,000/month | 3 months |
| QA (part-time) | $2,000/month | 3 months |
| Technical Writer (part-time) | $1,000/month | 3 months |
| Infrastructure (Dev/Test) | $500/month | 3 months |
| Tools & Services | $200/month | 3 months |
| **Total** | **$18,700/month × 3** | **$56,100** |

**Buffer (15%)**: $8,400
**Total with Buffer**: **$65,000**

---

## 12. Comparison: MVP vs Full Version

| Feature | MVP v0.1 | Full v3.0 |
|---------|----------|-----------|
| **Development Time** | 3 months | 10 months |
| **Budget** | $65,000 | $215,500 |
| **MQTT Broker** | Embedded only | Embedded + External |
| **Database** | SQLite only | SQLite + PostgreSQL + MongoDB |
| **Cache** | In-memory only | In-memory + Redis + Memcached |
| **Authentication** | Basic | Advanced (JWT, OAuth) |
| **User Management** | Single user | Multi-user, RBAC |
| **Monitoring** | Basic dashboard | Advanced monitoring & alerting |
| **API** | None | RESTful API |
| **High Availability** | No | Yes (clustering) |
| **SaaS Features** | No | Yes |
| **Target** | Small deployments | Enterprise + Small |

---

## 13. Next Steps

### Immediate (Week 1)

1. **Repository Setup**
   - Create GitHub repository
   - Set up project structure
   - Configure CI/CD

2. **Development Environment**
   - Set up local Go environment
   - Set up React + Vite
   - Configure hot reload

3. **Initial Commits**
   - Basic Go server
   - Embedded UI
   - Hello World page

### Short-term (Month 1)

4. **Core Backend**
   - MQTT broker integration
   - Database setup
   - Cache implementation

5. **Basic UI**
   - Dashboard skeleton
   - Device list
   - Topic list

---

## 14. Conclusion

### MVP Philosophy

**The goal is to build something SMALL that WORKS, not something BIG that doesn't.**

**What Matters**:
- ✅ Single executable
- ✅ Embedded MQTT broker
- ✅ Works out of the box
- ✅ Easy to configure
- ✅ Basic monitoring

**What Doesn't Matter (Yet)**:
- ❌ External databases
- ❌ Cloud integration
- ❌ Enterprise features
- ❌ Advanced monitoring
- ❌ Multi-tenancy

### Success Criteria

**MVP is successful if**:
1. User can download single executable
2. Run it on their machine
3. Complete onboarding wizard
4. Connect MQTT devices
5. See messages in dashboard
6. Everything works without external dependencies

With focused execution over 3 months, we can deliver a **working MVP** that proves the concept and provides value to early adopters.

---

**Document Version**: MVP v0.1
**Last Updated**: January 9, 2026
**Status**: Draft - Pending Review
**Development Time**: 3 months
**Budget**: $65,000

---

## Appendix

### A. MVP Checklist

**Must Have for MVP**:
- [ ] Single executable binary
- [ ] Embedded MQTT broker
- [ ] Local database (SQLite)
- [ ] In-memory cache
- [ ] Web admin interface
- [ ] Setup mode
- [ ] Runner mode
- [ ] Onboarding wizard
- [ ] Basic configuration (config.yaml)
- [ ] Dashboard with metrics
- [ ] Device management
- [ ] Topic browser
- [ ] Message viewer
- [ ] Documentation

### B. Future Features (Post-MVP)

**Phase 2** (6 months after MVP):
- External database support
- REST API
- Advanced monitoring
- User authentication

**Phase 3** (12 months after MVP):
- High availability
- Clustering
- Cloud integration
- Mobile app

### C. References

1. mochi-mqtt documentation
2. SQLite documentation
3. go-cache documentation
4. React documentation
5. Vite documentation

---

**Prepared by**: Development Team
**Approved by**: [Pending]
**Review Date**: [Pending]
**Target Release**: Q2 2026
