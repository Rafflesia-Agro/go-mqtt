# Project Plan: Universal MQTT Gateway Server (UM-Gateway) v4.0

## Executive Summary

**Project Name**: Universal MQTT Gateway Server (UM-Gateway)
**Version**: 4.0 (Internal Database for Application Configuration)
**Status**: Proposal/Planning
**Target Release**: Q2 2027
**Organization**: Rafflesia Agro Technology Division

---

## Version History

| Version | Date | Changes | Author |
|---------|------|---------|--------|
| **v4.0** | 2026-01-09 | **Internal Database for Configuration**: Replacing config.yaml with internal.db<br>• Completely removing config.yaml<br>• Using embedded SQLite database (internal.db) for application configuration<br>• Separating internal.db (application config) and database.db (system data)<br>• PocketBase-like approach for configuration management<br>• Setup mode and runner mode remain the same<br>• API for configuration management (instead of file editing)<br>• Automatic migration from previous versions<br>• Budget: $225,000 (10 months) | Development Team |
| **v3.0** | 2026-01-09 | **Technology Stack Change & Clarity Improvement**: React-based UI<br>• Replacing SvelteKit with React + Vite<br>• Adding TanStack Query for server state management<br>• Adding TanStack Router for routing<br>• Clarifying Setup Mode = config.yaml editor<br>• Clarifying Runner Mode = production operation (config read-only)<br>• Clear separation: config.yaml (infrastructure) vs database (runtime data)<br>• Budget: $215,500 (10 months) | Development Team |

---

## 1. Background

### Evolution from v3 to v4

**v3 Approach**: Using config.yaml for infrastructure configuration
**v4 Approach**: Using embedded database (internal.db) for all configuration

### v3 Limitations (config.yaml)

1. **File Editing Issues**: Users must manually edit YAML files or use complex UI
2. **Validation Issues**: YAML validation can be complex with unclear error messages
3. **Migration Issues**: Version upgrades require manual file manipulation
4. **No Versioning**: No configuration change history
5. **Backup Issues**: Separate backup for config and database
6. **Access Issues**: File permissions can be problematic on some OS
7. **Consistency Issues**: No ACID to ensure configuration consistency

### Key Improvements in v4

**1. Replacing config.yaml with internal.db**

**internal.db** (Application Configuration):
- **Purpose**: Store all application and infrastructure settings
- **Location**: `./data/internal.db`
- **Format**: Embedded SQLite database
- **Access**: Only via API (safe and structured)
- **Backup**: Included in routine database backup
- **Migration**: Automatic with database schema
- **Transactional**: ACID compliant

**database.db** (System Data):
- **Purpose**: Store MQTT runtime data (messages, devices, topics, etc.)
- **Location**: `./data/database.db` (default, configurable)
- **Format**: Embedded SQLite database with WAL mode
- **Adapter**: Configurable for other databases (PostgreSQL, MongoDB, InfluxDB)
- **Access**: Via API for CRUD operations
- **Backup**: Automatic/periodic

**2. PocketBase-like Approach**

Adopting concepts from PocketBase:
- Single file executable with embedded database
- Admin UI for configuration management
- REST API for all operations
- Automatic schema migration
- Collections for structured data

**3. API for Configuration Management**

Instead of editing YAML files:
- `GET /api/config` - Get all configuration
- `PUT /api/config/:key` - Update configuration value
- `POST /api/config/validate` - Validate configuration
- `POST /api/config/reset` - Reset to defaults
- `GET /api/config/history` - Change history
- `POST /api/config/rollback` - Rollback to previous version

**4. Setup Mode and Runner Mode**

**Setup Mode**:
- Can change configuration in internal.db
- Form-based UI for settings (not YAML editor)
- Real-time validation
- Test connection
- Save & Restart

**Runner Mode**:
- Configuration in internal.db becomes read-only
- Cannot change application settings
- Normal production operation
- Dashboard monitoring

### Market Opportunity

