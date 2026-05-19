import {
  IconX,
  IconInfoCircle,
  IconPlayerRecord,
  IconShare,
  IconClock,
  IconTerminal2,
} from '@tabler/icons-react';
import { useUIStore } from '@/stores/uiStore';
import { useTerminalStore } from '@/stores/terminalStore';
import { useAuthStore } from '@/stores/authStore';

function SectionHeader({ children }: { children: React.ReactNode }) {
  return (
    <h3 className="mb-2 text-[10px] font-semibold uppercase tracking-wider text-terminal-output/60">
      {children}
    </h3>
  );
}

function InfoRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex items-start justify-between gap-2 py-0.5">
      <span className="text-xs text-terminal-output/70">{label}</span>
      <span className="text-right font-mono text-xs text-terminal-text">{value}</span>
    </div>
  );
}

export function RightPanel() {
  const setRightPanelOpen = useUIStore((s) => s.setRightPanelOpen);
  const activeTabId = useTerminalStore((s) => s.activeTabId);
  const terminals = useTerminalStore((s) => s.terminals);
  const user = useAuthStore((s) => s.user);

  const activeTerminal = terminals.find((t) => t.id === activeTabId);

  return (
    <aside className="flex w-56 shrink-0 flex-col border-l border-terminal-border bg-terminal-bg-toolbar">
      {/* Header */}
      <div className="flex items-center justify-between border-b border-terminal-border px-3 py-2">
        <div className="flex items-center gap-1.5">
          <IconInfoCircle size={14} className="text-terminal-output" />
          <span className="text-xs font-medium uppercase tracking-wide text-terminal-output">
            Session Info
          </span>
        </div>
        <button
          onClick={() => setRightPanelOpen(false)}
          className="rounded p-0.5 text-terminal-output hover:bg-terminal-border hover:text-terminal-text"
          aria-label="Close panel"
        >
          <IconX size={14} />
        </button>
      </div>

      <div className="flex-1 overflow-y-auto p-3 space-y-4">
        {activeTerminal ? (
          <>
            {/* Terminal info */}
            <div>
              <SectionHeader>Terminal</SectionHeader>
              <div className="space-y-1">
                <InfoRow label="ID" value={activeTerminal.id.slice(0, 8)} />
                <InfoRow label="Owner" value={user?.username ?? '—'} />
                <InfoRow label="Shell" value={activeTerminal.shell} />
                <InfoRow label="CWD" value={activeTerminal.cwd} />
              </div>
            </div>

            {/* Processes */}
            <div>
              <SectionHeader>Processes</SectionHeader>
              <div className="space-y-1">
                <div className="flex items-center gap-2 rounded bg-terminal-bg/50 px-2 py-1">
                  <IconTerminal2 size={12} className="text-terminal-prompt" />
                  <span className="flex-1 truncate font-mono text-xs text-terminal-text">
                    {activeTerminal.shell.split('/').pop()}
                  </span>
                  <span className="text-[10px] text-terminal-output/50">PID 1</span>
                </div>
              </div>
            </div>

            {/* Activity */}
            <div>
              <SectionHeader>Activity</SectionHeader>
              <div className="space-y-1.5">
                <div className="flex items-start gap-2">
                  <IconClock size={12} className="mt-0.5 shrink-0 text-terminal-output/50" />
                  <div>
                    <p className="font-mono text-xs text-terminal-text">Session started</p>
                    <p className="text-[10px] text-terminal-output/50">
                      {new Date(activeTerminal.createdAt).toLocaleTimeString()}
                    </p>
                  </div>
                </div>
              </div>
            </div>

            {/* Quick actions */}
            <div>
              <SectionHeader>Quick Actions</SectionHeader>
              <div className="space-y-1">
                <button className="flex w-full items-center gap-2 rounded border border-terminal-border px-2 py-1.5 text-xs text-terminal-output transition-colors hover:bg-terminal-border hover:text-terminal-text">
                  <IconPlayerRecord size={14} className="text-status-red" />
                  Start recording
                </button>
                <button className="flex w-full items-center gap-2 rounded border border-terminal-border px-2 py-1.5 text-xs text-terminal-output transition-colors hover:bg-terminal-border hover:text-terminal-text">
                  <IconShare size={14} className="text-terminal-path" />
                  Share (readonly)
                </button>
              </div>
            </div>
          </>
        ) : (
          <div className="flex flex-col items-center justify-center py-8 text-center">
            <IconInfoCircle size={24} className="mb-2 text-terminal-output/30" />
            <p className="text-xs text-terminal-output/50">No active session</p>
          </div>
        )}
      </div>
    </aside>
  );
}
