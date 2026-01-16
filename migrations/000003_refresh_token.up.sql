create table refreshtoken(
	email text unique not null,
	refreshToken text unique not null
);

