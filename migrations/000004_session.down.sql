ALTER TABLE session
DROP CONSTRAINT IF EXISTS session_user_id_fkey;

DROP INDEX IF EXISTS idx_session_token;
DROP INDEX IF EXISTS idx_session_user_id;

DROP TABLE IF EXISTS session;
