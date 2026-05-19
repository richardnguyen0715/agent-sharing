import type { ReactNode } from 'react';
import { Topbar } from './Topbar';
import { Sidebar } from './Sidebar';
import { RightPanel } from './RightPanel';
import { useUIStore } from '@/stores/uiStore';

interface AppLayoutProps {
  children: ReactNode;
}

export function AppLayout({ children }: AppLayoutProps) {
  const sidebarOpen = useUIStore((s) => s.sidebarOpen);
  const rightPanelOpen = useUIStore((s) => s.rightPanelOpen);

  return (
    <div className="flex h-screen flex-col bg-terminal-bg text-terminal-text">
      <Topbar />

      <div className="flex flex-1 overflow-hidden">
        {sidebarOpen && <Sidebar />}

        <main className="flex flex-1 flex-col overflow-hidden">
          {children}
        </main>

        {rightPanelOpen && <RightPanel />}
      </div>
    </div>
  );
}
