# Project Brief: Universal MQTT Gateway Server (UM-Gateway) v3.0

## Executive Summary

**Project Name**: Universal MQTT Gateway Server (UM-Gateway)
**Version**: 3.0 (React-based UI with Enhanced Configuration Clarity)
**Status**: Proposal/Planning
**Target Release**: Q1 2027
**Organization**: Rafflesia Agro Technology Division

---

## Version History

| Version | Date | Changes | Author |
|---------|------|---------|--------|
| **v3.0** | 2026-01-09 | **Technology Stack Change & Clarity Improvements**: React-based UI<br>• Replaced SvelteKit with React + Vite<br>• Added TanStack Query for server state management<br>• Added TanStack Router for routing<br>• Clarified Setup Mode = config.yaml editor<br>• Clarified Runner Mode = production operation (config read-only)<br>• Clear separation: config.yaml (infrastructure) vs database (runtime data)<br>• Enhanced documentation on what goes where<br>• Budget: $215,500 (10 months) | Development Team |
| **v2.0** | 2026-01-09 | **Major Architecture Change**: Embedded MQTT broker implementation<br>• Added embedded MQTT broker using mochi-mqtt library<br>• Introduced dual-mode operation (setup/runner)<br>• Changed configuration from database to YAML config.yaml<br>• Replaced external cache with in-memory cache (go-cache)<br>• SQLite WAL mode as default (optimized for concurrency)<br>• Single executable architecture with process coordination<br>• New code examples for embedded broker, WAL mode, cache adapters<br>• Budget: $238,150 (10 months) | Development Team |
| **v1.0** | 2026-01-09 | **Web Admin Dashboard**: Visual configuration approach<br>• Added SvelteKit admin dashboard for configuration<br>• SQLite database with onboarding wizard<br>• Replaced JSON configuration with web UI forms<br>• External MQTT broker connection model<br>• Budget: $178,650 (reduced from v0) | Development Team |
| **v0.0** | 2026-01-09 | **Original Proposal**: JSON-based configuration<br>• Initial concept with JSON file definitions<br>• External MQTT broker connection required<br>• PostgreSQL/MongoDB configurable via JSON<br>• Focus on farm telemetry use case<br>• Budget: $283,200 | Development Team |

---

## 1. Background

### Evolution from v2 to v3

**v2 Approach**: SvelteKit-based UI with unclear separation between setup and runner modes
**v3 Approach**: React-based UI with TanStack ecosystem, crystal-clear mode definitions and data separation

### Limitations of v2

1. **Unclear Mode Purpose**: Setup mode purpose was vague - what exactly can be configured?
2. **Confusion About config.yaml**: Not clear what data belongs in config.yaml vs database
3. **Framework Choice**: SvelteKit less familiar to most developers compared to React
4. **State Management**: No clear server state management strategy
5. **Infrastructure vs Runtime Data**: No clear guidance on what goes where

### Key Improvements in v3

**1. Technology Stack Change**
- **React + Vite**: Faster build, better ecosystem, larger talent pool
- **TanStack Query**: Powerful server state management, caching, synchronization
- **TanStack Router**: Type-safe routing with excellent developer experience

**2. Crystal-Clear Mode Definitions**

**Setup Mode** (Configuration Mode):
- **Purpose**: Edit `config.yaml` infrastructure settings
- **When**: First run, or when changing infrastructure configuration
- **What You Can Change**:
  - Database connection settings (type, path, credentials)
  - MQTT broker settings (embedded/external, ports, limits)
  - Cache configuration (type, TTL, size)
  - Web server settings (ports, authentication)
  - Logging configuration (level, format, output)
  - Data retention policies
  - File paths (data, logs, PID)
- **What You Cannot Change**: Runtime data (topics, subscriptions, messages, device states)
- **How**: Web UI with YAML editor or form-based interface
- **Safety**: Changes validated before saving, backup created automatically

**Runner Mode** (Production Mode):
- **Purpose**: Normal production operation with embedded MQTT broker active
- **When**: Day-to-day operation after initial configuration
- **What Happens**:
  - `config.yaml` is **read-only** (cannot be edited)
  - Embedded MQTT broker runs and accepts connections
  - Web UI shows monitoring dashboard (read-only config view)
  - All runtime operations use database, not config.yaml
  - Cannot switch to setup mode without stopping process
- **What You Can Do**: Monitor system, view logs, check metrics, restart services
- **What You Cannot Do**: Edit infrastructure configuration

**3. Clear Data Separation**

