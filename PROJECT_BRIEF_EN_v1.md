# Project Brief: Universal MQTT Gateway Server (UM-Gateway) v1.0

## Executive Summary

**Project Name**: Universal MQTT Gateway Server (UM-Gateway)
**Version**: 1.0 (Self-Hosted with Web Admin)
**Status**: Proposal/Planning
**Target Release**: Q3 2026
**Organization**: Rafflesia Agro Technology Division

---

## Version History

| Version | Date | Changes | Author |
|---------|------|---------|--------|
| **v1.0** | 2026-01-09 | **Web Admin Dashboard**: Visual configuration approach<br>• Added SvelteKit admin dashboard for configuration<br>• SQLite database with onboarding wizard<br>• Replaced JSON configuration with web UI forms<br>• External MQTT broker connection model<br>• Budget: $178,650 (reduced from v0)<br>• Self-contained deployment like PocketBase<br>• Zero-config setup with automatic database initialization | Development Team |
| **v0.0** | 2026-01-09 | **Original Proposal**: JSON-based configuration<br>• Initial concept with JSON file definitions<br>• External MQTT broker connection required<br>• PostgreSQL/MongoDB configurable via JSON<br>• Focus on farm telemetry use case<br>• Budget: $283,200<br>• Manual configuration and deployment | Development Team |

---

## 1. Background

### Current State
The existing Go-MQTT server (v1.0) was developed as a specialized solution for Rafflesia Agro's smart poultry farming operations. While successful in its domain, the current implementation suffers from several limitations:

1. **Tight Coupling**: Business logic is hard-coded for farm telemetry use case
2. **JSON Configuration Complexity**: Requires technical knowledge to edit JSON files
3. **No User Interface**: All configuration done through text files
4. **Manual Deployment**: Each change requires server restart
5. **Limited Accessibility**: No web-based management interface

### Evolution from v0 to v1
**v0 Approach**: JSON-based configuration files that require manual editing and server restarts
**v1 Approach**: Self-contained web admin dashboard (like PocketBase) with embedded SQLite database for zero-config setup

### Market Opportunity
The IoT landscape needs a **self-hosted, user-friendly MQTT gateway** that non-technical users can deploy and manage:

- **Small Teams**: Want to avoid cloud SaaS costs but need easy management
- **IT Teams**: Prefer self-hosted solutions with web interfaces over file-based config
- **Startups**: Need quick deployment without learning complex configuration syntax
- **Edge Deployments**: Require local management without external dependencies

---

## 2. Business Model

### Value Proposition

**For Non-Technical Users**:
- Zero-config deployment - just run and access web interface
- Visual setup instead of editing JSON files
- Onboarding wizard walks through first-time setup
- No need to learn configuration syntax

**For System Integrators**:
- Faster deployment (minutes vs hours)
- Easy client demonstrations
- Remote management through web UI
- Client self-service reduces support burden

**For Rafflesia Agro**:
- Lower barrier to entry increases adoption
- Differentiates from JSON-config competitors
- Professional appearance with admin dashboard
- Potential for hosted version with same UX

### Revenue Streams

1. **Free Self-Hosted**: Basic version with SQLite (100% free, MIT License)
2. **Pro License** ($199 one-time): Advanced features, multi-database support
3. **Enterprise License** ($999/year): Priority support, white-label options
4. **Cloud Service** ($29/month): Hosted version with same UI/UX
5. **Support Contracts**: Custom integrations and consultations

### Pricing Strategy

| Edition | Target Market | Pricing Model | Features |
|---------|--------------|---------------|----------|
| Community | Hobbyists, students, small projects | Free (MIT) | SQLite, web admin, basic MQTT gateway |
| Pro | SMEs, professional use | $199 one-time | PostgreSQL/MongoDB, advanced features, email support |
| Enterprise | Large organizations | $999/year | Multi-user, RBAC, priority support, SLA |
| Cloud | All markets | $29/month | Fully hosted, auto-backups, 99.9% SLA |

---

## 3. Target Audience

### Primary Users

**1. Non-Technical Business Owners**
- Need IoT monitoring but lack coding skills
- Want "it just works" solution
- Comfortable with web forms and dashboards
- Pain point: Cannot edit JSON configuration files

**2. Small IT Teams**
- Manage internal IoT projects
- Prefer web UI over config files
- Need to delegate to non-technical staff
- Pain point: Team members have varying technical skills

