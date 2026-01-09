# Project Brief: Universal MQTT Gateway Server (UM-Gateway) v2.0

## Executive Summary

**Project Name**: Universal MQTT Gateway Server (UM-Gateway)
**Version**: 2.0 (Single Executable with Embedded MQTT)
**Status**: Proposal/Planning
**Target Release**: Q4 2026
**Organization**: Rafflesia Agro Technology Division

---

## Version History

| Version | Date | Changes | Author |
|---------|------|---------|--------|
| **v2.0** | 2026-01-09 | **Major Architecture Change**: Embedded MQTT broker implementation<br>• Added embedded MQTT broker using mochi-mqtt library<br>• Introduced dual-mode operation (setup/runner)<br>• Changed configuration from database to YAML config.yaml<br>• Replaced external cache with in-memory cache (go-cache)<br>• SQLite WAL mode as default (optimized for concurrency)<br>• Single executable architecture with process coordination<br>• New code examples for embedded broker, WAL mode, cache adapters<br>• Budget: $238,150 (10 months) | Development Team |
| **v1.0** | 2026-01-09 | **Web Admin Dashboard**: Visual configuration approach<br>• Added SvelteKit admin dashboard for configuration<br>• SQLite database with onboarding wizard<br>• Replaced JSON configuration with web UI forms<br>• External MQTT broker connection model<br>• Budget: $178,650 (reduced from v0) | Development Team |
| **v0.0** | 2026-01-09 | **Original Proposal**: JSON-based configuration<br>• Initial concept with JSON file definitions<br>• External MQTT broker connection required<br>• PostgreSQL/MongoDB configurable via JSON<br>• Focus on farm telemetry use case<br>• Budget: $283,200 | Development Team |

---

## 1. Background

### Evolution from v1 to v2

**v1 Approach**: Web admin dashboard with SQLite, separate MQTT broker connection
**v2 Approach**: Truly self-contained single executable with embedded MQTT broker, dual-mode operation

### Limitations of v1

1. **External MQTT Broker Dependency**: Still requires external MQTT broker (Mosquitto, HiveMQ, etc.)
2. **No True Offline Operation**: Cannot function without external services
3. **Setup vs Runtime Confusion**: Configuration changes require re-running setup
4. **Limited Deployment Scenarios**: Cannot run in air-gapped environments easily
4. **Cache/State Dependencies**: Redis required for production use

### Market Opportunity

The IoT landscape needs a **truly self-contained MQTT gateway** that:
- **Edge Deployment**: Run on edge devices without external dependencies
- **Air-Gapped Environments**: Function in isolated networks
- **Simplified Operations**: Single executable, no separate infrastructure
- **Embedded Use Cases**: Deploy on devices with limited resources

---

## 2. Business Model

### Value Proposition

**For Edge Deployments**:
- Single binary includes everything needed
- No separate MQTT broker to manage
- Reduced infrastructure complexity
- Lower operational costs

**For Air-Gapped Environments**:
- Complete offline operation
- No external dependencies
- Built-in MQTT broker
- Full control over data

**For Development/Testing**:
- Quick local setup without external services
- Easy to reset and reconfigure
- Setup mode for initial configuration
- Runner mode for production

### Revenue Streams

1. **Free Self-Hosted**: Single binary with embedded MQTT (MIT License)
2. **Pro License** ($299 one-time): External database adapters, advanced features
3. **Enterprise License** ($1,499/year): Priority support, custom builds
4. **Edge appliances**: Pre-configured hardware + software bundles

### Pricing Strategy

| Edition | Target Market | Pricing Model | Features |
|---------|--------------|---------------|----------|
| Community | Hobbyists, edge computing | Free (MIT) | Embedded MQTT, SQLite WAL, in-memory cache, setup mode |
| Professional | SMEs, production | $299 one-time | External DB adapters, persistent cache, API access |
| Enterprise | Large organizations | $1,499/year | Multiple gateways, clustering, priority support |

---

## 3. Target Audience

### Primary Users

**1. Edge Computing Engineers**
- Deploy IoT gateways on edge devices
- Need minimal footprint
- Require offline operation
- Pain point: External dependencies increase complexity

**2. Air-Gapped Environment Operators**
- Industrial IoT in isolated networks
- Cannot access external services
- Need complete control
- Pain point: Managing multiple separate services

**3. Development Teams**
- Need local testing environment
- Want quick setup/teardown
- Require configuration flexibility
- Pain point: Setting up multiple services for development

**4. System Integrators**
- Deploy customer premise solutions
- Want simplified operations
- Need easy reconfiguration
- Pain point: Managing external MQTT brokers per customer

---

## 4. Problem Statement

### Core Problems

**Problem 1: External MQTT Broker Dependency**
- Must deploy and manage separate MQTT broker (Mosquitto, HiveMQ, EMQX)
- Additional operational complexity
- More failure points
- Current approach: External broker connection required

**Problem 2: No True Offline Operation**
- v1 still requires external MQTT broker
- Cannot function in air-gapped networks without setup
- External dependencies reduce reliability
- Current approach: Client mode only, no broker capability

**Problem 3: Configuration vs Runtime Confusion**
- No clear separation between setup and runtime
- Configuration changes during runtime cause instability
- No dedicated setup mode
- Current approach: Web admin always available, can cause issues

**Problem 4: Cache/State Management Complexity**
- Redis required for production-grade caching
- External dependency for temporary data
- Difficult to deploy in resource-constrained environments
- Current approach: In-memory cache only or external Redis

