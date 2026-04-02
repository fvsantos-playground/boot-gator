-- +goose Up
CREATE VIEW feeds_with_user AS
SELECT f.id, f.name, f.url, f.created_at, f.updated_at, u.name AS user_name
FROM feeds f
JOIN users u ON f.user_id = u.id;

-- +goose Down
DROP VIEW feeds_with_user;
