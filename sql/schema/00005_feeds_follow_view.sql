-- +goose Up
CREATE VIEW feeds_follow_expanded AS
SELECT ff.id, ff.created_at, ff.updated_at, ff.user_id, ff.feed_id, u.name AS user_name, f.name AS feed_name, f.url AS feed_url
FROM feed_follows ff
JOIN users u ON ff.user_id = u.id
JOIN feeds f ON ff.feed_id = f.id;

-- +goose Down
DROP VIEW feeds_follow_expanded;
