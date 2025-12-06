-- Remove passed_test_count and failed_test_count columns from vocabularies table
ALTER TABLE vocabularies
DROP COLUMN IF EXISTS passed_test_count,
DROP COLUMN IF EXISTS failed_test_count;

-- Drop the indexes
DROP INDEX IF EXISTS idx_vocabularies_passed_test_count;
DROP INDEX IF EXISTS idx_vocabularies_failed_test_count;