**config.yaml** (Infrastructure Configuration - Static):
- **Purpose**: Defines HOW the gateway runs (infrastructure settings)
- **Location**: File system (same directory as executable)
- **Format**: YAML (human-readable text file)
- **Editing**: Setup mode only
- **Changes Require**: Process restart
- **Contains**:
  ```yaml
  # Operation mode
  mode: runner  # setup | runner

  # Database configuration
  database:
    type: sqlite
    connection:
      path: ./data/gateway.db
      wal_mode: true
      max_open_conns: 25
    cache:
      type: memory
      ttl: 3600

  # MQTT broker configuration
  mqtt:
    type: embedded
    embedded:
      listen_address: :1883
      max_connections: 1000

  # Web server configuration
  web:
    address: :8080
    admin_api:
      enabled: true

  # Logging configuration
  logging:
    level: info
    format: json
  ```

**Database** (Runtime Data - Dynamic):
- **Purpose**: Stores WHAT the gateway processes (actual data)
- **Location**: File system (SQLite: ./data/gateway.db)
- **Format**: Binary database file
- **Access**: Both setup and runner modes (via API)
- **Changes**: Real-time, no restart required
- **Contains**:
  - Device registry (connected devices, their metadata)
  - Topic definitions (what topics exist)
  - Subscription rules (who subscribes to what)
  - Message history (actual MQTT messages)
  - Device states (current status of each device)
  - Metrics and statistics (performance data)
  - User accounts (if authentication enabled)
  - API keys and tokens (if applicable)

**Key Differences**:

| Aspect | config.yaml | Database |
|--------|-------------|----------|
| **Purpose** | Infrastructure settings | Runtime data |
| **Example** | "Listen on port 1883" | "Device X published message Y" |
| **Editing** | Setup mode only | Both modes (via API) |
| **Changes Require** | Process restart | No restart |
| **Format** | YAML text | Binary database |
| **Location** | ./config.yaml | ./data/gateway.db |
| **Backup** | Manual or setup mode | Automatic/periodic |
| **Version Control** | Yes (Git-friendly) | No (binary data) |

### Market Opportunity

The IoT landscape needs a **truly self-contained MQTT gateway** with:
- **Clear Configuration**: Crystal-clear separation of infrastructure vs runtime data
- **Familiar Technology**: React-based UI, easier to find developers
- **Better State Management**: TanStack Query for reliable server state sync
- **Developer Experience**: Type-safe routing, excellent tooling

---

## 2. Business Model

### Value Proposition

**For Edge Deployments**:
- Single binary includes everything needed
- Clear separation: config.yaml (infrastructure) vs database (data)
- No confusion about what to edit and when
- React-based UI = larger talent pool

**For Air-Gapped Environments**:
- Complete offline operation
- No external dependencies
- Built-in MQTT broker
- Full control over data

**For Development/Testing**:
- Quick local setup without external services
- Easy to reset and reconfigure
- Setup mode for infrastructure changes
- Runner mode for production testing

### Revenue Streams

1. **Free Self-Hosted**: Single binary with embedded MQTT (MIT License)
2. **Pro License** ($299 one-time): External database adapters, advanced features
3. **Enterprise License** ($1,499/year): Priority support, custom builds
4. **Edge appliances**: Pre-configured hardware + software bundles

### Pricing Strategy

| Edition | Target Market | Pricing Model | Features |
|---------|--------------|---------------|----------|
| Community | Hobbyists, edge computing | Free (MIT) | Embedded MQTT, SQLite WAL, in-memory cache, setup mode, React UI |
| Professional | SMEs, production | $299 one-time | External DB adapters, persistent cache, API access, priority support |
| Enterprise | Large organizations | $1,499/year | Multiple gateways, clustering, priority support, custom builds |

---

## 3. Target Audience

### Primary Users

**1. Edge Computing Engineers**
- Deploy IoT gateways on edge devices
- Need minimal footprint
- Require offline operation
- Pain point: Unclear what to configure where

**2. Air-Gapped Environment Operators**
- Industrial IoT in isolated networks
- Cannot access external services
- Need complete control
- Pain point: Confusion between config and data

**3. Development Teams**
- Need local testing environment
- Want quick setup/teardown
- Prefer React ecosystem
- Pain point: Unclear configuration boundaries

**4. System Integrators**
- Deploy customer premise solutions
- Want simplified operations
- Need easy reconfiguration
- Pain point: Not knowing what requires restart

---

## 4. Problem Statement

### Core Problems

**Problem 1: Unclear Mode Purpose**
- v2 didn't clearly explain what Setup Mode does
- Users confused: "Can I change everything in setup mode?"
- Current approach: Vague "configuration changes" description
- Solution: Clear definition: Setup Mode = config.yaml editor ONLY

