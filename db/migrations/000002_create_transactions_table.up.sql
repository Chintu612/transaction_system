CREATE TABLE IF NOT EXISTS transactions(
    id BIGINT PRIMARY KEY,
    amount DOUBLE PRECISION NOT NULL,
    type VARCHAR(50) NOT NULL,
    parent_id BIGINT,
    FOREIGN KEY (parent_id) REFERENCES transactions(id)
);