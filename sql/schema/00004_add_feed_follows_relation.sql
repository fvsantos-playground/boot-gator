-- +goose Up
CREATE TABLE public.feed_follows (
	id int GENERATED ALWAYS AS IDENTITY NOT NULL,
	created_at date NOT NULL,
	updated_at date NOT NULL,
	user_id uuid NOT NULL,
	feed_id uuid NOT NULL,
	CONSTRAINT feed_follows_pk PRIMARY KEY (id),
	CONSTRAINT feed_follows_users_fk FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE,
	CONSTRAINT feed_follows_feeds_fk FOREIGN KEY (feed_id) REFERENCES public.feeds(id) ON DELETE CASCADE,
    CONSTRAINT unique_follows UNIQUE (user_id,feed_id)
);


-- +goose Down
DROP TABLE public.feed_follows;