**Problem 5: Multi-Process Coordination**
- Web server, API server, MQTT client running in same process
- Potential resource contention
- Graceful shutdown complexity
- Current approach: Single process with goroutines, not well optimized

### Impact

- **Deployment Complexity**: Multiple services to deploy vs single binary
- **Infrastructure Cost**: Additional servers/brokers needed
- **Failure Points**: External dependencies increase failure scenarios
- **Operational Overhead**: Managing multiple components
- **Edge Limitations**: Cannot deploy in truly resource-constrained environments

---

## 5. Proposed Solution

### Vision: Truly Self-Contained MQTT Gateway

A **single executable binary** that includes:
- **Embedded MQTT Broker** (using mochi-mqtt library)
- **Dual-Mode Operation**: Setup mode and Runner mode
- **SQLite in WAL Mode**: Optimized for concurrent reads/writes
- **In-Memory Cache**: Built-in cache with optional external adapters
- **YAML Configuration**: Clear separation of setup and runtime config
- **Process Coordination**: Optimized single-process multi-server architecture

### Key Innovations

**1. Embedded MQTT Broker**
- Uses `mochi-mqtt` Go library
- Full MQTT 3.1.1 and 5.0 broker support
- No external broker dependency
- Can still connect to external brokers if needed

**2. Dual-Mode Operation**
```
[Setup Mode]  ←→  Configure YAML file, database, MQTT settings
     ↓
[Runner Mode] ←→  Production operation with embedded broker
```

**3. SQLite in WAL Mode**
- Write-Ahead Logging for better concurrency
- Better performance than standard SQLite
- Sufficient for moderate workloads
- Optional upgrade to external databases

**4. In-Memory Cache with Adapters**
- Built-in cache using `github.com/patrickmn/go-cache`
- Pluggable cache adapters (Redis, Memcached available)
- TTL-based expiration
- Configurable cache sizes

**5. config.yaml Configuration**
```yaml
mode: runner  # setup | runner
model: embedded  # embedded | external

database:
  type: sqlite  # sqlite | postgres | mongodb | influxdb
  connection:
    path: ./data/gateway.db
    wal_mode: true
    max_open_conns: 25
    cache:
      type: memory  # memory | redis | memcached
      ttl: 3600
      max_size: 10000

mqtt:
  type: embedded  # embedded | external
  embedded:
    listen_address: :1883
    max_connections: 1000
    max_message_size: 256KB
  external:
    url: tcp://localhost:1883
    client_id: um-gateway
    username: ""
    password: ""

paths:
  data: ./data
  logs: ./logs
  pid: ./um-gateway.pid

retention:
  enabled: true
  default_ttl: 2592000  # 30 days
  cleanup_interval: 3600

web:
  address: :8080
  admin_api:
    enabled: true
  public_api:
    enabled: true

logging:
  level: info
  format: json
```

### Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│              Single Executable Binary (um-gateway)          │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌───────────────────────────────────────────────────────┐  │
│  │               Main Process (Go)                        │  │
│  │                                                        │  │
│  │  ┌────────────┐  ┌────────────┐  ┌────────────┐  │  │
│  │  │ Embedded  │  │ HTTP/Web  │  │ In-Memory │  │  │
│  │  │ MQTT      │  │ Server     │  │ Cache      │  │  │
│  │  │ Broker     │  │            │  │            │  │  │
│  │  │ (mochi-mqtt)│  │            │  │            │  │  │
│  │  └────────────┘  └────────────┘  └────────────┘  │  │
│  │                                                        │  │
│  │  ┌────────────┐  ┌────────────┐                    │  │
│  │  │ YAML       │  │ SQLite     │                    │  │
│  │  │ Config     │  │ (WAL Mode) │                    │  │
│  │  │ Loader     │  │            │                    │  │
│  │  └────────────┘  └────────────┘                    │  │
│  │                                                        │  │
│  │  ┌────────────┐  ┌────────────┐                    │  │
│  │  │ Setup      │  │ Mode       │                    │  │
│  │  │ Mode       │  │ Switcher   │                    │  │
│  │  │ Manager    │  │            │                    │  │
│  │  └────────────┘  └────────────┘                    │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                               │
│  ┌───────────────────────────────────────────────────────┐  │
│  │              Storage & Cache Layer                   │  │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────────┐      │  │
│  │  │ SQLite   │  │ Memory   │  │ External     │      │  │
│  │  │ (WAL)    │  │ Cache    │  │ Adapters    │      │  │
│  │  │          │  │ (go-cache)│  │ (Optional)  │      │  │
│  │  └──────────┘  └──────────┘  └──────────────┘      │  │
│  └───────────────────────────────────────────────────────┘  │
│                                                               │
│  ┌───────────────────────────────────────────────────────┐  │
│  │         Web Admin Dashboard (Embedded SPA)            │  │
│  │  - Setup Mode UI (first run or explicit)              │  │
│  │  - Configuration management (YAML editing)             │  │
│  │  - Monitoring dashboard                               │  │
│  │  - Logs viewer                                         │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

---

## 6. Project Objectives

### Primary Objectives (Must Have)

**O1: Implement Embedded MQTT Broker**
- Integrate `mochi-mqtt` library for full broker functionality
- Support MQTT 3.1.1 and 5.0 protocols
- Configurable listener address and limits
- Optional external broker connection (bridging)

