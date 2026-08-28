CREATE TABLE IF NOT EXISTS characters (
    user_id VARCHAR(64) PRIMARY KEY,
    level INT NOT NULL DEFAULT 1,
    current_xp INT NOT NULL DEFAULT 0,
    hp NUMERIC(3, 1) NOT NULL DEFAULT 6.0,
    strength INT NOT NULL DEFAULT 10,
    endurance INT NOT NULL DEFAULT 10,
    discipline INT NOT NULL DEFAULT 10,
    balance INT NOT NULL DEFAULT 10,
    current_streak INT NOT NULL DEFAULT 0,
    last_active_date DATE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);