-- Add status column to vocabularies table
ALTER TABLE vocabularies
ADD COLUMN status VARCHAR(50) DEFAULT 'learning' CHECK (status IN ('learning', 'memorized'));

-- Create index on status for better query performance
CREATE INDEX IF NOT EXISTS idx_vocabularies_status ON vocabularies(status);