**3. Startup Founders**
- Building IoT MVP quickly
- Need to iterate on configuration
- Don't want to maintain complex infrastructure
- Pain point: Time-to-market pressure, limited resources

**4. MSPs (Managed Service Providers)**
- Manage deployments for multiple clients
- Need consistent management interface
- Want to delegate configuration to clients
- Pain point: Managing diverse client deployments efficiently

### Secondary Users

**5. Students & Educators**
- Learning IoT concepts
- Need visual feedback
- Want to understand without deep technical knowledge

---

## 4. Problem Statement

### Core Problems

**Problem 1: JSON Configuration is Too Technical**
- JSON syntax errors break deployments
- No validation until server restart
- Difficult for non-technical users
- No visual feedback on configuration structure

**Problem 2: No User Interface for Management**
- All configuration requires SSH/file access
- Cannot manage remotely without server access
- No visual overview of subscriptions
- Difficult to troubleshoot without technical skills

**Problem 3: File-Based Configuration is Error-Prone**
- Manual file editing leads to syntax errors
- No configuration versioning
- Difficult to rollback changes
- No audit trail of configuration changes

**Problem 4: High Barrier to Entry**
- Requires knowledge of MQTT, databases, JSON
- No guided setup process
- Documentation-heavy learning curve
- No visual feedback during configuration

**Problem 5: Remote Management is Difficult**
- Cannot configure without server access
- No collaboration on configuration
- No multi-user management
- Difficult to manage multiple deployments

### Impact

- **Adoption Barrier**: Technical complexity prevents widespread use
- **Support Burden**: Simple configuration errors require technical support
- **Deployment Time**: Hours/days vs minutes for web-based solutions
- **User Frustration**: Technical users frustrated with UI limitations, non-technical users cannot use it
- **Lost Opportunities**: Market prefers web-managed solutions (evidenced by Home Assistant, Node-RED success)

---

## 5. Proposed Solution

### Vision: PocketBase for MQTT Gateway

A **self-hosted MQTT gateway with embedded web admin dashboard** that provides:
- **Zero-configuration deployment** - single binary, no setup required
- **Web-based administration** - everything managed through browser
- **SQLite database** - embedded, no external database needed for basic use
- **Visual configuration** - forms, wizards, drag-and-drop interfaces
- **Onboarding experience** - guided first-time setup

### Key Innovations

**1. Embedded Web Admin Dashboard**
- Built with SvelteKit (like PocketBase)
- Served directly by Go binary
- Single-page application for instant feedback
- Responsive design for mobile/tablet access

**2. First-Run Onboarding Experience**
- Setup wizard on first access
- Admin account creation
- Initial MQTT broker configuration
- First topic subscription setup
- Database selection (SQLite default, optional PostgreSQL/MongoDB)

**3. SQLite for Configuration & Data**
- Embedded database (zero external dependencies)
- Stores all configuration, users, subscriptions
- Optional upgrade to PostgreSQL/MongoDB for scale
- Automatic migrations on version updates

**4. Visual Topic Builder**
- Form-based topic subscription creation
- Real-time validation
- Preview of subscription configuration
- Test subscription immediately

**5. Self-Contained Architecture**
- Single binary includes web UI assets
- No separate web server needed
- No external dependencies for basic usage
- Easy deployment - just run the executable

### Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    UM-Gateway Server                         │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ MQTT         │  │ HTTP/WebSocket│  │ Embedded     │      │
│  │ Subscriber   │  │ API Server   │  │ Web Server   │      │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘      │
│         │                 │                 │               │
│  ┌──────▼─────────────────▼─────────────────▼───────┐      │
│  │           Message Processing Engine               │      │
│  │  - Topic Routing (Database-configured)          │      │
│  │  - Schema Validation                             │      │
│  │  - Data Transformation                           │      │
│  │  - State Management                              │      │
│  └──────┬────────────────────────────────────────────┘      │
│         │                                                  │
│  ┌──────▼────────────────────────────────────────────┐     │
│  │          Storage Layer (Pluggable)                │     │
│  │  ┌──────────────┐  ┌──────────────┐              │     │
│  │  │ SQLite       │  │ PostgreSQL   │              │     │
│  │  │ (Embedded)   │  │ (Optional)   │              │     │
│  │  └──────────────┘  └──────────────┘              │     │
│  │  ┌──────────────┐  ┌──────────────┐              │     │
│  │  │ MongoDB      │  │ InfluxDB     │              │     │
│  │  │ (Optional)   │  │ (Optional)   │              │     │
│  │  └──────────────┘  └──────────────┘              │     │
│  └──────────────────────────────────────────────────┘     │
│                                                           │
│  ┌──────────────────────────────────────────────────┐    │
│  │         Admin Dashboard (SvelteKit SPA)          │    │
│  │  - Onboarding Wizard                              │    │
│  │  - Topic Management UI                            │    │
│  │  - Configuration Builder                          │    │
│  │  - Real-time Monitoring                           │    │
│  │  - User Management                                │    │
│  └──────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

