import { TerminalToolbar } from './TerminalToolbar';
import { TerminalScreen } from './TerminalScreen';
import { StatusBar } from '@/components/layout/StatusBar';
import { useTerminal } from '@/hooks/useTerminal';
import { IconTerminal2 } from '@tabler/icons-react';

interface TerminalPaneProps {
  terminalId: string | null;
  active?: boolean;
}

export function TerminalPane({ terminalId, active = true }: TerminalPaneProps) {
  const { terminalRef, isConnected, dimensions } = useTerminal({
    terminalId,
    active,
  });

  if (!terminalId) {
    return (
      <div className="flex h-full flex-col items-center justify-center bg-terminal-bg">
        <IconTerminal2 size={48} className="mb-4 text-terminal-output/20" />
        <p className="text-sm text-terminal-output/50">No terminal selected</p>
        <p className="mt-1 text-xs text-terminal-output/30">
          Open a terminal from the sidebar or create a new one
        </p>
      </div>
    );
  }

  // Keep mounted when inactive so the xterm buffer + WebSocket stay alive.
  // Hiding via CSS (not unmounting) preserves scrollback and the live prompt.
  return (
    <div className={active ? 'flex h-full flex-col' : 'hidden'}>
      <TerminalToolbar terminalId={terminalId} />
      <TerminalScreen terminalRef={terminalRef} />
      <StatusBar
        isConnected={isConnected}
        terminalId={terminalId}
        dimensions={dimensions}
      />
    </div>
  );
}
