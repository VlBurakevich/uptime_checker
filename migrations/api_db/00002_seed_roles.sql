-- +goose Up
INSERT INTO roles (name) VALUES ('user'), ('admin')
ON CONFLICT (name) DO NOTHING;

-- +goose Down
DELETE FROM roles WHERE name IN ('user', 'admin');