**Problem 2: Confusion About Data Placement**
- Not clear what goes in config.yaml vs database
- Users ask: "Where do I store device metadata?"
- Current approach: Ambiguous data separation
- Solution: Crystal-clear rule: Infrastructure (config.yaml) vs Runtime Data (database)

**Problem 3: Framework Familiarity**
- SvelteKit less familiar than React
- Smaller talent pool
- Current approach: SvelteKit for web UI
- Solution: React + Vite with TanStack ecosystem

**Problem 4: State Management**
- No clear server state management
- Manual synchronization between UI and backend
- Current approach: Basic API calls
- Solution: TanStack Query for automatic caching, synchronization, retries

**Problem 5: Routing Complexity**
- Type safety issues with routing
- Manual parameter handling
- Current approach: Basic routing
- Solution: TanStack Router for type-safe routing

### Impact

- **Configuration Errors**: Users editing wrong files/places
- **Unnecessary Restarts**: Changing database data thinking it requires restart
- **Developer Availability**: Harder to find SvelteKit developers
- **State Inconsistency**: UI showing stale data
- **Poor Developer Experience**: Manual type checking, route parameter bugs

---

## 5. Proposed Solution

### Vision: Crystal-Clear Configuration with Modern React Stack

A **single executable binary** with:
- **React + Vite UI**: Modern, familiar technology stack
- **TanStack Query**: Automatic server state management, caching, synchronization
- **TanStack Router**: Type-safe routing with excellent DX
- **Crystal-Clear Mode Definitions**: Setup Mode = config.yaml editor, Runner Mode = production
- **Clear Data Separation**: Infrastructure (config.yaml) vs Runtime Data (database)

### Key Innovations

**1. React + Vite + TanStack Stack**

```bash
# Frontend stack
- React 18+ (UI library)
- Vite 5+ (build tool)
- TypeScript (type safety)
- TanStack Query (server state management)
- TanStack Router (routing)
- TailwindCSS (styling)

# Backend stack
- Go 1.24+
- mochi-mqtt (embedded MQTT broker)
- SQLite WAL mode (default database)
- go-cache (in-memory cache)
```

**2. Crystal-Clear Mode Definitions**

**Setup Mode Flow**:
```
1. User starts gateway with --mode=setup OR first run detected
2. Web UI launches in setup mode at http://localhost:8080/setup
3. User sees config.yaml editor with validation
4. User changes infrastructure settings (database, MQTT, cache, etc.)
5. User clicks "Save & Restart"
6. System validates config.yaml
7. System creates backup of old config.yaml
8. System writes new config.yaml
9. System restarts and switches to runner mode
```

**Runner Mode Flow**:
```
1. User starts gateway normally (or after setup)
2. config.yaml is loaded and becomes READ-ONLY
3. Embedded MQTT broker starts on configured port
4. Web UI launches in runner mode at http://localhost:8080
5. User sees monitoring dashboard
6. User can view config.yaml (read-only)
7. User cannot edit config.yaml
8. User interacts with runtime data via database (API calls)
```

**3. Clear Data Separation Rules**

**Rule 1: Infrastructure Settings → config.yaml**

Ask: "Does changing this require restarting the gateway?"
- YES → Goes in config.yaml
- NO → Goes in database

Examples:
- ✅ "Which database to use?" → config.yaml (requires restart)
- ✅ "Which port to listen on?" → config.yaml (requires restart)
- ✅ "How many cache items?" → config.yaml (requires restart)
- ❌ "What topics exist?" → database (no restart needed)
- ❌ "Which devices are connected?" → database (no restart needed)

**Rule 2: Runtime Data → Database**

Ask: "Is this data that changes during normal operation?"
- YES → Goes in database
- NO → Goes in config.yaml

Examples:
- ✅ "MQTT messages" → database (changes constantly)
- ✅ "Device states" → database (changes during operation)
- ✅ "Topic subscriptions" → database (changes during operation)
- ❌ "Database file path" → config.yaml (static setting)
- ❌ "Log level" → config.yaml (static setting)

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
│  │  │ - YAML editor with validation                  │  │  │
│  │  │ - Form-based configuration                     │  │  │
│  │  │ - Test connection button                       │  │  │
│  │  │ - Save & Restart button                        │  │  │
│  │  └────────────────────────────────────────────────┘  │  │
│  │                                                        │  │
│  │  ┌────────────────────────────────────────────────┐  │  │
│  │  │ Runner Mode UI (/)                              │  │  │
│  │  │ - Monitoring dashboard                          │  │  │
│  │  │ - Read-only config viewer                      │  │  │
│  │  │ - Real-time metrics (TanStack Query)            │  │  │
│  │  │ - Logs viewer                                   │  │  │
│  │  └────────────────────────────────────────────────┘  │  │
│  │                                                        │  │
│  │  ┌────────────────────────────────────────────────┐  │  │
│  │  │ TanStack Router (Type-Safe Routing)             │  │  │
│  │  │ TanStack Query (Server State Management)        │  │  │
│  │  └────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

