CREATE EXTENSION IF NOT EXISTS "pgcrypto";

ALTER TABLE users
ADD COLUMN user_uuid uuid not null default gen_random_uuid();


