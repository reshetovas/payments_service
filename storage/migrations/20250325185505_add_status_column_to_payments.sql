-- +goose Up
ALTER TABLE payments ADD COLUMN state TEXT DEFAULT 'NEW';
ALTER TABLE payments ADD COLUMN attempts INTEGER DEFAULT 0;

-- +goose Down
ALTER TABLE payments DROP COLUMN state;
