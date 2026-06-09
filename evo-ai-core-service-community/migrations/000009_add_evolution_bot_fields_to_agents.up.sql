ALTER TABLE evo_core_agents ADD COLUMN IF NOT EXISTS evolution_bot_id UUID;
ALTER TABLE evo_core_agents ADD COLUMN IF NOT EXISTS evolution_bot_sync BOOLEAN NOT NULL DEFAULT false;

CREATE INDEX IF NOT EXISTS idx_agents_evolution_bot_id ON evo_core_agents(evolution_bot_id);
CREATE INDEX IF NOT EXISTS idx_agents_evolution_bot_sync ON evo_core_agents(evolution_bot_sync);