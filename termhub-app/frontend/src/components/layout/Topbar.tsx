import { useState, useRef, useEffect } from 'react';
import {
  IconTerminal2,
  IconMenu2,
  IconPlus,
  IconX,
  IconChevronDown,
  IconLogout,
  IconSettings,
} from '@tabler/icons-react';
import { useUIStore } from '@/stores/uiStore';
import { useTerminalStore } from '@/stores/terminalStore';
import { useAuthStore } from '@/stores/authStore';
import { Avatar } from '@/components/ui/Avatar';
import { Badge } from '@/components/ui/Badge';

function StatusDot({ status }: { status: string }) {
  const color = status === 'active' ? 'bg-status-amber' : 'bg-status-green';
  return <span className={`inline-block h-2 w-2 rounded-full ${color}`} />;
}

export function Topbar() {
  const toggleSidebar = useUIStore((s) => s.toggleSidebar);
  const openModal = useUIStore((s) => s.openModal);
  const tabs = useTerminalStore((s) => s.tabs);
  const activeTabId = useTerminalStore((s) => s.activeTabId);
  const setActiveTab = useTerminalStore((s) => s.setActiveTab);
  const removeTab = useTerminalStore((s) => s.removeTab);
  const user = useAuthStore((s) => s.user);
  const logout = useAuthStore((s) => s.logout);

  const [userMenuOpen, setUserMenuOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);

  // Close user menu on outside click
  useEffect(() => {
    if (!userMenuOpen) return;
    const handler = (e: MouseEvent) => {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) {
        setUserMenuOpen(false);
      }
    };
    document.addEventListener('mousedown', handler);
    return () => document.removeEventListener('mousedown', handler);
  }, [userMenuOpen]);

  return (
    <header className="flex h-12 shrink-0 items-center border-b border-terminal-border bg-terminal-bg-toolbar">
      {/* Left: hamburger + brand */}
      <div className="flex items-center gap-2 px-3">
        <button
          onClick={toggleSidebar}
          className="rounded p-1 text-terminal-output hover:bg-terminal-border hover:text-terminal-text"
          aria-label="Toggle sidebar"
        >
          <IconMenu2 size={18} />
        </button>
        <div className="flex items-center gap-1.5">
          <IconTerminal2 size={20} className="text-brand" />
          <span className="text-sm font-semibold">TermHub</span>
        </div>
      </div>

      {/* Center: tab bar */}
      <div className="flex flex-1 items-end gap-0 overflow-x-auto px-1">
        {tabs.map((tab) => {
          const isActive = tab.terminalId === activeTabId;
          return (
            <button
              key={tab.terminalId}
              onClick={() => setActiveTab(tab.terminalId)}
              className={`group relative flex items-center gap-1.5 border-b-2 px-3 py-2.5 text-xs transition-colors ${
                isActive
                  ? 'border-brand text-terminal-text'
                  : 'border-transparent text-terminal-output hover:text-terminal-text'
              }`}
            >
              <StatusDot status={tab.status} />
              <span className="max-w-[120px] truncate">{tab.title}</span>
              <span
                role="button"
                tabIndex={0}
                onClick={(e) => {
                  e.stopPropagation();
                  removeTab(tab.terminalId);
                }}
                onKeyDown={(e) => {
                  if (e.key === 'Enter') {
                    e.stopPropagation();
                    removeTab(tab.terminalId);
                  }
                }}
                className="ml-0.5 rounded p-0.5 opacity-0 transition-opacity hover:bg-terminal-border group-hover:opacity-100"
                aria-label={`Close ${tab.title}`}
              >
                <IconX size={12} />
              </span>
            </button>
          );
        })}

        {/* New tab button */}
        <button
          onClick={() => openModal('newTerminal')}
          className="flex items-center rounded p-1.5 text-terminal-output hover:bg-terminal-border hover:text-terminal-text"
          aria-label="New terminal"
        >
          <IconPlus size={16} />
        </button>
      </div>

      {/* Right: user section */}
      {user && (
        <div className="relative px-3" ref={menuRef}>
          <button
            onClick={() => setUserMenuOpen((v) => !v)}
            className="flex items-center gap-2 rounded px-2 py-1 text-sm hover:bg-terminal-border"
          >
            <Avatar name={user.username} size="sm" />
            <span className="hidden text-xs text-terminal-text sm:inline">{user.username}</span>
            <Badge>{user.role}</Badge>
            <IconChevronDown size={14} className="text-terminal-output" />
          </button>

          {userMenuOpen && (
            <div className="absolute right-0 top-full z-50 mt-1 w-48 rounded-lg border border-terminal-border bg-terminal-bg-toolbar py-1 shadow-xl shadow-black/40">
              <button
                onClick={() => {
                  setUserMenuOpen(false);
                  openModal('settings');
                }}
                className="flex w-full items-center gap-2 px-3 py-2 text-sm text-terminal-output hover:bg-terminal-border hover:text-terminal-text"
              >
                <IconSettings size={16} />
                Settings
              </button>
              <div className="my-1 border-t border-terminal-border" />
              <button
                onClick={() => {
                  setUserMenuOpen(false);
                  logout();
                }}
                className="flex w-full items-center gap-2 px-3 py-2 text-sm text-terminal-error hover:bg-terminal-border"
              >
                <IconLogout size={16} />
                Sign out
              </button>
            </div>
          )}
        </div>
      )}
    </header>
  );
}
