CREATE INDEX idx_transaction_id ON transactions (id);

CREATE INDEX idx_transaction_type ON transactions (type);

CREATE INDEX idx_transaction_id_parent_id ON transactions (id, parent_id);