**O2: Implement Dual-Mode Operation**
- Setup mode: Initial configuration and migrations
- Runner mode: Production operation
- Mode switching through config.yaml or command-line flag
- Protection against mode switching in production

**O3: SQLite WAL Mode Optimization**
- Enable Write-Ahead Logging for better concurrency
- Configure connection pooling
- Optimize for embedded use case
- Optional external database adapters

**O4: In-Memory Cache with Pluggable Adapters**
- Built-in cache using `go-cache`
- TTL-based expiration
- Pluggable cache adapters (Redis, Memcached)
- Configurable cache sizes and policies

**O5: Process Coordination**
- Optimize goroutine usage for web, API, MQTT servers
- Graceful shutdown handling
- Resource management and limits
- Health checks for all components

### Secondary Objectives (Should Have)

**O6: YAML Configuration Management**
- Clear, human-readable config.yaml
- Hot-reload support in setup mode
- Validation before applying changes
- Configuration migration helpers

**O7: Setup Mode Web UI**
- Configuration editor with YAML validation
- Database migration wizard
- Mode switch confirmation
- Configuration backup/restore

**O8: External Database Adapters**
- PostgreSQL adapter
- MongoDB adapter
- InfluxDB adapter
- One-click migration from SQLite

**O9: External Cache Adapters**
- Redis adapter
- Memcached adapter
- Cache warming strategies
- Cache statistics and monitoring

### Tertiary Objectives (Nice to Have)

**O10: High Availability Mode**
- Multi-gateway clustering
- Data replication
- Failover mechanisms

**O11: Performance Optimization**
- Connection pooling
- Batch optimization
- Memory profiling and tuning

---

## 7. Project Scope & Boundaries

### In Scope (What We Will Build)

#### Core Functionality

1. **Embedded MQTT Broker**
   - Full broker implementation using mochi-mqtt
   - MQTT 3.1.1 and 5.0 support
   - Configurable listeners
   - Client authentication
   - Message persistence (optional)

2. **Dual-Mode System**
   - Setup mode: First-run or explicit activation
   - Runner mode: Production operation
   - Mode detection and validation
   - Safe mode switching

3. **SQLite WAL Mode Database**
   - Write-Ahead Logging enabled
   - Connection pooling
   - Automatic migrations
   - Backup/restore utilities

4. **In-Memory Cache**
   - go-cache implementation
   - TTL expiration
   - Size-based eviction
   - Statistics API

5. **YAML Configuration**
   - config.yaml structure
   - Validation on load
   - Hot-reload in setup mode
   - Configuration versioning

6. **Process Management**
   - Single process with optimized goroutines
   - Graceful shutdown
   - PID file management
   - Signal handling (SIGTERM, SIGINT)

#### Web UI Features

**Setup Mode UI**:
- YAML configuration editor
- Database migration wizard
- Mode switch confirmation
- System status checks

**Runner Mode UI**:
- Monitoring dashboard
- Logs viewer
- Metrics visualization
- Configuration read-only view

### Out of Scope (What We Won't Build - Initially)

1. **Hardware/Firmware**
   - Device firmware
   - OTA updates
   - Hardware provisioning

2. **Advanced MQTT Features** (Phase 2)
   - Clustered MQTT brokers
   - Bridge configurations
   - Advanced security (SSL/TLS setup)

3. **High Availability** (Phase 3)
   - Multi-gateway clustering
   - Data replication
   - Automatic failover

### Boundaries & Constraints

#### Technical Constraints
- **Language**: Go 1.24+
- **MQTT Library**: mochi-mqtt (embedded) or Eclipse Paho (external)
- **Minimum Hardware**: 1 CPU core, 1GB RAM (WAL mode requires more)
- **Network**: Optional connectivity for external brokers/databases

#### Time Constraints
- **Phase 1 (MVP)**: 10 months
- **Phase 2 (External Adapters)**: +3 months
- **Phase 3 (HA Features)**: +6 months

#### Budget Constraints
- **Development Team**: 2-3 full-time developers
- **Infrastructure**: $1,500/month for development/testing
- **Third-party Services**: Free tier initially

#### Resource Constraints
- **Development Team**: 1-2 Full-stack developers (Go + SvelteKit), 1 DevOps engineer
- **Domain Experts**: MQTT specialists, cache experts
- **Support**: Part-time technical writer

---

## 8. Technology Stack

### Core Technologies

#### Backend (Go)
| Component | Technology | Rationale |
|-----------|-----------|-----------|
| Language | Go 1.24+ | Performance, concurrency, single binary |
| MQTT Broker | mochi-mqtt | Embedded Go MQTT broker |
| MQTT Client (external) | Eclipse Paho | For external broker connections |
| HTTP Router | Chi v5 | Lightweight, idiomatic Go |
| YAML | gopkg.in/yaml.v3 | YAML parsing and serialization |
| SQLite | mattn/go-sqlite3 | WAL mode support |
| Cache | patrickmn/go-cache | In-memory caching with TTL |
| WebSocket | gorilla/websocket | Real-time UI updates |

#### Embedded MQTT Broker Configuration

```go
import mqtt "github.com/mochi-mqtt/mqtt/v2"

// Embedded broker setup
func NewEmbeddedBroker(config MQTTConfig) (*mqtt.Server, error) {
    server := mqtt.NewServer(nil)

    // Configure TCP listener
    tcp := mqtt.NewTCPListener(config.ListenAddress, nil)
    server.AddListener(tcp)

    return server, nil
}
```

