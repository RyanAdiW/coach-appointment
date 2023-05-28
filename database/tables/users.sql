CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.users (
	id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	name varchar(100) NOT NULL,
	email varchar(100) NOT NULL UNIQUE,
	password varchar(255) NOT NULL,
	role varchar(20) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
	updated_at TIMESTAMP NULL
);