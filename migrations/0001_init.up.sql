CREATE TABLE IF NOT EXISTS users(
    id UUID default gen_random_uuid() PRIMARY KEY,
    email varchar(256) unique
);

CREATE TABLE IF NOT EXISTS folders(
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name TEXT NOT NULL,
    parent_id UUID DEFAULT NULL,
    folder_client_path TEXT
);

ALTER TABLE IF EXISTS folders
ADD CONSTRAINT fk_parent
FOREIGN KEY (parent_id)
REFERENCES folders(id);

CREATE OR REPLACE FUNCTION set_default_folder_client_path()
RETURNS TRIGGER AS $$
BEGIN
    NEW.folder_client_path := NEW.id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_folder_client_path_default
BEFORE INSERT ON folders
FOR EACH ROW
WHEN (NEW.folder_client_path IS NULL)
EXECUTE FUNCTION set_default_folder_client_path();

CREATE type file_status AS ENUM ('OK', 'UPLOADING', 'ERROR');

CREATE TABLE IF NOT EXISTS file_meta_data(
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    folder_id UUID NOT NULL REFERENCES folders(id),
    size_bytes INT NOT NULL,
    mime TEXT NOT NULL,
    name TEXT NOT NULL,
    signed_link TEXT default null,
    status file_status default 'UPLOADING'
);
