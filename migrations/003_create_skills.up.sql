-- Create skills table with categories
CREATE TABLE skills (
    id SERIAL PRIMARY KEY,
    category VARCHAR(100) NOT NULL,
    name VARCHAR(255) NOT NULL,
    level VARCHAR(50), -- beginner, intermediate, advanced, expert
    years_experience INTEGER,
    order_index INTEGER DEFAULT 0, -- For ordering within category
    is_featured BOOLEAN DEFAULT FALSE, -- Highlight key skills
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure no duplicate skills within same category
    UNIQUE(category, name)
);

-- Create trigger for updated_at
CREATE TRIGGER update_skills_updated_at BEFORE UPDATE
    ON skills FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create indexes for efficient querying
CREATE INDEX idx_skills_category ON skills(category);
CREATE INDEX idx_skills_category_order ON skills(category, order_index);
CREATE INDEX idx_skills_featured ON skills(is_featured) WHERE is_featured = TRUE;
CREATE INDEX idx_skills_level ON skills(level);

-- Create check constraint for valid skill levels
ALTER TABLE skills ADD CONSTRAINT chk_skill_level 
    CHECK (level IN ('beginner', 'intermediate', 'advanced', 'expert') OR level IS NULL);