-- Create education table for education and certifications
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
    credential_id VARCHAR(255), -- For certifications
    credential_url VARCHAR(500), -- Link to verify certification
    expiry_date DATE, -- For certifications that expire
    order_index INTEGER DEFAULT 0,
    is_featured BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create trigger for updated_at
CREATE TRIGGER update_education_updated_at BEFORE UPDATE
    ON education FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create indexes
CREATE INDEX idx_education_type ON education(type);
CREATE INDEX idx_education_year ON education(year_completed DESC, year_started DESC);
CREATE INDEX idx_education_institution ON education(institution);
CREATE INDEX idx_education_status ON education(status);
CREATE INDEX idx_education_order ON education(type, order_index);
CREATE INDEX idx_education_featured ON education(is_featured) WHERE is_featured = TRUE;