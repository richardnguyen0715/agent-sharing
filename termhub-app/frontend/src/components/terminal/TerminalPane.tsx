import { Suspense, lazy, useState } from 'react';
import { TerminalToolbar, type TerminalPaneView } from './TerminalToolbar';
import { TerminalScreen } from './TerminalScreen';
import { StatusBar } from '@/components/layout/StatusBar';
import { useTerminal } from '@/hooks/useTerminal';
import { IconTerminal2, IconLoader2 } from '@tabler/icons-react';

// Monaco is a multi-MB dependency — only fetch it once someone opens the
// Editor tab, not as part of the base terminal bundle.
const EditorView = lazy(() => import('@/components/editor/EditorView'));

interface TerminalPaneProps {
  terminalId: string | null;
  active?: boolean;
}

export function TerminalPane({ terminalId, active = true }: TerminalPaneProps) {
  const { terminalRef, isConnected, dimensions } = useTerminal({
    terminalId,
    active,
  });
  const [view, setView] = useState<TerminalPaneView>('terminal');
  const [editorMounted, setEditorMounted] = useState(false);

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

  const handleViewChange = (next: TerminalPaneView) => {
    setView(next);
    if (next === 'editor') setEditorMounted(true);
  };

  // Keep mounted when inactive so the xterm buffer + WebSocket stay alive.
  // Hiding via CSS (not unmounting) preserves scrollback and the live prompt.
  // The editor tab follows the same rule once opened, so unsaved edits
  // survive switching back to the terminal tab.
  return (
    <div className={active ? 'flex h-full flex-col' : 'hidden'}>
      <TerminalToolbar terminalId={terminalId} view={view} onViewChange={handleViewChange} />
      <div className={view === 'terminal' ? 'flex min-h-0 flex-1 flex-col' : 'hidden'}>
        <TerminalScreen terminalRef={terminalRef} />
      </div>
      {editorMounted && (
        <div className={view === 'editor' ? 'flex min-h-0 flex-1' : 'hidden'}>
          <Suspense
            fallback={
              <div className="flex flex-1 items-center justify-center gap-2 text-sm text-terminal-output/60">
                <IconLoader2 size={16} className="animate-spin" /> Loading editor…
              </div>
            }
          >
            <EditorView terminalId={terminalId} />
          </Suspense>
        </div>
      )}
      {view === 'terminal' && (
        <StatusBar isConnected={isConnected} terminalId={terminalId} dimensions={dimensions} />
      )}
    </div>
  );
}
