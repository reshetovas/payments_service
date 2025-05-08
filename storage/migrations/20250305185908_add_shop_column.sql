-- +goose Up
ALTER TABLE payments ADD COLUMN shop_id INTEGER;
ALTER TABLE payments ADD COLUMN address TEXT;

CREATE TABLE items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    payment_id INTEGER,
    name TEXT NOT NULL,
    price REAL NOT NULL,
    quantity INTEGER NOT NULL,
    FOREIGN KEY (payment_id) REFERENCES payments(id) ON DELETE CASCADE
);

-- +goose Down
ALTER TABLE payments DROP COLUMN shop_id;
ALTER TABLE payments DROP COLUMN address;
DROP TABLE items;

