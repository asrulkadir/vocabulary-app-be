-- Convert vocabularies table id and user_id from INTEGER to UUID
-- Drop foreign key first
ALTER TABLE vocabularies DROP CONSTRAINT vocabularies_user_id_fkey;

-- Create new uuid columns
ALTER TABLE vocabularies ADD COLUMN id_uuid UUID DEFAULT gen_random_uuid();
ALTER TABLE vocabularies ADD COLUMN user_id_uuid UUID;

-- Make id_uuid unique
ALTER TABLE vocabularies ADD CONSTRAINT vocabularies_id_uuid_unique UNIQUE(id_uuid);

-- Populate user_id_uuid from users table
UPDATE vocabularies v SET user_id_uuid = u.id FROM users u WHERE v.user_id = u.id;

-- Make NOT NULL after population
ALTER TABLE vocabularies ALTER COLUMN id_uuid SET NOT NULL;
ALTER TABLE vocabularies ALTER COLUMN user_id_uuid SET NOT NULL;

-- Drop old columns
ALTER TABLE vocabularies DROP CONSTRAINT vocabularies_pkey;
ALTER TABLE vocabularies DROP COLUMN id;
ALTER TABLE vocabularies DROP COLUMN user_id;

-- Rename new columns
ALTER TABLE vocabularies RENAME COLUMN id_uuid TO id;
ALTER TABLE vocabularies RENAME COLUMN user_id_uuid TO user_id;

-- Add primary key and foreign key
ALTER TABLE vocabularies ADD PRIMARY KEY (id);
ALTER TABLE vocabularies ADD CONSTRAINT vocabularies_user_id_fkey 
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Recreate indexes
DROP INDEX IF EXISTS idx_vocabularies_user_id;
CREATE INDEX idx_vocabularies_user_id ON vocabularies(user_id);
