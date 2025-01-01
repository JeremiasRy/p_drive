CREATE TABLE IF NOT EXISTS file_meta_data(
    id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    folder_path TEXT,
    size_bytes INT NOT NULL,
    mime TEXT NOT NULL,
    name TEXT NOT NULL
);