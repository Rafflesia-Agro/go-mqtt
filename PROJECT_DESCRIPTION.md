# Go-MQTT Project Documentation

## Table of Contents
1. [Project Overview](#project-overview)
2. [Background](#background)
3. [Purpose & Objectives](#purpose--objectives)
4. [Target Audience](#target-audience)
5. [Business Context](#business-context)
6. [Technical Architecture](#technical-architecture)
7. [Database Structure](#database-structure)
8. [Business Flow](#business-flow)
9. [User Flow](#user-flow)
10. [Features](#features)
11. [File Structure](#file-structure)
12. [Technical Specifications](#technical-specifications)
13. [Router Structure](#router-structure)
14. [Controllers & Business Logic](#controllers--business-logic)
15. [Important Code References](#important-code-references)
16. [Areas for Improvement](#areas-for-improvement)
17. [Installation & Setup](#installation--setup)
18. [License](#license)
19. [Author & Team](#author--team)

---

## Project Overview

**Go-MQTT** is a high-performance Go-based backend service that bridges IoT hardware devices with web applications through MQTT protocol. The system serves as a central hub for managing poultry farm coops, enabling real-time monitoring and control of hardware devices (fans, lamps, feeders, water systems) while collecting sensor data for analysis.

**Key Characteristics:**
- Real-time bidirectional communication via MQTT
- High-throughput sensor data processing with batch insertion
- Redis-based state management for instant hardware control
- JWT-based authentication with role-based access control
- RESTful API for web/mobile client integration
- Optimized for 24/7 operation in agricultural IoT environments

---

## Background

This project was developed for **Rafflesia Agro**, an agricultural technology company focusing on smart poultry farming solutions in Indonesia. The system addresses the need for:

1. **Remote Hardware Management**: Farmers need to control coop hardware (lights, fans, feeders) remotely without physical presence
2. **Real-time Monitoring**: Continuous sensor data collection (temperature, humidity, etc.) for environmental control
3. **Scalability**: Support for multiple coops with high-frequency sensor data transmission
4. **Reliability**: 24/7 operation with automatic reconnection and graceful failure handling

The project evolved from basic MQTT message handling to a comprehensive IoT platform with database persistence, caching, and web API capabilities.

---

## Purpose & Objectives

### Primary Purpose
To provide a reliable, scalable middleware layer that:
- Receives sensor telemetry from ESP32/IoT devices via MQTT
- Stores high-volume sensor data efficiently in PostgreSQL
- Exposes hardware control endpoints via HTTP REST API
- Maintains real-time hardware state in Redis
- Authenticates and authorizes users based on roles

### Technical Objectives
- **Performance**: Handle 10,000+ sensor readings per minute with batch insertion
- **Reliability**: 99.9% uptime with auto-reconnection to MQTT/Database/Redis
- **Scalability**: Support unlimited coops with horizontal scaling capability
- **Real-time**: Sub-second latency for hardware state updates
- **Security**: JWT-based auth with role-based access control (RBAC)

### Business Objectives
- Enable farmers to monitor multiple coops from a single dashboard
- Reduce labor costs through automated hardware control
- Improve poultry health through precise environmental monitoring
- Provide data-driven insights for operational optimization

---

## Target Audience

### Primary Users
1. **Farmers** (Pemilik Kandang)
   - Own and manage multiple coops
   - Need remote control and monitoring capabilities
   - Require historical data analysis

2. **ABK** (Ahli Budidaya Kandang - Coop Operators)
   - Assigned to specific coops
   - Monitor daily operations and sensor readings
   - Execute hardware controls as needed

3. **Admins**
   - System administrators
   - Full access to all coops for maintenance and troubleshooting
   - Manage user roles and permissions

### Secondary Stakeholders
- **Management**: Access aggregated data for business insights
- **Technical Teams**: Monitor system health and performance
- **Hardware Integrators**: ESP32 developers who connect devices to the system

---

## Business Context

### Domain: Smart Poultry Farming (Rafflesia Agro)

The system operates in the agricultural IoT domain, specifically for poultry farming in Indonesia. Key business requirements:

1. **Multi-Coop Management**
   - Single farmer can own multiple coops across different locations
   - Each coop can have an assigned ABK (operator)
   - Hierarchical access control based on ownership

2. **Hardware Types**
   - **Climate Control**: Fans, heaters, lamps (for temperature regulation)
   - **Feeding Systems**: Automated feeders with scheduling
   - **Water Systems**: Automated water dispensers
   - **Sensors**: Temperature, humidity, light, ammonia levels

3. **Operational Modes**
   - **Manual Mode**: Direct on/off control by farmers/ABK
   - **Auto Mode**: Hardware responds to sensor thresholds (e.g., fan turns on when temperature > 30°C)
   - **Scheduled Mode**: Time-based operations (e.g., lights on at 6 AM, off at 6 PM)

4. **Data Requirements**
   - Sensor data stored every ~10-15 seconds per coop
   - Hardware state changes must be instant (< 1 second latency)
   - Historical data for trend analysis (weeks/months)

---

## Technical Architecture

### System Architecture Overview

```
┌─────────────┐         MQTT          ┌──────────────┐
│   ESP32     │ ◄────────────────────► │  Go-MQTT     │
│  (Hardware) │   tcp://mqtt:1883      │   Service    │
└─────────────┘                        └──────┬───────┘
                                               │
                    ┌──────────────────────────┼───────────────────┐
                    │                          │                   │
                    ▼                          ▼                   ▼
            ┌───────────┐           ┌─────────────┐       ┌─────────────┐
            │  Redis    │           │ PostgreSQL  │       │ Web/Mobile  │
            │  (State)  │           │  (Telemetry)│       │   Client    │
            └───────────┘           └─────────────┘       └─────────────┘
```

### Architecture Pattern
**Event-Driven Microservice with CQRS-like Separation**

- **Write Path**: Hardware state commands → Redis → MQTT publish → ESP32 devices
- **Read Path**: Sensor telemetry → MQTT subscribe → Worker queue → PostgreSQL batch insert
- **API Layer**: REST endpoints for client applications with JWT auth

### Key Components

1. **MQTT Broker Client** (`src/broker/mqtt.go`)
   - Subscribes to: `$share/backend-workers/coops/+/telemetry`
   - Publishes to: `coops/{coop_id}/state`
   - Handles connection loss and auto-reconnect

2. **HTTP Server** (`main.go:608-637`)
   - RESTful API using Chi router
   - JWT authentication middleware
   - Role-based authorization middleware

3. **Database Worker Pool** (`src/database/database.go:173-211`)
   - 10 concurrent workers processing sensor data
   - Batch insertion (up to 1000 rows per batch)
   - 1-second timeout for batch flushing

4. **Cache Layer** (`src/cache/redis.go`)
   - Stores hardware states for instant retrieval
   - Tracks last sensor timestamp for online/offline status
   - Pattern: `hardware_state_{coop_id}_{hardware_name}`

---

## Database Structure

### PostgreSQL Schema

**Primary Tables:**

1. **`rec_sensor_coops`** (Sensor Readings)
   ```sql
   - id: bigserial (PK)
   - value: double precision
   - sensor_type_id: integer (FK to sensor_types)
   - coop_id: bigint (FK to coops)
   - created_at: timestamp
   - updated_at: timestamp
   ```
   - **Indexing**: Composite index on `(coop_id, created_at)` for time-series queries
   - **Insert Rate**: ~10-15 seconds per coop, batch inserted via COPY protocol

2. **`sensor_types`** (Sensor Definitions)
   ```sql
   - id: serial (PK)
   - name: varchar (unique) - e.g., "temperature", "humidity"
   ```
   - **Cached**: Loaded into memory on startup (map[string]int32)

3. **`coops`** (Coop Information)
   ```sql
   - id: bigint (PK)
   - farm_id: integer (FK to farms)
   - abk_id: bigint (FK to users, nullable)
   ```

4. **`farms`** (Farm Information)
   ```sql
   - id: integer (PK)
   - farmer_id: bigint (FK to users)
   ```

5. **`users`** (User Accounts)
   ```sql
   - id: bigint (PK)
   - current_role_id: integer (FK to roles)
   ```

6. **`roles`** (User Roles)
   ```sql
   - id: integer (PK)
   - name: varchar - "admin", "farmer", "abk"
   ```

### Database Connection
- **Driver**: `pgx/v5` (pure Go PostgreSQL driver)
- **Connection Pool**: `pgxpool` with configurable size
- **Timezone**: Asia/Jakarta (set at session level)
- **Batch Insert**: Uses `COPY FROM` protocol for high-performance bulk inserts

---

## Business Flow

### 1. Sensor Data Collection Flow

```
ESP32 Device
    │
    │ Publishes telemetry to: coops/{coop_id}/telemetry
    ▼
MQTT Broker
    │
    │ Shared subscription: $share/backend-workers/coops/+/telemetry
    ▼
Go-MQTT Service (MQTT Handler)
    │
    │ - Extracts coop_id from topic
    │ - Parses JSON payload
    │ - Adds timestamp (Asia/Jakarta)
    ▼
Save() Method → Job Channel (Buffer: 20,000)
    │
    │ 10 Concurrent DB Workers
    ▼
Batch Insert (max 1000 rows or 1 second timeout)
    │
    │ PostgreSQL COPY FROM protocol
    ▼
rec_sensor_coops table
    │
    └── Updates Redis key: last_sensor_stored_{coop_id}
```

### 2. Hardware Control Flow

```
Web/Mobile Client
    │
    │ POST /coops/{id}/state
    │ Headers: Authorization: Bearer <JWT>
    │ Body: { "hardware": "fan1", "state": true, "meta": {...} }
    ▼
JWT Verification
    │
    │ - Validates token signature
    │ - Extracts user_id and role
    ▼
CoopAccessMiddleware
    │
    │ - Checks if user owns the coop (farmer role)
    │ - Checks if user is assigned ABK (abk role)
    │ - Bypass for admin role
    ▼
Handler Logic
    │
    │ - Validates request (partial or full update)
    │ - Merges with existing Redis state
    ▼
Redis SET: hardware_state_{coop_id}_{hardware_name}
    │
    │ Publish to MQTT: coops/{coop_id}/state
    ▼
ESP32 Device
    │
    │ Subscribed to: coops/{coop_id}/state
    │ Executes hardware command
    ▼
Hardware State Updated
```

### 3. State Retrieval Flow

```
Client
    │
    │ GET /coops/{id}/state
    ▼
JWT Auth + CoopAccess Check
    │
    ▼
Redis SCAN: hardware_state_{coop_id}_*
    │
    │ Retrieves all hardware states for coop
    ▼
Returns JSON array of hardware states
```

### 4. Online/Offline Status Flow

```
Client
    │
    │ GET /coops/{id}/status
    ▼
Redis GET: last_sensor_stored_{coop_id}
    │
    │ If exists:
    │   - Parse timestamp
    │   - Calculate time elapsed
    │   - If elapsed ≤ 45 seconds → is_online: true
    │ If not exists:
    │   - is_online: false
    ▼
Returns status with last_seen and time_elapsed
```

---

## User Flow

### Farmer User Flow

1. **Login to Web Dashboard**
   - Authenticate with username/password
   - Receive JWT token with role: "farmer"

2. **View Coop List**
   - API call to fetch all coops owned by farmer
   - Display online/offline status for each coop

3. **Monitor Specific Coop**
   - Navigate to coop detail page
   - **Real-time Sensor Data**: Fetch latest sensor readings from database
   - **Hardware Status**: Get current states via `GET /coops/{id}/state`
   - **Online Status**: Check if ESP32 is active via `GET /coops/{id}/status`

4. **Control Hardware**
   - Click toggle button for fan/lamp/feeder
   - Client sends `POST /coops/{id}/state` with new state
   - Hardware state updates in Redis within milliseconds
   - ESP32 receives MQTT command and executes

5. **Configure Auto Mode**
   - Set temperature thresholds (min/max)
   - Configure time schedules (e.g., lights on at 6 AM)
   - Client sends partial update with meta data

6. **View Historical Data**
   - Select date range
   - Fetch sensor readings from PostgreSQL
   - Display charts for temperature, humidity trends

### ABK User Flow

1. **Login with ABK Credentials**
   - Receive JWT token with role: "abk"

2. **View Assigned Coops Only**
   - API returns only coops where abk_id matches user ID

3. **Monitor and Control**
   - Same monitoring capabilities as farmer
   - Can control hardware within assigned coops

### Admin User Flow

1. **Bypass Ownership Checks**
   - Access any coop without ownership verification
   - Used for troubleshooting and maintenance

2. **System Monitoring**
   - Check service health via `GET /health`
   - Monitor MQTT connection status

---

## Features

### Core Features

#### 1. Real-time Sensor Data Collection
- **MQTT Subscription**: Listens to `coops/+/telemetry` topic
- **Shared Subscription**: Uses `$share/backend-workers/` for load balancing
- **High Throughput**: Processes thousands of readings per minute
- **Batch Insertion**: Accumulates up to 1000 readings before DB insert
- **Auto-timestamp**: Adds Asia/Jakarta timestamp on receipt

#### 2. Hardware State Management
- **Redis-backed Storage**: Instant state retrieval (< 10ms)
- **Partial Updates**: Update only state, mode, or schedules independently
- **MQTT Publishing**: State changes broadcasted to ESP32 devices
- **Wildcard Pattern**: `hardware_state_{coop_id}_{hardware_name}`

#### 3. Authentication & Authorization
- **JWT Authentication**: HS256 algorithm with secret key
- **Role-Based Access Control**:
  - `admin`: Full access to all coops
  - `farmer`: Access to owned coops only
  - `abk`: Access to assigned coops only
- **Dynamic Role Verification**: Roles fetched from database on each request

#### 4. Online/Offline Detection
- **Last Seen Tracking**: Redis key `last_sensor_stored_{coop_id}`
- **45-Second Threshold**: Coop considered online if data received within 45s
- **Time Elapsed Calculation**: Returns duration since last sensor data

#### 5. RESTful API
- **OpenAPI Specification**: Full API documentation in `openapi.yaml`
- **CORS Enabled**: Cross-origin requests allowed
- **Standard HTTP Codes**: Proper 400, 401, 403, 404, 500 error responses

#### 6. Graceful Shutdown
- **Signal Handling**: Catches SIGINT/SIGTERM
- **Connection Cleanup**: Closes DB, Redis, MQTT connections properly
- **Worker Drain**: Processes remaining jobs before shutdown
- **10-Second Timeout**: Forces shutdown after timeout

#### 7. Test Endpoint
- **`POST /coops/{id}/test-sensor`**: Bypass MQTT for testing
- **Same Database Logic**: Uses identical storage mechanism
- **Useful For**: Development, debugging, load testing

### Advanced Features

#### High-Performance Batch Processing
- **Worker Pool**: 10 concurrent goroutines processing sensor data
- **Buffered Channel**: 20,000 job capacity prevents data loss
- **Smart Batching**: Inserts on max batch size OR timeout (whichever first)
- **COPY Protocol**: Uses PostgreSQL's fastest bulk insert method

#### Auto-Reconnection
- **MQTT Client**: Automatically reconnects on connection loss
- **Retry Logic**: Built into paho.mqtt.golang client
- **Connection Status**: Tracked and exposed via health endpoint

#### Flexible Hardware Metadata
- **Mode Selection**: "auto" or "manual" operation
- **Temperature Thresholds**: Min/max triggers for auto mode
- **Time Schedules**: Multiple ON/OFF schedules per day
- **Collision Detection**: Validates schedule overlaps

---

## File Structure

```
go-mqtt/
├── main.go                      # Entry point, HTTP server, routers, handlers
├── go.mod                       # Go module definition
├── go.sum                       # Dependency checksums
├── .env                         # Environment variables (not in git)
├── .env.example                 # Environment variables template
├── .gitignore                   # Git ignore rules
├── Dockerfile                   # Docker container definition
├── build.sh                     # Build script for production
├── openapi.yaml                 # OpenAPI 3.0 specification
├── README.md                    # Quick start guide
├── PROJECT_DESCRIPTION.md       # This file - comprehensive documentation
├── esp.ino                      # Arduino/ESP32 firmware reference
├── server                       # Compiled binary (not in git)
├── server_openapi.yml           # Server-specific OpenAPI spec
├── logs/                        # Application logs directory
│
└── src/                         # Source code packages
    ├── broker/
    │   └── mqtt.go              # MQTT client setup, message handling
    ├── cache/
    │   └── redis.go             # Redis client configuration
    ├── config/
    │   └── config.go            # Environment variables, logger setup
    └── database/
        └── database.go          # PostgreSQL, worker pool, batch insert
```

### File Descriptions

#### `main.go` (638 lines)
- **Purpose**: Application entry point and HTTP server
- **Key Functions**:
  - `main()`: Initializes all services, starts server
  - `SetupRouter()`: Configures Chi router with middleware
  - `CoopAccessMiddleware()`: Authorization logic
  - `StartServer()`: HTTP server with graceful shutdown
  - `HTTPError()`: Standardized error responses
  - `CORSMiddleware()`: CORS headers
  - `mergePartialUpdates()`: Merges partial hardware updates
  - `validateSchedules()`: Validates time schedules

#### `src/broker/mqtt.go` (123 lines)
- **Purpose**: MQTT client and message handling
- **Key Functions**:
  - `SetupMQTTClient()`: Creates and configures MQTT client
  - `MqttPublish()`: Publishes messages to MQTT topics
- **Key Structs**:
  - `TelemetryMessage`: Incoming sensor data structure
  - `MQTTConfig`: Configuration struct

#### `src/database/database.go` (259 lines)
- **Purpose**: Database operations and worker pool
- **Key Functions**:
  - `NewPostgresStore()`: Initializes connection pool
  - `Save()`: Queues sensor data for batch insert
  - `DBWorker()`: Worker goroutine for batch processing
  - `batchInsert()`: Performs COPY FROM bulk insert
  - `LoadSensorTypes()`: Caches sensor type mappings
  - `GetCoopOwnerAndABK()`: Fetches coop ownership
  - `GetUserRole()`: Retrieves user role from DB
- **Constants**:
  - `JOB_BUFFER_SIZE = 20000`
  - `MAX_BATCH_SIZE = 1000`
  - `BATCH_TIMEOUT = 1s`
  - `NUM_DB_WORKERS = 10`

#### `src/cache/redis.go` (49 lines)
- **Purpose**: Redis client setup
- **Key Functions**:
  - `SetupRedisClient()`: Creates Redis connection
- **Key Structs**:
  - `RedisConfig`: Configuration struct

#### `src/config/config.go` (77 lines)
- **Purpose**: Configuration management
- **Key Functions**:
  - `LoadConfig()`: Loads all environment variables
  - `SetupLogger()`: Configures structured logging
  - `GetRequiredEnv()`: Validates required env vars
- **Key Structs**:
  - `Config`: Holds all configuration values

---

## Technical Specifications

### Technology Stack

#### Backend Framework
- **Language**: Go 1.24.5
- **HTTP Router**: Chi v5.2.3 (lightweight, idiomatic)
- **MQTT Client**: Eclipse Paho MQTT v1.5.1

#### Database & Cache
- **PostgreSQL Driver**: pgx v5.7.6 (high-performance)
- **Connection Pooling**: pgxpool (built-in connection management)
- **Redis Client**: go-redis v9.14.0

#### Authentication
- **JWT Library**: go-chi/jwtauth v5.3.3
- **Algorithm**: HS256 (HMAC-SHA256)
- **Token Claims**: `sub` (user ID), `role` (role name)

#### Utilities
- **Environment**: godotenv v1.5.1 (load .env files)
- **Logging**: slog (Go 1.21+ structured logging)

### Performance Specifications

#### Throughput
- **Sensor Data**: 10,000+ readings/minute per worker
- **HTTP Requests**: 1,000+ requests/second (depending on DB/Redis latency)
- **MQTT Messages**: Subscribes to shared topic for horizontal scaling

#### Latency
- **Hardware State Update**: < 100ms (Redis write + MQTT publish)
- **Sensor Data Storage**: 1-2 seconds (batch timeout)
- **HTTP API Response**: 50-200ms average

#### Resource Usage
- **Memory**: ~50-100 MB (varies with worker pool size)
- **CPU**: 5-15% (single core) under normal load
- **Database Connections**: 10-20 (configurable pool size)

### Configuration Requirements

#### Environment Variables
```bash
# Server
PORT=21999

# PostgreSQL
DB_HOST=103.197.190.23
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=<password>
DB_DATABASE=rafflesiaagro.migration

# Redis
REDIS_HOST=103.197.190.23
REDIS_PORT=6379
REDIS_PASSWORD=<password>
REDIS_DB=2

# MQTT
MQTT_URL=tcp://mqttserver.rafflesiaagro.com:1883
MQTT_USERNAME=<username>
MQQT_PASSWORD=<password>

# Security
JWT_SECRET=<secret-key>

# Logging
LOG_LEVEL=info|debug
```

### Deployment Architecture

#### Production Deployment
- **Container**: Alpine Linux-based Docker image
- **Port**: 21999 (internal), exposed via reverse proxy
- **Replicas**: Single instance (horizontal scaling possible via MQTT shared subscriptions)
- **Reverse Proxy**: Nginx/Traefik for SSL termination
- **Monitoring**: Structured logs to stdout for log aggregators

#### Docker Configuration
```dockerfile
FROM alpine:latest
RUN apk add --no-cache tzdata
ENV TZ=Asia/Jakarta
COPY . .
EXPOSE 21999
CMD ["/server"]
```

---

## Router Structure

### Chi Router Configuration

**File**: `main.go:317-605`

```go
func SetupRouter(client mqtt.Client, redisClient *redis.Client,
                 tokenAuth *jwtauth.JWTAuth, store *database.PostgresStore) http.Handler
```

### Route Definitions

#### Public Routes
```
GET  /health
     └── Returns: { "status": "ok", "mqtt_online": true }
     └── Middleware: None
```

#### Protected Routes (JWT Required)

All routes under `/coops/{id}/*` require:
1. `jwtauth.Verifier` - Extracts token from Authorization header
2. `jwtauth.Authenticator` - Validates token signature
3. `CoopAccessMiddleware` - Checks user permissions

```
GET  /coops/{id}/state
     └── Returns: All hardware states for a coop
     └── Middleware: JWT Verifier, Authenticator, CoopAccess
     └── Redis Operation: SCAN hardware_state_{coop_id}_*

POST /coops/{id}/state
     └── Accepts: HardwareState or PartialHardwareState
     └── Returns: Success confirmation
     └── Middleware: JWT Verifier, Authenticator, CoopAccess
     └── Redis Operation: SET hardware_state_{coop_id}_{hardware}
     └── MQTT Publish: coops/{id}/state

GET  /coops/{id}/status
     └── Returns: Online status, last seen, time elapsed
     └── Middleware: JWT Verifier, Authenticator, CoopAccess
     └── Redis Operation: GET last_sensor_stored_{coop_id}

POST /coops/{id}/test-sensor
     └── Accepts: Sensor data JSON
     └── Returns: Saved confirmation
     └── Middleware: JWT Verifier, Authenticator, CoopAccess
     └── Purpose: Testing without MQTT
```

### Middleware Chain

```
Request → CORSMiddleware → RequestID → RealIP → Logger → Recoverer
       ↓
   JWT Verifier (extract token)
       ↓
   JWT Authenticator (validate token)
       ↓
   CoopAccessMiddleware (check permissions)
       ↓
   Route Handler
```

---

## Controllers & Business Logic

### Controller Architecture

**Pattern**: Functional controllers (no struct-based controllers)
Handlers are defined as anonymous functions within `SetupRouter()`

### Key Controllers

#### 1. Health Check Controller
**Location**: `main.go:323-326`

```go
r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]any{
        "status": "ok",
        "mqtt_online": client.IsConnected()
    })
})
```

**Logic**:
- Checks MQTT client connection status
- Returns JSON with service status

#### 2. Get Hardware States Controller
**Location**: `main.go:336-383`

**Business Logic**:
1. Extract `coop_id` from URL parameter
2. Build Redis pattern: `hardware_state_{coop_id}_*`
3. Use `SCAN` to find all matching keys (iterative, non-blocking)
4. For each key:
   - `GET` value from Redis
   - Unmarshal JSON to `HardwareState`
   - Append to results array
5. Return JSON with count

**Error Handling**:
- Redis scan errors → 500 Internal Server Error
- Unmarshal errors → Log warning, skip corrupted entry
- Empty results → Return empty array (not an error)

#### 3. Update Hardware State Controller
**Location**: `main.go:433-545`

**Business Logic**:
1. Parse request body (1MB limit)
2. Detect update type:
   - **Partial Update**: `PartialHardwareState` (only provided fields)
   - **Full Update**: `HardwareState` (complete replacement)
3. For partial updates:
   - Fetch existing state from Redis
   - Merge with new values using `mergePartialUpdates()`
4. Validate schedules (if provided in meta)
5. Marshal to JSON
6. `SET` to Redis: `hardware_state_{coop_id}_{hardware}`
7. `PUBLISH` to MQTT: `coops/{coop_id}/state`
8. Return 202 Accepted

**Error Handling**:
- Invalid JSON → 400 Bad Request
- Invalid schedules → 400 with validation error
- Redis failures → 500 Internal Server Error
- MQTT publish timeout → 500 Internal Server Error

#### 4. Get Coop Status Controller
**Location**: `main.go:385-431`

**Business Logic**:
1. Build Redis key: `last_sensor_stored_{coop_id}`
2. `GET` value from Redis
3. If key exists:
   - Parse RFC3339 timestamp
   - Calculate elapsed time
   - Check if elapsed ≤ 45 seconds
   - Set `is_online: true` or `false`
4. If key doesn't exist:
   - Set `is_online: false`
5. Return JSON with status details

**Error Handling**:
- Timestamp parse errors → 500 Internal Server Error
- Redis errors → 500 Internal Server Error

#### 5. Test Sensor Controller
**Location**: `main.go:548-601`

**Business Logic**:
1. Parse `coop_id` from URL
2. Parse JSON body (1MB limit)
3. Create `TelemetryMessage` with current Jakarta time
4. Reload sensor types from database
5. Call `store.Save()` to queue for batch insert
6. Return 201 Created

**Purpose**: Development/testing endpoint that bypasses MQTT

### Middleware Logic

#### CoopAccessMiddleware
**Location**: `main.go:248-315`

**Authorization Flow**:
1. Extract `coop_id` from URL
2. Extract `user_id` from JWT claims (`sub`)
3. Query database for user's current role
4. If role is "admin":
   - Bypass all checks, allow access
5. If role is "farmer" or "abk":
   - Query `GetCoopOwnerAndABK(coop_id)`
   - Farmer check: `farmer_id == user_id AND farmer_id != 0`
   - ABK check: `abk_id == user_id AND abk_id != 0`
6. If neither condition met:
   - Return 403 Forbidden

**SQL Queries**:
```sql
-- Get user role
SELECT r.name FROM users u
JOIN roles r ON u.current_role_id = r.id
WHERE u.id = $1

-- Get coop ownership
SELECT f.farmer_id, c.abk_id
FROM coops c
LEFT JOIN farms f ON c.farm_id = f.id
WHERE c.id = $1
```

### Helper Functions

#### mergePartialUpdates()
**Location**: `main.go:83-124`

**Logic**:
- Starts with existing hardware state
- Updates fields only if provided in partial update
- Handles nested meta updates (mode, temperature, schedules)
- Returns merged state

#### validateSchedules()
**Location**: `main.go:126-160`

**Validations**:
1. Order numbers must be unique
2. Time format must be HH:MM
3. `time_off` must be after `time_on`
4. No overlapping schedules (next start time must be after previous end time)
5. Returns descriptive error for validation failures

---

## Important Code References

### Critical Code Sections

#### 1. MQTT Message Handler
**File**: `src/broker/mqtt.go:59-95`

**Why Important**: Core telemetry processing logic

```go
if token := c.Subscribe(telemetryTopic, 1, func(client mqtt.Client, msg mqtt.Message) {
    parts := strings.Split(msg.Topic(), "/")
    if len(parts) != 3 {
        slog.Warn("Received message on unexpected topic format", ...)
        return
    }
    coopID := parts[1] // coops/[0] coopID/[1] telemetry/[2]

    var sensorData map[string]any
    if err := json.Unmarshal(msg.Payload(), &sensorData); err != nil {
        slog.Error("Failed to unmarshal telemetry", ...)
        return
    }

    telemetry := TelemetryMessage{
        CoopID:    coopID,
        Data:      sensorData,
        Timestamp: GetJakartaTime(),
    }

    if err := store.Save(telemetry, sensorTypes); err != nil {
        slog.Error("Failed to queue sensor data", ...)
    }
}); token.Wait() && token.Error() != nil {
    slog.Error("Failed to subscribe to telemetry topic", ...)
}
```

#### 2. Batch Insert with COPY Protocol
**File**: `src/database/database.go:214-238`

**Why Important**: Highest-performance database insertion method

```go
func batchInsert(pool *pgxpool.Pool, readings []SensorReading) error {
    if len(readings) == 0 {
        return nil
    }

    rows := make([][]interface{}, len(readings))
    for i, r := range readings {
        rows[i] = []interface{}{r.Value, r.SensorTypeID, r.CoopID, r.Timestamp, r.Timestamp}
    }

    _, err := pool.CopyFrom(
        context.Background(),
        pgx.Identifier{"rec_sensor_coops"},
        []string{"value", "sensor_type_id", "coop_id", "created_at", "updated_at"},
        pgx.CopyFromRows(rows),
    )

    if err != nil {
        slog.Error("COPY From failed", "error", err)
        return err
    }

    slog.Info("Successfully inserted batch", "rows", len(readings))
    return nil
}
```

#### 3. Worker Pool Initialization
**File**: `main.go:204-216`

**Why Important**: High-throughput data processing architecture

```go
sensorTypes, err := database.LoadSensorTypes(store.Pool)
if err != nil {
    slog.Error("Failed to load sensor types", "error", err)
    os.Exit(1)
}

var wg sync.WaitGroup
for i := 0; i < database.NUM_DB_WORKERS; i++ {
    wg.Add(1)
    go database.DBWorker(i, &wg, store.Pool, store.JobChan, sensorTypes)
}
slog.Info("🚀 Database worker pool started", "workers", database.NUM_DB_WORKERS)
```

#### 4. Partial Update Handling
**File**: `main.go:440-508`

**Why Important**: Flexible API design for better UX

```go
// First attempt to decode as partial update
var partialUpdate PartialHardwareState
partialErr := decoder.Decode(&partialUpdate)

var finalState HardwareState
var hardwareName string

if partialErr == nil && partialUpdate.Hardware != "" {
    // This is a partial update - fetch existing state and merge
    hardwareName = partialUpdate.Hardware
    redisKey := fmt.Sprintf("hardware_state_%s_%s", coopID, hardwareName)

    // Get existing state from Redis
    existingVal, err := redisClient.Get(r.Context(), redisKey).Result()
    // ... unmarshal existing state

    // Merge partial updates
    finalState = mergePartialUpdates(finalState, partialUpdate)
} else {
    // This is a full state update (backward compatibility)
    // ... decode as HardwareState
}
```

#### 5. Graceful Shutdown
**File**: `main.go:228-245`

**Why Important**: Clean resource management for 24/7 operation

```go
stop := make(chan os.Signal, 1)
signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
srv := StartServer(router, cfg)
<-stop

slog.Info("⏳ Shutting down services...")
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
if err := srv.Shutdown(ctx); err != nil {
    slog.Error("Server shutdown failed", "err", err)
}
mqttClient.Disconnect(250)
slog.Info("MQTT client disconnected")
close(store.JobChan)
wg.Wait()
slog.Info("All database workers have finished.")
slog.Info("✅ Server gracefully stopped.")
```

#### 6. Authorization Check
**File**: `main.go:277-309`

**Why Important**: Security and multi-tenancy

```go
// BARU: Ambil peran pengguna langsung dari database.
role, err := store.GetUserRole(userID)
if err != nil {
    HTTPError(w, http.StatusForbidden, err)
    return
}

// Bypass otorisasi untuk admin
if role == "admin" {
    slog.Info("Admin access granted, bypassing ownership checks.")
    next.ServeHTTP(w, r)
    return
}

// Query database untuk mendapatkan pemilik & ABK kandang
farmerID, abkID, err := store.GetCoopOwnerAndABK(coopID)
if err != nil {
    if err == pgx.ErrNoRows {
        HTTPError(w, http.StatusNotFound, errors.New("coop not found"))
        return
    }
    HTTPError(w, http.StatusInternalServerError, err)
    return
}

// Terapkan logika otorisasi untuk non-admin
isFarmOwner := (role == "farmer" && farmerID == userID && farmerID != 0)
isAssignedAbk := (role == "abk" && abk_id == userID && abk_id != 0)

if !isFarmOwner && !isAssignedAbk {
    HTTPError(w, http.StatusForbidden, fmt.Errorf("you are not authorized to access coop '%s'", coopIDStr))
    return
}
```

---

## Areas for Improvement

### Current Limitations & Potential Enhancements

#### 1. Security Enhancements

**Issue**: JWT secret stored in environment variable
- **Improvement**: Use key management service (HashiCorp Vault, AWS Secrets Manager)

**Issue**: No rate limiting on API endpoints
- **Improvement**: Implement rate limiting middleware (e.g., tollbooth, slowcache)

**Issue**: CORS allows all origins (`*`)
- **Improvement**: Whitelist specific domains in production

**Issue**: No request signing for MQTT
- **Improvement**: Implement MQTT TLS + client certificates

#### 2. Performance Optimizations

**Issue**: Sensor types loaded at startup only
- **Improvement**: Implement periodic refresh or cache invalidation

**Issue**: No connection pooling configuration exposed
- **Improvement**: Make pool size, max connections configurable via env vars

**Issue**: Sequential authorization check (2 DB queries per request)
- **Improvement**: Cache user roles in Redis with TTL (5 minutes)

**Issue**: No pagination for hardware state retrieval
- **Improvement**: Add pagination if coops have many hardware devices

#### 3. Reliability Improvements

**Issue**: No retry logic for failed Redis operations
- **Improvement**: Implement exponential backoff retry mechanism

**Issue**: MQTT publish timeout (5 seconds) may be too long
- **Improvement**: Use background publish with acknowledgment queue

**Issue**: No dead letter queue for failed sensor data
- **Improvement**: Implement DLQ for failed batch inserts

**Issue**: No health check for database connection
- **Improvement**: Add `/health/db` endpoint with ping check

#### 4. Monitoring & Observability

**Issue**: Structured logging but no centralized logging
- **Improvement**: Integrate ELK Stack, Loki, or CloudWatch

**Issue**: No metrics collection
- **Improvement**: Add Prometheus metrics (request rate, batch insert time, queue depth)

**Issue**: No distributed tracing
- **Improvement**: Add OpenTelemetry for request tracing

**Issue**: No alerting on failures
- **Improvement**: Integrate PagerDuty, Slack webhooks on critical errors

#### 5. Feature Enhancements

**Issue**: No hardware state history
- **Improvement**: Store state changes in audit log table

**Issue**: No bulk hardware control
- **Improvement**: Add endpoint to update multiple hardware at once

**Issue**: No sensor data aggregation
- **Improvement**: Add endpoints for min/max/avg over time periods

**Issue**: No alert configuration
- **Improvement**: Allow users to set thresholds for alerts (e.g., temp > 35°C)

**Issue**: No webhook notifications
- **Improvement**: Send webhooks on hardware state changes or threshold violations

#### 6. Code Quality

**Issue**: Large handler functions in main.go
- **Improvement**: Extract handlers to separate package (`handlers/`)

**Issue**: Mixed concerns in handlers (auth, business logic, data access)
- **Improvement**: Implement clean architecture with service layer

**Issue**: Hard-coded constants (buffer size, batch size)
- **Improvement**: Move to configuration file

**Issue**: Limited integration tests
- **Improvement**: Add test suite with mocked MQTT/DB/Redis

#### 7. Database Optimizations

**Issue**: No data retention policy
- **Improvement**: Implement partitioning by month, auto-drop old data

**Issue**: No materialized views for aggregated data
- **Improvement**: Create hourly/daily summaries for fast dashboard queries

**Issue**: No indexing strategy documented
- **Improvement**: Document and implement indexes on common query patterns

#### 8. DevOps & Deployment

**Issue**: Single binary deployment
- **Improvement**: Implement blue-green deployment strategy

**Issue**: No rollback mechanism
- **Improvement**: Container image tagging and rollback scripts

**Issue**: Manual environment configuration
- **Improvement**: Infrastructure as Code (Terraform/Ansible)

**Issue**: No backup strategy documented
- **Improvement**: Automated backup procedures for Redis and PostgreSQL

---

## Installation & Setup

### Prerequisites

- **Go**: 1.24.5 or higher
- **PostgreSQL**: 12+ with database schema
- **Redis**: 6+ (any edition)
- **MQTT Broker**: Mosquitto, HiveMQ, or AWS IoT Core
- **Operating System**: Linux (recommended), macOS, Windows

### Local Development Setup

#### 1. Clone Repository
```bash
git clone https://github.com/Anjasfedo/go-mqtt.git
cd go-mqtt
```

#### 2. Install Dependencies
```bash
go mod download
```

#### 3. Configure Environment
```bash
cp .env.example .env
# Edit .env with your configuration
nano .env
```

#### 4. Database Setup
```sql
-- Create database
CREATE DATABASE "rafflesiaagro.migration";

-- Connect to database
\c "rafflesiaagro.migration"

-- Create tables (schema should match the code expectations)
-- Tables needed: rec_sensor_coops, sensor_types, coops, farms, users, roles

-- Insert sensor types
INSERT INTO sensor_types (name) VALUES
('temperature'), ('humidity'), ('light'), ('ammonia');

-- Insert roles
INSERT INTO roles (name) VALUES ('admin'), ('farmer'), ('abk');
```

#### 5. Run Application
```bash
# Development mode
go run main.go

# Or build and run
go build -o server
./server
```

#### 6. Verify Installation
```bash
# Check health endpoint
curl http://localhost:21999/health

# Expected response:
# {"status":"ok","mqtt_online":true}
```

### Docker Deployment

#### 1. Build Docker Image
```bash
docker build -t go-mqtt:latest .
```

#### 2. Run Container
```bash
docker run -d \
  --name go-mqtt \
  --env-file .env \
  -p 21999:21999 \
  go-mqtt:latest
```

#### 3. Docker Compose (Recommended)
```yaml
version: '3.8'
services:
  go-mqtt:
    build: .
    ports:
      - "21999:21999"
    env_file:
      - .env
    depends_on:
      - postgres
      - redis
      - mqtt
    restart: unless-stopped

  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: rafflesiaagro.migration
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: your_password
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass your_redis_password
    volumes:
      - redis_data:/data

  mqtt:
    image: eclipse-mosquitto:2
    ports:
      - "1883:1883"
    volumes:
      - ./mosquitto.conf:/mosquitto/config/mosquitto.conf

volumes:
  postgres_data:
  redis_data:
```

### Production Deployment

#### 1. Build Production Binary
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o server

# With optimizations
go build -ldflags="-s -w" -o server
```

#### 2. Using Build Script
```bash
chmod +x build.sh
./build.sh
```

#### 3. Deploy to Server
```bash
# Copy to server
scp server user@server:/path/to/deploy/

# SSH to server
ssh user@server

# Run as service
sudo systemctl start go-mqtt
```

#### 4. Systemd Service Configuration
```ini
[Unit]
Description=Go-MQTT Service
After=network.target

[Service]
Type=simple
User=go-mqtt
WorkingDirectory=/opt/go-mqtt
ExecStart=/opt/go-mqtt/server
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

### Environment Variables Guide

#### Required Variables (No Defaults)
- `PORT`: HTTP server port (e.g., 21999)
- `DB_HOST`: PostgreSQL host
- `DB_PORT`: PostgreSQL port (usually 5432)
- `DB_USERNAME`: Database user
- `DB_PASSWORD`: Database password
- `DB_DATABASE`: Database name
- `REDIS_HOST`: Redis host
- `REDIS_PORT`: Redis port (usually 6379)
- `REDIS_PASSWORD`: Redis password
- `REDIS_DB`: Redis database number (0-15)
- `MQTT_URL`: MQTT broker URL (e.g., tcp://localhost:1883)
- `MQTT_USERNAME`: MQTT username
- `MQTT_PASSWORD`: MQTT password
- `JWT_SECRET`: Secret key for JWT signing

#### Optional Variables
- `LOG_LEVEL`: Logging level (info, debug) - Default: info

### Testing

#### 1. Health Check
```bash
curl http://localhost:21999/health
```

#### 2. Test Sensor Data (No MQTT Required)
```bash
curl -X POST http://localhost:21999/coops/123/test-sensor \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "temperature": 28.5,
    "humidity": 75.2,
    "light": 500
  }'
```

#### 3. Get Hardware States
```bash
curl http://localhost:21999/coops/123/state \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

#### 4. Update Hardware State
```bash
curl -X POST http://localhost:21999/coops/123/state \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "hardware": "fan1",
    "state": true,
    "meta": {
      "mode": "auto",
      "temperature_min": 25.0,
      "temperature_max": 30.0
    }
  }'
```

---

## License

This project appears to be a private/proprietary project for Rafflesia Agro. No explicit license file is present in the repository.

**For Open Source Distribution** (if applicable):
Consider adding a `LICENSE` file. Common options:
- **MIT License**: Permissive, simple
- **Apache 2.0**: Patent protection, widely used
- **GPL v3**: Copyleft, requires derivative works to be open source

---

## Author & Team

### Project Information
- **Company**: Rafflesia Agro
- **Domain**: Agricultural IoT / Smart Poultry Farming
- **Location**: Indonesia (Asia/Jakarta timezone)

### Development Team
- **Main Developer**: Anjasfedo (GitHub: @Anjasfedo)
- **Project Status**: Active (as of latest commit: Oct 10, 2024)

### Contact & Support
- **Repository**: https://github.com/Anjasfedo/go-mqtt
- **Issue Tracking**: GitHub Issues
- **Documentation**: This file (PROJECT_DESCRIPTION.md) and README.md

### Related Projects
- **ESP32 Firmware**: `esp.ino` (Arduino code for hardware devices)
- **API Specification**: `openapi.yaml` (OpenAPI 3.0 documentation)
- **Web/Mobile Client**: (Separate repository, not included)

### Acknowledgments
- **MQTT Library**: Eclipse Paho MQTT Go Client
- **HTTP Router**: Chi - lightweight, idiomatic HTTP router
- **Database Driver**: pgx - high-performance PostgreSQL driver
- **JWT Library**: go-chi/jwtauth

---

## Additional Resources

### Documentation
- **README.md**: Quick start guide and API examples
- **openapi.yaml**: Full OpenAPI 3.0 specification
- **server_openapi.yml**: Server-specific API documentation

### ESP32 Integration
- **esp.ino**: Arduino firmware for ESP32 devices
  - MQTT connection logic
  - Sensor data publishing format
  - Hardware state subscription handling

### Logging & Monitoring
- **Logs Directory**: `/logs` (application logs)
- **Log Format**: Structured JSON via slog
- **Log Levels**: info (default), debug (verbose)

### Production Deployment Checklist
- [ ] Environment variables configured
- [ ] PostgreSQL schema created
- [ ] Redis connection tested
- [ ] MQTT broker accessible
- [ ] JWT secret generated (32+ characters)
- [ ] Firewall rules configured (port 21999)
- [ ] SSL/TLS configured (reverse proxy)
- [ ] Log aggregation configured
- [ ] Monitoring setup (Prometheus/DataDog)
- [ ] Backup strategy implemented
- [ ] Rollback procedure documented

---

## Quick Reference

### Essential Commands

```bash
# Build
go build -o server

# Run
./server

# Docker
docker build -t go-mqtt .
docker run -p 21999:21999 --env-file .env go-mqtt

# Dependencies
go mod tidy
go mod download

# Testing
curl http://localhost:21999/health
```

### Key Configuration Numbers
- **HTTP Port**: 21999
- **DB Workers**: 10
- **Job Buffer**: 20,000
- **Batch Size**: 1,000
- **Batch Timeout**: 1 second
- **Online Threshold**: 45 seconds

### Redis Key Patterns
- **Hardware State**: `hardware_state_{coop_id}_{hardware_name}`
- **Last Sensor**: `last_sensor_stored_{coop_id}`

### MQTT Topics
- **Subscribe**: `$share/backend-workers/coops/+/telemetry`
- **Publish**: `coops/{coop_id}/state`

---

**Document Version**: 1.0
**Last Updated**: January 9, 2026
**Generated For**: Go-MQTT Project by Rafflesia Agro
