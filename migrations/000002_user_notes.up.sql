create table user_notes(
	user_id uuid references users_table(user_id),
	title text unique not null,
	myspace text,
	created_at timestamptz default now()
);

create index my_index_getnotes on user_notes(user_id, created_at);
