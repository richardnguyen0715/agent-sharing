import { useEffect, useState, useRef } from 'react';

interface StatusBarProps {
  isConnected: boolean;
  terminalId: string | null;
  dimensions: { cols: number; rows: number };
}

function formatDuration(seconds: number): string {
  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  const s = seconds % 60;
  return [h, m, s].map((v) => String(v).padStart(2, '0')).join(':');
}

export function StatusBar({ isConnected, terminalId, dimensions }: StatusBarProps) {
  const [elapsed, setElapsed] = useState(0);
  const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);

  useEffect(() => {
    if (isConnected) {
      setElapsed(0);
      intervalRef.current = setInterval(() => {
        setElapsed((prev) => prev + 1);
      }, 1000);
    } else {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
        intervalRef.current = null;
      }
    }
    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
      }
    };
  }, [isConnected, terminalId]);

  return (
    <footer className="flex h-6 shrink-0 items-center justify-between border-t border-terminal-border bg-terminal-bg-toolbar px-3 font-mono text-[11px] text-terminal-output">
      <div className="flex items-center gap-4">
        <span className="flex items-center gap-1.5">
          <span
            className={`inline-block h-2 w-2 rounded-full ${
              isConnected ? 'bg-status-green' : 'bg-status-red'
            }`}
          />
          {isConnected ? 'connected' : 'disconnected'}
        </span>

        {terminalId && (
          <span className="text-terminal-output/60">
            pty: {terminalId.slice(0, 8)}
          </span>
        )}

        {isConnected && (
          <span className="text-terminal-output/60">{formatDuration(elapsed)}</span>
        )}
      </div>

      <div className="flex items-center gap-4">
        {terminalId && (
          <span className="text-terminal-output/60">
            {dimensions.cols}×{dimensions.rows}
          </span>
        )}
      </div>
    </footer>
  );
}
