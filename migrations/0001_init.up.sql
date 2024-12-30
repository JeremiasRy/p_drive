CREATE TABLE IF NOT EXISTS file_meta_data(
    path TEXT PRIMARY KEY,
    size_bytes INT NOT NULL,
    mime TEXT NOT NULL
);