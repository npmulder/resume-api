# Resume API - Task Tracking

## Project Overview
Building a Go-based REST API to serve resume data from PostgreSQL, with Kubernetes deployment for homelab infrastructure. This project serves as a learning vehicle for Go development best practices and enterprise patterns.

## Task Status Legend
- 🔴 **Not Started** - Task not yet begun
- 🟡 **In Progress** - Currently working on task
- 🟢 **Completed** - Task finished and verified
- 🔵 **Review** - Task completed, needs review/testing
- ⚪ **Blocked** - Task blocked by dependency or issue

---

## Phase 1: Project Foundation ✅

### 1.1 Project Setup
- [🟢] Update .gitignore to exclude resume folder
- [🟢] Initialize Go module with proper naming
- [🟢] Create comprehensive directory structure
- [🟢] Create task tracking document (this file)

### 1.2 Documentation Framework
- [🔴] Create system design document
- [🔴] Create development guide
- [🔴] Update CLAUDE.md with project specifics

---

## Phase 2: Core Development Setup

### 2.1 Database Design (Est: 2-3 hours)
- [🔴] **Design PostgreSQL schema**
  - Profile table (personal info, summary)
  - Experiences table (work history, dates, descriptions)
  - Skills table (categorized skills with JSONB)
  - Achievements table 
  - Education table
  - Certifications table
  - **Learning**: PostgreSQL data types, indexing, relationships

- [🔴] **Create database migrations**
  - Install golang-migrate tool
  - Write up/down migration files
  - **Learning**: Database migration patterns in Go

- [🔴] **Create seed data script**
  - Parse resume data into database format
  - Create SQL insert scripts
  - **Learning**: Data transformation, SQL scripting

### 2.2 Go Project Dependencies (Est: 1 hour)
- [🔴] **Add core dependencies to go.mod**
  - HTTP router (gin or echo)
  - PostgreSQL driver (pgx)
  - Configuration (viper)
  - Logging (slog or logrus)
  - Testing (testify)
  - **Learning**: Go dependency management, popular libraries

### 2.3 Configuration System (Est: 1-2 hours)
- [🔴] **Create config package**
  - Environment variable loading
  - Configuration validation
  - Default values
  - **Learning**: Go struct tags, environment handling

- [🔴] **Create .env.example file**
  - Database connection settings
  - Server configuration
  - **Learning**: Environment-based configuration

---

## Phase 3: Database Layer

### 3.1 Database Connection (Est: 2 hours)
- [🔴] **Implement database package**
  - Connection pooling with pgx
  - Health check functions
  - Migration runner
  - **Learning**: Connection pooling, database health checks

### 3.2 Repository Pattern (Est: 3-4 hours)
- [🔴] **Create repository interfaces**
  - ProfileRepository interface
  - ExperienceRepository interface
  - SkillRepository interface
  - **Learning**: Go interfaces, dependency injection

- [🔴] **Implement PostgreSQL repositories**
  - CRUD operations for each entity
  - Error handling patterns
  - Context usage
  - **Learning**: SQL in Go, error wrapping, context patterns

### 3.3 Models and DTOs (Est: 2 hours)
- [🔴] **Create domain models**
  - Profile, Experience, Skill structs
  - JSON marshaling tags
  - Validation tags
  - **Learning**: Struct tags, JSON marshaling, validation

---

## Phase 4: Business Logic Layer

### 4.1 Service Layer (Est: 2-3 hours)
- [🔴] **Create service interfaces**
  - ResumeService interface
  - Business logic abstraction
  - **Learning**: Service layer patterns, business logic separation

- [🔴] **Implement service logic**
  - Data aggregation
  - Business rules
  - Error handling
  - **Learning**: Clean architecture, error handling strategies

---

## Phase 5: HTTP Layer

### 5.1 HTTP Handlers (Est: 3-4 hours)
- [🔴] **Create handler package**
  - REST endpoint handlers
  - Request/response DTOs
  - HTTP status codes
  - **Learning**: HTTP handling in Go, REST conventions

