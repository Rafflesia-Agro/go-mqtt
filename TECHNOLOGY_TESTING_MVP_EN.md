# Technology Testing Plan: UM-Gateway MVP v0.1

## Executive Summary

**Project Name**: UM-Gateway MVP Technology Testing
**Purpose**: Verify all technologies/libraries before full development
**Timeline**: 2 weeks (before main development starts)
**Target**: Complete all proof-of-concept (PoC) tests

---

## Version History

| Version | Date | Changes | Author |
|---------|------|---------|--------|
| **v1.0** | 2026-01-09 | Initial testing plan<br>• mochi-mqtt broker PoC<br>• SQLite WAL mode testing<br>• go-cache integration<br>• React + Vite in Go<br>• YAML configuration<br>• All critical technology verification | Development Team |

---

## 1. Purpose

### Why Technology Testing?

Before starting the 3-month MVP development, we need to **verify that all chosen technologies work together**. This prevents:
- ❌ Incompatible libraries
- ❌ Unexpected technical blockers
- ❌ Performance issues
- ❌ Architecture redesign mid-development

### Testing Goals

1. **Verify Feasibility**: Prove each technology can do what we need
2. **Test Integration**: Ensure technologies work together
3. **Measure Performance**: Establish baseline metrics
4. **Identify Issues**: Discover problems early (when they're cheap to fix)
5. **Build Confidence**: Know the tech stack works before committing

### Success Criteria

- [ ] All PoCs run successfully
- [ ] All technologies integrate without conflicts
- [ ] Performance meets minimum requirements
- [ ] No critical blockers discovered
- [ ] Documentation for each PoC created

---

## 2. Testing Overview

### Testing Timeline

**Week 1: Core Backend Technologies**
- Days 1-2: mochi-mqtt embedded broker
- Days 3-4: SQLite WAL mode
- Days 5-7: go-cache integration

**Week 2: Frontend & Integration**
- Days 1-2: React + Vite setup
- Days 3-4: Embed React in Go binary
- Days 5-6: YAML configuration
- Day 7: Integration testing & documentation

### Test Environment

**Hardware**:
- Development machine (any OS)
- Minimum: 2 CPU cores, 4GB RAM

**Software**:
- Go 1.24+
- Node.js 20+
- Git
- Modern browser (Chrome/Firefox)

**Tools**:
- IDE (VS Code, GoLand, etc.)
- Postman/cURL (for API testing)
- MQTT client (MQTTX, mosquitto_pub/sub)

---

## 3. Test Plan Details

### Test 1: mochi-mqtt Embedded Broker

**Purpose**: Verify mochi-mqtt can work as embedded broker

**What to Test**:
- [ ] Create basic MQTT broker
- [ ] TCP listener on port 1883
- [ ] Client authentication (username/password)
- [ ] Publish/Subscribe functionality
- [ ] Message persistence (optional)
- [ ] Concurrent connections (10+ clients)
- [ ] Graceful shutdown

**Success Criteria**:
- ✅ Broker starts and listens on port 1883
- ✅ MQTT client can connect with username/password
- ✅ Client can publish to topic
- ✅ Client can subscribe to topic
- ✅ Subscriber receives published messages
- ✅ 10+ concurrent clients work without issues
- ✅ Broker shuts down gracefully

**Expected Code Structure**:

```go
package main

import (
    "log"
    mqtt "github.com/mochi-mqtt/mqtt/v2"
    "github.com/mochi-mqtt/mqtt/hooks/auth"
    "github.com/mochi-mqtt/mqtt/listeners"
)

func main() {
    // Create MQTT server
    server := mqtt.NewServer(nil)

    // Add authentication hook
    authHook := auth.NewHook("username", "password")
    server.AddHook(authHook)

    // Create TCP listener
    tcp := listeners.NewTCPListener("localhost:1883", nil)
    server.AddListener(tcp)

    // Start server
    log.Println("Starting MQTT broker on :1883")
    go server.Serve()

    // Keep running
    select {}
}
```

**Test Procedure**:
1. Create new directory: `tests/mochi-mqtt-poc`
2. Initialize Go module
3. Install mochi-mqtt
4. Copy expected code
5. Run: `go run main.go`
6. Use MQTTX to connect to `localhost:1883`
7. Publish message to topic `test/topic`
8. Subscribe to `test/topic`
9. Verify message received
10. Test with 10+ concurrent connections
11. Stop server (Ctrl+C) and verify graceful shutdown

**Estimated Time**: 4-6 hours

**Deliverables**:
- Working PoC code
- Test results documentation
- Performance notes (memory usage, CPU)
- Known issues or limitations

---

### Test 2: SQLite WAL Mode

**Purpose**: Verify SQLite with WAL mode works for concurrent access

**What to Test**:
- [ ] Create SQLite database
- [ ] Enable WAL mode
- [ ] Concurrent reads
- [ ] Concurrent writes
- [ ] Connection pooling
- [ ] Schema migrations
- [ ] Backup/restore

**Success Criteria**:
- ✅ Database created successfully
- ✅ WAL mode enabled
- ✅ Multiple concurrent reads work
- ✅ Concurrent writes don't block reads
- ✅ Connection pool works
- ✅ Migrations run without errors
- ✅ Backup can be created

**Expected Code Structure**:

```go
package main

import (
    "database/sql"
    "log"
    "sync"
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    // Open database
    db, err := sql.Open("sqlite3", "./test.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Enable WAL mode
    _, err = db.Exec("PRAGMA journal_mode=WAL")
    if err != nil {
        log.Fatal(err)
    }

    // Create table
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS messages (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            topic TEXT,
            payload TEXT,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    `)
    if err != nil {
        log.Fatal(err)
    }

    // Test concurrent writes
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(n int) {
            defer wg.Done()
            _, err := db.Exec(
                "INSERT INTO messages (topic, payload) VALUES (?, ?)",
                "test/topic",
                "message-"+string(rune(n)),
            )
            if err != nil {
                log.Printf("Write error: %v", err)
            }
        }(i)
    }
    wg.Wait()

    // Verify data
    rows, _ := db.Query("SELECT COUNT(*) FROM messages")
    defer rows.Close()
    var count int
    rows.Next()
    rows.Scan(&count)
    log.Printf("Total messages: %d", count)
}
```

**Test Procedure**:
1. Create directory: `tests/sqlite-wal-poc`
2. Initialize Go module
3. Install sqlite3 driver
4. Copy expected code
5. Run: `go run main.go`
6. Verify WAL mode enabled (check `.db-wal` and `.db-shm` files)
7. Test concurrent reads (multiple goroutines reading)
8. Test concurrent writes (multiple goroutines writing)
9. Verify connection pooling works
10. Test backup creation
11. Document performance

**Estimated Time**: 3-4 hours

**Deliverables**:
- Working PoC code
- WAL mode verification results
- Concurrency test results
- Performance notes
- Backup/restore procedure

---

### Test 3: go-cache Integration

**Purpose**: Verify go-cache works for in-memory caching

**What to Test**:
- [ ] Create cache instance
- [ ] Set values with TTL
- [ ] Get values
- [ ] Automatic expiration
- [ ] Delete values
- [ ] Flush all
- [ ] Thread safety

**Success Criteria**:
- ✅ Cache created successfully
- ✅ Values can be set with TTL
- ✅ Values can be retrieved
- ✅ Items expire after TTL
- ✅ Values can be deleted
- ✅ Cache can be flushed
- ✅ Concurrent access is safe

**Expected Code Structure**:

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/patrickmn/go-cache"
)

func main() {
    // Create cache with default TTL of 5 minutes
    // Cleanup interval of 10 minutes
    c := cache.New(5*time.Minute, 10*time.Minute)

    // Set value
    c.Set("foo", "bar", cache.DefaultExpiration)

    // Get value
    if x, found := c.Get("foo"); found {
        fmt.Println("Found foo:", x.(string))
    }

    // Set with short TTL for testing
    c.Set("temp", "expires soon", 2*time.Second)

    // Wait for expiration
    time.Sleep(3 * time.Second)

    // Check if expired
    if _, found := c.Get("temp"); !found {
        fmt.Println("temp expired as expected")
    }

    // Test concurrent access
    done := make(chan bool)
    for i := 0; i < 10; i++ {
        go func(n int) {
            key := fmt.Sprintf("key%d", n)
            c.Set(key, n, cache.DefaultExpiration)
            if x, found := c.Get(key); found {
                log.Printf("Goroutine %d: value = %v", n, x)
            }
            done <- true
        }(i)
    }

    // Wait for all goroutines
    for i := 0; i < 10; i++ {
        <-done
    }

    // Item count
    fmt.Printf("Cache items: %d\n", c.ItemCount())
}
```

