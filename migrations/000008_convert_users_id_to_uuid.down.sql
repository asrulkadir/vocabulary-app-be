-- Rollback: Convert users and vocabularies back to SERIAL integers
-- Drop foreign key
ALTER TABLE vocabularies DROP CONSTRAINT IF EXISTS vocabularies_user_id_fkey;

-- Drop indexes
DROP INDEX IF EXISTS idx_vocabularies_user_id;
DROP INDEX IF EXISTS idx_vocabularies_id;

-- Create temporary int column for users
ALTER TABLE users ADD COLUMN id_int BIGINT;
UPDATE users SET id_int = ROW_NUMBER() OVER (ORDER BY id ASC);

-- Store uuid mapping before dropping
ALTER TABLE users ADD COLUMN id_uuid_backup UUID;
UPDATE users SET id_uuid_backup = id;

-- Drop old uuid id and rename
ALTER TABLE users DROP CONSTRAINT users_pkey;
ALTER TABLE users DROP COLUMN id;
ALTER TABLE users RENAME COLUMN id_int TO id;
ALTER TABLE users ADD PRIMARY KEY (id);

-- Create temporary int column for vocabularies.user_id
ALTER TABLE vocabularies ADD COLUMN user_id_int BIGINT;

-- Populate vocabularies.user_id from users using the uuid backup
UPDATE vocabularies SET user_id_int = users.id 
FROM users WHERE users.id_uuid_backup = vocabularies.user_id;

-- Drop uuid columns and rename integer columns
ALTER TABLE vocabularies DROP COLUMN user_id;
ALTER TABLE vocabularies RENAME COLUMN user_id_int TO user_id;

ALTER TABLE users DROP COLUMN id_uuid_backup;

-- Recreate foreign key
ALTER TABLE vocabularies ADD CONSTRAINT vocabularies_user_id_fkey 
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Recreate indexes
CREATE INDEX idx_vocabularies_user_id ON vocabularies(user_id);
CREATE INDEX idx_vocabularies_id ON vocabularies(id);
