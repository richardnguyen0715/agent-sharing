import { useState, useEffect, useCallback, useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import {
  IconDashboard,
  IconRefresh,
  IconDownload,
  IconUserPlus,
  IconTerminal2,
  IconUsers,
  IconAlertTriangle,
  IconPlayerRecord,
  IconSearch,
  IconEye,
  IconPlayerStop,
  IconLoader2,
  IconAlertCircle,
  IconInfoCircle,
  IconCircleCheck,
  IconCircleX,
  IconArrowLeft,
} from '@tabler/icons-react';
import { Avatar } from '@/components/ui/Avatar';
import * as api from '@/services/api';
import type { Terminal, User, AuditLog } from '@/types';

// ── Types ────────────────────────────────────────────────────────────

type TerminalFilter = 'all' | 'connected' | 'idle' | 'warnings';

interface AuditSeverity {
  level: 'warn' | 'info' | 'ok' | 'err';
  label: string;
}

// ── Helpers ──────────────────────────────────────────────────────────

function classifyAudit(action: string): AuditSeverity {
  const a = action.toLowerCase();
  if (a.includes('block') || a.includes('denied') || a.includes('warn'))
    return { level: 'warn', label: 'Warning' };
  if (a.includes('disconnect') || a.includes('error') || a.includes('fail'))
    return { level: 'err', label: 'Error' };
  if (a.includes('create') || a.includes('record') || a.includes('start') || a.includes('success'))
    return { level: 'ok', label: 'Success' };
  return { level: 'info', label: 'Info' };
}

const AUDIT_ICON_CLASSES: Record<string, string> = {
  warn: 'bg-[#BA7517]/20 text-[#D29922]',
  info: 'bg-terminal-path/20 text-terminal-path',
  ok: 'bg-[#1D9E75]/20 text-[#3FB950]',
  err: 'bg-[#E24B4A]/20 text-[#F85149]',
};

function AuditIcon({ level }: { level: string }) {
  const cls = `flex h-6 w-6 shrink-0 items-center justify-center rounded-full ${AUDIT_ICON_CLASSES[level] ?? ''}`;
  switch (level) {
    case 'warn':
      return <span className={cls}><IconAlertTriangle size={13} /></span>;
    case 'err':
      return <span className={cls}><IconCircleX size={13} /></span>;
    case 'ok':
      return <span className={cls}><IconCircleCheck size={13} /></span>;
    default:
      return <span className={cls}><IconInfoCircle size={13} /></span>;
  }
}

const ROLE_BADGE: Record<string, string> = {
  admin: 'bg-[#EEEDFE] text-[#3C3489]',
  operator: 'bg-[#E6F1FB] text-[#0C447C]',
  viewer: 'bg-[#F1EFE8] text-[#444441]',
};

function RoleBadge({ role }: { role: string }) {
  const cls = ROLE_BADGE[role] ?? 'bg-terminal-border text-terminal-text';
  return (
    <span className={`inline-flex items-center rounded-full px-2 py-0.5 text-[11px] font-medium capitalize ${cls}`}>
      {role}
    </span>
  );
}

function StatusDot({ status }: { status: string }) {
  const color =
    status === 'active' ? 'bg-[#1D9E75]' : status === 'idle' ? 'bg-[#BA7517]' : 'bg-[#E24B4A]';
  return <span className={`inline-block h-2 w-2 rounded-full ${color}`} />;
}

function formatUptime(createdAt: string): string {
  const diff = Math.max(0, Math.floor((Date.now() - new Date(createdAt).getTime()) / 1000));
  if (diff < 60) return `${diff}s`;
  if (diff < 3600) return `${Math.floor(diff / 60)}m`;
  const h = Math.floor(diff / 3600);
  const m = Math.floor((diff % 3600) / 60);
  return `${h}h ${m}m`;
}

function formatTime(dateStr: string): string {
  return new Date(dateStr).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });
}

