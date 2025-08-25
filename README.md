# HTMX + Go Web Application

A modern web application built with Go 1.25's standard `net/http` library, HTMX 2.0.6 for dynamic interactions, Tailwind CSS v4.1.12 for styling, and Templ v0.3.943 for type-safe HTML templates.

## Features

- **Go Standard Library**: Built with `net/http` for maximum performance and minimal dependencies
- **HTMX 2.0.6**: Dynamic interactions without writing JavaScript
- **Tailwind CSS v4.1.12**: Zero-config utility-first CSS framework with standalone CLI
- **Templ v0.3.943**: Type-safe HTML templates for Go
- **PostgreSQL**: Modern PostgreSQL database with connection pooling
- **Hyperscript v0.9.14**: Enhanced HTML interactivity
- **Live Reload**: Hot reload for Go code, templates, and CSS during development
- **Task Runner v3**: Automated development workflow with modern Taskfile
- **Production Ready**: Dockerized with graceful shutdown

## Technology Stack

### Backend
- Go 1.25 with `net/http` standard library
- PostgreSQL with jackc/pgx/v5 driver and connection pooling
- Templ v0.3.943 for type-safe HTML templates
- Air for hot reloading during development
- Taskfile v3 for task automation

### Frontend
- HTMX 2.0.6 for dynamic interactions
- Tailwind CSS v4.1.12 with standalone CLI (zero-config)
- Hyperscript v0.9.14 for enhanced interactivity
- Live reload via Templ proxy

## Quick Start

### Prerequisites

- Go 1.25+ (includes container-aware GOMAXPROCS)
- PostgreSQL 13+ (or Docker/Docker Compose)
- [Task](https://taskfile.dev) (optional but recommended)
- curl (for downloading Tailwind CLI)

**No Node.js required!** Uses Tailwind v4.1.12 standalone CLI.

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

2. Start PostgreSQL database:
   ```bash
   # Using Docker Compose (recommended)
   docker-compose up postgres -d
   
   # Or start your local PostgreSQL instance
   # Make sure database 'htmx_learn' exists
   ```

3. Set database URL (optional, uses default if not set):
   ```bash
   export DATABASE_URL="postgres://user:password@localhost:5432/htmx_learn?sslmode=disable"
   ```

4. Start development environment:
   ```bash
   task dev
   ```

5. Open your browser to http://localhost:7331 (Templ proxy with live reload)

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
# Using Docker Compose (includes PostgreSQL)
docker-compose up --build

# Or build and run manually
docker build -t htmx-learn .
docker run -p 8080:8080 -e DATABASE_URL="postgres://user:password@host.docker.internal:5432/htmx_learn?sslmode=disable" htmx-learn
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

- `DATABASE_URL`: PostgreSQL connection string (default: postgres://user:password@localhost:5432/htmx_learn?sslmode=disable)
- `PORT`: Server port (default: 8080)

### Development Tools

- **Air**: Configured inline via Taskfile for Go hot reloading
- **Templ v0.3.943**: Live reload proxy for browser auto-refresh
- **Tailwind v4.1.12**: Zero-config setup with custom components
- **Task v3**: Modern task runner with YAML configuration

## Contributing

1. Follow the existing code style
2. Run `task fmt` before committing
3. Ensure tests pass with `task test`
4. Update documentation as needed

## License

MIT License