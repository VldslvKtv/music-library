CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    group_id INT REFERENCES groups(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    UNIQUE (group_id, name)
);