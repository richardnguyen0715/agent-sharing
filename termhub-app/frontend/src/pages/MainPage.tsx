import { useEffect } from 'react';
import { AppLayout } from '@/components/layout/AppLayout';
import { TerminalPane } from '@/components/terminal/TerminalPane';
import { useTerminalStore } from '@/stores/terminalStore';

export function MainPage() {
  const fetchTerminals = useTerminalStore((s) => s.fetchTerminals);
  const tabs = useTerminalStore((s) => s.tabs);
  const activeTabId = useTerminalStore((s) => s.activeTabId);

  useEffect(() => {
    void fetchTerminals();
  }, [fetchTerminals]);

  return (
    <AppLayout>
      {tabs.length === 0 ? (
        <TerminalPane terminalId={null} />
      ) : (
        <div className="relative h-full">
          {tabs.map((tab) => (
            <div key={tab.terminalId} className="absolute inset-0">
              <TerminalPane
                terminalId={tab.terminalId}
                active={tab.terminalId === activeTabId}
              />
            </div>
          ))}
        </div>
      )}
    </AppLayout>
  );
}
