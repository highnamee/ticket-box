-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_categories_deleted_at ON categories (deleted_at);
CREATE INDEX IF NOT EXISTS idx_categories_slug ON categories (slug);

ALTER TABLE tickets 
    ADD COLUMN IF NOT EXISTS category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS venue VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS date VARCHAR(100) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS time VARCHAR(50) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS image_url TEXT NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS tags TEXT[] NOT NULL DEFAULT '{}',
    ADD COLUMN IF NOT EXISTS featured BOOLEAN NOT NULL DEFAULT false;

CREATE INDEX IF NOT EXISTS idx_tickets_category_id ON tickets (category_id);
CREATE INDEX IF NOT EXISTS idx_tickets_featured ON tickets (featured);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_tickets_featured;
DROP INDEX IF EXISTS idx_tickets_category_id;

ALTER TABLE tickets 
    DROP COLUMN IF EXISTS featured,
    DROP COLUMN IF EXISTS tags,
    DROP COLUMN IF EXISTS image_url,
    DROP COLUMN IF EXISTS time,
    DROP COLUMN IF EXISTS date,
    DROP COLUMN IF EXISTS venue,
    DROP COLUMN IF EXISTS category_id;

DROP TABLE IF EXISTS categories;
-- +goose StatementEnd