---

## 6. Project Objectives

### Primary Objectives (Must Have)

**O1: Implement React + Vite Frontend**
- Replace SvelteKit with React + Vite
- TypeScript for type safety
- TailwindCSS for styling
- Single-page application embedded in binary

**O2: Integrate TanStack Query**
- Server state management
- Automatic caching and revalidation
- Background refetching
- Optimistic updates
- Error handling and retries

**O3: Integrate TanStack Router**
- Type-safe routing
- Route parameters with type inference
- Nested routes and layouts
- Code splitting
- Search parameter handling

**O4: Clarify Setup Mode Purpose**
- Setup Mode = config.yaml editor ONLY
- Clear messaging: "Configure infrastructure settings"
- YAML editor with validation
- Form-based alternative interface
- Test configuration before saving

**O5: Clarify Runner Mode Purpose**
- Runner Mode = production operation
- config.yaml is READ-ONLY
- Monitoring dashboard
- Read-only config viewer
- Cannot edit infrastructure settings

**O6: Document Data Separation Clearly**
- Infrastructure → config.yaml
- Runtime data → database
- Decision tree: "Does it require restart?"
- Examples and rules
- Visual diagrams

### Secondary Objectives (Should Have)

**O7: Enhanced Configuration Validation**
- Real-time YAML validation
- Connection testing
- Schema validation
- Error messages with suggestions
- Configuration diff viewer

**O8: Improved Error Handling**
- User-friendly error messages
- Error recovery suggestions
- Configuration error details
- Log file integration
- Error reporting

**O9: Better Developer Experience**
- Hot reload in development
- TypeScript strict mode
- ESLint + Prettier
- Storybook for components
- Comprehensive documentation

### Tertiary Objectives (Nice to Have)

**O10: Configuration Import/Export**
- Export config.yaml
- Import from file
- Configuration templates
- Preset configurations

**O11: Configuration History**
- Track config.yaml changes
- Rollback to previous versions
- Change diff viewer
- Audit log

---

## 7. Project Scope & Boundaries

### In Scope (What We Will Build)

#### Core Functionality

1. **React + Vite Frontend**
   - SPA embedded in Go binary
   - TypeScript for type safety
   - TailwindCSS for styling
   - Responsive design

2. **TanStack Query Integration**
   - Query hooks for all API calls
   - Mutation hooks for data changes
   - Automatic caching and revalidation
   - Error and loading states

3. **TanStack Router Integration**
   - Type-safe routes
   - Route protection (setup vs runner)
   - Nested layouts
   - Search parameters

4. **Setup Mode UI**
   - YAML editor with syntax highlighting
   - Form-based configuration interface
   - Real-time validation
   - Test configuration button
   - Save & Restart button
   - Configuration backup

5. **Runner Mode UI**
   - Monitoring dashboard
   - Read-only config viewer
   - Real-time metrics
   - Logs viewer
   - Device management

6. **Clear Documentation**
   - Data separation rules
   - Mode definitions
   - Configuration examples
   - Decision trees

### Out of Scope (What We Won't Build - Initially)

1. **Advanced Features** (Phase 2)
   - Configuration history/versioning
   - Configuration import/export
   - Advanced monitoring and alerting

2. **High Availability** (Phase 3)
   - Multi-gateway clustering
   - Data replication
   - Automatic failover

### Boundaries & Constraints

#### Technical Constraints
- **Frontend**: React 18+, Vite 5+, TypeScript 5+
- **Backend**: Go 1.24+
- **MQTT Library**: mochi-mqtt
- **Minimum Hardware**: 1 CPU core, 1GB RAM

#### Time Constraints
- **Phase 1 (MVP)**: 10 months
- **Phase 2 (Advanced Features)**: +3 months
- **Phase 3 (HA Features)**: +6 months

#### Budget Constraints
- **Development Team**: 2-3 full-time developers
- **Infrastructure**: $1,500/month for development/testing
- **Third-party Services**: Free tier initially

---

## 8. Technology Stack

### Frontend Stack

| Component | Technology | Version | Rationale |
|-----------|-----------|---------|-----------|
| UI Library | React | 18.3+ | Largest ecosystem, familiar to most developers |
| Build Tool | Vite | 5.0+ | Fast HMR, optimized builds, great DX |
| Language | TypeScript | 5.3+ | Type safety, better developer experience |
| State Management | TanStack Query | 5.0+ | Powerful server state management, caching |
| Routing | TanStack Router | 1.0+ | Type-safe routing, excellent DX |
| Styling | TailwindCSS | 3.4+ | Utility-first, fast development |
| Code Quality | ESLint, Prettier | Latest | Consistent code formatting, linting |

