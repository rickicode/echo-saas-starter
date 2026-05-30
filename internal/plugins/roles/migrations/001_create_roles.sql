CREATE TABLE IF NOT EXISTS roles (
    id INTEGER PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    description TEXT DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO roles (id, name, description) VALUES (1, 'super_admin', 'Super Administrator with full access') ON CONFLICT DO NOTHING;
INSERT INTO roles (id, name, description) VALUES (2, 'admin', 'Administrator with elevated access') ON CONFLICT DO NOTHING;
INSERT INTO roles (id, name, description) VALUES (3, 'user', 'Standard user') ON CONFLICT DO NOTHING;
