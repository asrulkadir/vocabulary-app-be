-- Add test_count column to vocabularies table
ALTER TABLE vocabularies
ADD COLUMN test_count INTEGER DEFAULT 0;

-- Create index on test_count for better query performance
CREATE INDEX IF NOT EXISTS idx_vocabularies_test_count ON vocabularies(test_count);
