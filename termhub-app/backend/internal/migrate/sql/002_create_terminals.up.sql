CREATE TABLE IF NOT EXISTS terminals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL DEFAULT 'Terminal',
    cwd VARCHAR(1024) NOT NULL DEFAULT '~',
    shell VARCHAR(255) NOT NULL DEFAULT '/bin/bash',
    status VARCHAR(50) NOT NULL DEFAULT 'idle'
        CHECK (status IN ('active', 'idle', 'terminated')),
    color_tag VARCHAR(7) NOT NULL DEFAULT '#00FF00',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_terminals_owner_id ON terminals (owner_id);
CREATE INDEX idx_terminals_status ON terminals (status);
