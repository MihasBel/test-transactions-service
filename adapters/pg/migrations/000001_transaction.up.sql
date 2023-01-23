
CREATE TABLE transactions (
                           id UUID PRIMARY KEY,
                           created_at TIMESTAMP DEFAULT now() NOT NULL,
                           user_id UUID NOT NULL,
                           amount BIGINT NOT NULL,
                           status SMALLINT DEFAULT 0 NOT NULL,
                           description CHAR (128)
);
CREATE INDEX transactions_user_id_idx ON transactions USING hash (user_id);