#### Cache Implementation

```go
import "github.com/patrickmn/go-cache"

// In-memory cache
type Cache struct {
    store *cache.Cache
    ttl   time.Duration
}

func NewCache(defaultTTL time.Duration, maxSize int) *Cache {
    return &Cache{
        store: cache.New(defaultTTL, 10*time.Minute),
        ttl:   defaultTTL,
    }
}

// Pluggable cache adapter interface
type CacheAdapter interface {
    Get(key string) (interface{}, error)
    Set(key string, value interface{}, ttl time.Duration) error
    Delete(key string) error
    Clear() error
}
```

#### SQLite WAL Mode

```go
import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

// Enable WAL mode
func enableWALMode(db *sql.DB) error {
    _, err := db.Exec("PRAGMA journal_mode=WAL")
    if err != nil {
        return err
    }

    _, err = db.Exec("PRAGMA synchronous=NORMAL")
    if err != nil {
        return err
    }

    _, err = db.Exec("PRAGMA cache_size=-10000") // 10MB cache
    return err
}
```

### Configuration Structure

#### config.yaml Schema

```yaml
# Operation Mode
mode: runner  # setup | runner

# Model Type
model: embedded  # embedded | external

# Database Configuration
database:
  # Type: sqlite, postgres, mongodb, influxdb
  type: sqlite
  connection:
    # SQLite specific
    path: ./data/gateway.db
    wal_mode: true
    max_open_conns: 25
    max_idle_conns: 5
    conn_max_lifetime: 300s

    # Cache configuration
    cache:
      # Type: memory, redis, memcached
      type: memory
      ttl: 3600
      max_size: 10000

      # Redis (if type is redis)
      redis:
        network: tcp
        address: localhost:6379
        password: ""
        db: 0

      # Memcached (if type is memcached)
      memcached:
        address: localhost:11211
        timeout: 1000

# MQTT Configuration
mqtt:
  # Type: embedded, external
  type: embedded

  # Embedded broker settings
  embedded:
    listen_address: :1883
    max_connections: 1000
    max_message_size: 256KB
    keepalive: 60s

  # External broker (if type is external)
  external:
    url: tcp://localhost:1883
    client_id: um-gateway
    username: ""
    password: ""
    clean_session: true
    auto_reconnect: true
    keepalive: 60s

# Paths Configuration
paths:
  data: ./data       # Data directory
  logs: ./logs       # Logs directory
  pid: ./um-gateway.pid  # PID file location

# Data Retention
retention:
  enabled: true
  default_ttl: 2592000  # 30 days in seconds
  cleanup_interval: 3600   # Cleanup interval in seconds

# Web Server Configuration
web:
  address: :8080
  admin_api:
    enabled: true
    authentication:
      type: basic  # basic, jwt, none
  public_api:
    enabled: true
    rate_limit:
      enabled: true
      requests_per_minute: 1000

# Logging Configuration
logging:
  level: info  # debug, info, warn, error
  format: json  # json, text
  output: stdout  # stdout, file
  file:
    path: ./logs/gateway.log
    max_size: 100MB
    max_backups: 10
    max_age: 30
```

---

## 9. Features

### 9.1 Core Features

#### F1: Dual-Mode Operation

**Setup Mode**:
- Activated on first run or explicit flag (`--mode=setup`)
- Web UI allows configuration changes
- YAML configuration editing
- Database migrations
- Mode switching to runner mode

**Runner Mode**:
- Normal production operation
- Configuration is read-only
- Embedded MQTT broker active
- Cannot switch to setup mode without stopping process

**Mode Switching**:

```go
// Mode detection
func DetermineMode(config *Config) string {
    // Check command line flag
    if *setupFlag {
        return "setup"
    }

    // Check config file
    if config.Mode == "" {
        // First run - default to setup
        if isFirstRun() {
            return "setup"
        }
        return "runner"
    }

    return config.Mode
}

// Mode enforcement
func RunInMode(config *Config) error {
    mode := DetermineMode(config)

    switch mode {
    case "setup":
        return runSetupMode(config)
    case "runner":
        return runRunnerMode(config)
    default:
        return fmt.Errorf("unknown mode: %s", mode)
    }
}
```

---

#### F2: Embedded MQTT Broker

**Implementation using mochi-mqtt**:

```go
package broker

import (
    mqtt "github.com/mochi-mqtt/mqtt/v2"
    "github.com/mochi-mqtt/mqtt/hooks"
)

type EmbeddedBroker struct {
    server   *mqtt.Server
    config   *MQTTConfig
    hooks    *hooks.Hooks
}

func NewEmbeddedBroker(config *MQTTConfig) (*EmbeddedBroker, error) {
    // Create server
    server := mqtt.NewServer(nil)

    // Setup hooks
    h := hooks.NewHooks()
    h.OnConnect = onClientConnect
    h.OnMessage = onMessageReceived
    h.OnDisconnect = onClientDisconnect

    server.AddHook(h)

    // Configure TCP listener
    tcp := mqtt.NewTCPListener(config.ListenAddress, nil)
    server.AddListener(tcp)

    // Configure WebSocket listener (optional)
    ws := mqtt.NewWebSocketListener(config.WSAddress, nil)
    server.AddListener(ws)

    return &EmbeddedBroker{
        server: server,
        config: config,
        hooks:  h,
    }, nil
}

func (b *EmbeddedBroker) Start() error {
    go b.server.Serve()
    return nil
}

func (b *EmbeddedBroker) Stop() {
    b.server.Close()
}
```

