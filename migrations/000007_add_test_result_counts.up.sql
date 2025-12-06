-- Add passed_test_count and failed_test_count columns to vocabularies table
ALTER TABLE vocabularies
ADD COLUMN passed_test_count INTEGER DEFAULT 0,
ADD COLUMN failed_test_count INTEGER DEFAULT 0;

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_vocabularies_passed_test_count ON vocabularies(passed_test_count);
CREATE INDEX IF NOT EXISTS idx_vocabularies_failed_test_count ON vocabularies(failed_test_count);
