CREATE TABLE IF NOT EXISTS refresh_tokens (
    user_id VARCHAR(255) PRIMARY KEY,
    token TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

create table if not exists users (
     user_id uuid primary key default gen_random_uuid(),
     email text not null unique,
     password_hash text not null,
     created_at timestamptz not null default now(),
     updated_at timestamptz not null default now()
);