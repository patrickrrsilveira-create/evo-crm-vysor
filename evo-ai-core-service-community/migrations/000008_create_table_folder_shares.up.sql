CREATE TABLE IF NOT EXISTS evo_core_folder_shares (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    folder_id UUID REFERENCES evo_core_folders(id) ON DELETE CASCADE,
    shared_by_user_id UUID, -- Removed foreign key reference to external users table
    shared_with_email VARCHAR(255) NOT NULL,
    shared_with_user_id UUID, -- Removed foreign key reference to external users table
    permission_level VARCHAR(5) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_permission_level
        CHECK (permission_level::text = ANY (ARRAY['read', 'write']))
);

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'evo_core_folder_shares'
        AND column_name = 'folder_id'
    ) THEN
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
    
            WHERE tablename = 'evo_core_folder_shares'
            AND indexname = 'idx_evo_core_folder_shares_folder_id'
        ) THEN
            CREATE INDEX idx_evo_core_folder_shares_folder_id ON evo_core_folder_shares (folder_id);
        END IF;
    END IF;

    IF EXISTS (

        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'evo_core_folder_shares'
        AND column_name = 'shared_by_user_id'
    ) THEN
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE tablename = 'evo_core_folder_shares'
            AND indexname = 'idx_evo_core_folder_shares_shared_by_user_id'
        ) THEN
            CREATE INDEX idx_evo_core_folder_shares_shared_by_user_id ON evo_core_folder_shares (shared_by_user_id);
        END IF;
    END IF;

    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'evo_core_folder_shares'
        AND column_name = 'shared_with_user_id'
    ) THEN
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE tablename = 'evo_core_folder_shares'
            AND indexname = 'idx_evo_core_folder_shares_shared_with_user_id'
        ) THEN
            CREATE INDEX idx_evo_core_folder_shares_shared_with_user_id ON evo_core_folder_shares (shared_with_user_id);
        END IF;
    END IF;
END $$;
