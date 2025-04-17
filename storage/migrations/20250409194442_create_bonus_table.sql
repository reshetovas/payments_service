-- +goose Up
CREATE TABLE bonuses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    payment_id INTEGER,
    amount REAL NOT NULL,
    FOREIGN KEY (payment_id) REFERENCES payments(id)
);

-- +goose Down
DROP TABLE bonuses;
