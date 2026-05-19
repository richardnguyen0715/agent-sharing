import { useState, type FormEvent } from 'react';
import { IconX, IconTerminal2, IconLoader2 } from '@tabler/icons-react';
import { useUIStore } from '@/stores/uiStore';
import { useTerminalStore } from '@/stores/terminalStore';
import { Toggle } from '@/components/ui/Toggle';

const SHELL_OPTIONS = [
  { value: '/bin/bash', label: 'Bash' },
  { value: '/bin/zsh', label: 'Zsh' },
  { value: '/bin/sh', label: 'POSIX sh' },
];

const DIR_OPTIONS = [
  { value: '~', label: '~ (Home)' },
  { value: '/workspace', label: '/workspace' },
  { value: '/tmp', label: '/tmp' },
];

interface ColorOption {
  name: string;
  value: string;
  className: string;
}

const COLOR_OPTIONS: ColorOption[] = [
  { name: 'Default', value: 'default', className: 'bg-terminal-output' },
  { name: 'Blue', value: 'blue', className: 'bg-terminal-path' },
  { name: 'Green', value: 'green', className: 'bg-terminal-prompt' },
  { name: 'Amber', value: 'amber', className: 'bg-status-amber' },
  { name: 'Coral', value: 'coral', className: 'bg-terminal-error' },
  { name: 'Pink', value: 'pink', className: 'bg-[#D2A8FF]' },
];

