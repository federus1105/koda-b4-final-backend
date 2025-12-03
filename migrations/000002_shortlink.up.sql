CREATE TABLE shortlink (
    id SERIAL PRIMARY KEY,
    account_id INT NULL,
    short_code VARCHAR(20) NOT NULL UNIQUE,
    original_url TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    expired_at TIMESTAMP,
    total_click BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

ALTER TABLE shortlink
ADD CONSTRAINT shortlink_account_id_fkey
FOREIGN KEY (account_id) REFERENCES account(id) ON DELETE CASCADE;

CREATE INDEX idx_shortlink_short_code ON shortlink(short_code);
CREATE INDEX idx_shortlink_account_id ON shortlink(account_id);
CREATE INDEX idx_shortlink_created_at ON shortlink(created_at);
