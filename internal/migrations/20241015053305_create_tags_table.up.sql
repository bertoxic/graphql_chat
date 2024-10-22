-- Create tags table
CREATE TABLE IF NOT EXISTS tags(
                                    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                    name VARCHAR(255) UNIQUE NOT NULL,
                                    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create index for tags
CREATE INDEX idx_tags_name ON tags(name);
