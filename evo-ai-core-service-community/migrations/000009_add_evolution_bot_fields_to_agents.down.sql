DROP INDEX IF EXISTS idx_agents_evolution_bot_sync;
DROP INDEX IF EXISTS idx_agents_evolution_bot_id;

ALTER TABLE evo_core_agents DROP IF EXISTS COLUMN evolution_bot_sync;
ALTER TABLE evo_core_agents DROP IF EXISTS COLUMN evolution_bot_id;