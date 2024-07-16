DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'thread_t') THEN
		CREATE TYPE thread_t AS ENUM ('direct_message', 'server_thread', 'group_message');
	END IF;
END $$;

DO $$
BEGIN
	IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role') THEN
		CREATE TYPE role AS ENUM ('user', 'moderator', 'admin');
	END IF;
END $$;


CREATE TABLE IF NOT EXISTS threads(
	id uuid NOT NULL,
	thread_type thread_t NOT NULL,
	PRIMARY KEY(id),
	server_id uuid,
	name VARCHAR(50),
	voice_enabled BOOLEAN,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz
);


CREATE TABLE IF NOT EXISTS thread_users(
	thread_id uuid NOT NULL REFERENCES threads(id),
	user_id uuid NOT NULL,
	PRIMARY KEY(thread_id, user_id),
	user_role role NOT NULL DEFAULT 'user'
);

CREATE TABLE IF NOT EXISTS messages(
	id uuid NOT NULL,
	PRIMARY KEY(id),
	thread_id uuid NOT NULL REFERENCES threads(id),
	sender_id uuid NOT NULL,
	message VARCHAR(256),
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz
);

CREATE TABLE IF NOT EXISTS reactions(
	message_id uuid NOT NULL REFERENCES messages(id),
	user_id uuid NOT NULL,
	PRIMARY KEY(message_id, user_id),
	reaction VARCHAR(4)
);

CREATE TABLE IF NOT EXISTS reads(
	message_id uuid NOT NULL REFERENCES messages(id),
	user_id uuid NOT NULL,
	PRIMARY KEY(message_id, user_id)
);