The IoT landscape needs **MQTT gateway with modern configuration management**:
- **No File Editing**: All configuration via API/Web UI
- **Transactional**: Consistent and reliable configuration
- **Versioning**: Change history and rollback
- **Easy Migration**: Smooth upgrades
- **Modern Approach**: Like PocketBase, Strapi, etc.

---

## 2. Business Model

### Value Proposition

**For Edge Deployments**:
- No configuration file editing needed
- All settings via modern Web UI
- Safe and transactional configuration
- Change history and rollback

**For Air-Gapped Environments**:
- Fully offline operation
- Embedded database (no external database needed)
- Simple backup and restore
- Easy configuration cloning

**For Development/Testing**:
- Reset configuration with one click
- Snapshot and restore state
- Reproducible environments
- Automatic version migration

### Revenue Streams

1. **Free Self-Hosted**: Single binary with embedded databases (MIT License)
2. **Pro License** ($299 one-time): External database adapters, advanced features
3. **Enterprise License** ($1,499/year): Priority support, custom builds

### Pricing Strategy

| Edition | Target Market | Pricing Model | Features |
|---------|--------------|---------------|----------|
| Community | Hobbyists, edge computing | Free (MIT) | Internal.db, embedded database.db, setup mode, React UI |
| Professional | SMBs, production | $299 one-time | External DB adapters, API access, priority support |
| Enterprise | Large organizations | $1,499/year | Multiple gateways, clustering, priority support |

---

## 3. Target Audience

### Primary Users

**1. Edge Computing Engineers**
- Deploy IoT gateways on edge devices
- Don't want to edit YAML files
- Need transactional configuration
- Pain point: YAML parsing errors

**2. Air-Gapped Environment Operators**
- Industrial IoT in isolated networks
- Need easy configuration management
- Cannot edit files remotely
- Pain point: Difficult file access

**3. Development Teams**
- Need local testing environment
- Want quick configuration reset
- Need state snapshotting
- Pain point: Difficult configuration reproduction

**4. System Integrators**
- Deploy solutions at customer sites
- Want UI-based configuration, not files
- Need change history
- Pain point: No audit trail

---

## 4. Problem Statement

### Core Problems

**Problem 1: Configuration File Editing**
- v3 uses config.yaml that must be edited
- Users must understand YAML syntax
- Indentation and format errors are hard to debug
- Current approach: YAML editor with validation
- Solution: API/Web UI for all configuration (no file editing)

**Problem 2: No Configuration History**
- No trace of configuration changes
- Difficult to rollback to previous settings
- Current approach: No history
- Solution: History table in internal.db

**Problem 3: Difficult Configuration Migration**
- Version upgrades require manual config.yaml editing
- Not automatic
- Current approach: Manual migration guides
- Solution: Automatic migration with database schema

**Problem 4: Separate Backups**
- config.yaml and database.db must be backed up separately
- Not consistent
- Current approach: Manual backup scripts
- Solution: Single backup for both databases

**Problem 5: Complex Configuration Validation**
- YAML validation can be complex
- Error messages unclear
- Current approach: Basic YAML validation
- Solution: API-level validation with clear error messages

### Impact

- **Configuration Errors**: Users make YAML syntax mistakes
- **Difficult Upgrades**: Users afraid to upgrade due to manual migration
- **No Audit**: Can't see who changed what and when
- **Inconsistent Backups**: Config and database out of sync
- **Poor DX**: Difficult configuration debugging

---

## 5. Proposed Solution

### Vision: Database-First Configuration

**Single executable binary** with:
- **internal.db**: Embedded SQLite for application configuration
- **database.db**: Embedded SQLite for system data (with adapter for other DBs)
- **API/Web UI**: All configuration via API (no file editing)
- **Setup Mode**: Form-based configuration UI
- **Runner Mode**: Read-only configuration
- **History & Rollback**: Version control for configuration
- **Automatic Migration**: Database schema for upgrades

### Key Innovations