### Backend Stack

| Component | Technology | Version | Rationale |
|-----------|-----------|---------|-----------|
| Language | Go | 1.24+ | Performance, concurrency, single binary |
| MQTT Broker | mochi-mqtt | 2.0+ | Embedded Go MQTT broker |
| Database | SQLite | 3.40+ | WAL mode support, single file |
| Cache | go-cache | Latest | In-memory caching with TTL |
| HTTP Router | Chi | 5.0+ | Lightweight, idiomatic Go |
| YAML | yaml.v3 | Latest | YAML parsing and serialization |

### TanStack Query Integration

```typescript
// Example: Query hook for fetching metrics
import { useQuery } from '@tanstack/react-query';

function useMetrics() {
  return useQuery({
    queryKey: ['metrics'],
    queryFn: async () => {
      const response = await fetch('/api/metrics');
      if (!response.ok) throw new Error('Failed to fetch metrics');
      return response.json();
    },
    refetchInterval: 5000, // Auto-refresh every 5 seconds
  });
}

// Usage in component
function MetricsDashboard() {
  const { data, error, isLoading } = useMetrics();

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div>
      <h2>Metrics</h2>
      <p>Connected Clients: {data.connectedClients}</p>
      <p>Messages/sec: {data.messagesPerSec}</p>
    </div>
  );
}
```

### TanStack Router Integration

```typescript
// Example: Type-safe routes
import { createRootRoute, createRoute, createRouter } from '@tanstack/react-router';

// Setup mode route (only accessible in setup mode)
const setupRoute = createRoute({
  path: '/setup',
  component: SetupMode,
  beforeLoad: async ({ location }) => {
    // Check if in setup mode
    const mode = await fetch('/api/mode').then(r => r.json());
    if (mode !== 'setup') {
      throw redirect({ to: '/' });
    }
  },
});

// Runner mode route (only accessible in runner mode)
const rootRoute = createRootRoute({
  component: RunnerMode,
});

// Router configuration
const router = createRouter({
  routeTree: rootRoute.addChildren([setupRoute]),
});

// Type-safe navigation
function NavigateToSetup() {
  const navigate = useNavigate();
  return (
    <button onClick={() => navigate({ to: '/setup' })}>
      Go to Setup
    </button>
  );
}
```

### Configuration Structure

#### config.yaml (Infrastructure - Static)

```yaml
# =============================================================================
# UM-Gateway Configuration File
# =============================================================================
# This file contains INFRASTRUCTURE SETTINGS ONLY
# These settings define HOW the gateway runs
# Changes to this file require restarting the gateway
#
# DO NOT store runtime data here (use database instead)
# =============================================================================

# Operation Mode
# - setup: Configuration mode (edit this file via web UI)
# - runner: Production mode (this file becomes read-only)
mode: runner

# Database Configuration
database:
  # Database type: sqlite, postgres, mongodb, influxdb
  type: sqlite

  connection:
    # SQLite specific settings
    path: ./data/gateway.db
    wal_mode: true
    max_open_conns: 25
    max_idle_conns: 5
    conn_max_lifetime: 300s

    # Cache configuration
    cache:
      # Cache type: memory, redis, memcached
      type: memory
      ttl: 3600
      max_size: 10000

# MQTT Broker Configuration
mqtt:
  # Broker type: embedded, external
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

# Paths Configuration
paths:
  data: ./data
  logs: ./logs
  pid: ./um-gateway.pid

# Web Server Configuration
web:
  address: :8080
  admin_api:
    enabled: true
    authentication:
      type: basic

# Data Retention Policy
retention:
  enabled: true
  default_ttl: 2592000  # 30 days
  cleanup_interval: 3600

# Logging Configuration
logging:
  level: info
  format: json
  output: stdout
```

#### Database Schema (Runtime Data - Dynamic)

```sql
-- Devices table (runtime data: which devices exist)
CREATE TABLE devices (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  type TEXT NOT NULL,
  metadata TEXT, -- JSON
  first_seen DATETIME DEFAULT CURRENT_TIMESTAMP,
  last_seen DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Topics table (runtime data: what topics exist)
CREATE TABLE topics (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  description TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Subscriptions table (runtime data: who subscribes to what)
CREATE TABLE subscriptions (
  id TEXT PRIMARY KEY,
  device_id TEXT NOT NULL,
  topic_id TEXT NOT NULL,
  qos INTEGER DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (device_id) REFERENCES devices(id),
  FOREIGN KEY (topic_id) REFERENCES topics(id)
);

-- Messages table (runtime data: actual MQTT messages)
CREATE TABLE messages (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  topic TEXT NOT NULL,
  payload TEXT NOT NULL,
  qos INTEGER DEFAULT 0,
  retained BOOLEAN DEFAULT FALSE,
  published_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Device states table (runtime data: current status)
CREATE TABLE device_states (
  device_id TEXT PRIMARY KEY,
  status TEXT NOT NULL,
  last_payload TEXT,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (device_id) REFERENCES devices(id)
);
```

