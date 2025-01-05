CREATE TABLE IF NOT EXISTS users(
    id UUID default gen_random_uuid() PRIMARY KEY,
    email varchar(256) unique
);

CREATE TABLE IF NOT EXISTS folders(
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name TEXT NOT NULL,
    parent_id UUID DEFAULT NULL
);

ALTER TABLE IF EXISTS folders
ADD CONSTRAINT fk_parent
FOREIGN KEY (parent_id)
REFERENCES folders(id);

CREATE TABLE IF NOT EXISTS file_meta_data(
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    folder_id UUID NOT NULL REFERENCES folders(id),
    size_bytes INT NOT NULL,
    mime TEXT NOT NULL,
    name TEXT NOT NULL,
    signed_link TEXT default null
);
