# Resume API - System Design Document

## Overview
A Go-based REST API that serves resume/CV data from a PostgreSQL database, designed as a personal portfolio project and learning vehicle for Go development best practices.

## Architecture Principles

### Clean Architecture
The project follows Clean Architecture principles with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                    HTTP Layer (Handlers)                    │
├─────────────────────────────────────────────────────────────┤
│                   Business Layer (Services)                 │
├─────────────────────────────────────────────────────────────┤
│                  Data Layer (Repositories)                  │
├─────────────────────────────────────────────────────────────┤
│                     Database (PostgreSQL)                   │
└─────────────────────────────────────────────────────────────┘
```

### Dependency Flow
- **Handlers** → **Services** → **Repositories** → **Database**
- Dependencies flow inward through interfaces
- Each layer only knows about the layer directly below it

## System Components

### 1. HTTP Layer (`internal/handlers/`)
**Responsibility**: Handle HTTP requests and responses

- **ResumeHandler**: Main handler for resume endpoints
- **HealthHandler**: Health check endpoints
- **Middleware**: CORS, logging, recovery, timeouts

**Key Patterns**:
- Dependency injection of services
- Request validation
- Error handling with appropriate HTTP status codes
- JSON marshaling/unmarshaling

### 2. Business Layer (`internal/services/`)
**Responsibility**: Business logic and data orchestration

- **ResumeService**: Aggregates and processes resume data
- **HealthService**: System health checks

**Key Patterns**:
- Service interfaces for testability
- Business rule enforcement
- Error wrapping and context propagation

### 3. Data Layer (`internal/repository/`)
**Responsibility**: Data access and persistence

- **ProfileRepository**: Personal information CRUD
- **ExperienceRepository**: Work experience CRUD  
- **SkillRepository**: Skills and categories CRUD
- **AchievementRepository**: Achievements CRUD
- **EducationRepository**: Education and certifications CRUD

**Key Patterns**:
- Repository pattern with interfaces
- SQL query encapsulation
- Connection pool management
- Context-aware operations

### 4. Models (`internal/models/`)
**Responsibility**: Data structures and domain objects

**Domain Models**:
```go
type Profile struct {
    ID       int       `json:"id" db:"id"`
    Name     string    `json:"name" db:"name"`
    Title    string    `json:"title" db:"title"`
    Email    string    `json:"email" db:"email"`
    Phone    string    `json:"phone" db:"phone"`
    Location string    `json:"location" db:"location"`
    LinkedIn string    `json:"linkedin" db:"linkedin"`
    Summary  string    `json:"summary" db:"summary"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Experience struct {
    ID          int       `json:"id" db:"id"`
    Company     string    `json:"company" db:"company"`
    Position    string    `json:"position" db:"position"`
    StartDate   time.Time `json:"start_date" db:"start_date"`
    EndDate     *time.Time `json:"end_date" db:"end_date"`
    Description string    `json:"description" db:"description"`
    Order       int       `json:"order" db:"order"`
}

type Skill struct {
    ID       int    `json:"id" db:"id"`
    Category string `json:"category" db:"category"`
    Name     string `json:"name" db:"name"`
    Level    string `json:"level" db:"level"`
    Order    int    `json:"order" db:"order"`
}
```

## Database Design

### Schema Overview
```sql
-- Core profile information
CREATE TABLE profiles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(50),
    location VARCHAR(255),
    linkedin VARCHAR(255),
    summary TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Work experience
CREATE TABLE experiences (
    id SERIAL PRIMARY KEY,
    company VARCHAR(255) NOT NULL,
    position VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    description TEXT,
    order_index INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Skills with categories
CREATE TABLE skills (
    id SERIAL PRIMARY KEY,
    category VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    level VARCHAR(50),
    order_index INTEGER DEFAULT 0
);

-- Achievements
CREATE TABLE achievements (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    order_index INTEGER DEFAULT 0
);

-- Education and certifications
CREATE TABLE education (
    id SERIAL PRIMARY KEY,
    institution VARCHAR(255) NOT NULL,
    degree VARCHAR(255),
    field VARCHAR(255),
    year INTEGER,
    description TEXT,
    type VARCHAR(50) CHECK (type IN ('education', 'certification')),
    order_index INTEGER DEFAULT 0
);
```

### Indexing Strategy
```sql
-- Performance indexes
CREATE INDEX idx_experiences_company ON experiences(company);
CREATE INDEX idx_experiences_order ON experiences(order_index);
CREATE INDEX idx_skills_category ON skills(category);
CREATE INDEX idx_skills_order ON skills(category, order_index);
CREATE INDEX idx_education_type ON education(type);
CREATE INDEX idx_education_order ON education(type, order_index);
```

## API Design

### Base URL
- Development: `http://localhost:8080/api/v1`
- Production: `https://api.neilmulder.dev/api/v1`

### Endpoints

#### Profile Information
```
GET /api/v1/profile
Response: Profile object with personal information
```

#### Work Experience
```
GET /api/v1/experiences
Response: Array of Experience objects, ordered by date (newest first)

Query Parameters:
- limit: Maximum number of results (default: 10)
- company: Filter by company name
```

#### Skills
```
GET /api/v1/skills
Response: Object with skills grouped by category

Query Parameters:
- category: Filter by skill category
- level: Filter by skill level
```

#### Achievements
```
GET /api/v1/achievements
Response: Array of Achievement objects
```

#### Education
```
GET /api/v1/education
Response: Object with education and certifications

Query Parameters:
- type: Filter by 'education' or 'certification'
```

#### Health Check
```
GET /health
Response: System health status including database connectivity
```

### Error Handling
Standard HTTP status codes with consistent error response format:

```json
{
    "error": {
        "code": "RESOURCE_NOT_FOUND",
        "message": "The requested resource was not found",
        "details": "Profile with ID 123 does not exist"
    },
    "timestamp": "2025-01-11T10:30:00Z"
}
```

## Go-Specific Design Patterns

### 1. Interface-Based Architecture
```go
// Service interfaces for dependency injection
type ResumeService interface {
    GetProfile(ctx context.Context) (*models.Profile, error)
    GetExperiences(ctx context.Context, filters ExperienceFilters) ([]*models.Experience, error)
    GetSkills(ctx context.Context, category string) (map[string][]*models.Skill, error)
}

// Repository interfaces for data access
type ProfileRepository interface {
    GetProfile(ctx context.Context) (*models.Profile, error)
    UpdateProfile(ctx context.Context, profile *models.Profile) error
}
```

### 2. Context-Aware Operations
All operations accept `context.Context` for:
- Request timeouts
- Cancellation handling
- Request tracing
- Metadata propagation

### 3. Error Wrapping
```go
func (r *profileRepository) GetProfile(ctx context.Context) (*models.Profile, error) {
    var profile models.Profile
    err := r.db.GetContext(ctx, &profile, "SELECT * FROM profiles LIMIT 1")
    if err != nil {
        return nil, fmt.Errorf("failed to get profile: %w", err)
    }
    return &profile, nil
}
```

### 4. Configuration Management
```go
type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Log      LogConfig      `mapstructure:"log"`
}

// Environment-based configuration with defaults
func LoadConfig() (*Config, error) {
    viper.SetDefault("server.port", 8080)
    viper.SetDefault("server.timeout", "30s")
    // ... load from environment
}
```

## Deployment Architecture

### Local Development
```
┌─────────────┐    ┌─────────────┐
│   Go API    │    │ PostgreSQL  │
│  (Port 8080)│◄──►│ (Port 5432) │
└─────────────┘    └─────────────┘
```

### Kubernetes Production
```
┌─────────────────────────────────────────────────────────┐
│                       Ingress                           │
│                  (TLS Termination)                      │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                   Service                               │
│              (Load Balancer)                            │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                  Deployment                             │
│              (Multiple Pods)                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐     │
│  │   Pod 1     │  │   Pod 2     │  │   Pod N     │     │
│  │   Go API    │  │   Go API    │  │   Go API    │     │
│  └─────────────┘  └─────────────┘  └─────────────┘     │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│               PostgreSQL Service                        │
│              (External Database)                        │
└─────────────────────────────────────────────────────────┘
```

## Performance Considerations

### Connection Pooling
- PostgreSQL connection pool: 10-25 connections
- Connection lifetime: 1 hour
- Idle timeout: 10 minutes

### Caching Strategy (Future)
- Redis for frequently accessed data
- Cache TTL: 15 minutes for profile data
- Cache invalidation on data updates

### Response Times
- Target: < 100ms for all endpoints
- Timeout: 30 seconds for requests
- Database query timeout: 10 seconds

## Security Considerations

### Data Protection
- No authentication required (public resume data)
- Input validation on all parameters
- SQL injection prevention via parameterized queries
- XSS prevention via JSON encoding

### Infrastructure Security
- TLS encryption in production
- Rate limiting: 100 requests/minute per IP
- Health check endpoint access control
- Container security best practices

## Monitoring and Observability

### Logging
- Structured JSON logging
- Request ID tracing
- Error logging with stack traces
- Performance metrics

### Health Checks
- Database connectivity
- Application health
- Dependency status

### Future Metrics
- Request duration histograms
- Error rate monitoring
- Database connection pool metrics

## Testing Strategy

### Unit Tests
- Repository layer: Mock database connections
- Service layer: Mock repositories
- Handler layer: Mock services and HTTP testing

### Integration Tests
- End-to-end API testing
- Database integration with test containers
- Configuration testing

### Test Coverage
- Target: 80%+ code coverage
- Critical path: 95%+ coverage

This design provides a solid foundation for learning Go while building a production-ready API that showcases modern development practices.