**Client Connection Handling**:

```go
func onClientConnect(cl mqtt.Client, pk packets.Packet) {
    // Authentication logic
    clientID := cl.ClientInfo().ClientID
    username := cl.ClientInfo().Username
    password := cl.ClientInfo().Password

    // Validate credentials
    if !authenticateClient(clientID, username, password) {
        cl.Disconnect(0x80) // Not authorized
        return
    }

    // Log connection
    log.Info("Client connected", "client_id", clientID)
}

func onMessageReceived(cl mqtt.Client, pk packets.Packet) {
    topic := pk.TopicName
    payload := pk.Payload

    // Process message
    handleMessage(topic, payload)
}

func onClientDisconnect(cl mqtt.Client, err error) {
    clientID := cl.ClientInfo().ClientID
    log.Info("Client disconnected", "client_id", clientID, "error", err)
}
```

---

#### F3: SQLite WAL Mode

**Optimization for Concurrent Access**:

```go
package database

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

type SQLiteDatabase struct {
    db *sql.DB
}

func NewSQLiteDatabase(config SQLiteConfig) (*SQLiteDatabase, error) {
    // Open database
    db, err := sql.Open("sqlite3", config.Path)
    if err != nil {
        return nil, err
    }

    // Enable WAL mode for better concurrency
    if err := enableWALMode(db); err != nil {
        return nil, err
    }

    // Configure connection pool
    db.SetMaxOpenConns(config.MaxOpenConns)
    db.SetMaxIdleConns(config.MaxIdleConns)
    db.SetConnMaxLifetime(config.ConnMaxLifetime)

    return &SQLiteDatabase{db: db}, nil
}

func enableWALMode(db *sql.DB) error {
    settings := []string{
        "PRAGMA journal_mode=WAL",
        "PRAGMA synchronous=NORMAL",
        "PRAGMA cache_size=-10000",        // -10000 means use max available
        "PRAGMA temp_store=memory",
        "PRAGMA mmap_size=268435456",        // 256MB
    }

    for _, setting := range settings {
        if _, err := db.Exec(setting); err != nil {
            return fmt.Errorf("failed to execute %s: %w", setting, err)
        }
    }

    return nil
}
```

**Performance Benefits**:
- Better concurrency (multiple readers + single writer)
- Reduced disk I/O
- Faster commit times
- Better for embedded/edge use cases

---

#### F4: In-Memory Cache with Adapters

**Base Cache Implementation**:

```go
package cache

import (
    "github.com/patrickmn/go-cache"
    "time"
)

type MemoryCache struct {
    store *cache.Cache
    ttl   time.Duration
}

func NewMemoryCache(defaultTTL time.Duration, maxSize int) *MemoryCache {
    return &MemoryCache{
        store: cache.New(defaultTTL, 10*time.Minute),
        ttl:   defaultTTL,
    }
}

func (c *MemoryCache) Get(key string) (interface{}, error) {
    val, found := c.store.Get(key)
    if !found {
        return nil, ErrKeyNotFound
    }
    return val, nil
}

func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) error {
    if ttl == 0 {
        ttl = c.ttl
    }
    c.store.Set(key, value, ttl)
    return nil
}

func (c *MemoryCache) Delete(key string) error {
    c.store.Delete(key)
    return nil
}

func (c *MemoryCache) Clear() error {
    c.store.Flush()
    return nil
}

// Statistics
func (c *MemoryCache) Stats() CacheStats {
    stats := c.store.ItemCount()
    return CacheStats{
        ItemCount: stats,
        Size:       c.estimateSize(),
    }
}
```

**Redis Adapter**:

```go
package cache

import (
    "github.com/redis/go-redis/v9"
    "context"
    "time"
)

type RedisCache struct {
    client *redis.Client
    ttl    time.Duration
}

func NewRedisCache(addr, password string, db int, ttl time.Duration) (*RedisCache, error) {
    client := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })

    if err := client.Ping(context.Background()).Err(); err != nil {
        return nil, err
    }

    return &RedisCache{
        client: client,
        ttl:    ttl,
    }, nil
}

func (c *RedisCache) Get(key string) (interface{}, error) {
    val, err := c.client.Get(context.Background(), key).Result()
    if err == redis.Nil {
        return nil, ErrKeyNotFound
    }
    return val, err
}

func (c *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
    if ttl == 0 {
        ttl = c.ttl
    }
    return c.client.Set(context.Background(), key, value, ttl).Err()
}

func (c *RedisCache) Delete(key string) error {
    return c.client.Del(context.Background(), key).Err()
}

func (c *RedisCache) Clear() error {
    // Warning: This will delete all keys in the database!
    return c.client.FlushDB(context.Background()).Err()
}
```

---

#### F5: Process Coordination

**Single Process with Optimized Goroutines**:

