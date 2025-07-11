-- Create experiences table for work history
CREATE TABLE experiences (
    id SERIAL PRIMARY KEY,
    company VARCHAR(255) NOT NULL,
    position VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE, -- NULL indicates current position
    description TEXT,
    highlights TEXT[], -- Array of key achievements/responsibilities
    order_index INTEGER DEFAULT 0, -- For ordering display
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger for updated_at
CREATE TRIGGER update_experiences_updated_at BEFORE UPDATE
    ON experiences FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create indexes for common queries
CREATE INDEX idx_experiences_company ON experiences(company);
CREATE INDEX idx_experiences_dates ON experiences(start_date DESC, end_date DESC);
CREATE INDEX idx_experiences_order ON experiences(order_index);
CREATE INDEX idx_experiences_current ON experiences(end_date) WHERE end_date IS NULL;