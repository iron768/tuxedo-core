# Tuxedo Core - API Server

Backend API server for the Tuxedo scene editor. Serves scene files, assets, and provides file management for CPPS Yukon development.

## Features

- 📁 Scene file management (list, read, write)
- 🖼️ Asset serving from yukon project
- 🔄 WebSocket support for live reload
- 🌐 CORS-enabled for local development
- ⚙️ JSON-based configuration
- 📝 Request logging middleware

## Prerequisites

- Go 1.21 or higher
- Yukon project in adjacent directory (`../yukon`)

## Installation

```bash
# Install Go dependencies
go mod download
```

## Configuration

The server uses `config.json` for configuration. Create it or modify the default:

```json
{
  "server": {
    "port": "3000",
    "host": "0.0.0.0",
    "allowOrigins": "*"
  },
  "project": {
    "yukonPath": "../yukon",
    "scenesPath": "src/scenes",
    "assetsPath": "assets"
  },
  "logging": {
    "enabled": true,
    "level": "info",
    "format": "json"
  }
}
```

### Configuration Options

**Server Settings:**
- `port`: Server port (default: 3000)
- `host`: Bind address (default: 0.0.0.0)
- `allowOrigins`: CORS allowed origins (default: *)

**Project Paths:**
- `yukonPath`: Path to yukon project root
- `scenesPath`: Relative path to scenes within yukon
- `assetsPath`: Relative path to assets within yukon

**Logging:**
- `enabled`: Enable/disable logging
- `level`: Log level (debug, info, warn, error)
- `format`: Log format (json, text)

## Running

```bash
# Run with hot reload (development)
go run main.go

# Build and run
go build
./tuxedo-core  # or tuxedo-core.exe on Windows
```

The server will start on `http://localhost:3000` by default.

## Project Structure

```
tuxedo-core/
├── config/              # Configuration package
│   └── config.go        # Config loader and types
├── handlers/            # HTTP request handlers
│   ├── assets.go        # Asset endpoints
│   ├── scenes.go        # Scene CRUD operations
│   ├── project.go       # Project info endpoints
│   └── websocket.go     # WebSocket handler
├── middleware/          # HTTP middleware
│   ├── cors.go          # CORS handling
│   └── logger.go        # Request logging
├── models/              # Data models
│   └── scene.go         # Scene types
├── config.json          # Configuration file
├── go.mod               # Go modules
└── main.go              # Entry point
```

## API Endpoints

### Scenes

**GET** `/api/scenes`
- List all available scenes
- Returns array of scene metadata

**GET** `/api/scenes/{path}`
- Get specific scene file
- `path`: Scene path (e.g., `rooms/town/Town`)
- Returns scene JSON

**PUT** `/api/scenes/{path}`
- Update scene file
- Request body: Scene JSON
- Returns updated scene

**POST** `/api/scenes`
- Create new scene
- Request body: Scene JSON with path
- Returns created scene

### Prefabs

**GET** `/api/prefab/{id}`
- Get prefab definition by ID
- `id`: Prefab UUID (e.g., `d3866883-7507-4f66-a7e3-bc9a896c4a22`)
- Searches in `shared_prefabs` directory
- Returns prefab scene JSON with `sceneType: "PREFAB"`

### Assets

**GET** `/assets/{path}`
- Serve asset files from yukon
- Supports images, JSON, and other static files
- Example: `/assets/media/rooms/town/town-pack.json`

**GET** `/api/assets`
- List available assets
- Returns array of asset metadata

### Project

**GET** `/api/project`
- Get project information
- Returns project stats and structure

### WebSocket

**GET** `/api/ws`
- WebSocket connection for live updates
- Supports file watching and hot reload

## Development

### Testing Endpoints

```bash
# List scenes
curl http://localhost:3000/api/scenes

# Get specific scene
curl http://localhost:3000/api/scenes/rooms/town/Town

# Test asset serving
curl http://localhost:3000/assets/media/rooms/town/town-pack.json
```

## Logging

The server logs all HTTP requests:

```
2024/01/15 12:34:56 Starting Tuxedo Core server...
2024/01/15 12:34:56 Configuration:
2024/01/15 12:34:56   - Port: 3000
2024/01/15 12:34:56   - Assets Path: ../yukon/assets
2024/01/15 12:34:56   - Scenes Path: ../yukon/src/scenes
2024/01/15 12:34:56 Serving assets from: ../yukon/assets
2024/01/15 12:34:56 Tuxedo Core server listening on 0.0.0.0:3000
2024/01/15 12:34:58 GET /api/scenes - 200 OK (12ms)
```

## Troubleshooting

### Port already in use
```bash
# Change port in config.json
{
  "server": {
    "port": "3001"
  }
}
```

### Cannot find yukon directory
```bash
# Update yukonPath in config.json
{
  "project": {
    "yukonPath": "/absolute/path/to/yukon"
  }
}
```

### CORS errors
```bash
# Update allowOrigins in config.json for specific domain
{
  "server": {
    "allowOrigins": "http://localhost:8080"
  }
}
```

### Assets not loading
- Verify yukon directory structure
- Check file permissions
- Ensure assets directory exists: `yukon/assets/`
- Verify pack files are valid JSON

## Building for Production

```bash
# Build binary
go build -o tuxedo-core

# Build with optimizations
go build -ldflags="-s -w" -o tuxedo-core

# Cross-compile for different platforms
GOOS=windows GOARCH=amd64 go build -o tuxedo-core.exe
GOOS=linux GOARCH=amd64 go build -o tuxedo-core-linux
GOOS=darwin GOARCH=amd64 go build -o tuxedo-core-mac
```

## Dependencies

- `github.com/gorilla/mux` - HTTP router
- `github.com/gorilla/websocket` - WebSocket support

## Adding New Endpoints

1. Create handler in `handlers/` directory
2. Define route in `main.go`
3. Add CORS and logging middleware
4. Update this README with endpoint documentation

Example:
```go
// handlers/myhandler.go
func MyHandler(w http.ResponseWriter, r *http.Request) {
    // Your handler code
}

// main.go
api.HandleFunc("/my-endpoint", handlers.MyHandler).Methods("GET")
```

## Security Notes

- CORS is wide open by default (`*`) for development
- No authentication/authorization implemented
- File paths are validated to prevent directory traversal
- Suitable for local development only

## Contributing

1. Follow Go coding conventions
2. Add error handling for all file operations
3. Log important operations and errors
4. Update API documentation for new endpoints
5. Test with the tuxedo client

## Related Projects

- **tuxedo**: Frontend scene editor (Vue 3 + TypeScript)
- **yukon**: Yukon game client (Phaser 3)

## License

See main project LICENSE file.