- [🔴] **Implement endpoints**
  - GET /api/v1/profile
  - GET /api/v1/experiences
  - GET /api/v1/skills
  - GET /api/v1/achievements
  - GET /api/v1/education
  - GET /health
  - **Learning**: RESTful API design, HTTP best practices

### 5.2 Middleware (Est: 2-3 hours)
- [🔴] **Create middleware package**
  - CORS middleware
  - Logging middleware
  - Recovery middleware
  - Request timeout middleware
  - **Learning**: Middleware patterns, HTTP middleware chains

### 5.3 Main Application (Est: 1-2 hours)
- [🔴] **Create cmd/api/main.go**
  - Dependency injection
  - Graceful shutdown
  - Signal handling
  - **Learning**: Application lifecycle, graceful shutdown patterns

---

## Phase 6: Testing

### 6.1 Unit Tests (Est: 4-5 hours)
- [🔴] **Repository tests**
  - Table-driven tests
  - Mock database connections
  - **Learning**: Table-driven testing, mocking in Go

- [🔴] **Service tests**
  - Business logic testing
  - Mock repositories
  - **Learning**: Interface mocking, testify/mock

- [🔴] **Handler tests**
  - HTTP testing
  - Request/response validation
  - **Learning**: HTTP testing in Go, httptest package

### 6.2 Integration Tests (Est: 3-4 hours)
- [🔴] **Database integration tests**
  - Test containers or test database
  - End-to-end data flow
  - **Learning**: Integration testing strategies

---

## Phase 7: DevOps and Deployment

### 7.1 Containerization (Est: 2-3 hours)
- [🔴] **Create Dockerfile**
  - Multi-stage build
  - Security best practices
  - Non-root user
  - **Learning**: Docker best practices, Go containerization

### 7.2 Kubernetes Deployment (Est: 3-4 hours)
- [🔴] **Create Kubernetes manifests**
  - Deployment configuration
  - Service configuration
  - ConfigMap and Secrets
  - Health probes
  - **Learning**: Kubernetes deployment patterns

### 7.3 Local Development (Est: 1-2 hours)
- [🔴] **Create docker-compose.yml**
  - PostgreSQL container
  - API container
  - Development environment
  - **Learning**: Local development with containers

---

## Phase 8: Advanced Features (Stretch Goals)

### 8.1 Observability (Est: 2-3 hours)
- [🔴] **Add structured logging**
  - Request tracing
  - Error logging
  - Performance metrics
  - **Learning**: Structured logging, observability patterns

### 8.2 Performance (Est: 2-3 hours)
- [🔴] **Add caching layer**
  - Redis integration
  - Cache strategies
  - **Learning**: Caching patterns, Redis in Go

### 8.3 Security (Est: 2-3 hours)
- [🔴] **Security middleware**
  - Rate limiting
  - Input validation
  - Security headers
  - **Learning**: API security, rate limiting

---

## Learning Objectives Tracker

### Go Language Features Covered
- [ ] Interfaces and dependency injection
- [ ] Context package usage
- [ ] Error handling and wrapping
- [ ] Struct tags and JSON marshaling
- [ ] Table-driven testing
- [ ] Goroutines and channels (if needed)
- [ ] HTTP server patterns
- [ ] Configuration management

### Enterprise Patterns
- [ ] Clean Architecture
- [ ] Repository pattern
- [ ] Service layer pattern
- [ ] Dependency injection
- [ ] Middleware patterns
- [ ] Graceful shutdown
- [ ] Health checks

### DevOps Skills
- [ ] Docker containerization
- [ ] Kubernetes deployment
- [ ] Configuration management
- [ ] Database migrations
- [ ] CI/CD preparation

---

## Current Sprint: Foundation Phase
**Sprint Goal**: Complete project setup and documentation framework

**Active Tasks**:
- Create system design document
- Create development guide  
- Update CLAUDE.md

**Next Sprint**: Database design and core dependencies