DROP INDEX IF EXISTS idx_evo_core_api_keys_name_active_unique;

CREATE UNIQUE INDEX idx_evo_core_api_keys_name_unique ON evo_core_api_keys (name);
