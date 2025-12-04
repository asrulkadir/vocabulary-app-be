-- Add unique constraint on user_id and word to prevent duplicate vocabularies per user
ALTER TABLE vocabularies
ADD CONSTRAINT unique_user_word UNIQUE(user_id, word);

-- Create index for the unique constraint (PostgreSQL does this automatically, but making it explicit)
CREATE INDEX IF NOT EXISTS idx_vocabularies_user_word ON vocabularies(user_id, word);