**Test Procedure**:
1. Create directory: `tests/go-cache-poc`
2. Initialize Go module
3. Install go-cache
4. Copy expected code
5. Run: `go run main.go`
6. Verify basic set/get works
7. Verify TTL expiration works
8. Test concurrent access (100+ goroutines)
9. Verify thread safety
10. Document memory usage

**Estimated Time**: 2-3 hours

**Deliverables**:
- Working PoC code
- TTL verification results
- Concurrency test results
- Memory usage notes

---

### Test 4: React + Vite Frontend

**Purpose**: Verify React + Vite works for admin UI

**What to Test**:
- [ ] Create React + Vite project
- [ ] TypeScript configuration
- [ ] Basic routing
- [ ] API calls (fetch)
- [ ] Build for production
- [ ] TanStack Query integration
- [ ] Responsive design

**Success Criteria**:
- ✅ Project created successfully
- ✅ TypeScript works
- ✅ Basic routing works
- ✅ Can call API endpoints
- ✅ Production build succeeds
- ✅ TanStack Query fetches data
- ✅ UI is responsive

**Expected Code Structure**:

```bash
# Create project
npm create vite@latest mqtt-gateway-poc -- --template react-ts
cd mqtt-gateway-poc

# Install dependencies
npm install

# Install additional packages
npm install @tanstack/react-query
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p
```

