-- Rollback: Convert vocabularies table id and user_id from UUID back to INTEGER
-- Drop foreign key
ALTER TABLE vocabularies DROP CONSTRAINT IF EXISTS vocabularies_user_id_fkey;

-- Create new integer columns
ALTER TABLE vocabularies ADD COLUMN id_serial SERIAL;
ALTER TABLE vocabularies ADD COLUMN user_id_serial INTEGER;

-- Populate user_id_serial (this is a simplified approach - in real scenario you'd need mapping)
UPDATE vocabularies SET user_id_serial = (SELECT id FROM users WHERE users.id = vocabularies.user_id LIMIT 1);

-- Drop old uuid columns
ALTER TABLE vocabularies DROP CONSTRAINT vocabularies_pkey;
ALTER TABLE vocabularies DROP COLUMN id;
ALTER TABLE vocabularies DROP COLUMN user_id;

-- Rename new columns
ALTER TABLE vocabularies RENAME COLUMN id_serial TO id;
ALTER TABLE vocabularies RENAME COLUMN user_id_serial TO user_id;

-- Add primary key and foreign key
ALTER TABLE vocabularies ADD PRIMARY KEY (id);
ALTER TABLE vocabularies ADD CONSTRAINT vocabularies_user_id_fkey 
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Recreate indexes
DROP INDEX IF EXISTS idx_vocabularies_user_id;
CREATE INDEX idx_vocabularies_user_id ON vocabularies(user_id);