```go
package main

import (
    "sync"
    "context"
    "os"
    "os/signal"
    "syscall"
    "time"
)

type GatewayServer struct {
    config      *Config
    mqttBroker  *broker.EmbeddedBroker
    webServer  *WebServer
    apiServer  *APIServer
    cache       cache.Cache
    database   *database.Database
    wg          sync.WaitGroup
    shutdownCtx context.Context
    shutdown    func()
}

func NewGatewayServer(configPath string) (*GatewayServer, error) {
    // Load YAML configuration
    config, err := LoadConfig(configPath)
    if err != nil {
        return nil, err
    }

    // Initialize cache
    cache, err := initializeCache(config.Database.Cache)
    if err != nil {
        return nil, err
    }

    // Initialize database
    db, err := initializeDatabase(config.Database)
    if err != nil {
        return nil, err
    }

    // Initialize embedded MQTT broker
    broker, err := initializeMQTTBroker(config.MQTT, db, cache)
    if err != nil {
        return nil, err
    }

    // Create shutdown context
    shutdownCtx, shutdown := context.WithCancel(context.Background())

    return &GatewayServer{
        config:     config,
        mqttBroker: broker,
        cache:      cache,
        database:   db,
        shutdownCtx: shutdownCtx,
        shutdown:   shutdown,
    }, nil
}

func (s *GatewayServer) Start() error {
    // Start embedded MQTT broker
    s.wg.Add(1)
    go func() {
        defer s.wg.Done()
        if err := s.mqttBroker.Start(); err != nil {
            log.Error("MQTT broker error", "error", err)
        }
    }()

    // Start web server (admin dashboard)
    s.wg.Add(1)
    go func() {
        defer s.wg.Done()
        s.webServer.Start(s.shutdownCtx)
    }()

    // Start API server
    s.wg.Add(1)
    go func() {
        defer s.wg.Done()
        s.apiServer.Start(s.shutdownCtx)
    }()

    // Start background jobs
    s.wg.Add(1)
    go func() {
        defer s.wg.Done()
        s.runBackgroundJobs(s.shutdownCtx)
    }()

    return nil
}

func (s *GatewayServer) runBackgroundJobs(ctx context.Context) {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            log.Info("Background jobs shutting down")
            return
        case <-ticker.C:
            // Perform cleanup tasks
            s.performCleanup()

            // Perform cache stats
            stats := s.cache.Stats()
            log.Debug("Cache stats", "stats", stats)
        }
    }
}

func (s *GatewayServer) Stop() {
    log.Info("Shutting down gateway...")

    // Signal shutdown
    s.shutdown()

    // Stop embedded MQTT broker
    s.mqttBroker.Stop()

    // Wait for all goroutines (with timeout)
    done := make(chan struct{})
    go func() {
        s.wg.Wait()
        close(done)
    }()

    select {
    case <-done:
        log.Info("All services stopped gracefully")
    case <-time.After(30 * time.Second):
        log.Warn("Shutdown timed out, forcing exit")
    }

    log.Info("Gateway stopped")
}

func (s *GatewayServer) performCleanup() {
    // Delete expired data
    s.database.CleanupExpiredData(s.config.Retention.DefaultTTL)

    // Flush cache if needed
    // s.cache.Flush()
}
```

---

### 9.2 Advanced Features

#### F6: Setup Mode Web UI

**Configuration Editor**:

```
┌─────────────────────────────────────────────────────────────┐
│  Setup Mode - Configuration Editor                           │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  Mode: ● Setup Mode  ○ Runner Mode                           │
│                                                               │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  config.yaml                                        │   │
│  │  ┌───────────────────────────────────────────────┐  │   │
│  │  │ mode: runner                                     │  │   │
│  │  │                                                │  │   │
│  │  │ model: embedded                                 │  │   │
│  │  │                                                │  │   │
│  │  │ database:                                      │  │   │
│  │  │   type: sqlite                                 │  │   │
│  │  │   connection:                                  │  │   │
│  │  │     path: ./data/gateway.db                  │  │   │
│  │  │     wal_mode: true                             │  │   │
│  │  │     cache:                                     │  │   │
│  │  │       type: memory                             │  │   │
│  │  │       ttl: 3600                               │  │   │
│  │  │                                                │  │   │
│  │  │ mqtt:                                          │  │   │
│  │  │   type: embedded                               │  │   │
│  │  │   embedded:                                    │  │   │
│  │  │     listen_address: :1883                     │  │   │
│  │  │     max_connections: 1000                      │  │   │
│  │  │                                                │  │   │
│  │  │ retention:                                     │  │   │
│  │  │   enabled: true                                │  │   │
│  │  │   default_ttl: 2592000                         │  │   │
│  │  │                                                │  │   │
│  │  └───────────────────────────────────────────────┘  │   │
│  │                                                      │   │
│  │  [Validate]  [Save]  [Cancel]                        │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  [Switch to Runner Mode] [Export Config] [Backup Data]        │
└─────────────────────────────────────────────────────────────┘
```

---

#### F7: Database Migration Wizard

**SQLite → PostgreSQL**:

```
┌─────────────────────────────────────────────────────────────┐
│  Database Migration Wizard                                     │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  Current Database: SQLite (WAL mode)                         │
│  Target Database: ○ PostgreSQL  ● MongoDB  ● InfluxDB     │
│                                                               │
│  Migration Data:                                             │
│  ☑ Configuration data                                       │
│  ☑ Time-series data (may take hours)                          │
│  ☑ Cache data                                               │
│  ☐ User accounts (Pro feature)                               │
│                                                               │
│  PostgreSQL Connection:                                       │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Host: [localhost]                                   │   │
│  │ Port: [5432]                                        │   │
│  │ Database: [um_gateway_prod]                         │   │
│  │ Username: [um_gateway]                              │   │
│  │ Password: [••••••••••]                              │   │
│  │                                                      │   │
│  │ [Test Connection] ✓ Connected successfully      │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  Migration Options:                                          │
│  ☑ Stop MQTT broker during migration                       │
│  ☑ Create backup before migration                          │
│  ☑ Verify data after migration                             │
│                                                               │
│  Estimated time: ~2 hours for 1.2M records                    │
│                                                               │
│  [Start Migration]  [Cancel]                                 │
└─────────────────────────────────────────────────────────────┘
```