**src/App.tsx**:
```typescript
import { useQuery } from '@tanstack/react-query';

function App() {
  const { data, isLoading, error } = useQuery({
    queryKey: ['test'],
    queryFn: async () => {
      const response = await fetch('http://localhost:8080/api/test');
      if (!response.ok) throw new Error('Network error');
      return response.json();
    },
  });

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div className="p-4">
      <h1>MQTT Gateway PoC</h1>
      <pre>{JSON.stringify(data, null, 2)}</pre>
    </div>
  );
}

export default App;
```

**Test Procedure**:
1. Create directory: `tests/react-vite-poc`
2. Create Vite + React project
3. Install dependencies
4. Create basic components
5. Add TanStack Query
6. Create mock API endpoint (use json-server or simple Go server)
7. Test API calls
8. Test TypeScript compilation
9. Build for production: `npm run build`
10. Preview build: `npm run preview`
11. Test responsiveness (browser devTools)

**Estimated Time**: 4-5 hours

**Deliverables**:
- Working PoC code
- Build artifacts
- Performance notes (bundle size, load time)
- TypeScript configuration

---

### Test 5: Embed React in Go Binary

**Purpose**: Verify React app can be embedded in Go binary

**What to Test**:
- [ ] Build React app for production
- [ ] Embed static files in Go
- [ ] Serve React app from Go
- [ ] API proxying
- [ ] SPA routing support
- [ ] Hot reload in development

**Success Criteria**:
- ✅ React app builds successfully
- ✅ Static files embedded in Go binary
- ✅ Go serves React app
- ✅ API calls work
- ✅ SPA routing works
- ✅ Single binary executable

**Expected Code Structure**:

**Go Backend**:
```go
package main

import (
    "embed"
    "io/fs"
    "log"
    "net/http"
)

//go:embed frontend/dist
var frontendFS embed.FS

func main() {
    // Get subdirectory for frontend
    distFS, err := fs.Sub(frontendFS, "frontend/dist")
    if err != nil {
        log.Fatal(err)
    }

    // Serve frontend
    http.Handle("/", http.FileServer(http.FS(distFS)))

    // API endpoint
    http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(`{"status":"ok"}`))
    })

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

**Test Procedure**:
1. Create directory: `tests/embed-react-poc`
2. Create React app (use previous PoC)
3. Build React app: `npm run build`
4. Create Go backend with embed
5. Place React build in `frontend/dist`
6. Run Go server: `go run main.go`
7. Open browser to `http://localhost:8080`
8. Verify React app loads
9. Verify API calls work
10. Test SPA routing (navigate between pages)
11. Build Go binary: `go build -o gateway.exe`
12. Test standalone binary
13. Measure binary size