---

## 9. Features

### 9.1 Core Features

#### F1: React + Vite Frontend

**Project Structure**:

```
frontend/
├── src/
│   ├── main.tsx           # Entry point
│   ├── App.tsx            # Root component
│   ├── routes/            # Route definitions
│   │   ├── __root.tsx     # Root layout
│   │   ├── index.tsx      # Runner mode route
│   │   └── setup.tsx      # Setup mode route
│   ├── pages/             # Page components
│   │   ├── SetupMode.tsx  # Setup mode page
│   │   ├── RunnerMode.tsx # Runner mode page
│   │   └── Dashboard.tsx  # Dashboard page
│   ├── components/        # Reusable components
│   │   ├── YAMLEditor.tsx
│   │   ├── MetricsCard.tsx
│   │   └── LogViewer.tsx
│   ├── hooks/             # Custom hooks
│   │   ├── useMetrics.ts
│   │   ├── useConfig.ts
│   │   └── useMode.ts
│   ├── lib/               # Utilities
│   │   ├── api.ts         # API client
│   │   └── query.ts       # TanStack Query config
│   └── styles/            # Styles
│       └── index.css      # Tailwind CSS
├── index.html
├── vite.config.ts
├── tsconfig.json
└── package.json
```

**TanStack Query Configuration**:

```typescript
// frontend/src/lib/query.ts
import { QueryClient } from '@tanstack/react-query';

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5000,        // Data fresh for 5 seconds
      gcTime: 1000 * 60 * 5, // Cache for 5 minutes
      retry: 3,              // Retry failed requests 3 times
      refetchOnWindowFocus: true,
    },
    mutations: {
      retry: 1,
    },
  },
});
```

---

#### F2: Setup Mode UI

**Purpose**: Edit config.yaml infrastructure settings

**Components**:

```typescript
// frontend/src/pages/SetupMode.tsx
import { useMutation, useQuery } from '@tanstack/react-query';

export function SetupMode() {
  // Fetch current config
  const { data: config, isLoading } = useQuery({
    queryKey: ['config'],
    queryFn: () => fetch('/api/config').then(r => r.json()),
  });

  // Save config mutation
  const saveConfig = useMutation({
    mutationFn: (newConfig: string) =>
      fetch('/api/config', {
        method: 'POST',
        body: newConfig,
        headers: { 'Content-Type': 'text/yaml' },
      }),
    onSuccess: () => {
      alert('Config saved! Restarting gateway...');
      setTimeout(() => window.location.reload(), 2000);
    },
  });

  if (isLoading) return <div>Loading configuration...</div>;

  return (
    <div className="setup-mode">
      <h1>Setup Mode - Infrastructure Configuration</h1>
      <p className="warning">
        ⚠️ You are editing INFRASTRUCTURE SETTINGS (config.yaml)
        <br />
        Changes will require restarting the gateway
      </p>

      <YAMLEditor
        value={config.yaml}
        onChange={(newValue) => setConfig(newValue)}
      />

      <div className="actions">
        <button onClick={() => saveConfig.mutate(config)}>
          Save & Restart
        </button>
        <button onClick={() => testConfig.mutate(config)}>
          Test Configuration
        </button>
      </div>
    </div>
  );
}
```

---

#### F3: Runner Mode UI

**Purpose**: Production operation monitoring

**Components**:

```typescript
// frontend/src/pages/RunnerMode.tsx
import { useQuery } from '@tanstack/react-query';

export function RunnerMode() {
  // Fetch metrics (auto-refresh every 5 seconds)
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
      <h1>Runner Mode - Production Dashboard</h1>
      <p className="info">
        ℹ️ Gateway is running in PRODUCTION mode
        <br />
        Configuration is read-only. To make changes, switch to setup mode.
      </p>

      <MetricsCard metrics={metrics} />
      <ConfigViewer config={config} />
      <LogViewer />
    </div>
  );
}
```

---

#### F4: Clear Data Separation Documentation

**Decision Tree**:

