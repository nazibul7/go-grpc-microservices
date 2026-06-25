CREATE TYPE user_role AS ENUM ('user', 'admin');

CREATE TABLE
    users (
        id UUID PRIMARY KEY, -- same id as user_db.users, set by UserService
        email TEXT NOT NULL UNIQUE,
        password_hash TEXT NOT NULL,
        role user_role NOT NULL DEFAULT 'user',
        created_at TIMESTAMPTZ NOT NULL DEFAULT now ()
    );