CREATE TABLE IF NOT EXISTS evo_core_mcp_servers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name          VARCHAR(255)             NOT NULL,
    description   TEXT,
    config_type   VARCHAR(10)             NOT NULL,
    config_json   JSON                     NOT NULL,
    environments  JSON                     NOT NULL,
    tools         JSON                     NOT NULL,
    type          VARCHAR(10)              NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT check_mcp_server_config_type CHECK (
        config_type::TEXT = ANY (ARRAY['studio', 'sse'])
    ),
    CONSTRAINT check_mcp_server_type CHECK (
        type::TEXT = ANY (ARRAY['official', 'community'])
    )
);

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.tables
        WHERE table_name = 'evo_core_mcp_servers'
    ) AND EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'evo_core_mcp_servers'
        AND column_name = 'name'
    ) THEN
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE tablename = 'evo_core_mcp_servers'
            AND indexname = 'idx_evo_core_mcp_servers_name'
        ) THEN
            CREATE UNIQUE INDEX idx_evo_core_mcp_servers_name ON evo_core_mcp_servers (name);
        END IF;
    END IF;
END
$$;