ALTER TABLE click
DROP CONSTRAINT IF EXISTS click_shortlink_id_fkey;

DROP INDEX IF EXISTS idx_click_shortlink_id;
DROP INDEX IF EXISTS idx_click_created_at;
DROP TABLE IF EXISTS click;
