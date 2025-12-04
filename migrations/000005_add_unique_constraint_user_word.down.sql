-- Remove unique constraint on user_id and word
ALTER TABLE vocabularies
DROP CONSTRAINT IF EXISTS unique_user_word;

-- Drop the index
DROP INDEX IF EXISTS idx_vocabularies_user_word;
