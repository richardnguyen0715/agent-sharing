CREATE TABLE IF NOT EXISTS terminal_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    terminal_id UUID NOT NULL REFERENCES terminals(id) ON DELETE CASCADE,
    websocket_id VARCHAR(255) NOT NULL,
    connected_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    disconnected_at TIMESTAMPTZ
);

CREATE INDEX idx_terminal_sessions_terminal_id ON terminal_sessions (terminal_id);
CREATE INDEX idx_terminal_sessions_connected_at ON terminal_sessions (connected_at);
