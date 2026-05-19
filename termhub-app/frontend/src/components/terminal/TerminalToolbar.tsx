import { useState } from 'react';
import {
  IconX,
  IconMaximize,
  IconInfoCircle,
  IconCopy,
  IconPlayerRecord,
  IconTerminal2,
  IconLoader2,
} from '@tabler/icons-react';
import { useTerminalStore } from '@/stores/terminalStore';
import { useUIStore } from '@/stores/uiStore';

interface TerminalToolbarProps {
  terminalId: string;
}

export function TerminalToolbar({ terminalId }: TerminalToolbarProps) {
  const deleteTerminal = useTerminalStore((s) => s.deleteTerminal);
  const terminals = useTerminalStore((s) => s.terminals);
  const toggleRightPanel = useUIStore((s) => s.toggleRightPanel);
  const [killing, setKilling] = useState(false);

  const terminal = terminals.find((t) => t.id === terminalId);
  const shellName = terminal?.shell.split('/').pop() ?? 'shell';
  const cwd = terminal?.cwd ?? '~';

  const handleKill = async () => {
    if (killing) return;
    setKilling(true);
    try {
      await deleteTerminal(terminalId);
    } catch {
      setKilling(false);
    }
  };

  return (
    <div className="flex h-9 shrink-0 items-center justify-between border-b border-terminal-border bg-terminal-bg-toolbar px-3">
      <div className="flex items-center gap-2">
        <IconTerminal2 size={14} className="text-terminal-prompt" />
        <span className="font-mono text-xs text-terminal-output">
          {shellName} — <span className="text-terminal-path">{cwd}</span>
        </span>
      </div>
      <div className="flex items-center gap-0.5">
        <button
          className="rounded p-1 text-terminal-output hover:bg-terminal-border hover:text-terminal-text"
          aria-label="Copy"
        >
          <IconCopy size={14} />
        </button>
        <button
          className="rounded p-1 text-terminal-output hover:bg-terminal-border hover:text-terminal-text"
          aria-label="Record"
        >
          <IconPlayerRecord size={14} />
        </button>
        <button
          onClick={toggleRightPanel}
          className="rounded p-1 text-terminal-output hover:bg-terminal-border hover:text-terminal-text"
          aria-label="Terminal info"
        >
          <IconInfoCircle size={14} />
        </button>
        <button
          className="rounded p-1 text-terminal-output hover:bg-terminal-border hover:text-terminal-text"
          aria-label="Maximize"
        >
          <IconMaximize size={14} />
        </button>
        <button
          onClick={handleKill}
          disabled={killing}
          className="rounded p-1 text-terminal-output hover:bg-terminal-error hover:text-white disabled:opacity-50"
          aria-label="Kill terminal"
          title="Kill terminal"
        >
          {killing ? <IconLoader2 size={14} className="animate-spin" /> : <IconX size={14} />}
        </button>
      </div>
    </div>
  );
}
