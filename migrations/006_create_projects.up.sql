-- Create projects table for notable projects like homelab
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    short_description VARCHAR(500), -- For summary views
    technologies JSONB, -- Store array of technologies used
    github_url VARCHAR(500),
    demo_url VARCHAR(500),
    start_date DATE,
    end_date DATE, -- NULL for ongoing projects
    status VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active', 'completed', 'archived', 'planned')),
    is_featured BOOLEAN DEFAULT FALSE,
    order_index INTEGER DEFAULT 0,
    key_features TEXT[], -- Array of key features/highlights
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger for updated_at
CREATE TRIGGER update_projects_updated_at BEFORE UPDATE
    ON projects FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create indexes
CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_projects_dates ON projects(start_date DESC, end_date DESC);
CREATE INDEX idx_projects_featured ON projects(is_featured) WHERE is_featured = TRUE;
CREATE INDEX idx_projects_order ON projects(order_index);

-- Create GIN index for JSONB technologies column for efficient searching
CREATE INDEX idx_projects_technologies ON projects USING GIN (technologies);