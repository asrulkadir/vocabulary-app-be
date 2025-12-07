-- Convert users table id from SERIAL to UUID
-- Create new uuid column for users
ALTER TABLE users ADD COLUMN id_uuid UUID DEFAULT gen_random_uuid();

-- Make it unique
ALTER TABLE users ADD CONSTRAINT users_id_uuid_unique UNIQUE(id_uuid);

-- Drop foreign key constraint
ALTER TABLE vocabularies DROP CONSTRAINT vocabularies_user_id_fkey;

-- Create new uuid column for vocabularies.user_id
ALTER TABLE vocabularies ADD COLUMN user_id_uuid UUID;

-- Populate uuid columns from old id columns
UPDATE users SET id_uuid = gen_random_uuid() WHERE id_uuid IS NULL;
UPDATE vocabularies SET user_id_uuid = (
  SELECT id_uuid FROM users WHERE users.id = vocabularies.user_id
);

-- Make uuid columns NOT NULL
ALTER TABLE users ALTER COLUMN id_uuid SET NOT NULL;
ALTER TABLE vocabularies ALTER COLUMN user_id_uuid SET NOT NULL;

-- Drop old columns and rename new ones
ALTER TABLE vocabularies DROP COLUMN user_id;
ALTER TABLE vocabularies RENAME COLUMN user_id_uuid TO user_id;

ALTER TABLE users DROP CONSTRAINT users_pkey;
ALTER TABLE users DROP COLUMN id;
ALTER TABLE users RENAME COLUMN id_uuid TO id;
ALTER TABLE users ADD PRIMARY KEY (id);

-- Recreate foreign key constraint
ALTER TABLE vocabularies ADD CONSTRAINT vocabularies_user_id_fkey 
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Recreate indexes
CREATE INDEX idx_vocabularies_user_id ON vocabularies(user_id);
CREATE INDEX idx_vocabularies_id ON vocabularies(id);
