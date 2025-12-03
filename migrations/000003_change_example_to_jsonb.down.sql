-- Revert example column back to TEXT
ALTER TABLE vocabularies
ALTER COLUMN example TYPE TEXT USING CASE 
  WHEN example::text = '[]' THEN NULL
  WHEN jsonb_array_length(example) > 0 THEN example->>0
  ELSE NULL
END,
ALTER COLUMN example DROP DEFAULT;
