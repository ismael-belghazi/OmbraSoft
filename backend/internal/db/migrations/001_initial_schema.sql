
-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Series table
CREATE TABLE IF NOT EXISTS series (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    source_site VARCHAR(255),
    source_url VARCHAR(255) NOT NULL,
    last_chapter_number INT,
    last_checked_at TIMESTAMP
);

-- Bookmarks table
CREATE TABLE IF NOT EXISTS bookmarks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    series_id UUID NOT NULL REFERENCES series(id),
    last_read_chapter INT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
