-- Create achievements table for key accomplishments
CREATE TABLE achievements (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100), -- performance, security, leadership, innovation, etc.
    impact_metric VARCHAR(255), -- "30% reduction", "40% improvement", etc.
    year_achieved INTEGER,
    order_index INTEGER DEFAULT 0,
    is_featured BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger for updated_at
CREATE TRIGGER update_achievements_updated_at BEFORE UPDATE
    ON achievements FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create indexes
CREATE INDEX idx_achievements_category ON achievements(category);
CREATE INDEX idx_achievements_year ON achievements(year_achieved DESC);
CREATE INDEX idx_achievements_order ON achievements(order_index);
CREATE INDEX idx_achievements_featured ON achievements(is_featured) WHERE is_featured = TRUE;