---

#### F8: Runner Mode Dashboard

**Production Monitoring**:

```
┌─────────────────────────────────────────────────────────────┐
│  Runner Mode - Monitoring Dashboard                            │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  System Status                                                │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐ │
│  │  ✓ MQTT  │  │  ✓ Web    │  │  ✓ API    │  │  ✓ DB     │ │
│  │  Broker  │  │  Server   │  │  Server   │  │  (WAL)   │ │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘ │
│                                                               │
│  Broker Statistics                                            │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Connected Clients: 47                                  │   │
│  │ Messages/min: 1,234                                   │   │
│  │ Uptime: 45 days, 6 hours                                │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  Cache Statistics                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Type: Memory (go-cache)                               │   │
│  │ Items: 5,234                                            │   │
│  │ Hit Rate: 87.3%                                        │   │
│  │ Memory: ~45 MB                                          │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  Database Statistics                                          │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Type: SQLite (WAL mode)                               │   │
│  │ Size: 127 MB                                           │   │
│  │ Records: 1,234,567                                       │   │
│  │ Connections: 5/25                                       │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  [View Logs]  [View Metrics]  [Enter Setup Mode]              │
└─────────────────────────────────────────────────────────────┘
```

---

## 10. Development Roadmap

### Phase 1: Foundation (Months 1-5)

**Sprint 1-2: Project Setup & Architecture**
- [ ] Monorepo structure (Go + SvelteKit)
- [ ] YAML configuration loading
- [ ] Mode detection and switching logic
- [ ] CI/CD pipeline

**Sprint 3-4: Embedded MQTT Broker**
- [ ] mochi-mqtt integration
- [ ] Client authentication
- [ ] Message persistence (optional)
- [ ] WebSocket support

**Sprint 5-6: SQLite WAL Mode**
- [ ] SQLite with WAL mode
- [ ] Connection pooling
- [ ] Schema migrations
- [ ] Backup/restore utilities

**Sprint 7-8: In-Memory Cache**
- [ ] go-cache integration
- [ ] TTL expiration
- [ ] Cache statistics API
- [ ] Memory management

**Sprint 9-10: Dual-Mode Implementation**
- [ ] Setup mode logic
- [ ] Runner mode logic
- [ ] Mode switching
- [ ] Protection mechanisms

**Milestone**: Alpha - Single executable with embedded broker

---

### Phase 2: Web UI & Features (Months 6-8)

**Sprint 11-12: Setup Mode UI**
- [ ] YAML configuration editor
- [ ] Validation feedback
- [ ] Mode switch confirmation
- [ ] Configuration backup/restore

**Sprint 13-14: Runner Mode Dashboard**
- [ ] Monitoring dashboard
- [ ] Logs viewer
- [ ] Metrics visualization
- [ ] Status indicators

**Sprint 15-16: Process Coordination**
- [ ] Optimized goroutine management
- [ ] Graceful shutdown
- [ ] Resource monitoring
- [ ] Health checks

**Milestone**: Beta - Feature-complete v2

---

### Phase 3: Polish & Production (Months 9-10)

**Sprint 17-18: Testing & Quality**
- [ ] Load testing (embedded broker)
- [ ] Concurrency testing (WAL mode)
- [ ] Memory leak testing
- [ ] Security audit

**Sprint 19-20: Documentation & Deployment**
- [ ] User guide
- [ ] Installation guide
- [ ] Architecture documentation
- [ ] Cross-platform builds

**Milestone**: v2.0 General Availability

---

## 11. Success Metrics

### Technical Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| Binary Size | < 60 MB | Build artifacts |
| Startup Time | < 3 seconds | Automated tests |
| Memory Usage (Idle) | < 100 MB | Resource monitoring |
| Messages/Second (Embedded) | 5K msg/sec | Benchmark tests |
| MQTT Clients | 1,000 concurrent | Load testing |
| Cache Hit Rate | > 85% | Cache statistics |
| Database Writes (SQLite WAL) | 1K writes/sec | Benchmark tests |

### Business Metrics

| Metric | Target | Timeline |
|--------|--------|----------|
| Active Installations | 1,000 | 6 months post-launch |
| GitHub Stars | 2,000 | 6 months post-launch |
| Edge Deployments | 200 | 12 months post-launch |
| Pro License Sales | 100 | 12 months post-launch |

### Quality Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| Mode Switch Success Rate | 100% | Analytics |
| Setup Completion Rate | > 95% | Analytics |
| Mean Time Between Failures | 720 hours (30 days) | Uptime monitoring |
| Documentation Coverage | 100% features | Manual review |

---

## 12. Risk Analysis

### Technical Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| mochi-mqtt library limitations | Medium | Medium | Fallback to Eclipse Paho, extensive testing |
| SQLite WAL mode not sufficient | Medium | Low | External database adapters from start |
| Single-process resource contention | High | Medium | Careful goroutine management, resource limits |
| Cache memory exhaustion | Medium | Medium | Size limits, LRU eviction |

