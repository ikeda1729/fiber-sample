CREATE TABLE users (
	id bigserial NOT NULL unique,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	username text NOT null unique,
	email text NOT null unique,
	"password" text NOT null,
	names text NULL
);

CREATE INDEX idx_users_deleted_at ON users (deleted_at timestamptz_ops);


CREATE TABLE tweets (
	id bigserial NOT NULL unique,
	created_at timestamptz NULL,
	updated_at timestamptz NULL,
	deleted_at timestamptz NULL,
	"content" text NOT NULL,
	user_id int8 NOT NULL
);
CREATE INDEX idx_tweets_deleted_at ON tweets (deleted_at timestamptz_ops);


CREATE TABLE user_followees (
	id int primary key NOT NULL,
	user_id int8 NOT NULL,
	followee_id int8 NOT NULL
);