---

## 6. Project Objectives

### Primary Objectives (Must Have)

**O1: Develop Self-Contained Web Admin Dashboard**
- SvelteKit-based SPA embedded in Go binary
- Responsive design for all screen sizes
- Dark/light theme support
- Multi-language support (EN, ID initially)

**O2: Implement First-Run Onboarding Experience**
- Setup wizard for initial configuration
- Admin account creation with secure password
- MQTT broker connection testing
- First topic subscription creation
- Success confirmation with next steps

**O3: SQLite Database for Configuration**
- Embedded SQLite for config storage
- Automatic schema migrations
- Optional upgrade path to PostgreSQL/MongoDB
- Backup/restore functionality

**O4: Visual Topic Subscription Builder**
- Form-based topic subscription creation
- Real-time MQTT topic validation
- Transformation rule builder (visual)
- Storage destination selector
- Test subscription before activating

**O5: User Authentication & Authorization**
- Admin account creation on first run
- Additional user management (Pro feature)
- Role-based access control (Admin, Editor, Viewer)
- Session management with timeout

### Secondary Objectives (Should Have)

**O6: Real-Time Monitoring Dashboard**
- Live message throughput display
- Active subscription status
- Connection health indicators
- Storage capacity monitoring

**O7: Configuration Import/Export**
- Export configuration as JSON
- Import configuration from file
- Configuration templates
- Share configurations between deployments

**O8: Multi-Database Support (Pro)**
- Upgrade from SQLite to PostgreSQL
- MongoDB for flexible schemas
- InfluxDB for time-series optimization
- One-click migration wizard

**O9: API for Automation**
- RESTful API matching web UI functionality
- API key generation
- API documentation within admin UI
- Code examples for common operations

### Tertiary Objectives (Nice to Have)

**O10: Mobile App**
- React Native or Flutter mobile app
- Push notifications for alerts
- Quick status checks
- Simplified mobile UI

**O11: Plugin System**
- Custom data processors
- Custom storage adapters
- Custom authentication providers
- Plugin marketplace

---

## 7. Project Scope & Boundaries

### In Scope (What We Will Build)

#### Core Functionality
1. **Embedded Web Server**
   - Serve SvelteKit SPA from embedded filesystem
   - WebSocket support for real-time updates
   - Automatic HTTPS with Let's Encrypt (optional)

2. **Admin Dashboard**
   - Setup wizard for first-time configuration
   - Topic subscription management
   - MQTT broker connection management
   - Storage destination configuration
   - User management (Pro feature)

3. **SQLite Database**
   - Configuration storage
   - User authentication data
   - Topic subscriptions
   - Audit logs
   - Optional time-series data storage

4. **MQTT Processing**
   - Subscribe to topics configured through UI
   - Message routing based on database config
   - Data transformation
   - State management

5. **Storage Options**
   - SQLite (embedded, default)
   - PostgreSQL (Pro feature)
   - MongoDB (Pro feature)
   - InfluxDB (Pro feature)

#### Web UI Features
- Onboarding wizard
- Dashboard home
- Topic management
- Storage management
- User management (Pro)
- Settings page
- System status

#### Deployment & Operations
- Single binary deployment
- Docker images
- Systemd service files
- Cross-platform builds (Linux, macOS, Windows)
- Auto-update mechanism (optional)

### Out of Scope (What We Won't Build - Initially)

1. **Hardware/Firmware**
   - Device firmware (ESP32, Arduino, etc.)
   - Hardware provisioning
   - OTA updates

2. **Advanced Analytics**
   - Data visualization dashboards (can integrate with Grafana)
   - Business intelligence
   - Predictive analytics

