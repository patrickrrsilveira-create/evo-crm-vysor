-- Alter the 'key' column type from VARCHAR(255) to TEXT to support encrypted keys
-- Encrypted keys can be much longer than 255 characters
ALTER TABLE evo_core_api_keys 
ALTER COLUMN key TYPE TEXT;

