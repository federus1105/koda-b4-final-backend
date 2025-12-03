ALTER TABLE shortlink DROP CONSTRAINT IF EXISTS shortlink_account_id_fkey;

-- Drop indexes
DROP INDEX IF EXISTS idx_shortlink_short_code;
DROP INDEX IF EXISTS idx_shortlink_account_id;
DROP INDEX IF EXISTS idx_shortlink_created_at;

-- Drop table
DROP TABLE IF EXISTS shortlink;
