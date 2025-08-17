DROP TRIGGER IF EXISTS update_holidays_updated_at;
DROP INDEX IF EXISTS idx_holidays_date_type;
DROP INDEX IF EXISTS idx_holidays_is_active;
DROP INDEX IF EXISTS idx_holidays_type;
DROP INDEX IF EXISTS idx_holidays_date;
DROP TABLE IF EXISTS holidays;
