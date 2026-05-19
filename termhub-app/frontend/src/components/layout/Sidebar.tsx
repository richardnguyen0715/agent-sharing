import {
  IconPlus,
  IconTerminal2,
  IconPlayerRecord,
  IconFolder,
  IconLock,
  IconBrandGithub,
  IconCloud,
  IconShieldLock,
  IconDashboard,
} from '@tabler/icons-react';
import { useNavigate } from 'react-router-dom';
import { useUIStore } from '@/stores/uiStore';
import { useTerminalStore } from '@/stores/terminalStore';
import { useAuthStore } from '@/stores/authStore';

function StatusDot({ status }: { status: string }) {
  const color =
    status === 'active'
      ? 'bg-status-amber'
      : status === 'terminated'
        ? 'bg-status-red'
        : 'bg-status-green';
  return <span className={`inline-block h-2 w-2 shrink-0 rounded-full ${color}`} />;
}

interface CredentialItemProps {
  icon: React.ReactNode;
  label: string;
  status: 'connected' | 'disconnected';
}

function CredentialItem({ icon, label, status }: CredentialItemProps) {
  return (
    <div className="flex items-center gap-2 px-2 py-1 text-xs text-terminal-output">
      {icon}
      <span className="flex-1">{label}</span>
      <span
        className={`inline-block h-1.5 w-1.5 rounded-full ${
          status === 'connected' ? 'bg-status-green' : 'bg-terminal-border'
        }`}
      />
    </div>
  );
}

export function Sidebar() {
  const openModal = useUIStore((s) => s.openModal);
  const terminals = useTerminalStore((s) => s.terminals);
  const activeTabId = useTerminalStore((s) => s.activeTabId);
  const addTab = useTerminalStore((s) => s.addTab);
  const user = useAuthStore((s) => s.user);
  const navigate = useNavigate();

  const handleTerminalClick = (terminal: (typeof terminals)[number]) => {
    addTab({
      terminalId: terminal.id,
      title: terminal.title,
      status: terminal.status === 'active' ? 'active' : 'idle',
      colorTag: terminal.colorTag,
    });
  };

  return (
    <aside className="flex w-52 shrink-0 flex-col border-r border-terminal-border bg-terminal-bg-toolbar">
      {/* My Terminals */}
      <div className="flex items-center justify-between border-b border-terminal-border px-3 py-2">
        <span className="text-xs font-medium uppercase tracking-wide text-terminal-output">
          My Terminals
        </span>
        <button
          onClick={() => openModal('newTerminal')}
          className="rounded p-0.5 text-terminal-output hover:bg-terminal-border hover:text-terminal-text"
          aria-label="New terminal"
        >
          <IconPlus size={14} />
        </button>
      </div>

      <div className="flex-1 overflow-y-auto">
        {/* Terminal list */}
        <div className="space-y-0.5 px-1.5 py-1.5">
          {terminals.length === 0 && (
            <p className="px-2 py-3 text-center text-xs text-terminal-output/60">
              No terminals yet
            </p>
          )}
          {terminals
            .filter((t) => t.status !== 'terminated')
            .map((terminal) => {
              const isActive = terminal.id === activeTabId;
              return (
                <button
                  key={terminal.id}
                  onClick={() => handleTerminalClick(terminal)}
                  className={`flex w-full items-center gap-2 rounded px-2 py-1.5 text-xs transition-colors ${
                    isActive
                      ? 'bg-brand/15 text-terminal-text'
                      : 'text-terminal-output hover:bg-terminal-border/50 hover:text-terminal-text'
                  }`}
                >
                  <IconTerminal2 size={14} className="shrink-0" />
                  <span className="flex-1 truncate text-left">{terminal.title}</span>
                  <StatusDot status={terminal.status} />
                </button>
              );
            })}
        </div>

        {/* Shared section */}
        <div className="border-t border-terminal-border px-3 py-2">
          <span className="text-xs font-medium uppercase tracking-wide text-terminal-output">
            Shared
          </span>
        </div>
        <div className="space-y-0.5 px-1.5 pb-1.5">
          <div className="flex items-center gap-2 rounded px-2 py-1.5 text-xs text-terminal-output">
            <IconFolder size={14} className="shrink-0" />
            <span className="flex-1 truncate">/workspace/shared</span>
          </div>
          <div className="flex items-center gap-2 rounded px-2 py-1.5 text-xs text-terminal-output">
            <IconLock size={14} className="shrink-0" />
            <span className="flex-1 truncate">.env readonly</span>
          </div>
        </div>

        {/* Actions */}
        <div className="border-t border-terminal-border px-3 py-2">
          <span className="text-xs font-medium uppercase tracking-wide text-terminal-output">
            Actions
          </span>
        </div>
        <div className="space-y-0.5 px-1.5 pb-1.5">
          <button
            onClick={() => openModal('newTerminal')}
            className="flex w-full items-center gap-2 rounded px-2 py-1.5 text-xs text-terminal-output hover:bg-terminal-border/50 hover:text-terminal-text"
          >
            <IconPlus size={14} />
            New terminal
          </button>
          <button
            onClick={() => navigate('/recordings')}
            className="flex w-full items-center gap-2 rounded px-2 py-1.5 text-xs text-terminal-output hover:bg-terminal-border/50 hover:text-terminal-text"
          >
            <IconPlayerRecord size={14} />
            Recordings
          </button>
          {user?.role === 'admin' && (
            <button
              onClick={() => navigate('/admin')}
              className="flex w-full items-center gap-2 rounded px-2 py-1.5 text-xs text-terminal-output hover:bg-terminal-border/50 hover:text-terminal-text"
            >
              <IconDashboard size={14} />
              Admin dashboard
            </button>
          )}
        </div>
      </div>

      {/* Host credentials */}
      <div className="border-t border-terminal-border">
        <div className="px-3 py-2">
          <span className="text-xs font-medium uppercase tracking-wide text-terminal-output">
            Host Credentials
          </span>
        </div>
        <div className="pb-2">
          <CredentialItem
            icon={<IconBrandGithub size={13} />}
            label="GitHub"
            status="connected"
          />
          <CredentialItem
            icon={<IconCloud size={13} />}
            label="AWS"
            status="connected"
          />
          <CredentialItem
            icon={<IconShieldLock size={13} />}
            label="VPN"
            status="disconnected"
          />
        </div>
      </div>
    </aside>
  );
}
