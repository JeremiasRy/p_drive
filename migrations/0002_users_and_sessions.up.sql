CREATE TABLE IF NOT EXISTS users(
    id UUID default gen_random_uuid() PRIMARY KEY,
    email varchar(256) unique
);