CREATE TABLE IF NOT EXISTS folders(
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name TEXT NOT NULL,
    parent_id UUID NULL
);

ALTER TABLE IF EXISTS folders
ADD CONSTRAINT fk_parent
FOREIGN KEY (parent_id)
REFERENCES folders(id);