**1. Dual Database Architecture**

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
│  │  │  • Application settings                         │ │  │
│  │  │  • MQTT broker configuration                    │ │  │
│  │  │  • Database settings                            │ │  │
│  │  │  • Web server settings                          │ │  │
│  │  │  • Authentication credentials                   │ │  │
│  │  │  • Logging configuration                        │ │  │
│  │  │  • Change history                               │ │  │
│  │  │  • Schema migrations                            │ │  │
│  │  └────────────────────────────────────────────────┘ │  │
│  │                                                        │  │
│  │  ┌────────────────────────────────────────────────┐ │  │
│  │  │      System Database (database.db - Default)    │ │  │
│  │  │      Runtime Data (with Adapter Pattern)        │ │  │
│  │  │                                                  │ │  │
│  │  │  • Device registry                              │ │  │
│  │  │  • Topics and subscriptions                     │ │  │
│  │  │  • MQTT message history                         │ │  │
│  │  │  • Device states                                │ │  │
│  │  │  • Metrics and statistics                       │ │  │
│  │  │  • User accounts (if any)                       │ │  │
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
│  │  │ - Real-time metrics                             │  │  │
│  │  │ - Log viewer                                    │  │  │
│  │  └────────────────────────────────────────────────┘  │  │
│  │                                                        │  │
│  │  ┌────────────────────────────────────────────────┐  │  │
│  │  │ TanStack Router & Query                         │  │  │
│  │  └────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

**2. Internal Database Schema (internal.db)**

```sql
-- Application settings table (key-value store)
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

-- Change history table
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

-- Credentials table (encrypted)
CREATE TABLE credentials (
    id TEXT PRIMARY KEY,
    service TEXT NOT NULL,  -- 'mqtt', 'database', etc.
    username TEXT NOT NULL,
    password_hash TEXT NOT NULL,  -- Hashed
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Schema migrations table
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

**3. API for Configuration**

```go
// API endpoints for configuration

// GET /api/config - Get all configuration
func GetConfig(c *gin.Context) {
    config := configService.GetAll()
    c.JSON(200, config)
}

// GET /api/config/:category - Get configuration by category
func GetConfigByCategory(c *gin.Context) {
    category := c.Param("category")
    config := configService.GetByCategory(category)
    c.JSON(200, config)
}

// PUT /api/config/:key - Update configuration value
func UpdateConfig(c *gin.Context) {
    var req UpdateConfigRequest
    if err := c.BindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // Validate
    if err := configService.Validate(req.Key, req.Value); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // Update with history
    if err := configService.Set(req.Key, req.Value, req.Reason); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"status": "ok"})
}

// POST /api/config/validate - Validate configuration
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

// GET /api/config/history - Change history
func GetConfigHistory(c *gin.Context) {
    history := configService.GetHistory()
    c.JSON(200, history)
}

// POST /api/config/rollback/:version - Rollback to previous version
func RollbackConfig(c *gin.Context) {
    version := c.Param("version")
    if err := configService.Rollback(version); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"status": "rolled_back"})
}

