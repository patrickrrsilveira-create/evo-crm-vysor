-- Add 'external' agent type to the CHECK constraint
ALTER TABLE evo_core_agents
DROP CONSTRAINT IF EXISTS check_agent_type;

ALTER TABLE evo_core_agents
ADD CONSTRAINT check_agent_type CHECK (
    type::TEXT = ANY (
        ARRAY[
            'llm'::VARCHAR,
            'sequential'::VARCHAR,
            'parallel'::VARCHAR,
            'loop'::VARCHAR,
            'a2a'::VARCHAR,
            'workflow'::VARCHAR,
            'crew_ai'::VARCHAR,
            'task'::VARCHAR,
            'external'::VARCHAR
        ]::TEXT[]
    )
);
