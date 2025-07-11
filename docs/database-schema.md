# Database Schema Documentation

## Overview
PostgreSQL database schema for the Resume API, designed to store comprehensive resume/CV data with flexibility and performance in mind.

## Schema Design Principles

### 1. Normalization
- Each entity (profile, experience, skills, etc.) has its own table
- Minimal data duplication
- Clear relationships between entities

### 2. Flexibility
- JSONB columns for semi-structured data (technologies, features)
- TEXT arrays for lists (highlights, key_features)
- Optional fields for varying data completeness

### 3. Performance
- Strategic indexing on commonly queried fields
- Composite indexes for multi-column queries
- GIN indexes for JSONB and array columns

### 4. Data Integrity
- Check constraints for enumerated values
- Unique constraints where appropriate
- NOT NULL constraints for required fields

## Tables

### profiles
Core personal information and summary.

```sql
CREATE TABLE profiles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(50),
    location VARCHAR(255),
    linkedin VARCHAR(255),
    github VARCHAR(255),
    summary TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Key Features:**
- Single profile record (only one resume owner)
- Email uniqueness enforced
- Automatic timestamp management
- Support for multiple contact methods

**Indexes:**
- `idx_profiles_email` - Fast email lookups

### experiences
Work history and employment details.

```sql
CREATE TABLE experiences (
    id SERIAL PRIMARY KEY,
    company VARCHAR(255) NOT NULL,
    position VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE, -- NULL indicates current position
    description TEXT,
    highlights TEXT[], -- Array of key achievements
    order_index INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Key Features:**
- Open-ended employment (NULL end_date for current roles)
- Highlights stored as PostgreSQL array
- Manual ordering support
- Date range queries supported

**Indexes:**
- `idx_experiences_company` - Company-based filtering
- `idx_experiences_dates` - Date range queries (DESC order)
- `idx_experiences_order` - Display order
- `idx_experiences_current` - Quick current position lookup

### skills
Technical and professional skills with categorization.

```sql
CREATE TABLE skills (
    id SERIAL PRIMARY KEY,
    category VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    level VARCHAR(50), -- beginner, intermediate, advanced, expert
    years_experience INTEGER,
    order_index INTEGER DEFAULT 0,
    is_featured BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(category, name)
);
```

**Key Features:**
- Categorized skills (Languages, Databases, etc.)
- Skill level tracking with validation
- Years of experience quantification
- Featured skills for highlighting
- No duplicate skills per category

**Indexes:**
- `idx_skills_category` - Category-based grouping
- `idx_skills_category_order` - Ordered skills within category
- `idx_skills_featured` - Quick featured skills access
- `idx_skills_level` - Skill level filtering

**Constraints:**
- `chk_skill_level` - Valid skill levels only

### achievements
Key accomplishments and metrics.

```sql
CREATE TABLE achievements (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100), -- performance, security, leadership, etc.
    impact_metric VARCHAR(255), -- "30% reduction", "40% improvement"
    year_achieved INTEGER,
    order_index INTEGER DEFAULT 0,
    is_featured BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Key Features:**
- Quantifiable impact metrics
- Categorized achievements
- Year-based tracking
- Featured achievements for prominence

**Indexes:**
- `idx_achievements_category` - Category grouping
- `idx_achievements_year` - Chronological ordering
- `idx_achievements_order` - Manual ordering
- `idx_achievements_featured` - Featured achievements

### education
Education and professional certifications.

```sql
CREATE TABLE education (
    id SERIAL PRIMARY KEY,
    institution VARCHAR(255) NOT NULL,
    degree_or_certification VARCHAR(255) NOT NULL,
    field_of_study VARCHAR(255),
    year_completed INTEGER,
    year_started INTEGER,
    description TEXT,
    type VARCHAR(50) NOT NULL CHECK (type IN ('education', 'certification')),
    status VARCHAR(50) DEFAULT 'completed' CHECK (status IN ('completed', 'in_progress', 'planned')),
    credential_id VARCHAR(255),
    credential_url VARCHAR(500),
    expiry_date DATE,
    order_index INTEGER DEFAULT 0,
    is_featured BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Key Features:**
- Unified education and certification storage
- Progress tracking (completed, in_progress, planned)
- Credential verification support
- Expiration tracking for certifications
- Date range support

**Indexes:**
- `idx_education_type` - Education vs. certification separation
- `idx_education_year` - Chronological ordering
- `idx_education_institution` - Institution-based queries
- `idx_education_status` - Status filtering
- `idx_education_order` - Type-specific ordering
- `idx_education_featured` - Featured items

### projects
Notable projects and portfolio items.

```sql
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    short_description VARCHAR(500),
    technologies JSONB, -- Array of technologies
    github_url VARCHAR(500),
    demo_url VARCHAR(500),
    start_date DATE,
    end_date DATE, -- NULL for ongoing
    status VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active', 'completed', 'archived', 'planned')),
    is_featured BOOLEAN DEFAULT FALSE,
    order_index INTEGER DEFAULT 0,
    key_features TEXT[],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Key Features:**
- JSONB for flexible technology storage
- Multiple URL types (GitHub, demo)
- Project lifecycle tracking
- Key features as array
- Status-based filtering

**Indexes:**
- `idx_projects_status` - Status filtering
- `idx_projects_dates` - Date-based ordering
- `idx_projects_featured` - Featured projects
- `idx_projects_order` - Manual ordering
- `idx_projects_technologies` - GIN index for JSONB search

## Common Patterns

### Timestamps
All tables include `created_at` and `updated_at` with automatic maintenance via triggers.

### Ordering
Most tables include `order_index` for manual display ordering.

### Featured Items
Boolean `is_featured` flags for highlighting important items.

### Status Tracking
Where applicable, status enums track item lifecycle.

## Query Examples

### Get Current Position
```sql
SELECT company, position, start_date 
FROM experiences 
WHERE end_date IS NULL 
ORDER BY start_date DESC 
LIMIT 1;
```

### Skills by Category
```sql
SELECT category, array_agg(name ORDER BY order_index) as skills
FROM skills 
GROUP BY category
ORDER BY category;
```

### Featured Items Across Tables
```sql
-- Featured skills
SELECT 'skill' as type, name as title, category as detail FROM skills WHERE is_featured = TRUE
UNION ALL
-- Featured achievements  
SELECT 'achievement' as type, title, category as detail FROM achievements WHERE is_featured = TRUE
UNION ALL
-- Featured projects
SELECT 'project' as type, name as title, status as detail FROM projects WHERE is_featured = TRUE;
```

### Technology Search in Projects
```sql
SELECT name, description, technologies
FROM projects 
WHERE technologies @> '["Kubernetes"]'
AND status = 'active';
```

## Performance Considerations

### Index Strategy
- Primary keys (auto-indexed)
- Foreign keys (when added)
- Frequently filtered columns
- Sort columns (with DESC where appropriate)
- JSONB columns with GIN indexes
- Partial indexes for boolean flags

### Query Optimization
- Use specific indexes for common queries
- Leverage partial indexes for sparse data
- Consider covering indexes for read-heavy operations

## Migration Strategy

### Migration Files
- Sequential numbering (001, 002, etc.)
- Paired up/down migrations
- Atomic operations where possible
- Rollback safety

### Data Loading
- Seed scripts for initial data
- COPY commands for bulk loading
- Transaction safety

This schema provides a solid foundation for the Resume API while maintaining flexibility for future enhancements.