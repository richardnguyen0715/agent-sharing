import { useEffect } from 'react';
import { AppLayout } from '@/components/layout/AppLayout';
import { TerminalPane } from '@/components/terminal/TerminalPane';
import { useTerminalStore } from '@/stores/terminalStore';

export function MainPage() {
  const fetchTerminals = useTerminalStore((s) => s.fetchTerminals);
  const activeTabId = useTerminalStore((s) => s.activeTabId);

  useEffect(() => {
    void fetchTerminals();
  }, [fetchTerminals]);

  return (
    <AppLayout>
      <TerminalPane terminalId={activeTabId} />
    </AppLayout>
  );
}
