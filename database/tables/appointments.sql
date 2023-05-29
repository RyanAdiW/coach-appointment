CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE public.appointments (
	id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	user_id VARCHAR NOT NULL,
    coach_name VARCHAR NOT NULL,
    status varchar(50) NOT NULL,
    appointment_start TIMESTAMP WITH TIME ZONE,
    appointment_end TIMESTAMP WITH TIME ZONE,
	created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
	updated_at TIMESTAMP NULL
);