### Business Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Users prefer external broker | Low | Medium | Keep external broker option |
| Embedded broker limitations | Medium | High | Clear communication of capabilities |
| Single binary deployment concerns | Low | Low | Security scans, code signing |

### Operational Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Process hangs on shutdown | Medium | Low | Timeout mechanisms, force kill option |
| Mode switching in production | High | Low | Explicit confirmation, safety checks |
| Data loss during mode switch | High | Low | Automatic backups, validation |

---

## 13. Resource Requirements

### Team Structure

**Core Team (Phase 1-3)**
- 2x Full-Stack Developers (Go + SvelteKit)
- 1x DevOps Engineer
- 1x QA Engineer (part-time)
- 1x Technical Writer (part-time)

**Budget Estimate**

| Category | Cost (Monthly) | Duration |
|----------|----------------|---------|
| Development Team Salaries | $18,000 | 10 months |
| Infrastructure (Dev/Test) | $1,500 | 10 months |
| Tools & Services | $400 | Ongoing |
| Documentation & Design | $1,000 | 8 months |
| Contingency (15%) | $2,915 | - |
| **Total Phase 1-3** | **$23,815/mo × 10** | **$238,150** |

---

## 14. Comparison: v1 vs v2

| Aspect | v1 (Web Admin + SQLite) | v2 (Single Executable + Embedded MQTT) |
|--------|---------------------------|-----------------------------------------|
| **MQTT Broker** | External connection required | Embedded (mochi-mqtt) |
| **Deployment** | Requires separate MQTT broker | Single binary |
| **Offline Operation** | Requires external services | Truly offline-capable |
| **Configuration** | SQLite database + Web UI | YAML config file + Web UI |
| **Cache** | External (Redis) recommended | In-memory (go-cache) |
| **Setup Process** | Onboarding wizard | Setup mode + Runner mode |
| **Use Case** | General IoT deployments | Edge, air-gapped, embedded |
| **Complexity** | Lower (separate services) | Higher (single process coordination) |
| **Resource Usage** | Higher (multiple processes) | Lower (optimized single process) |

---

## 15. Next Steps

### Immediate Actions (Week 1-2)

1. **Technical Research**
   - Evaluate mochi-mqtt library capabilities
   - Test SQLite WAL mode performance
   - Prototype go-cache implementation

2. **Architecture Design**
   - Design dual-mode system
   - Plan process coordination
   - Define goroutine management

3. **Team Formation**
   - Recruit developers with Go expertise
   - Define collaboration processes

### Short-term Actions (Month 1)

4. **Prototype Development**
   - Embedded MQTT broker prototype
   - Dual-mode switching
   - YAML configuration loading

5. **Development Environment**
   - Monorepo setup
   - CI/CD for cross-platform builds
   - Testing framework

---

## 16. Conclusion

The v2 approach represents the **ultimate self-contained MQTT gateway**:

### Key Advantages over v1

**True Self-Containment**
- ✅ Embedded MQTT broker (no external dependency)
- ✅ In-memory caching (no Redis needed)
- ✅ Single executable deployment
- ✅ Works in air-gapped environments

**Simplified Operations**
- ✅ No separate MQTT broker to manage
- ✅ Clear setup/runner mode separation
- ✅ YAML configuration (human-readable)
- ✅ Mode switching for maintenance

**Better Performance**
- ✅ SQLite WAL mode for better concurrency
- ✅ In-memory cache for faster lookups
- ✅ Optimized single-process architecture
- ✅ Reduced inter-process communication overhead

**Targeted Use Cases**
- ✅ Edge computing deployments
- ✅ Air-gapped environments
- ✅ Resource-constrained devices
- ✅ Development and testing

### Evolution Path

**v0** → **v1** → **v2**:
- v0: JSON configuration, technical users only
- v1: Web admin, SQLite, external MQTT broker
- **v2: Single executable, embedded MQTT, YAML config, dual-mode**

The v2 version addresses the ultimate goal: **a single executable that provides complete MQTT gateway functionality** without any external dependencies, perfect for edge, air-gapped, and embedded deployments.

With focused execution over 10 months, we can deliver a **production-ready v2.0** that serves the untapped market of truly self-contained IoT gateways.

---

**Document Version**: 2.0
**Last Updated**: January 9, 2026
**Status**: Draft - Pending Review
**Supersedes**: PROJECT_BRIEF_EN_v1.md

---

## Appendix

### A. Glossary

- **WAL Mode**: Write-Ahead Logging - SQLite optimization for concurrency
- **Dual-Mode**: Setup mode (configuration) and Runner mode (production)
- **Embedded Broker**: MQTT broker running within the same process
- **In-Memory Cache**: Cache stored in RAM using go-cache
- **Single Executable**: All components bundled in one binary file

### B. References

1. mochi-mqtt library documentation
2. go-cache documentation
3. SQLite WAL mode documentation
4. PocketBase (inspiration for single executable)
5. Edge computing best practices

### C. Related Documents

- `PROJECT_BRIEF_EN_v0.md` - Original JSON-based approach
- `PROJECT_BRIEF_EN_v1.md` - Web admin approach
- `PROJECT_BRIEF_ID_v0.md` - Original approach (Bahasa Indonesia)
- `PROJECT_BRIEF_ID_v1.md` - Web admin approach (Bahasa Indonesia)

---

**Prepared by**: Development Team
**Approved by**: [Pending]
**Review Date**: [Pending]
