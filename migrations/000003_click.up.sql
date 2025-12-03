CREATE TABLE click (
    id SERIAL PRIMARY KEY,
    shortlink_id INT NOT NULL,
    user_agent TEXT,
    ip INET,
    referer TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

ALTER TABLE click
ADD CONSTRAINT click_shortlink_id_fkey
FOREIGN KEY (shortlink_id) REFERENCES shortlink(id) ON DELETE CASCADE;

CREATE INDEX idx_click_shortlink_id ON click(shortlink_id);
CREATE INDEX idx_click_created_at ON click(created_at);
