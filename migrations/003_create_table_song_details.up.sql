CREATE TABLE IF NOT EXISTS song_details (
    id SERIAL PRIMARY KEY,
    song_id INT REFERENCES songs(id) UNIQUE ON DELETE CASCADE,
    release_date DATE,
    text TEXT,
    link TEXT
);