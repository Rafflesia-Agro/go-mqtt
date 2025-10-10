# go-mqtt

A Go-based MQTT broker with HTTP API for hardware state management and monitoring.

## API Endpoints

### Hardware State Management

#### GET `/coops/{id}/state`

Retrieves all hardware states for a specific coop from Redis.

**Headers:**
- `Authorization: Bearer <JWT_TOKEN>`

**Response Example:**
```json
{
  "ok": true,
  "data": {
    "coop_id": "123",
    "hardware_states": [
      {
        "hardware": "lamp",
        "state": true,
        "meta": {
          "mode": "auto",
          "temperature_min": 25.0,
          "temperature_max": 30.0,
          "schedules": [
            {
              "order": 1,
              "time_on": "06:00",
              "time_off": "08:00"
            },
            {
              "order": 2,
              "time_on": "17:00",
              "time_off": "19:00"
            }
          ]
        }
      },
      {
        "hardware": "fan",
        "state": false,
        "meta": {
          "mode": "manual",
          "temperature_min": 20.0,
          "temperature_max": 28.0,
          "schedules": []
        }
      },
      {
        "hardware": "lamp_2",
        "state": true,
        "meta": null
      },
      {
        "hardware": "water",
        "state": true,
        "meta": null
      }
    ],
    "count": 3
  }
}
```

**Empty Response (no hardware states found):**
```json
{
  "ok": true,
  "data": {
    "coop_id": "123",
    "hardware_states": [],
    "count": 0
  }
}
```

#### POST `/coops/{id}/state`

Updates hardware state for a specific coop. Supports both full and partial updates.

**Headers:**
- `Authorization: Bearer <JWT_TOKEN>`
- `Content-Type: application/json`

**Full State Update (backward compatibility):**
```json
{
  "hardware": "feeder1",
  "state": true,
  "meta": {
    "mode": "auto",
    "temperature_min": 25.0,
    "temperature_max": 30.0,
    "schedules": [
      {
        "order": 1,
        "time_on": "06:00",
        "time_off": "18:00"
      }
    ]
  }
}
```

**Partial Update Examples:**

*Update only state:*
```json
{
  "hardware": "feeder1",
  "state": false
}
```

*Update only mode:*
```json
{
  "hardware": "heater1",
  "meta": {
    "mode": "manual"
  }
}
```

*Update only temperature settings:*
```json
{
  "hardware": "heater1",
  "meta": {
    "temperature_min": 22.0,
    "temperature_max": 27.0
  }
}
```

*Update only schedules:*
```json
{
  "hardware": "light1",
  "meta": {
    "schedules": [
      {
        "order": 1,
        "time_on": "06:00",
        "time_off": "18:00"
      },
      {
        "order": 2,
        "time_on": "19:00",
        "time_off": "21:00"
      }
    ]
  }
}
```

**Response Example:**
```json
{
  "ok": true,
  "message": "hardware state updated successfully",
  "data": {
    "hardware": "feeder1",
    "coop_id": "123"
  }
}
```

### Health Check

#### GET `/health`

Returns the service status and MQTT connection status.

**Response Example:**
```json
{
  "status": "ok",
  "mqtt_online": true
}
```

## Data Models

### HardwareState
```json
{
  "hardware": "string",
  "state": "boolean",
  "meta": "HardwareMeta|null"
}
```

### HardwareMeta
```json
{
  "mode": "string",
  "temperature_min": "number|null",
  "temperature_max": "number|null",
  "schedules": "TimeSchedule[]"
}
```

### TimeSchedule
```json
{
  "order": "integer",
  "time_on": "string", // HH:MM format
  "time_off": "string" // HH:MM format
}
```

## Redis Key Pattern

Hardware states are stored in Redis using the pattern:
```
hardware_state_{coop_id}_{hardware_name}
```

Example: `hardware_state_123_feeder1`

## Error Responses

All endpoints return consistent error responses:

```json
{
  "ok": false,
  "error": "error message"
}
```

Common HTTP status codes:
- `400 Bad Request` - Invalid input data
- `401 Unauthorized` - Missing or invalid JWT token
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Coop not found
- `500 Internal Server Error` - Server error

## Authentication & Authorization

- JWT-based authentication required for all protected endpoints
- Role-based access control (admin, farmer, abk)
- Users can only access coops they own or are assigned to