// ── Metric Card ──────────────────────────────────────────────────────

interface MetricCardProps {
  icon: React.ReactNode;
  label: string;
  value: string | number;
  sub: string;
  subColor?: string;
}

function MetricCard({ icon, label, value, sub, subColor = 'text-terminal-output' }: MetricCardProps) {
  return (
    <div className="rounded-lg border border-terminal-border bg-terminal-bg-toolbar p-4">
      <div className="mb-2 flex items-center gap-2 text-terminal-output">
        {icon}
        <span className="text-xs font-medium uppercase tracking-wide">{label}</span>
      </div>
      <p className="text-2xl font-bold text-terminal-text">{value}</p>
      <p className={`mt-0.5 text-xs ${subColor}`}>{sub}</p>
    </div>
  );
}

// ── Main Component ───────────────────────────────────────────────────

export function AdminPage() {
  const navigate = useNavigate();

  // Data
  const [terminals, setTerminals] = useState<Terminal[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [auditLogs, setAuditLogs] = useState<AuditLog[]>([]);

  // UI state
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filter, setFilter] = useState<TerminalFilter>('all');
  const [search, setSearch] = useState('');
  const [refreshing, setRefreshing] = useState(false);

  // Build user lookup
  const userMap = useMemo(() => {
    const map = new Map<string, User>();
    for (const u of users) map.set(u.id, u);
    return map;
  }, [users]);

  // Fetch all data
  const fetchData = useCallback(async (isRefresh = false) => {
    if (isRefresh) setRefreshing(true);
    else setLoading(true);
    setError(null);

    try {
      const [t, u, a] = await Promise.all([
        api.adminListTerminals(),
        api.adminListUsers(),
        api.adminGetAuditLogs(100),
      ]);
      setTerminals(t);
      setUsers(u);
      setAuditLogs(a);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load admin data');
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  }, []);

  // Initial load + auto-refresh
  useEffect(() => {
    void fetchData();
    const interval = setInterval(() => void fetchData(true), 10_000);
    return () => clearInterval(interval);
  }, [fetchData]);

  // Metrics
  const activeTerminals = terminals.filter((t) => t.status === 'active');
  const idleTerminals = terminals.filter((t) => t.status === 'idle');
  const liveTerminals = terminals.filter((t) => t.status !== 'terminated');
  const warningLogs = auditLogs.filter((l) => classifyAudit(l.action).level === 'warn');
  const blockedCount = auditLogs.filter((l) => l.action.toLowerCase().includes('block')).length;
  // Recording: use a heuristic — count terminals that have a "recording" audit action
  const recordingIds = new Set(
    auditLogs
      .filter((l) => l.action.toLowerCase().includes('record') && l.terminalId)
      .map((l) => l.terminalId),
  );
  const recordingCount = Math.min(recordingIds.size, liveTerminals.length);
  const recordingPct = liveTerminals.length > 0 ? Math.round((recordingCount / liveTerminals.length) * 100) : 0;

  // Filtered terminals
  const filteredTerminals = useMemo(() => {
    let list = terminals;

    // Status filter
    if (filter === 'connected') list = list.filter((t) => t.status === 'active');
    else if (filter === 'idle') list = list.filter((t) => t.status === 'idle');
    else if (filter === 'warnings') {
      const warnTermIds = new Set(warningLogs.map((l) => l.terminalId).filter(Boolean));
      list = list.filter((t) => warnTermIds.has(t.id));
    }

    // Search
    if (search.trim()) {
      const q = search.toLowerCase();
      list = list.filter((t) => {
        const owner = userMap.get(t.ownerId);
        return (
          t.title.toLowerCase().includes(q) ||
          t.id.toLowerCase().includes(q) ||
          t.cwd.toLowerCase().includes(q) ||
          (owner?.username.toLowerCase().includes(q) ?? false)
        );
      });
    }

    return list;
  }, [terminals, filter, search, warningLogs, userMap]);

  // Warnings for right panel
  const activeWarnings = warningLogs.slice(0, 5);

  // ── Loading / Error states ─────────────────────────────────────────

  if (loading) {
    return (
      <div className="flex h-screen items-center justify-center bg-terminal-bg">
        <div className="flex flex-col items-center gap-3">
          <IconLoader2 size={32} className="animate-spin text-brand" />
          <p className="text-sm text-terminal-output">Loading admin dashboard…</p>
        </div>
      </div>
    );
  }

  if (error && terminals.length === 0) {
    return (
      <div className="flex h-screen items-center justify-center bg-terminal-bg">
        <div className="flex flex-col items-center gap-3 text-center">
          <IconAlertCircle size={32} className="text-terminal-error" />
          <p className="text-sm text-terminal-error">{error}</p>
          <button
            onClick={() => void fetchData()}
            className="mt-2 rounded bg-brand px-4 py-2 text-sm font-medium text-white hover:bg-brand-hover"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  // ── Filter pills ──────────────────────────────────────────────────

  const FILTERS: { key: TerminalFilter; label: string; count: number }[] = [
    { key: 'all', label: 'All', count: terminals.length },
    { key: 'connected', label: 'Connected', count: activeTerminals.length },
    { key: 'idle', label: 'Idle', count: idleTerminals.length },
    { key: 'warnings', label: 'Warnings', count: warningLogs.length },
  ];

  return (
    <div className="flex h-screen flex-col bg-terminal-bg text-terminal-text">
      {/* ── Top bar ─────────────────────────────────────────────────── */}
      <header className="flex h-14 shrink-0 items-center justify-between border-b border-terminal-border bg-terminal-bg-toolbar px-5">
        <div className="flex items-center gap-3">
          <button
            onClick={() => navigate('/')}
            className="rounded p-1 text-terminal-output hover:bg-terminal-border hover:text-terminal-text"
            aria-label="Back to main"
          >
            <IconArrowLeft size={18} />
          </button>
          <IconDashboard size={20} className="text-brand" />
          <h1 className="text-sm font-semibold">Admin dashboard</h1>
          <span className="text-xs text-terminal-output">— live view</span>
        </div>
        <div className="flex items-center gap-2">
          <button
            onClick={() => void fetchData(true)}
            disabled={refreshing}
            className="flex items-center gap-1.5 rounded border border-terminal-border px-3 py-1.5 text-xs font-medium text-terminal-output transition-colors hover:bg-terminal-border hover:text-terminal-text disabled:opacity-50"
          >
            <IconRefresh size={14} className={refreshing ? 'animate-spin' : ''} />
            Refresh
          </button>
          <button className="flex items-center gap-1.5 rounded border border-terminal-border px-3 py-1.5 text-xs font-medium text-terminal-output transition-colors hover:bg-terminal-border hover:text-terminal-text">
            <IconDownload size={14} />
            Export logs
          </button>
          <button className="flex items-center gap-1.5 rounded bg-brand px-3 py-1.5 text-xs font-semibold text-white transition-colors hover:bg-brand-hover">
            <IconUserPlus size={14} />
            Add user
          </button>
        </div>
      </header>

      {/* ── Content ─────────────────────────────────────────────────── */}
      <div className="flex-1 overflow-y-auto p-5">
        {/* Error banner (non-blocking) */}
        {error && (
          <div className="mb-4 rounded-lg border border-terminal-error/30 bg-terminal-error/10 px-4 py-2.5 text-sm text-terminal-error">
            {error}
          </div>
        )}

        {/* ── Metrics row ───────────────────────────────────────────── */}
        <div className="mb-5 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <MetricCard
            icon={<IconTerminal2 size={16} />}
            label="Active terminals"
            value={activeTerminals.length}
            sub={`${liveTerminals.length} connected now`}
            subColor="text-[#3FB950]"
          />
          <MetricCard
            icon={<IconUsers size={16} />}
            label="Users online"
            value={users.length}
            sub={`of ${users.length} total`}
          />
          <MetricCard
            icon={<IconAlertTriangle size={16} />}
            label="Warnings today"
            value={warningLogs.length}
            sub={`${blockedCount} blocked`}
            subColor="text-[#D29922]"
          />
          <MetricCard
            icon={<IconPlayerRecord size={16} />}
            label="Session recording"
            value={`${recordingPct}%`}
            sub={`${recordingCount}/${liveTerminals.length} recording`}
            subColor="text-[#3FB950]"
          />
        </div>

        {/* ── Main body: table + right panel ────────────────────────── */}
        <div className="flex gap-5">
          {/* ── Terminals table ──────────────────────────────────────── */}
          <div className="min-w-0 flex-1 rounded-lg border border-terminal-border bg-terminal-bg-toolbar">
            {/* Table header */}
            <div className="flex items-center justify-between border-b border-terminal-border px-4 py-3">
              <div className="flex items-center gap-2">
                <h2 className="text-sm font-semibold">All terminals</h2>
                <span className="rounded-full bg-terminal-border px-2 py-0.5 text-[11px] font-medium text-terminal-output">
                  {filteredTerminals.length}
                </span>
              </div>
              <div className="relative">
                <IconSearch size={14} className="absolute left-2.5 top-1/2 -translate-y-1/2 text-terminal-output" />
                <input
                  type="text"
                  placeholder="Search terminals…"
                  value={search}
                  onChange={(e) => setSearch(e.target.value)}
                  className="rounded border border-terminal-border bg-terminal-bg py-1.5 pl-8 pr-3 text-xs text-terminal-text placeholder-terminal-output/50 outline-none transition-colors focus:border-brand"
                />
              </div>
            </div>

            {/* Filter pills */}
            <div className="flex gap-1.5 border-b border-terminal-border px-4 py-2">
              {FILTERS.map((f) => (
                <button
                  key={f.key}
                  onClick={() => setFilter(f.key)}
                  className={`rounded-full px-3 py-1 text-xs font-medium transition-colors ${
                    filter === f.key
                      ? 'bg-brand text-white'
                      : 'bg-terminal-border/50 text-terminal-output hover:bg-terminal-border hover:text-terminal-text'
                  }`}
                >
                  {f.label}
                  {f.count > 0 && (
                    <span className="ml-1.5 opacity-70">{f.count}</span>
                  )}
                </button>
              ))}
            </div>

            {/* Table */}
            <div className="overflow-x-auto">
              <table className="w-full text-left text-xs">
                <thead>
                  <tr className="border-b border-terminal-border text-[11px] uppercase tracking-wider text-terminal-output">
                    <th className="px-4 py-2.5 font-medium">Terminal</th>
                    <th className="px-4 py-2.5 font-medium">Owner</th>
                    <th className="px-4 py-2.5 font-medium">Role</th>
                    <th className="px-4 py-2.5 font-medium">Status</th>
                    <th className="px-4 py-2.5 font-medium">CWD</th>
                    <th className="px-4 py-2.5 font-medium">Uptime</th>
                    <th className="px-4 py-2.5 font-medium">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {filteredTerminals.length === 0 ? (
                    <tr>
                      <td colSpan={7} className="px-4 py-8 text-center text-terminal-output">
                        No terminals match the current filter.
                      </td>
                    </tr>
                  ) : (
                    filteredTerminals.map((t) => {
                      const owner = userMap.get(t.ownerId);
                      return (
                        <tr
                          key={t.id}
                          className="group border-b border-terminal-border/50 transition-colors hover:bg-terminal-bg/50"
                        >
                          <td className="px-4 py-2.5">
                            <div className="flex items-center gap-2">
                              <StatusDot status={t.status} />
                              <span className="font-mono text-terminal-text">{t.title}</span>
                            </div>
                          </td>
                          <td className="px-4 py-2.5">
                            <div className="flex items-center gap-2">
                              <Avatar name={owner?.username ?? 'U'} size="sm" />
                              <span className="text-terminal-text">{owner?.username ?? 'Unknown'}</span>
                            </div>
                          </td>
                          <td className="px-4 py-2.5">
                            <RoleBadge role={owner?.role ?? 'viewer'} />
                          </td>
                          <td className="px-4 py-2.5">
                            <span
                              className={`capitalize ${
                                t.status === 'active'
                                  ? 'text-[#3FB950]'
                                  : t.status === 'idle'
                                    ? 'text-[#D29922]'
                                    : 'text-[#F85149]'
                              }`}
                            >
                              {t.status}
                            </span>
                          </td>
                          <td className="px-4 py-2.5">
                            <span className="font-mono text-terminal-output">{t.cwd}</span>
                          </td>
                          <td className="px-4 py-2.5 font-mono text-terminal-output">
                            {formatUptime(t.createdAt)}
                          </td>
                          <td className="px-4 py-2.5">
                            <div className="flex items-center gap-1 opacity-0 transition-opacity group-hover:opacity-100">
                              <button
                                className="rounded p-1 text-terminal-output hover:bg-terminal-border hover:text-terminal-text"
                                aria-label={`View ${t.title}`}
                                title="View"
                              >
                                <IconEye size={14} />
                              </button>
                              <button
                                className="rounded p-1 text-terminal-output hover:bg-terminal-border hover:text-terminal-error"
                                aria-label={`Stop ${t.title}`}
                                title="Stop"
                              >
                                <IconPlayerStop size={14} />
                              </button>
                            </div>
                          </td>
                        </tr>
                      );
                    })
                  )}
                </tbody>
              </table>
            </div>
          </div>

          {/* ── Right panel ──────────────────────────────────────────── */}
          <div className="hidden w-[260px] shrink-0 flex-col gap-5 xl:flex">
            {/* Active warnings */}
            <div className="rounded-lg border border-terminal-border bg-terminal-bg-toolbar p-4">
              <h3 className="mb-3 text-xs font-semibold uppercase tracking-wider text-terminal-output">
                Active warnings
              </h3>
              {activeWarnings.length === 0 ? (
                <p className="text-xs text-terminal-output/50">No active warnings</p>
              ) : (
                <div className="space-y-2">
                  {activeWarnings.map((w) => (
                    <div key={w.id} className="rounded-md border border-[#BA7517]/30 bg-[#BA7517]/10 px-3 py-2">
                      <div className="flex items-center gap-1.5">
                        <IconAlertTriangle size={13} className="shrink-0 text-[#D29922]" />
                        <span className="text-xs font-medium text-[#D29922]">Command blocked</span>
                      </div>
                      <p className="mt-1 text-[11px] text-terminal-output">
                        {w.action}{w.terminalId ? ` · ${w.terminalId.slice(0, 8)}` : ''}
                      </p>
                    </div>
                  ))}
                </div>
              )}
            </div>

            {/* Audit log */}
            <div className="rounded-lg border border-terminal-border bg-terminal-bg-toolbar p-4">
              <h3 className="mb-3 text-xs font-semibold uppercase tracking-wider text-terminal-output">
                Audit log
              </h3>
              <div className="space-y-3">
                {auditLogs.slice(0, 15).map((log) => {
                  const sev = classifyAudit(log.action);
                  return (
                    <div key={log.id} className="flex items-start gap-2.5">
                      <AuditIcon level={sev.level} />
                      <div className="min-w-0 flex-1">
                        <p className="truncate text-xs text-terminal-text">{log.action}</p>
                        <p className="font-mono text-[10px] text-terminal-output/50">
                          {formatTime(log.createdAt)}
                        </p>
                      </div>
                    </div>
                  );
                })}
                {auditLogs.length === 0 && (
                  <p className="text-xs text-terminal-output/50">No audit logs yet</p>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
