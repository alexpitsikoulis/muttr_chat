CREATE TABLE IF NOT EXISTS dm_threads(
	id uuid NOT NULL,
	PRIMARY KEY(id),
	first_user_id uuid NOT NULL,
	second_user_id uuid NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz,
	CHECK (first_user_id < second_user_id)
);

CREATE TABLE IF NOT EXISTS direct_messages(
	id uuid NOT NULL,
	PRIMARY KEY(id),
	thread_id uuid NOT NULL REFERENCES dm_threads(id),
	sender_id uuid NOT NULL,
	message VARCHAR(2000) NOT NULL,
	is_read BOOLEAN DEFAULT false,
	reaction VARCHAR(4),
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz
);

CREATE TABLE IF NOT EXISTS server_threads(
	id uuid NOT NULL,
	PRIMARY KEY(id),
	server_id uuid NOT NULL,
	name VARCHAR(50),
	voice_enabled BOOLEAN DEFAULT true,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz
);

CREATE TABLE IF NOT EXISTS server_messages(
	id uuid NOT NULL,
	PRIMARY KEY(id),
	thread_id uuid NOT NULL REFERENCES server_threads(id),
	sender_id uuid NOT NULL,
	message VARCHAR(2000) NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz
);

CREATE TABLE IF NOT EXISTS server_message_reads(
	message_id uuid NOT NULL REFERENCES server_messages(id),
	user_id uuid NOT NULL,
	PRIMARY KEY(message_id, user_id),
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz
);

CREATE TABLE IF NOT EXISTS server_message_reactions(
	message_id uuid NOT NULL REFERENCES server_messages(id),
	user_id uuid NOT NULL,
	PRIMARY KEY(message_id, user_id),
	reaction VARCHAR(4) NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz
);

CREATE TABLE IF NOT EXISTS group_chats(
	id uuid,
	PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS group_chat_members(
	group_chat_id uuid NOT NULL REFERENCES group_chats(id),
	user_id uuid NOT NULL,
	PRIMARY KEY(group_chat_id, user_id)
);

CREATE TABLE IF NOT EXISTS group_chat_messages(
	id uuid NOT NULL,
	PRIMARY KEY(id),
	group_chat_id uuid NOT NULL REFERENCES group_chats(id),
	sender_id uuid NOT NULL,
	message VARCHAR(2000) NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz
);

CREATE TABLE IF NOT EXISTS group_chat_message_reads(
	message_id uuid NOT NULL REFERENCES group_chat_messages(id),
	user_id uuid NOT NULL,
	PRIMARY KEY(message_id, user_id),
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz
);

CREATE TABLE IF NOT EXISTS group_chat_message_reactions(
	message_id uuid NOT NULL REFERENCES group_chat_messages(id),
	user_id uuid NOT NULL,
	PRIMARY KEY(message_id, user_id),
	reaction VARCHAR(4) NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NOT NULL DEFAULT now(),
	deleted_at timestamptz
);