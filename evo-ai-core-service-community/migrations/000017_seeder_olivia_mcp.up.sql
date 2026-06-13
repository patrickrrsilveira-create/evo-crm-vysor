DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'evo_core_mcp_servers') THEN
        
        -- Insert Olivia MCP Server
        IF NOT EXISTS (SELECT 1 FROM evo_core_mcp_servers WHERE id = 'c1234567-89ab-cdef-0123-456789abcdef'::uuid) THEN
            INSERT INTO evo_core_mcp_servers (id, name, description, config_type, config_json, environments, tools, type, created_at, updated_at)
            VALUES (
                'c1234567-89ab-cdef-0123-456789abcdef'::uuid,
                'Olivia Assistant Tools',
                'Servidor MCP com as ferramentas da Agente Executiva Olivia (Gestão de Agenda, E-mails e Anotações).',
                'studio',
                '{"command": "python", "args": ["/olivia-mcp-server/main.py"], "env": {}}'::json,
                '{}'::json,
                '[
                    {"id": "check_calendar", "name": "check_calendar", "description": "Verifica a agenda da diretoria para um dia específico.", "tags": ["calendar", "schedule"], "examples": ["verifique minha agenda hoje", "quais minhas reuniões amanha"], "inputModes": ["text"], "outputModes": ["text"]},
                    {"id": "schedule_meeting", "name": "schedule_meeting", "description": "Agenda um novo compromisso na agenda da diretoria.", "tags": ["calendar", "schedule"], "examples": ["marque uma reunião com o joão as 14h"], "inputModes": ["text"], "outputModes": ["text"]},
                    {"id": "send_email", "name": "send_email", "description": "Envia um e-mail em nome da diretoria.", "tags": ["email", "communication"], "examples": ["envie um email para carlos cancelando a reunião"], "inputModes": ["text"], "outputModes": ["text"]},
                    {"id": "save_note", "name": "save_note", "description": "Salva uma anotação, ata de reunião ou relatório no arquivo da secretaria.", "tags": ["notes", "reports"], "examples": ["salve esta ata da reunião de marketing"], "inputModes": ["text"], "outputModes": ["text"]},
                    {"id": "read_reports", "name": "read_reports", "description": "Busca relatórios e anotações salvos pela secretaria.", "tags": ["notes", "reports"], "examples": ["busque a ata da última reunião", "quais os últimos relatórios salvos"], "inputModes": ["text"], "outputModes": ["text"]}
                ]'::json,
                'official',
                now(),
                now()
            );
        END IF;

        RAISE NOTICE 'Olivia MCP Server inserted into evo_core_mcp_servers table.', 8;
        
    ELSE
        RAISE NOTICE 'Table evo_core_mcp_servers not found. Seeder not executed.';
    END IF;
END $$;
