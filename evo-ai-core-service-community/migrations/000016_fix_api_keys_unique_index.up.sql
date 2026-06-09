-- Drop the too restrictive unique index on name
DROP INDEX IF EXISTS idx_evo_core_api_keys_name_unique;

-- Re-create a unique index on name but only for active keys
CREATE UNIQUE INDEX idx_evo_core_api_keys_name_active_unique ON evo_core_api_keys (name) WHERE (is_active = true);