export function NewTerminalModal() {
  const closeModal = useUIStore((s) => s.closeModal);
  const createTerminal = useTerminalStore((s) => s.createTerminal);

  const [title, setTitle] = useState('');
  const [shell, setShell] = useState('/bin/bash');
  const [cwd, setCwd] = useState('~');
  const [colorTag, setColorTag] = useState('default');
  const [autoRecord, setAutoRecord] = useState(false);
  const [allowSharing, setAllowSharing] = useState(false);
  const [persist, setPersist] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const shellLabel = SHELL_OPTIONS.find((s) => s.value === shell)?.label ?? 'shell';
  const effectiveTitle = title.trim() || 'New Terminal';

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);
    setIsSubmitting(true);

    try {
      await createTerminal({
        title: effectiveTitle,
        shell,
        cwd,
        colorTag,
      });
      closeModal();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create terminal');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60">
      <div
        className="w-full max-w-lg rounded-xl border border-terminal-border bg-terminal-bg-toolbar shadow-2xl shadow-black/50"
        role="dialog"
        aria-modal="true"
        aria-labelledby="new-terminal-title"
      >
        {/* Header */}
        <div className="flex items-center justify-between border-b border-terminal-border px-6 py-4">
          <div className="flex items-center gap-3">
            <div className="flex h-9 w-9 items-center justify-center rounded-lg bg-brand/15">
              <IconTerminal2 size={18} className="text-brand" />
            </div>
            <div>
              <h2 id="new-terminal-title" className="text-base font-semibold text-terminal-text">
                New Terminal
              </h2>
              <p className="text-xs text-terminal-output">Configure your terminal session</p>
            </div>
          </div>
          <button
            onClick={closeModal}
            className="rounded p-1 text-terminal-output hover:bg-terminal-border hover:text-terminal-text"
            aria-label="Close modal"
          >
            <IconX size={18} />
          </button>
        </div>

        {/* Body */}
        <form onSubmit={handleSubmit}>
          <div className="space-y-5 px-6 py-5">
            {/* Terminal name */}
            <div>
              <label htmlFor="terminal-title" className="mb-1.5 block text-sm font-medium text-terminal-output">
                Terminal name
                <span className="ml-1 text-terminal-output/50">(optional)</span>
              </label>
              <input
                id="terminal-title"
                type="text"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="My Terminal"
                className="w-full rounded-lg border border-terminal-border bg-terminal-bg px-3 py-2 text-sm text-terminal-text placeholder-terminal-output/40 outline-none transition-colors focus:border-brand focus:ring-1 focus:ring-brand/40"
              />
            </div>

            {/* Directory + Shell row */}
            <div className="grid grid-cols-2 gap-3">
              <div>
                <label htmlFor="terminal-cwd" className="mb-1.5 block text-sm font-medium text-terminal-output">
                  Starting directory
                </label>
                <select
                  id="terminal-cwd"
                  value={cwd}
                  onChange={(e) => setCwd(e.target.value)}
                  className="w-full rounded-lg border border-terminal-border bg-terminal-bg px-3 py-2 text-sm text-terminal-text outline-none transition-colors focus:border-brand focus:ring-1 focus:ring-brand/40"
                >
                  {DIR_OPTIONS.map((opt) => (
                    <option key={opt.value} value={opt.value}>{opt.label}</option>
                  ))}
                </select>
              </div>
              <div>
                <label htmlFor="terminal-shell" className="mb-1.5 block text-sm font-medium text-terminal-output">
                  Shell
                </label>
                <select
                  id="terminal-shell"
                  value={shell}
                  onChange={(e) => setShell(e.target.value)}
                  className="w-full rounded-lg border border-terminal-border bg-terminal-bg px-3 py-2 text-sm text-terminal-text outline-none transition-colors focus:border-brand focus:ring-1 focus:ring-brand/40"
                >
                  {SHELL_OPTIONS.map((opt) => (
                    <option key={opt.value} value={opt.value}>{opt.label}</option>
                  ))}
                </select>
              </div>
            </div>

            {/* Color tag */}
            <div>
              <label className="mb-2 block text-sm font-medium text-terminal-output">
                Color tag
              </label>
              <div className="flex gap-2">
                {COLOR_OPTIONS.map((opt) => (
                  <button
                    key={opt.value}
                    type="button"
                    onClick={() => setColorTag(opt.value)}
                    className={`flex flex-col items-center gap-1 rounded-lg border px-2.5 py-2 transition-colors ${
                      colorTag === opt.value
                        ? 'border-brand bg-brand/10'
                        : 'border-terminal-border hover:border-terminal-output/30'
                    }`}
                  >
                    <span className={`h-3 w-3 rounded-full ${opt.className}`} />
                    <span className="text-[10px] text-terminal-output">{opt.name}</span>
                  </button>
                ))}
              </div>
            </div>

            {/* Preview */}
            <div className="overflow-hidden rounded-lg border border-terminal-border bg-terminal-bg">
              <div className="px-3 py-2 font-mono text-xs">
                <span className="text-terminal-prompt">❯</span>{' '}
                <span className="text-terminal-path">{cwd}</span>{' '}
                <span className="text-terminal-output">({shellLabel})</span>
                <span className="ml-0.5 inline-block h-3.5 w-1.5 animate-pulse bg-terminal-cursor" />
              </div>
            </div>

            {/* Toggles */}
            <div className="space-y-3">
              <Toggle
                checked={autoRecord}
                onChange={setAutoRecord}
                label="Auto-record session"
              />
              <Toggle
                checked={allowSharing}
                onChange={setAllowSharing}
                label="Allow sharing"
              />
              <Toggle
                checked={persist}
                onChange={setPersist}
                label="Persist after disconnect"
              />
            </div>

            {/* Error */}
            {error && (
              <div className="rounded-lg border border-terminal-error/30 bg-terminal-error/10 px-3 py-2 text-sm text-terminal-error">
                {error}
              </div>
            )}
          </div>

          {/* Footer */}
          <div className="flex items-center justify-end gap-2 border-t border-terminal-border px-6 py-4">
            <button
              type="button"
              onClick={closeModal}
              disabled={isSubmitting}
              className="rounded-lg border border-terminal-border px-4 py-2 text-sm text-terminal-output transition-colors hover:bg-terminal-border hover:text-terminal-text disabled:opacity-50"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isSubmitting}
              className="flex items-center gap-2 rounded-lg bg-brand px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-brand-hover disabled:opacity-50"
            >
              {isSubmitting && <IconLoader2 size={14} className="animate-spin" />}
              Create terminal
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
