CREATE TABLE IF NOT EXISTS plans (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    description TEXT DEFAULT '',
    price_monthly REAL DEFAULT 0,
    price_yearly REAL DEFAULT 0,
    features TEXT DEFAULT '[]',
    is_active INTEGER DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO plans (id, name, slug, description, price_monthly, price_yearly, features, is_active) VALUES (1, 'Free', 'free', 'Basic plan for getting started', 0, 0, '["basic"]', 1) ON CONFLICT DO NOTHING;
INSERT INTO plans (id, name, slug, description, price_monthly, price_yearly, features, is_active) VALUES (2, 'Pro', 'pro', 'Professional plan with advanced features', 29, 290, '["basic","advanced","analytics"]', 1) ON CONFLICT DO NOTHING;
INSERT INTO plans (id, name, slug, description, price_monthly, price_yearly, features, is_active) VALUES (3, 'Enterprise', 'enterprise', 'Full-featured enterprise plan', 99, 990, '["basic","advanced","analytics","priority","custom"]', 1) ON CONFLICT DO NOTHING;
