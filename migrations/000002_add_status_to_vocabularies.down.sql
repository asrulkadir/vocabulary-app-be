-- Remove status column from vocabularies table
ALTER TABLE vocabularies
DROP COLUMN IF EXISTS status;

-- Drop the status index
DROP INDEX IF EXISTS idx_vocabularies_status;
