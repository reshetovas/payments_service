-- +goose Up
CREATE TABLE payments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    amount REAL NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    currency TEXT
);

-- +goose Down
DROP TABLE payments;
