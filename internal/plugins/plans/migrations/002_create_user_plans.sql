CREATE TABLE IF NOT EXISTS user_plans (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    plan_id INTEGER NOT NULL,
    status TEXT DEFAULT 'active',
    starts_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (plan_id) REFERENCES plans(id)
);
