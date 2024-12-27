CREATE TABLE IF NOT EXISTS users(
    id UUID default gen_random_uuid() PRIMARY KEY,
    email varchar(256) unique
);

ALTER TABLE IF EXISTS file_meta_data
ADD COLUMN user_id UUID;

ALTER TABLE IF EXISTS file_meta_data
ADD CONSTRAINT fk_file_users
FOREIGN KEY (user_id)
REFERENCES users(id)
ON DELETE CASCADE;