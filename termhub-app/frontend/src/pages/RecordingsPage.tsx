import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  IconPlayerRecord,
  IconPlayerPlay,
  IconArrowLeft,
  IconDownload,
  IconTrash,
  IconSearch,
  IconClock,
  IconTerminal2,
  IconUser,
} from '@tabler/icons-react';

// ── Placeholder recordings data ─────────────────────────────────────

interface Recording {
  id: string;
  terminalTitle: string;
  username: string;
  duration: number;
  date: string;
  fileSize: string;
}

const PLACEHOLDER_RECORDINGS: Recording[] = [
  {
    id: 'rec-001',
    terminalTitle: 'dev-server',
    username: 'alice',
    duration: 1823,
    date: '2025-01-15T10:30:00Z',
    fileSize: '2.4 MB',
  },
  {
    id: 'rec-002',
    terminalTitle: 'prod-deploy',
    username: 'bob',
    duration: 542,
    date: '2025-01-15T09:15:00Z',
    fileSize: '890 KB',
  },
  {
    id: 'rec-003',
    terminalTitle: 'db-migration',
    username: 'alice',
    duration: 3601,
    date: '2025-01-14T16:45:00Z',
    fileSize: '5.1 MB',
  },
  {
    id: 'rec-004',
    terminalTitle: 'api-debug',
    username: 'charlie',
    duration: 267,
    date: '2025-01-14T14:20:00Z',
    fileSize: '420 KB',
  },
  {
    id: 'rec-005',
    terminalTitle: 'build-pipeline',
    username: 'bob',
    duration: 1200,
    date: '2025-01-14T11:00:00Z',
    fileSize: '1.8 MB',
  },
];

// ── Helpers ──────────────────────────────────────────────────────────

function formatDuration(seconds: number): string {
  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  const s = seconds % 60;
  if (h > 0) return `${h}h ${m}m ${s}s`;
  if (m > 0) return `${m}m ${s}s`;
  return `${s}s`;
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString(undefined, {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
}

// ── Component ────────────────────────────────────────────────────────

export function RecordingsPage() {
  const navigate = useNavigate();
  const [search, setSearch] = useState('');

  const filtered = PLACEHOLDER_RECORDINGS.filter((r) => {
    if (!search.trim()) return true;
    const q = search.toLowerCase();
    return (
      r.terminalTitle.toLowerCase().includes(q) ||
      r.username.toLowerCase().includes(q)
    );
  });

  return (
    <div className="flex h-screen flex-col bg-terminal-bg text-terminal-text">
      {/* ── Header ────────────────────────────────────────────────── */}
      <header className="flex h-14 shrink-0 items-center justify-between border-b border-terminal-border bg-terminal-bg-toolbar px-5">
        <div className="flex items-center gap-3">
          <button
            onClick={() => navigate('/')}
            className="rounded p-1 text-terminal-output hover:bg-terminal-border hover:text-terminal-text"
            aria-label="Back to main"
          >
            <IconArrowLeft size={18} />
          </button>
          <IconPlayerRecord size={20} className="text-status-red" />
          <div>
            <h1 className="text-sm font-semibold">Session Recordings</h1>
            <p className="text-[11px] text-terminal-output">Browse and replay terminal sessions</p>
          </div>
        </div>
        <div className="relative">
          <IconSearch size={14} className="absolute left-2.5 top-1/2 -translate-y-1/2 text-terminal-output" />
          <input
            type="text"
            placeholder="Search recordings…"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="rounded border border-terminal-border bg-terminal-bg py-1.5 pl-8 pr-3 text-xs text-terminal-text placeholder-terminal-output/50 outline-none transition-colors focus:border-brand"
          />
        </div>
      </header>

      {/* ── Content ───────────────────────────────────────────────── */}
      <div className="flex-1 overflow-y-auto p-5">
        {/* Info banner */}
        <div className="mb-5 rounded-lg border border-terminal-border bg-terminal-bg-toolbar px-4 py-3">
          <p className="text-xs text-terminal-output">
            <span className="font-medium text-terminal-text">{PLACEHOLDER_RECORDINGS.length} recordings</span>
            {' '}available. Recording backend is under development — showing placeholder data.
          </p>
        </div>

        {/* Recording cards */}
        {filtered.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 text-center">
            <IconPlayerRecord size={32} className="mb-3 text-terminal-output/30" />
            <p className="text-sm text-terminal-output">No recordings match your search.</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 gap-3 md:grid-cols-2 xl:grid-cols-3">
            {filtered.map((rec) => (
              <div
                key={rec.id}
                className="group rounded-lg border border-terminal-border bg-terminal-bg-toolbar p-4 transition-colors hover:border-brand/40"
              >
                {/* Top row: terminal name + play button */}
                <div className="mb-3 flex items-start justify-between">
                  <div className="flex items-center gap-2">
                    <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-brand/15">
                      <IconTerminal2 size={16} className="text-brand" />
                    </div>
                    <div>
                      <p className="font-mono text-sm font-medium text-terminal-text">
                        {rec.terminalTitle}
                      </p>
                      <p className="text-[11px] text-terminal-output">{rec.id}</p>
                    </div>
                  </div>
                  <button
                    className="flex h-8 w-8 items-center justify-center rounded-full bg-brand text-white opacity-0 transition-opacity hover:bg-brand-hover group-hover:opacity-100"
                    aria-label={`Play recording ${rec.terminalTitle}`}
                    title="Play recording"
                  >
                    <IconPlayerPlay size={14} />
                  </button>
                </div>

                {/* Metadata */}
                <div className="space-y-1.5 text-xs text-terminal-output">
                  <div className="flex items-center gap-2">
                    <IconUser size={13} className="shrink-0" />
                    <span>{rec.username}</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <IconClock size={13} className="shrink-0" />
                    <span>{formatDuration(rec.duration)}</span>
                    <span className="text-terminal-output/40">·</span>
                    <span>{formatDate(rec.date)}</span>
                  </div>
                </div>

                {/* Footer */}
                <div className="mt-3 flex items-center justify-between border-t border-terminal-border/50 pt-3">
                  <span className="font-mono text-[11px] text-terminal-output/50">{rec.fileSize}</span>
                  <div className="flex items-center gap-1">
                    <button
                      className="rounded p-1 text-terminal-output hover:bg-terminal-border hover:text-terminal-text"
                      aria-label="Download"
                      title="Download"
                    >
                      <IconDownload size={14} />
                    </button>
                    <button
                      className="rounded p-1 text-terminal-output hover:bg-terminal-border hover:text-terminal-error"
                      aria-label="Delete"
                      title="Delete"
                    >
                      <IconTrash size={14} />
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