// POST /api/config/reset - Reset to defaults
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
  // Fetch current configuration
  const { data: config, isLoading } = useQuery({
    queryKey: ['config'],
    queryFn: () => fetch('/api/config').then(r => r.json()),
  });

  // Save config mutation
  const saveConfig = useMutation({
    mutationFn: (newConfig: Record<string, any>) =>
      fetch('/api/config', {
        method: 'POST',
        body: JSON.stringify(newConfig),
      }),
    onSuccess: () => {
      alert('Configuration saved! Restarting gateway...');
      setTimeout(() => window.location.reload(), 2000);
    },
  });

  if (isLoading) return <div>Loading...</div>;

  return (
    <div className="setup-mode">
      <h1>Setup Mode - Application Configuration</h1>

      <form onSubmit={(e) => {
        e.preventDefault();
        saveConfig.mutate(config);
      }}>
        {/* MQTT Configuration */}
        <Section title="MQTT Broker">
          <FormField
            label="Broker Type"
            type="select"
            value={config.mqtt.type}
            onChange={(v) => setConfig('mqtt.type', v)}
            options={[
              { value: 'embedded', label: 'Embedded' },
              { value: 'external', label: 'External' },
            ]}
          />
          <FormField
            label="Listen Port"
            type="number"
            value={config.mqtt.listen_port}
            onChange={(v) => setConfig('mqtt.listen_port', v)}
          />
        </Section>

        {/* Database Configuration */}
        <Section title="Database">
          <FormField
            label="Database Type"
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

        <button type="submit">Save & Restart</button>
      </form>
    </div>
  );
}
```

**5. Automatic Migration**

```go
// Migration from v3 to v4
func Migration_v3_to_v4() error {
    // Read old config.yaml if exists
    if oldConfigExists() {
        yamlConfig := readYAMLConfig()

        // Convert to database
        for key, value := range yamlConfig {
            configService.Set(key, value, "Migration from v3")
        }

        // Backup and remove config.yaml
        backupFile("config.yaml")
        removeFile("config.yaml")
    }

    return nil
}
```

---

## 6. Project Objectives

### Primary Objectives (Must Have)

**O1: Implement Internal Database (internal.db)**
- Replace config.yaml with embedded SQLite
- Schema for settings, credentials, history
- CRUD API for configuration
- Automatic migration from v3

**O2: Form-Based Setup Mode**
- Form-based UI for configuration (not YAML editor)
- Real-time validation
- Settings categorization
- Connection testing

**O3: Configuration History**
- Track all changes
- Rollback to previous versions
- Export/import configuration
- Audit log

**O4: Read-Only Runner Mode**
- Configuration cannot be changed
- Dashboard monitoring
- Configuration viewer

**O5: Adapter Pattern for Database**
- database.db as default (SQLite embedded)
- Adapters for PostgreSQL, MongoDB, InfluxDB
- Configuration in internal.db
- Runtime data in database.db (or external)

### Secondary Objectives (Should Have)

**O6: Enhanced Configuration Validation**
- Per-field validation
- Dependency checking
- Conflict detection
- Clear error messages

**O7: Snapshot & Restore**
- Complete configuration snapshot
- Restore from snapshot
- Configuration cloning

**O8: Advanced Settings UI**
- Search and filter settings
- Per-category reset
- Import/export configuration

---

## 7. Scope & Project Constraints

### In Scope (What We Will Build)

#### Core Functionality

1. **Internal Database (internal.db)**
   - Schema for settings, credentials, history
   - CRUD API
   - Automatic migration
   - Backup/restore

2. **System Database (database.db)**
   - Schema for device, topic, message, etc.
   - Adapter pattern for external DB
   - Default SQLite embedded

3. **Setup Mode**
   - Form-based UI
   - Real-time validation
   - Test connection
   - Save & Restart

4. **Runner Mode**
   - Read-only configuration
   - Dashboard monitoring
   - Configuration viewer

5. **Configuration History**
   - Change tracking
   - Rollback
   - Audit log

### Out of Scope (What We Won't Build - Initially)

1. **Advanced Features** (Phase 2)
   - Multi-user management
   - Role-based access control
   - Advanced authentication

2. **High Availability** (Phase 3)
   - Clustering
   - Data replication
   - Automatic failover

### Constraints & Limitations

#### Technical Constraints
- **Frontend**: React 18+, Vite 5+, TypeScript 5+
- **Backend**: Go 1.24+
- **MQTT Library**: mochi-mqtt
- **Database**: SQLite for internal.db and database.db (default)

#### Time Constraints
- **Phase 1 (v4.0)**: 10 months
- **Phase 2 (Advanced Features)**: +3 months

#### Budget Constraints
- **Development Team**: 2-3 full-time developers
- **Infrastructure**: $1,500/month

---

## 8. Technology Stack

### Frontend Stack

| Component | Technology | Version | Rationale |
|-----------|-----------|-------|-----------|
| UI Library | React | 18.3+ | Largest ecosystem |
| Build Tool | Vite | 5.0+ | Fast HMR, optimized builds |
| Language | TypeScript | 5.3+ | Type safety |
| State Management | TanStack Query | 5.0+ | Server state management |
| Routing | TanStack Router | 1.0+ | Type-safe routing |
| Forms | React Hook Form | Latest | Good form handling |

### Backend Stack

| Component | Technology | Version | Rationale |
|-----------|-----------|-------|-----------|
| Language | Go | 1.24+ | Performance, single binary |
| MQTT Broker | mochi-mqtt | 2.0+ | Embedded Go MQTT broker |
| Internal DB | SQLite | 3.40+ | Embedded, single file |
| System DB | SQLite | 3.40+ | Default, with adapters |
| Cache | go-cache | Latest | In-memory caching |
| HTTP Router | Chi | 5.0+ | Lightweight, idiomatic |

---

## 9. Development Roadmap

### Phase 1: Foundation (Months 1-4)

**Sprint 1-2: Project Setup & Architecture**
- [ ] Monorepo structure
- [ ] Setup React + Vite frontend
- [ ] Design internal.db schema
- [ ] Design database.db schema

**Sprint 3-4: Internal Database**
- [ ] Implement internal.db
- [ ] Settings CRUD API
- [ ] Configuration history
- [ ] Migration from v3

**Sprint 5-6: System Database**
- [ ] Implement database.db
- [ ] Adapter pattern
- [ ] Schema for device, topic, message

**Sprint 7-8: Frontend**
- [ ] Form-based setup mode
- [ ] Runner mode dashboard
- [ ] TanStack Query integration

### Phase 2: Features (Months 5-8)

**Sprint 9-12: Validation & Testing**
- [ ] Configuration validation
- [ ] Integration testing
- [ ] E2E testing
- [ ] Load testing

**Sprint 13-16: Polish**
- [ ] UI improvements
- [ ] Error handling
- [ ] Performance optimization
- [ ] Documentation

### Phase 3: Production (Months 9-10)

**Sprint 17-20: Deploy**
- [ ] Production build
- [ ] Cross-platform compilation
- [ ] Installation guide
- [ ] Release notes

---

## 10. Success Metrics

### Technical Metrics

| Metric | Target | Measurement |
|--------|--------|-------------|
| Binary Size | < 60 MB | Build artifacts |
| Startup Time | < 2 seconds | Manual testing |
| Memory Usage | < 100 MB idle | Monitoring |
| Setup Completion Rate | > 95% | Analytics |

---

## 11. Resource Requirements

### Budget Estimate

| Category | Cost (Monthly) | Duration |
|----------|----------------|---------|
| Development Team Salaries | $20,000 | 10 months |
| Infrastructure | $1,500 | 10 months |
| Tools & Services | $400 | Ongoing |
| Contingency (15%) | $3,150 | - |
| **Total** | **$22,500/month × 10** | **$225,000** |

---

## 12. Comparison: v3 vs v4

| Aspect | v3 (config.yaml) | v4 (internal.db) |
|--------|----------------|----------------|
| **Configuration** | YAML file | SQLite database |
| **Editing** | Manual file edit | API/Web UI |
| **History** | None | Yes, complete |
| **Migration** | Manual | Automatic |
| **Backup** | Separate | Single backup |
| **Validation** | YAML parsing | API validation |
| **Rollback** | Manual | Automatic |
| **Budget** | $215,500 | $225,000 |

---

## 13. Conclusion

Version v4 represents **a modern, database-first approach**:

### Key Advantages over v3

**Configuration Management**
- ✅ No file editing
- ✅ API/Web UI for all configuration
- ✅ Complete history
- ✅ Automatic rollback
- ✅ Automatic migration

**Reliability**
- ✅ Transactional (ACID)
- ✅ Consistent
- ✅ Unified backup
- ✅ API-level validation

**User Experience**
- ✅ Form-based UI (not YAML editor)
- ✅ Real-time validation
- ✅ Clear error messages
- ✅ Easy to use

With focused execution over 10 months, we can deliver **production-ready v4.0** with modern and reliable configuration management.

---

**Document Version**: 4.0
**Last Updated**: January 9, 2026
**Status**: Draft - Pending Review
**Replaces**: PROJECT_BRIEF_EN_v3.md

---

**Prepared by**: Development Team
**Approved by**: [Pending]
**Review Date**: [Pending]
