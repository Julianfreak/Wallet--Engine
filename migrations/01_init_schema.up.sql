CREATE TABLE accounts (
    id VARCHAR(50) PRIMARY KEY,
    owner VARCHAR(100) NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0,
    CONSTRAINT balance_non_negative CHECK (balance >= 0)
);

CREATE TABLE transactions (
    id VARCHAR(50) PRIMARY KEY,
    from_account_id VARCHAR(50) NOT NULL,
    to_account_id VARCHAR(50) NOT NULL,
    amount BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);