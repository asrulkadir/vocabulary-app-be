-- Make definition column nullable
ALTER TABLE vocabularies
ALTER COLUMN definition DROP NOT NULL;
