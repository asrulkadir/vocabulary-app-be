-- Remove test_count column from vocabularies table
ALTER TABLE vocabularies
DROP COLUMN IF EXISTS test_count;

-- Drop the test_count index
DROP INDEX IF EXISTS idx_vocabularies_test_count;
