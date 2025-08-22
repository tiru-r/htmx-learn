# HTMX + Go Web Application

A modern web application built with Go's standard `net/http` library, HTMX for dynamic interactions, Tailwind CSS v4.1 for styling, and Templ for type-safe HTML templates.

## Features

- **Go Standard Library**: Built with `net/http` for maximum performance and minimal dependencies
- **HTMX**: Dynamic interactions without writing JavaScript
- **Tailwind CSS v4.1**: Zero-config utility-first CSS framework with standalone CLI
- **Templ**: Type-safe HTML templates for Go
- **SQLite**: Modern SQLite database with WAL mode for persistence
- **Live Reload**: Hot reload for Go code, templates, and CSS during development
- **Task Runner**: Automated development workflow with Taskfile
- **Production Ready**: Dockerized with graceful shutdown

## Technology Stack

### Backend
- Go with `net/http` standard library
- SQLite with modernc.org/sqlite driver
- Templ for type-safe HTML templates
- Air for hot reloading during development
- Taskfile for task automation

### Frontend
- HTMX for dynamic interactions
- Tailwind CSS v4.1 with standalone CLI (zero-config)
- Hyperscript for enhanced interactivity
- Live reload via Templ proxy

## Quick Start

### Prerequisites

- Go 1.25+ (includes container-aware GOMAXPROCS)
- [Task](https://taskfile.dev) (optional but recommended)
- curl (for downloading Tailwind CLI)

**No Node.js required!** Uses Tailwind standalone CLI.

### Installation

1. Install dependencies:
   ```bash
   task install
   # or manually:
   # go install github.com/air-verse/air@latest
   # go install github.com/a-h/templ/cmd/templ@latest
   # curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64
   # chmod +x tailwindcss-linux-x64 && mv tailwindcss-linux-x64 tailwindcss
   ```

2. Start development environment:
   ```bash
   task dev
   ```

3. Open your browser to http://localhost:7331 (Templ proxy with live reload)

The application will be running with:
- Go server on http://localhost:8080
- Templ live reload proxy on http://localhost:7331
- Hot reload for Go code changes
- Auto-regeneration of Templ templates
- CSS watching and rebuilding

## Development Workflow

### Available Tasks

```bash
task                 # Show available tasks
task install         # Install all dependencies
task dev             # Start full development environment
task build           # Build the application
task run             # Run the application
task clean           # Clean build artifacts
task test            # Run tests
task fmt             # Format code
task lint            # Run linters
task prod-build      # Build for production
task db-init         # Initialize database directory
task db-reset        # Reset database (destroys all data)
task db-backup       # Backup database
```

### Development Commands

```bash
# Start development with live reload
task dev

# Build CSS manually
task css

# Watch CSS changes
task css-watch

# Generate Templ templates
task generate
```

## Project Structure

```
.
├── cmd/
│   └── htmx-learn/          # Application entry point
│       └── main.go
├── handlers/                # HTTP handlers
│   └── handlers.go
├── middleware/              # HTTP middleware
│   └── middleware.go
├── templates/               # Templ templates
│   ├── layouts/
│   │   └── base.templ
│   ├── components/
│   │   ├── counter.templ
│   │   └── dynamic.templ
│   └── pages/
│       ├── home.templ
│       ├── counter.templ
│       └── dynamic.templ
├── static/                  # Static assets
│   ├── css/
│   │   ├── input.css        # Tailwind CSS input
│   │   └── output.css       # Generated CSS (gitignored)
│   └── js/
├── tmp/                     # Build output (gitignored)
├── Taskfile.yml             # Task definitions
├── .air.toml                # Air configuration
├── package.json             # NPM dependencies
├── Dockerfile               # Production container
└── README.md
```

## Examples

The application includes several HTMX examples:

### Counter Example (`/counter`)
- HTMX GET/POST requests
- Target swapping
- Server-side state management

### Dynamic Content (`/dynamic`)
- Real-time content loading
- User management with CRUD operations
- Search with debounced input
- Loading indicators

## Production Deployment

### Docker

```bash
# Build production image
docker build -t htmx-learn .

# Run container
docker run -p 8080:8080 htmx-learn
```

### Manual Build

```bash
# Build for production
task prod-build

# Run binary
./tmp/htmx-learn
```

## Configuration

### Environment Variables

- `PORT`: Server port (default: 8080)

### Development Tools

- **Air**: Configured inline via Taskfile for Go hot reloading
- **Templ**: Live reload proxy for browser auto-refresh
- **Tailwind**: Zero-config setup with custom components

## Contributing

1. Follow the existing code style
2. Run `task fmt` before committing
3. Ensure tests pass with `task test`
4. Update documentation as needed

## License

MIT License