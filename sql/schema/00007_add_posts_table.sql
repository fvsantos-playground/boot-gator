-- +goose Up
CREATE TABLE public.posts (
	id int GENERATED ALWAYS AS IDENTITY NOT NULL,
	created_at timestamp NOT NULL,
	updated_at timestamp NOT NULL,
	title text NOT NULL,
	url text NOT NULL,
	description text NULL,
	published_at timestamp NULL,
	feed_id uuid NOT NULL,
	CONSTRAINT posts_pk PRIMARY KEY (id),
	CONSTRAINT posts_unique UNIQUE (url)
);

-- +goose Down
DROP TABLE IF EXISTS posts;
