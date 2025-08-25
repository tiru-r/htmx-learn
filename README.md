# 🚀 Modern HTMX + Go Web Application

A production-ready web application built with **Go 1.25**, **HTMX 2.0**, **PostgreSQL**, and modern security practices. Features comprehensive user management, real-time interactions, and enterprise-grade observability.

[![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![HTMX Version](https://img.shields.io/badge/HTMX-2.0-663399?style=flat-square)](https://htmx.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?style=flat-square&logo=postgresql)](https://postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker)](https://docker.com/)

## ✨ Key Features

### 🏗️ **Modern Architecture**
- **Go 1.25** with standard `net/http` library (no frameworks needed)
- **PostgreSQL 16** with `pgx/v5` driver and connection pooling
- **HTMX 2.0** for dynamic interactions without JavaScript complexity  
- **Templ** for type-safe HTML templates
- **Clean Architecture** with repository pattern and dependency injection

### 🔒 **Enterprise Security**
- **Environment-based Configuration** (no hardcoded credentials)
- **Configurable CORS** with origin validation
- **Comprehensive Security Headers** (CSP, HSTS, X-Frame-Options, etc.)
- **Rate Limiting** with IP-based throttling
- **Input Validation & Sanitization** with XSS protection
- **Structured Logging** with `log/slog` for security auditing

### 🛡️ **Production Resilience**
- **Circuit Breaker Pattern** for database operations
- **Health Check Endpoints** (`/health`, `/health/ready`, `/health/live`)
- **Graceful Shutdown** with configurable timeouts
- **Database Connection Pooling** with retry logic
- **Comprehensive Error Handling** with context propagation

### 📊 **User Management & Features**
- **CRUD Operations** for user management
- **Paginated Results** with efficient database queries
- **Search Functionality** with debounced input
- **Real-time Counter** with optimistic updates
- **Dynamic Content Loading** with HTMX

## 🗂️ **Project Structure**

```
htmx-learn/                    # 2,889 lines of Go code
├── cmd/htmx-learn/           # Application entry point
│   └── main.go               # Server setup, middleware chain, graceful shutdown
├── config/                   # Centralized configuration management
│   └── config.go             # Environment-based config with validation
├── handlers/                 # HTTP request handlers  
│   ├── handlers.go           # Main business logic handlers
│   └── helpers.go            # Template rendering and error handling utilities
├── middleware/               # HTTP middleware stack
│   └── middleware.go         # Security, logging, CORS, rate limiting
├── db/                       # Database layer
│   ├── db.go                 # Connection management & circuit breaker
│   ├── interfaces.go         # Repository interfaces
│   ├── models.go             # Data models and repository implementations
│   ├── pagination.go         # Generic pagination utilities
│   ├── pagination_test.go    # Pagination unit tests
│   └── schema.sql            # Database schema
├── validation/               # Input validation & security
│   ├── validation.go         # User input validation with XSS protection
│   └── validation_test.go    # Validation unit tests
├── circuitbreaker/           # Resilience patterns
│   └── circuitbreaker.go     # Circuit breaker implementation
├── templates/                # Type-safe HTML templates
│   ├── layouts/base.templ    # Base layout with meta tags & scripts
│   ├── pages/                # Full page templates
│   │   ├── home.templ        # Landing page
│   │   ├── counter.templ     # Interactive counter demo
│   │   └── dynamic.templ     # Dynamic content with user management
│   └── components/           # Reusable UI components
│       ├── counter.templ     # Counter widget with HTMX actions
│       ├── dynamic.templ     # User cards, search, time display
│       └── pagination.templ  # Pagination controls
├── static/                   # Static assets
│   ├── css/
│   │   ├── input.css         # Tailwind CSS configuration
│   │   └── output.css        # Generated CSS (auto-generated)
│   └── js/                   # JavaScript assets (if needed)
├── Dockerfile                # Multi-stage production container
├── docker-compose.yml        # Local development environment
├── Taskfile.yml              # Modern task runner configuration
└── go.mod                    # Go module with minimal dependencies
```

## 🔧 **Technology Stack**

### **Backend**
| Component | Version | Purpose |
|-----------|---------|---------|
| **Go** | 1.25 | High-performance backend with standard library |
| **PostgreSQL** | 16+ | Primary database with ACID compliance |
| **pgx** | v5.7.5 | High-performance PostgreSQL driver |
| **Templ** | v0.3.943 | Type-safe HTML templates |

### **Frontend** 
| Component | Version | Purpose |
|-----------|---------|---------|
| **HTMX** | 2.0 | Dynamic interactions without JavaScript |
| **Tailwind CSS** | v4.1.12 | Utility-first CSS framework |
| **Hyperscript** | v0.9.14 | Enhanced HTML interactivity |

### **DevOps & Tools**
| Tool | Purpose |
|------|---------|
| **Task** | Modern task runner (replaces Make) |
| **Air** | Hot reload for Go development |
| **Docker** | Containerization with multi-stage builds |
| **PostgreSQL** | Database with Docker Compose |

## 🚀 **Quick Start**

### **Prerequisites**
- **Go 1.25+** (includes container-aware GOMAXPROCS)
- **Task** (`brew install go-task/tap/go-task` or [install guide](https://taskfile.dev/installation/))
- **Docker & Docker Compose** (for database)

### **Installation**

1. **Clone and install dependencies:**
   ```bash
   git clone <repository-url>
   cd htmx-learn
   task install  # Installs Go tools + Tailwind CLI
   ```

2. **Start database:**
   ```bash
   task db-start  # Starts PostgreSQL in Docker
   ```

3. **Start development environment:**
   ```bash
   task dev  # Starts everything with live reload
   ```

4. **Open your browser:**
   - **App**: http://localhost:8080
   - **Live Reload Proxy**: http://localhost:7331 (recommended for development)

## 🌐 **API Endpoints**

### **Web Pages**
| Route | Method | Description |
|-------|--------|-------------|
| `/` | GET | Landing page with navigation |
| `/counter` | GET | Interactive counter demonstration |
| `/dynamic` | GET | User management with real-time features |

### **API Endpoints**
| Route | Method | Description |
|-------|--------|-------------|
| `/api/time` | GET | Current server time (HTMX demo) |
| `/api/users` | GET | List all users |
| `/api/users` | POST | Create new user |
| `/api/users/{id}` | DELETE | Delete user by ID |
| `/api/users/paginated` | GET | Paginated user list |
| `/api/search` | POST | Search users |
| `/api/search/paginated` | POST | Paginated search results |

### **Counter API**
| Route | Method | Description |
|-------|--------|-------------|
| `/counter/increment` | POST | Increment counter |
| `/counter/decrement` | POST | Decrement counter |
| `/counter/reset` | POST | Reset counter to zero |

### **Health Checks**
| Route | Method | Description |
|-------|--------|-------------|
| `/health` | GET | Comprehensive health check with database status |
| `/health/ready` | GET | Readiness probe for load balancers |
| `/health/live` | GET | Liveness probe for container orchestrators |

## ⚙️ **Configuration**

### **Environment Variables**

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | *required* | PostgreSQL connection string |
| `SECRET_KEY` | *required* | 32+ character secret for security |
| `PORT` | `8080` | Server port |
| `HOST` | `localhost` | Server host |
| `ENVIRONMENT` | `development` | Environment: development/staging/production |

#### **Database Configuration**
| Variable | Default | Description |
|----------|---------|-------------|
| `DB_MAX_CONNECTIONS` | `10` | Maximum database connections |
| `DB_MIN_CONNECTIONS` | `2` | Minimum database connections |
| `DB_CONN_MAX_LIFETIME` | `1h` | Connection maximum lifetime |

#### **Security Configuration**
| Variable | Default | Description |
|----------|---------|-------------|
| `ALLOWED_ORIGINS` | `http://localhost:8080,...` | Comma-separated CORS origins |
| `TRUSTED_PROXIES` | `127.0.0.1,::1` | Trusted proxy IP addresses |
| `RATE_LIMIT` | `100` | Requests per minute per IP |
| `RATE_LIMIT_WINDOW` | `1m` | Rate limiting time window |
| `RATE_LIMIT_BURST` | `20` | Burst capacity for rate limiting |

#### **Logging Configuration**
| Variable | Default | Description |
|----------|---------|-------------|
| `LOG_LEVEL` | `info` | Log level: debug/info/warn/error |
| `LOG_FORMAT` | `json` | Log format: json/text |

### **Example .env file**
```env
DATABASE_URL=postgres://user:password@localhost:5432/htmx_learn?sslmode=disable
SECRET_KEY=your-super-secret-key-32-chars-minimum
ENVIRONMENT=development
LOG_LEVEL=debug
ALLOWED_ORIGINS=http://localhost:8080,https://localhost:8080
RATE_LIMIT=100
```

## 🛠️ **Development Workflow**

### **Available Tasks**
```bash
task                    # Show all available tasks
task install           # Install all dependencies (Go tools + Tailwind)
task dev               # Start full development environment with live reload
task build             # Build the application
task run               # Build and run the application
task test              # Run all tests
task fmt               # Format Go code and Templ templates
task lint              # Run linters and static analysis
task clean             # Clean build artifacts

# Database operations
task db-start          # Start PostgreSQL in Docker
task db-stop           # Stop PostgreSQL
task db-reset          # Reset database (destroys all data!)
task db-logs           # View database logs
task db-shell          # Connect to database shell
task db-backup         # Backup database to backups/ directory

# CSS & Templates
task css               # Build CSS with Tailwind
task css-watch         # Watch and rebuild CSS
task generate          # Generate Templ templates
```

### **Development Features**
- **🔄 Hot Reload**: Go code changes trigger automatic rebuilds
- **📝 Template Watching**: Templ templates auto-regenerate
- **🎨 CSS Watching**: Tailwind CSS rebuilds on file changes  
- **📡 Live Reload Proxy**: Browser automatically refreshes on changes
- **🚀 Zero Config**: Works out of the box with sane defaults

## 🐳 **Production Deployment**

### **Docker Deployment**
```bash
# Full stack with PostgreSQL
docker-compose up --build -d

# Just the application (external database)
docker build -t htmx-learn .
docker run -p 8080:8080 \
  -e DATABASE_URL="postgres://user:pass@db:5432/htmx_learn?sslmode=disable" \
  -e SECRET_KEY="your-production-secret-key-32-chars" \
  -e ENVIRONMENT="production" \
  htmx-learn
```

### **Manual Deployment**
```bash
# Build for production
task prod-build

# Set environment variables
export DATABASE_URL="postgres://..."
export SECRET_KEY="..."
export ENVIRONMENT="production"

# Run binary
./tmp/htmx-learn
```

### **Production Checklist**
- ✅ Set strong `SECRET_KEY` (32+ characters)
- ✅ Configure `ALLOWED_ORIGINS` for your domains
- ✅ Set `ENVIRONMENT=production`
- ✅ Configure proper `DATABASE_URL` with connection pooling
- ✅ Set up monitoring for `/health` endpoints
- ✅ Configure reverse proxy (nginx/Caddy) for HTTPS
- ✅ Set resource limits in container orchestrator

## 📊 **Monitoring & Observability**

### **Health Check Endpoints**
```bash
# Comprehensive health check
curl http://localhost:8080/health
# Returns: {"status":"healthy","timestamp":"...","checks":{"database":{"status":"healthy","latency":"2ms"}}}

# Kubernetes readiness probe  
curl http://localhost:8080/health/ready

# Kubernetes liveness probe
curl http://localhost:8080/health/live
```

### **Structured Logging**
All logs are structured JSON for easy parsing:
```json
{
  "time": "2025-01-XX",
  "level": "INFO", 
  "msg": "HTTP Request",
  "method": "GET",
  "path": "/api/users", 
  "status": 200,
  "duration": "15ms",
  "remote_addr": "127.0.0.1"
}
```

### **Circuit Breaker Monitoring**
Database operations are protected by circuit breakers that log state transitions and provide statistics for monitoring external dependencies.

## 🧪 **Testing**

```bash
# Run all tests
task test

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -v ./validation
```

**Test Coverage:**
- ✅ **Database Pagination**: Comprehensive pagination logic testing
- ✅ **Input Validation**: XSS protection and sanitization tests  
- ✅ **Circuit Breaker**: Resilience pattern testing
- 🔄 **Integration Tests**: Coming soon

## 🔒 **Security Features**

This application implements defense-in-depth security:

### **Input Security**
- 🛡️ **XSS Protection**: All user inputs sanitized and validated
- 🔍 **SQL Injection Prevention**: Parameterized queries with pgx
- 📝 **Input Validation**: Comprehensive validation with custom error types

### **HTTP Security**  
- 🌐 **Secure CORS**: Configurable origin validation (no wildcards)
- 🔒 **Security Headers**: CSP, HSTS, X-Frame-Options, X-XSS-Protection
- 🚦 **Rate Limiting**: IP-based throttling with configurable limits
- 🔐 **HTTPS Ready**: Security headers configured for TLS

### **Application Security**
- 🗝️ **Environment Config**: No secrets in code
- 📊 **Structured Logging**: Security event auditing
- 🛡️ **Circuit Breakers**: Prevent cascade failures
- 🐳 **Container Security**: Non-root user, minimal attack surface

## 🤝 **Contributing**

1. **Code Style**: Follow Go conventions and run `task fmt`
2. **Testing**: Ensure tests pass with `task test` 
3. **Security**: No secrets in code, validate all inputs
4. **Documentation**: Update README for significant changes

## 📚 **Architecture Decisions**

### **Why Go Standard Library?**
- **Performance**: Direct control over HTTP handling
- **Simplicity**: No framework lock-in or magic
- **Reliability**: Battle-tested standard library
- **Security**: Smaller attack surface

### **Why HTMX?**
- **Simplicity**: HTML-driven interactions
- **Performance**: Minimal JavaScript overhead  
- **SEO Friendly**: Server-side rendering
- **Progressive Enhancement**: Works without JavaScript

### **Why PostgreSQL?**
- **ACID Compliance**: Data integrity guarantees
- **Performance**: Excellent for read-heavy workloads
- **JSON Support**: Flexible data structures when needed
- **Ecosystem**: Rich tooling and extensions

## 📄 **License**

MIT License - see LICENSE file for details.

---

**Built with ❤️ using Go 1.25, HTMX 2.0, and modern web standards.**