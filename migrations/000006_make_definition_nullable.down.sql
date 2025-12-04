-- Revert definition column to NOT NULL
ALTER TABLE vocabularies
ALTER COLUMN definition SET NOT NULL;
