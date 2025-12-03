-- Change example column to support multiple examples using JSON array
ALTER TABLE vocabularies
ALTER COLUMN example TYPE JSONB USING CASE 
  WHEN example IS NULL THEN '[]'::JSONB
  WHEN example = '' THEN '[]'::JSONB
  ELSE jsonb_build_array(example)
END,
ALTER COLUMN example SET DEFAULT '[]'::JSONB;