```
┌─────────────────────────────────────────────────────────────┐
│  Where should this data go?                                 │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  Q1: Does changing this require restarting the gateway?      │
│                                                              │
│     YES → config.yaml (Infrastructure Settings)               │
│     ┌─────────────────────────────────────────────────┐     │
│     │ Examples:                                        │     │
│     │ • Database type and connection settings          │     │
│     │ • MQTT broker settings (port, limits)            │     │
│     │ • Cache configuration                            │     │
│     │ • Web server settings                            │     │
│     │ • Logging configuration                          │     │
│     │ • File paths                                     │     │
│     └─────────────────────────────────────────────────┘     │
│                                                              │
│     NO → Database (Runtime Data)                             │
│     ┌─────────────────────────────────────────────────┐     │
│     │ Examples:                                        │     │
│     │ • Device registry                                │     │
│     │ • Topics and subscriptions                       │     │
│     │ • MQTT messages                                  │     │
│     │ • Device states                                  │     │
│     │ • Metrics and statistics                         │     │
│     │ • User accounts                                  │     │
│     └─────────────────────────────────────────────────┘     │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## 10. Development Roadmap

### Phase 1: Foundation (Months 1-5)

**Sprint 1-2: Project Setup & Architecture**
- [ ] Monorepo structure (Go + React + Vite)
- [ ] React + Vite frontend setup
- [ ] TanStack Query configuration
- [ ] TanStack Router setup
- [ ] TypeScript configuration
- [ ] CI/CD pipeline

**Sprint 3-4: Backend Core**
- [ ] YAML configuration loading
- [ ] Mode detection and switching
- [ ] Embedded MQTT broker integration
- [ ] SQLite WAL mode setup
- [ ] API endpoints for config and metrics

**Sprint 5-6: Frontend Core**
- [ ] Setup mode UI skeleton
- [ ] Runner mode UI skeleton
- [ ] TanStack Query hooks
- [ ] TanStack Router routes
- [ ] Basic layouts

**Sprint 7-8: Setup Mode Implementation**
- [ ] YAML editor component
- [ ] Configuration validation
- [ ] Test configuration button
- [ ] Save & restart functionality
- [ ] Configuration backup

**Sprint 9-10: Runner Mode Implementation**
- [ ] Monitoring dashboard
- [ ] Read-only config viewer
- [ ] Real-time metrics
- [ ] Logs viewer
- [ ] Auto-refresh with TanStack Query

**Milestone**: Alpha - Basic functionality complete

---

### Phase 2: Polish & Features (Months 6-8)

**Sprint 11-12: Enhanced UI**
- [ ] Responsive design
- [ ] Error handling
- [ ] Loading states
- [ ] Toast notifications
- [ ] Modal dialogs

**Sprint 13-14: Documentation**
- [ ] Data separation guide
- [ ] Mode definitions
- [ ] Configuration examples
- [ ] Decision trees
- [ ] API documentation

**Sprint 15-16: Testing**
- [ ] Unit tests (React components)
- [ ] Integration tests (API)
- [ ] E2E tests (Playwright)
- [ ] Load testing
- [ ] Security audit

**Milestone**: Beta - Feature-complete v3

---

### Phase 3: Production (Months 9-10)

**Sprint 17-18: Hardening**
- [ ] Error boundary implementation
- [ ] Graceful degradation
- [ ] Performance optimization
- [ ] Bundle size optimization
- [ ] Accessibility improvements

**Sprint 19-20: Deployment**
- [ ] Production builds
- [ ] Cross-platform compilation
- [ ] Installation guides
- [ ] User documentation
- [ ] Release notes

**Milestone**: v3.0 General Availability

---

## 11. Success Metrics

### Technical Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| Frontend Bundle Size | < 500 KB (gzipped) | Build artifacts |
| Initial Page Load | < 2 seconds | Lighthouse |
| Time to Interactive | < 3 seconds | Lighthouse |
| API Response Time | < 100ms | Automated tests |
| TanStack Query Cache Hit Rate | > 90% | Cache statistics |
| TypeScript Coverage | 100% | TSC --noEmit |

### Business Metrics

| Metric | Target | Timeline |
|--------|--------|----------|
| Active Installations | 1,000 | 6 months post-launch |
| GitHub Stars | 2,500 | 6 months post-launch |
| Average Session Duration | > 10 minutes | Analytics |
| Setup Completion Rate | > 95% | Analytics |

---

## 12. Risk Analysis

### Technical Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| TanStack Query complexity | Medium | Low | Comprehensive documentation, examples |
| TanStack Router learning curve | Medium | Low | Training, code examples |
| React bundle size | Low | Medium | Code splitting, lazy loading |
| TypeScript configuration issues | Low | Low | Strict config from start |

### Business Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Users prefer SvelteKit | Low | Low | React has larger ecosystem |
| TanStack adoption issues | Low | Low | Well-established libraries |
| Framework fatigue | Medium | Medium | Stable, proven technologies |

---

## 13. Resource Requirements

### Team Structure

**Core Team (Phase 1-3)**
- 2x Full-Stack Developers (Go + React)
- 1x Frontend Specialist (React + TanStack)
- 1x DevOps Engineer
- 1x QA Engineer (part-time)
- 1x Technical Writer (part-time)

**Budget Estimate**

| Category | Cost (Monthly) | Duration |
|----------|----------------|---------|
| Development Team Salaries | $20,000 | 10 months |
| Infrastructure (Dev/Test) | $1,500 | 10 months |
| Tools & Services | $400 | Ongoing |
| Documentation & Design | $1,000 | 8 months |
| Contingency (15%) | $2,865 | - |
| **Total Phase 1-3** | **$21,500/mo × 10** | **$215,500** |

---

## 14. Comparison: v2 vs v3

| Aspect | v2 (SvelteKit) | v3 (React + TanStack) |
|--------|----------------|----------------------|
| **Frontend Framework** | SvelteKit | React + Vite |
| **State Management** | Basic | TanStack Query |
| **Routing** | SvelteKit router | TanStack Router |
| **Type Safety** | Good | Excellent (TanStack Router) |
| **Mode Clarity** | Vague | Crystal-clear |
| **Data Separation** | Unclear | Documented rules |
| **Talent Pool** | Smaller | Larger (React) |
| **Bundle Size** | Smaller | Slightly larger |
| **Learning Curve** | Medium | Low (React popular) |
| **Budget** | $238,150 | $215,500 |

---

## 15. Next Steps

### Immediate Actions (Week 1-2)

1. **Technology Research**
   - Evaluate TanStack Query features
   - Test TanStack Router capabilities
   - Prototype React + Vite integration

2. **Architecture Design**
   - Define route structure
   - Plan query hooks
   - Design component hierarchy

3. **Team Formation**
   - Recruit React developers
   - TanStack training if needed

### Short-term Actions (Month 1)

4. **Prototype Development**
   - React + Vite frontend skeleton
   - TanStack Query integration
   - TanStack Router setup

5. **Documentation**
   - Data separation guide
   - Mode definitions
   - Decision tree

---

## 16. Conclusion

The v3 version represents the **most developer-friendly and clear approach**:

### Key Improvements over v2

**Technology Stack**
- ✅ React + Vite (larger ecosystem)
- ✅ TanStack Query (powerful state management)
- ✅ TanStack Router (type-safe routing)
- ✅ TypeScript (excellent DX)

**Clarity**
- ✅ Crystal-clear mode definitions
- ✅ Setup Mode = config.yaml editor
- ✅ Runner Mode = production operation
- ✅ Clear data separation rules
- ✅ Decision trees and examples

**Developer Experience**
- ✅ Familiar React ecosystem
- ✅ Type-safe routing
- ✅ Automatic state synchronization
- ✅ Better error handling
- ✅ Larger talent pool

**Lower Budget**
- ✅ $215,500 (vs $238,150 for v2)
- ✅ More React developers available
- ✅ Lower training costs

With focused execution over 10 months, we can deliver a **production-ready v3.0** with excellent clarity and developer experience.

---

**Document Version**: 3.0
**Last Updated**: January 9, 2026
**Status**: Draft - Pending Review
**Supersedes**: PROJECT_BRIEF_EN_v2.md

---

## Appendix

### A. Data Separation Examples

**config.yaml contains** (Infrastructure):
```yaml
database:
  type: sqlite           # Which database to use
  path: ./data/gateway.db  # Where to store it

mqtt:
  type: embedded         # Use embedded broker
  embedded:
    listen_address: :1883  # Which port to listen on
```

**Database contains** (Runtime Data):
```sql
-- Device "temp-sensor-1" exists
INSERT INTO devices (id, name, type) VALUES ('temp-sensor-1', 'Temperature Sensor', 'sensor');

-- Topic "sensors/temperature" exists
INSERT INTO topics (id, name) VALUES ('topic-1', 'sensors/temperature');

-- Device subscribed to topic
INSERT INTO subscriptions (device_id, topic_id) VALUES ('temp-sensor-1', 'topic-1');

-- Message published
INSERT INTO messages (topic, payload) VALUES ('sensors/temperature', '{"value": 25.5}');
```

### B. References

1. TanStack Query documentation
2. TanStack Router documentation
3. React documentation
4. Vite documentation
5. TypeScript documentation

---

**Prepared by**: Development Team
**Approved by**: [Pending]
**Review Date**: [Pending]
