create extension if not exists pgcrypto;

create table users_table (
	user_id uuid primary key default gen_random_uuid(),
	email text unique not null,
	password text not null,
	created_at timestamptz default now()
);