**Estimated Time**: 4-6 hours

**Deliverables**:
- Working PoC code
- Standalone binary
- Binary size measurement
- Embed procedure documentation
- Development workflow (hot reload)

---

### Test 6: YAML Configuration

**Purpose**: Verify YAML configuration loading works

**What to Test**:
- [ ] Create YAML config file
- [ ] Parse YAML in Go
- [ ] Validate configuration
- [ ] Default values
- [ ] Environment variables
- [ ] Hot reload (optional)
- [ ] Error handling

**Success Criteria**:
- ✅ YAML file created
- ✅ YAML parsed successfully
- ✅ Validation works
- ✅ Defaults applied correctly
- ✅ Environment variables override config
- ✅ Errors handled gracefully

**Expected Code Structure**:

**config.yaml**:
```yaml
mqtt:
  embedded:
    enabled: true
    listen_address: :1883
    username: admin
    password: secret
    max_connections: 100

database:
  type: sqlite
  path: ./data/gateway.db
  wal_mode: true

web:
  address: :8080
```

**config.go**:
```go
package main

import (
    "fmt"
    "log"
    "os"

    "gopkg.in/yaml.v3"
)

type Config struct {
    MQTT struct {
        Embedded struct {
            Enabled         bool   `yaml:"enabled"`
            ListenAddress   string `yaml:"listen_address"`
            Username        string `yaml:"username"`
            Password        string `yaml:"password"`
            MaxConnections  int    `yaml:"max_connections"`
        } `yaml:"embedded"`
    } `yaml:"mqtt"`

    Database struct {
        Type    string `yaml:"type"`
        Path    string `yaml:"path"`
        WALMode bool   `yaml:"wal_mode"`
    } `yaml:"database"`

    Web struct {
        Address string `yaml:"address"`
    } `yaml:"web"`
}

func LoadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    // Set defaults
    if config.MQTT.Embedded.ListenAddress == "" {
        config.MQTT.Embedded.ListenAddress = ":1883"
    }
    if config.Web.Address == "" {
        config.Web.Address = ":8080"
    }

    return &config, nil
}

func main() {
    config, err := LoadConfig("config.yaml")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("MQTT Listen: %s\n", config.MQTT.Embedded.ListenAddress)
    fmt.Printf("Web Address: %s\n", config.Web.Address)
    fmt.Printf("Database: %s\n", config.Database.Path)
}
```

**Test Procedure**:
1. Create directory: `tests/yaml-config-poc`
2. Create sample config.yaml
3. Create Go config loader
4. Test parsing valid YAML
5. Test parsing invalid YAML (verify error handling)
6. Test default values
7. Test environment variable override
8. Test missing config file
9. Test validation rules
10. Document configuration schema

**Estimated Time**: 2-3 hours

**Deliverables**:
- Working PoC code
- Configuration schema
- Validation rules
- Error handling examples

---

### Test 7: Integration Test (All Together)

**Purpose**: Verify all technologies work together

**What to Test**:
- [ ] Start embedded MQTT broker
- [ ] Start HTTP server with embedded React
- [ ] Connect MQTT client
- [ ] Publish/Subscribe messages
- [ ] Store messages in SQLite
- [ ] Cache recent messages
- [ ] Display in React UI
- [ ] Update config via YAML

**Success Criteria**:
- ✅ All services start without errors
- ✅ MQTT broker accepts connections
- ✅ Web UI loads successfully
- ✅ Messages flow through system
- ✅ Database stores messages
- ✅ Cache improves performance
- ✅ UI updates in real-time
- ✅ Config changes work

**Expected Architecture**:

```
┌─────────────────────────────────────────────┐
│         Go Process (Single Executable)       │
│                                             │
│  ┌──────────────┐  ┌──────────────┐       │
│  │ mochi-mqtt   │  │ HTTP Server  │       │
│  │ Broker       │  │              │       │
│  │              │  │ Serves:      │       │
│  │ Port: 1883   │  │ - React UI   │       │
│  └──────────────┘  │ - API        │       │
│        ↓           │              │       │
│  ┌──────────────┐  │              │       │
│  │ Hooks        │  │              │       │
│  │ (on message) │  │              │       │
│  └──────────────┘  │              │       │
│        ↓           └──────────────┘       │
│  ┌──────────────┐                          │
│  │ SQLite       │                          │
│  │ (WAL mode)   │                          │
│  └──────────────┘                          │
│        ↓                                   │
│  ┌──────────────┐                          │
│  │ go-cache     │                          │
│  │ (in-memory)  │                          │
│  └──────────────┘                          │
└─────────────────────────────────────────────┘
```

**Test Procedure**:
1. Create directory: `tests/integration-poc`
2. Copy all previous PoCs
3. Integrate into single application
4. Start all services
5. Connect with MQTT client
6. Publish test messages
7. Verify database storage
8. Verify cache hits
9. Open web UI
10. Verify real-time updates
11. Test config reload
12. Measure overall performance
13. Test graceful shutdown

**Estimated Time**: 8-12 hours

**Deliverables**:
- Complete integration PoC
- Performance metrics
- Known issues
- Recommendations
- Final architecture decision

---

## 4. Test Organization

### Directory Structure

```
tests/
├── README.md                           (this file overview)
├── mochi-mqtt-poc/                     (Test 1)
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   └── RESULTS.md
├── sqlite-wal-poc/                     (Test 2)
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   └── RESULTS.md
├── go-cache-poc/                       (Test 3)
│   ├── main.go
│   ├── go.mod
│   ├── go.sum
│   └── RESULTS.md
├── react-vite-poc/                     (Test 4)
│   ├── src/
│   ├── package.json
│   ├── vite.config.ts
│   └── RESULTS.md
├── embed-react-poc/                    (Test 5)
│   ├── main.go
│   ├── frontend/
│   │   └── (React build output)
│   └── RESULTS.md
├── yaml-config-poc/                    (Test 6)
│   ├── main.go
│   ├── config.yaml
│   └── RESULTS.md
└── integration-poc/                    (Test 7)
    ├── main.go
    ├── embedded/
    │   └── (React build)
    └── RESULTS.md
```

### Documentation Template

Each test should have a `RESULTS.md` file with:

```markdown
# Test Results: [Test Name]

## Summary
- **Date**: [Date]
- **Tester**: [Name]
- **Status**: ✅ PASS / ❌ FAIL

## What Was Tested
[List test objectives]

## Results
### Success Criteria
- [ ] Criterion 1 - ✅ PASS
- [ ] Criterion 2 - ✅ PASS
- [ ] Criterion 3 - ❌ FAIL

### Performance
- Metric 1: [Value]
- Metric 2: [Value]

### Issues Found
1. [Issue description]
2. [Issue description]

### Recommendations
1. [Recommendation]
2. [Recommendation]

## Conclusion
[Overall assessment - go/no-go decision]

## Next Steps
[What to do based on results]
```

---

## 5. Success Criteria (Overall)

### Go/No-Go Decision

**GO (Start Development) if**:
- ✅ All 7 tests pass
- ✅ No critical blockers
- ✅ Performance meets requirements
- ✅ Integration works
- ✅ Team is confident

**NO-GO (Re-evaluate) if**:
- ❌ Any test fails critically
- ❌ Performance is unacceptable
- ❌ Integration has major issues
- ❌ Better alternative found

### Minimum Requirements

**Functional Requirements**:
- All core features work as expected
- No show-stopper bugs
- Integration is stable

**Performance Requirements**:
- Binary size < 50 MB
- Startup time < 2 seconds
- Memory usage < 100 MB (idle)
- MQTT throughput > 1,000 msg/sec
- Web page load < 1 second

**Developer Experience Requirements**:
- Code is clear and maintainable
- Libraries have good documentation
- Debugging is straightforward
- Hot reload works in development

---

## 6. Timeline & Resources

### Timeline