3. **Multi-Tenancy** (Phase 2)
   - Multiple organizations in single instance
   - Resource quotas per tenant

4. **Edge Computing** (Phase 3)
   - Lightweight version
   - Offline operation mode

### Boundaries & Constraints

#### Technical Constraints
- **Language**: Go 1.24+ (backend), SvelteKit (frontend)
- **MQTT Version**: 3.1.1 and 5.0 support
- **Minimum Hardware**: 1 CPU core, 512MB RAM for basic usage
- **Network**: Connectivity to MQTT broker and optional external databases

#### Time Constraints
- **Phase 1 (MVP)**: 9 months
- **Phase 2 (Pro Features)**: +3 months
- **Phase 3 (Advanced)**: +6 months

#### Budget Constraints
- **Development Team**: 2-3 full-time developers
- **Infrastructure**: $1,000/month for development/testing
- **Third-party Services**: Free tier initially

#### Resource Constraints
- **Development Team**: 1-2 Full-stack developers (Go + SvelteKit), 1 UI/UX designer
- **Domain Experts**: IoT specialists for consultation
- **Support**: Part-time technical writer

---

## 8. Technology Stack

### Core Technologies

#### Backend
| Component | Technology | Rationale |
|-----------|-----------|-----------|
| Language | Go 1.24+ | Performance, concurrency, single binary |
| HTTP Router | Chi v5 | Lightweight, idiomatic Go |
| MQTT Client | Eclipse Paho MQTT v2 | Mature, MQTT 5.0 support |
| Embedded Files | embed package | Bundle web UI in binary |
| SQLite | mattn/go-sqlite3 | Pure Go SQLite driver |
| WebSocket | gorilla/websocket | Real-time UI updates |

#### Frontend (Admin Dashboard)
| Component | Technology | Rationale |
|-----------|-----------|-----------|
| Framework | SvelteKit | Modern, fast, great UX, small bundle size |
| UI Components | Skeleton UI | Svelte component library, beautiful defaults |
| State Management | Svelte stores | Built-in, simple, reactive |
| Forms | SvelteKit form actions | Native form handling |
| Charts | Chart.js / ApexCharts | Visualization |
| Icons | Lucide Svelte | Lightweight icon library |
| Styling | TailwindCSS | Rapid UI development |

#### Database Options
| Database | Use Case | Advantages |
|----------|----------|------------|
| **SQLite** | Default, embedded | Zero-config, single file, sufficient for moderate loads |
| **PostgreSQL** | Production scaling | SQL, ACID, mature ecosystem |
| **MongoDB** | Flexible schemas | No migrations, horizontal scaling |
| **InfluxDB** | Time-series optimization | Purpose-built, high write throughput |

#### Authentication & Security
| Component | Technology | Rationale |
|-----------|-----------|-----------|
| Password Hashing | bcrypt | Industry standard |
| Session Storage | SQLite (default) / Redis (optional) | Embedded or scalable |
| JWT | golang-jwt/jwt | Stateless auth for API |
| CSRF Protection | Built-in to SvelteKit | Security best practice |

#### Deployment
| Component | Technology | Rationale |
|-----------|-----------|-----------|
| Build | Go embed + SvelteKit static adapter | Single binary |
| Containerization | Docker | Industry standard |
| Cross-platform | Go cross-compilation | Single codebase, multiple platforms |

### Architecture Pattern

**1. Embedded Web UI**
```go
//go:embed all:dist
var uiFS embed.FS

func (s *Server) ServeDashboard(w http.ResponseWriter, r *http.Request) {
    http.FileServer(http.FS(uiFS)).ServeHTTP(w, r)
}
```