**Week 1**: Backend Technologies (Tests 1-3)
**Week 2**: Frontend & Integration (Tests 4-7)

### Resources

**People**:
- 1 Full-Stack Developer (Go + React)
- 10-15 hours per week
- Total: 20-30 hours

**Tools**:
- Development machine
- Go 1.24+
- Node.js 20+
- Git repository
- MQTT client (MQTTX)

**Budget**:
- Developer time: ~$1,500
- Tools: $0 (all free/open source)
- **Total**: ~$1,500

---

## 7. Risk Management

### Potential Risks

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| mochi-mqtt doesn't support needed features | HIGH | LOW | Have alternative: Eclipse Paho |
| SQLite WAL performance issues | MEDIUM | LOW | Can tune PRAGMA settings |
| Embed React significantly increases binary size | LOW | MEDIUM | Can optimize build, use compression |
| Integration issues between components | HIGH | MEDIUM | Integration test will reveal this |
| Poor documentation | MEDIUM | LOW | Create our own docs |

### Contingency Plans

**If mochi-mqtt fails**:
- Alternative: Use Eclipse Paho (client mode only, no embedded broker)
- Impact: No embedded broker, must use external Mosquitto
- Timeline: +1 week to implement alternative

**If SQLite WAL fails**:
- Alternative: Use standard SQLite (no WAL)
- Impact: Lower concurrency
- Timeline: Minimal change

**If embed React fails**:
- Alternative: Serve React separately (not embedded)
- Impact: Not single binary
- Timeline: +1 week to adjust architecture

---

## 8. Deliverables

### For Each Test

1. **Source Code**
   - Complete, working PoC
   - Well-commented
   - Follows best practices

2. **Results Documentation**
   - RESULTS.md file
   - Screenshots (if applicable)
   - Performance metrics
   - Issues and recommendations

3. **Setup Instructions**
   - How to run the PoC
   - Dependencies
   - Environment requirements

### Final Deliverables

1. **Complete Test Suite**
   - All 7 PoCs
   - Integration test
   - Documentation

2. **Summary Report**
   - Executive summary
   - Go/No-Go recommendation
   - Lessons learned
   - Next steps

3. **Architecture Decision**
   - Final tech stack confirmed
   - Any changes from original plan
   - Justification for decisions

---

## 9. Next Steps

### After Testing Complete

**If GO**:
1. Create main project repository
2. Set up development environment
3. Begin MVP development (3-month sprint)
4. Reference PoCs for implementation

**If NO-GO**:
1. Identify specific blockers
2. Research alternatives
3. Re-test problematic components
4. Make final architecture decision

---

## 10. Conclusion

This testing plan ensures we **validate our technology choices before committing to 3 months of development**. By spending 2 weeks and ~$1,500 now, we reduce the risk of costly rework or architecture changes later.

**Key Benefits**:
- ✅ Confidence in tech stack
- ✅ Early issue discovery
- ✅ Performance baseline
- ✅ Reusable code for MVP
- ✅ Team learning

**Success Means**:
- All technologies verified
- Integration working
- No blockers
- Ready to start MVP development

With proper testing, we can begin MVP development with confidence, knowing our chosen technologies will work as expected.

---

**Document Version**: 1.0
**Last Updated**: January 9, 2026
**Status**: Ready for Testing
**Testing Period**: 2 weeks
**Owner**: Development Team

---

## Appendix

### A. Quick Reference

**Repository**: [Link to test repository]
**Wiki**: [Link to internal wiki]
**Slack**: #mqtt-gateway-testing

### B. Resources

- [mochi-mqtt Documentation](https://github.com/mochi-mqtt/mqtt)
- [SQLite WAL Mode](https://www.sqlite.org/wal.html)
- [go-cache](https://github.com/patrickmn/go-cache)
- [React + Vite](https://vitejs.dev/guide/)
- [TanStack Query](https://tanstack.com/query/latest)

### C. Contact

**Questions**: Contact [Tech Lead]
**Issues**: Create GitHub issue
**Emergencies**: Slack #dev-urgent

---

**Prepared by**: Development Team
**Approved by**: [Pending]
**Review Date**: [Pending]