**2. Database-Driven Configuration**
```sql
-- SQLite schema
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

**3. Real-Time Updates**
```go
// WebSocket endpoint for dashboard
func (s *Server) DashboardWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, _ := upgrader.Upgrade(w, r, nil)
    for {
        // Push updates: message counts, connection status, etc.
    }
}
```

---

## 9. Features

### 9.1 Core Features

#### F1: First-Run Onboarding Wizard

**Description**: Guided setup experience for first-time users

**Screens**:

**Screen 1: Welcome**
```
┌─────────────────────────────────────┐
│  Welcome to UM-Gateway!             │
│                                     │
│  Let's get you set up in 3 steps   │
│                                     │
│  [Get Started]                      │
└─────────────────────────────────────┘
```

**Screen 2: Create Admin Account**
```
┌─────────────────────────────────────┐
│  Create Your Admin Account          │
│                                     │
│  Email: [________________]          │
│  Password: [________________]       │
│  Confirm: [________________]       │
│                                     │
│  [Back]  [Next]                     │
└─────────────────────────────────────┘
```

**Screen 3: Connect MQTT Broker**
```
┌─────────────────────────────────────┐
│  Connect to MQTT Broker             │
│                                     │
│  Broker Name: [My Broker]          │
│  Broker URL: [mqtt://localhost...]│
│  Port: [1883]                       │
│  Username: [________________]      │
│  Password: [________________]      │
│                                     │
│  [Test Connection]                  │
│  ✓ Connected successfully!         │
│                                     │
│  [Back]  [Next]                     │
└─────────────────────────────────────┘
```

**Screen 4: Create First Topic Subscription**
```
┌─────────────────────────────────────┐
│  Subscribe to Topics                │
│                                     │
│  Subscription Name: [Farm Sensors] │
│  Topic Pattern: [farms/+/sensors/#]│
│  QoS: [0 ▼]                         │
│                                     │
│  □ Enable data transformation       │
│  □ Add metadata                     │
│                                     │
│  Storage: SQLite ▼                  │
│                                     │
│  [Test Subscription]                │
│  ✓ Receiving messages...           │
│                                     │
│  [Back]  [Finish]                   │
└─────────────────────────────────────┘
```

**Screen 5: Success**
```
┌─────────────────────────────────────┐
│  You're All Set! 🎉                 │
│                                     │
│  Your MQTT gateway is now running  │
│                                     │
│  Active Subscriptions: 1           │
│  Messages Received: 127            │
│                                     │
│  [Go to Dashboard]                  │
│  [View Documentation]               │
└─────────────────────────────────────┘
```

---

#### F2: Admin Dashboard Home

**Layout**:
```
┌─────────────────────────────────────────────────────────────┐
│  UM-Gateway  [☰]  [Admin ▼]  [●]              [Search...]  │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   127       │  │     1       │  │   Running   │         │
│  │  Messages/s │  │  Active Sub │  │    Status   │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│                                                               │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │  Recent Activity                                        │ │
│  │  • New subscription: farm/#/sensors  (2 min ago)      │ │
│  │  • Configuration updated (1 hour ago)                  │ │
│  │  • User 'admin' logged in (3 hours ago)                │ │
│  └─────────────────────────────────────────────────────────┘ │
│                                                               │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │  Quick Actions                                          │ │
│  │  [+ Add Subscription]  [+ Add Broker]  [⚙ Settings]   │ │
│  └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

---

#### F3: Visual Topic Subscription Builder

**Form-Based Interface**:

```
┌─────────────────────────────────────────────────────────────┐
│  New Topic Subscription                           [×]        │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  Basic Information                                          │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Name *        [________________]                    │   │
│  │ Description   [________________]                    │   │
│  │ Enabled       ☑                                     │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  MQTT Configuration                                         │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Topic Pattern * [sensors/+/telemetry/#]            │   │
│  │                  ✓ Valid pattern                   │   │
│  │                                                      │   │
│  │ QoS Level     ○ None  ● At least once  ○ Exactly  │   │
│  │                                                      │   │
│  │ [Test Connection]  Receiving 23 messages/min       │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  Storage Destination                                         │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Destination   SQLite ▼                               │   │
│  │ Table         [sensor_readings]                     │   │
│  │                                                      │   │
│  │ [+ Create New Table]                                │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  Transformation (Optional)                                  │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ ☐ Enable data transformation                        │   │
│  │                                                      │   │
│  │ [+ Add Transformation Rule]                          │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  [Cancel]  [Save & Activate]                                │
└─────────────────────────────────────────────────────────────┘
```

---

#### F4: SQLite Database Management

**Visual Schema Builder**:

```
┌─────────────────────────────────────────────────────────────┐
│  Data Storage                                   [×]        │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  Current Databases                                           │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ 🔵 main (SQLite)                   [Manage] [Export] │   │
│  │    Tables: 5  Size: 127 MB  Records: 1.2M          │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  [+ Add Database]                                            │
│                                                               │
│  Tables in 'main'                                            │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ sensor_readings                    [Browse] [Edit]   │   │
│  │   Columns: id, timestamp, device_id, value, ...     │   │
│  │   Records: 1,234,567                                  │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ device_states                      [Browse] [Edit]   │   │
│  │   Columns: id, device_id, state, updated_at        │   │
│  │   Records: 45                                        │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  [+ Create Table]  [Backup Database]                        │
└─────────────────────────────────────────────────────────────┘
```

---

#### F5: Real-Time Monitoring Dashboard

```
┌─────────────────────────────────────────────────────────────┐
│  Monitoring Dashboard                                        │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  System Health                                               │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  │
│  │  ✓ MQTT  │  │  ✓ DB    │  │  ✓ API   │  │  ✓ Web   │  │
│  │ Connected│  │  Healthy │  │  Running │  │  Serving │  │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘  │
│                                                               │
│  Message Throughput (Last 24h)                              │
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
│  Active Subscriptions                                        │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ farm/#/sensors/#     45 msg/s  ● Active             │   │
│  │ building/+/temp/#     12 msg/s  ● Active             │   │
│  │ devices/+/state       3 msg/s   ● Active             │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

---

#### F6: User Management (Pro Feature)

```
┌─────────────────────────────────────────────────────────────┐
│  Users                                               [+ Add] │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  admin@example.com                    [Edit] [×]    │   │
│  │  Role: Admin  ●  Last active: 2 min ago            │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  user@company.com                      [Edit] [×]    │   │
│  │  Role: Editor  ○  Last active: 1 day ago            │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  Roles                                                        │
│  • Admin - Full access                                       │
│  • Editor - Can manage subscriptions, not users              │
│  • Viewer - Read-only access                                  │
└─────────────────────────────────────────────────────────────┘
```

---

### 9.2 Advanced Features (Pro)

#### F7: Database Upgrade Wizard

**SQLite → PostgreSQL**:

```
┌─────────────────────────────────────────────────────────────┐
│  Upgrade Database                                             │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  Current: SQLite (main)                     [Backup Now]    │
│  Upgrade to: ○ PostgreSQL  ● MongoDB  ○ InfluxDB           │
│                                                               │
│  PostgreSQL Connection                                       │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Host: [localhost]                                   │   │
│  │ Port: [5432]                                        │   │
│  │ Database: [um_gateway_prod]                         │   │
│  │ Username: [____________]                             │   │
│  │ Password: [____________]                             │   │
│  │                                                      │   │
│  │ [Test Connection] ✓ Connected                        │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  Migration Options                                          │
│  ☑ Migrate configuration data                              │
│  ☑ Migrate time-series data (may take hours)              │
│  ☐ Keep SQLite as backup                                   │
│                                                               │
│  Estimated time: ~2 hours for 1.2M records                    │
│                                                               │
│  [Start Migration]  [Cancel]                                 │
└─────────────────────────────────────────────────────────────┘
```

---

#### F8: Configuration Import/Export

```
┌─────────────────────────────────────────────────────────────┐
│  Import/Export Configuration                                 │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  Export Configuration                                         │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Include:                                               │   │
│  │ ☑ Topic subscriptions                                 │   │
│  │ ☑ MQTT brokers                                        │   │
│  │ ☑ Storage destinations                                │   │
│  │ ☐ Users (if exporting for backup)                     │   │
│  │                                                      │   │
│  │ [Export as JSON]  [Export as File]                   │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  Import Configuration                                         │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ Drag & drop file here or click to browse              │   │
│  │                                                      │   │
│  │ Configuration file: um-gateway-config.json           │   │
│  │                                                      │   │
│  │ [Preview]  [Import]                                   │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                               │
│  Templates                                                   │
│  • Smart Farm Starter Template                             │
│  • Smart Home Basic Template                               │
│  • Industrial Monitoring Template                         │
└─────────────────────────────────────────────────────────────┘
```

---

## 10. Development Roadmap

### Phase 1: Foundation (Months 1-4)

**Sprint 1-2: Project Setup**
- [ ] Repository structure (Go + SvelteKit monorepo)
- [ ] Development environment setup
- [ ] CI/CD pipeline
- [ ] Design system setup

**Sprint 3-4: Backend Core**
- [ ] SQLite database schema
- [ ] Configuration models
- [ ] MQTT client refactoring
- [ ] Embedded web server
- [ ] WebSocket support

**Sprint 5-6: Onboarding Experience**
- [ ] First-run detection
- [ ] Admin account creation
- [ ] MQTT broker connection wizard
- [ ] First subscription setup
- [ ] Success confirmation flow

**Sprint 7-8: Basic Dashboard**
- [ ] SvelteKit project setup
- [ ] Layout and navigation
- [ ] Home dashboard
- [ ] Authentication UI
- [ ] Settings page

**Milestone**: Alpha - Self-hosted gateway with onboarding

---

### Phase 2: Core Features (Months 5-7)

**Sprint 9-10: Topic Management**
- [ ] Topic subscription builder UI
- [ ] Topic pattern validation
- [ ] Storage destination selector
- [ ] Test subscription feature
- [ ] Subscription list view

**Sprint 11-12: Storage Management**
- [ ] SQLite database browser
- [ ] Table creation UI
- [ ] Data viewer (paginated)
- [ ] Backup/restore functionality
- [ ] Export to CSV/JSON

**Sprint 13-14: Real-Time Monitoring**
- [ ] WebSocket real-time updates
- [ ] Throughput charts
- [ ] Connection status indicators
- [ ] Active subscriptions view
- [ ] System health page

**Milestone**: Beta - Feature-complete self-hosted gateway

---

### Phase 3: Polish & Production (Months 8-9)

**Sprint 15-16: Testing & Quality**
- [ ] End-to-end testing
- [ ] Load testing (1K msg/sec)
- [ ] Security audit
- [ ] Performance optimization
- [ ] Cross-platform testing

**Sprint 17-18: Documentation**
- [ ] User guide (with screenshots)
- [ ] Installation guide
- [ ] Troubleshooting guide
- [ ] Video tutorials (onboarding, common tasks)
- [ ] API documentation

**Sprint 19-20: Deployment**
- [ ] Single binary builds (Linux, macOS, Windows)
- [ ] Docker images
- [ ] Installation scripts
- [ ] Auto-update mechanism
- [ ] Release process automation

**Milestone**: v1.0 General Availability

---

## 11. Success Metrics

### Technical Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| Time to First Subscription | < 5 minutes | Automated testing |
| Dashboard Load Time | < 2 seconds | Performance monitoring |
| Binary Size | < 50 MB | Build artifacts |
| Memory Usage | < 200 MB (idle) | Resource monitoring |
| Messages/Second (SQLite) | 10K msg/sec | Benchmark tests |

### Business Metrics

| Metric | Target | Timeline |
|--------|--------|----------|
| Active Installations | 500 | 6 months post-launch |
| GitHub Stars | 1,000 | 6 months post-launch |
| Pro License Sales | 50 | 6 months post-launch |
| Return Rate | < 5% | 6 months post-launch |
| Avg Session Duration | > 10 minutes | Analytics |

### Quality Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| Setup Completion Rate | > 90% | Analytics |
| Onboarding Drop-off | < 10% per step | Analytics |
| Support Requests/100 Users | < 5 | Support tracking |
| User Satisfaction | > 4.5/5 | Surveys |
| Documentation Completeness | 100% features | Manual review |

---

## 12. Risk Analysis

### Technical Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| SQLite performance limits | Medium | Medium | Document limits, offer PostgreSQL upgrade |
| Embedded UI increases binary size | Low | High | Compression, lazy loading |
| SvelteKit learning curve | Medium | Low | Team training, prototype early |
| Cross-platform build complexity | Medium | Medium | GitHub Actions matrix builds |

### Business Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Competitors add web UI | High | Medium | Fast development, focus on UX |
| Market prefers cloud over self-hosted | Medium | Low | Offer both options |
| Free tier support burden | High | Medium | Community forum, documentation |

### Operational Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| User data loss | High | Low | Auto-backup prompts, export features |
| Security vulnerabilities | High | Low | Security audits, dependency scanning |
| Difficult upgrades | Medium | Medium | Auto-migration, extensive testing |

---

## 13. Resource Requirements

### Team Structure

**Core Team (Phase 1-3)**
- 1x Full-Stack Developer (Go + SvelteKit)
- 1x Frontend Developer (SvelteKit focus)
- 1x UI/UX Designer (part-time)
- 1x QA Engineer (part-time)

**Stakeholders**
- Product Manager
- Technical Lead
- Domain Experts (IoT consultants)

### Budget Estimate

| Category | Cost (Monthly) | Duration |
|----------|----------------|---------|
| Development Team Salaries | $15,000 | 9 months |
| Infrastructure (Dev/Test) | $1,000 | 9 months |
| Tools & Services | $300 | Ongoing |
| Design Resources | $1,000 | 6 months |
| Contingency (15%) | $2,550 | - |
| **Total Phase 1-3** | **$19,850/mo × 9** | **$178,650** |

---

## 14. Comparison: v0 vs v1

| Aspect | v0 (JSON Config) | v1 (Web Admin) |
|--------|------------------|-----------------|
| Configuration | JSON files | Web UI + SQLite |
| First Setup | Manual config | Onboarding wizard |
| User Type | Technical | Technical + Non-technical |
| Deployment Time | 30-60 minutes | 2-5 minutes |
| Barrier to Entry | High | Low |
| Management | SSH/File access | Web browser |
| Collaboration | Git/Version control | Multi-user with RBAC |
| Visual Feedback | None (edit files) | Real-time validation |
| Error Prevention | Runtime only | Before submission |
| Learning Curve | Steep | Gentle |

---

## 15. Next Steps

### Immediate Actions (Week 1-2)

1. **Stakeholder Review**
   - Present v0 vs v1 comparison
   - Validate web admin approach
   - Secure budget approval

2. **Technical Proof-of-Concept**
   - SvelteKit embedded in Go binary
   - SQLite configuration storage
   - Basic onboarding flow

3. **Design System Setup**
   - Choose UI component library
   - Create design mockups
   - Define onboarding flow

4. **Team Formation**
   - Recruit SvelteKit developer
   - Define collaboration process

### Short-term Actions (Month 1)

5. **Development Environment**
   - Monorepo structure (Go + SvelteKit)
   - Hot reload for both backend and frontend
   - CI/CD for cross-platform builds

6. **Sprint Planning**
   - Detailed sprint breakdown
   - Architecture decisions
   - Technology lock-in

7. **Prototype**
   - Working onboarding wizard
   - Basic topic subscription UI
   - SQLite integration

---

## 16. Conclusion

The v1 approach with **embedded web admin dashboard** dramatically improves upon the v0 JSON configuration approach:

### Key Improvements

**Accessibility**
- ✅ Non-technical users can now use the system
- ✅ No need to learn JSON syntax
- ✅ Visual feedback throughout setup

**User Experience**
- ✅ Guided onboarding reduces errors
- ✅ Real-time validation prevents misconfiguration
- ✅ Web-based management from anywhere

**Time to Value**
- ✅ 2-5 minutes to first subscription vs 30-60 minutes
- ✅ Self-service reduces support burden
- ✅ Visual configuration is faster than editing files

**Market Position**
- ✅ Differentiates from JSON-config competitors
- ✅ Competes with cloud SaaS on UX
- ✅ Appeals to broader audience

The **SvelteKit + SQLite + Embedded Go** approach follows the successful pattern of **PocketBase**, creating a self-hosted solution that's as easy to use as cloud alternatives while maintaining data privacy and control.

With focused execution over 9 months, we can deliver a **production-ready v1.0** that addresses real market needs and establishes a foundation for long-term growth.

---

**Document Version**: 1.0
**Last Updated**: January 9, 2026
**Status**: Draft - Pending Review
**Supersedes**: PROJECT_BRIEF_EN_v0.md

---

## Appendix

### A. Glossary

- **Self-Hosted**: Software run on user's own infrastructure, not cloud
- **Embedded**: Bundled within single executable/binary
- **Onboarding**: Guided setup process for new users
- **Zero-Config**: No manual configuration required
- **SPA**: Single Page Application

### B. References

1. PocketBase (inspiration for embedded web UI)
2. Home Assistant (visual IoT configuration)
3. Node-RED (visual programming for IoT)
4. SvelteKit documentation
5. SQLite embedded Go drivers

### C. Related Documents

- `PROJECT_BRIEF_EN_v0.md` - Original JSON-based approach
- `PROJECT_BRIEF_ID_v0.md` - Original approach (Bahasa Indonesia)
- `PROJECT_DESCRIPTION.md` - Current system documentation

---

**Prepared by**: Development Team
**Approved by**: [Pending]
**Review Date**: [Pending]
