import { useCallback, useEffect, useState } from 'react';
import {
  IconChevronRight,
  IconChevronDown,
  IconFile,
  IconFolder,
  IconFilePlus,
  IconFolderPlus,
  IconTrash,
  IconPencil,
  IconLoader2,
} from '@tabler/icons-react';
import { cn } from '@/lib/utils';
import { listFileTree, createFileEntry, renameFileEntry, deleteFileEntry } from '@/services/api';
import type { FileTreeEntry } from '@/types';

interface FileTreeProps {
  terminalId: string;
  activePath: string | null;
  onOpenFile: (path: string) => void;
  onDeletedPath: (path: string) => void;
  onRenamedPath: (from: string, to: string) => void;
}

/** Inline text input used for both "new file/folder" and "rename" flows. */
function InlineInput({
  initialValue,
  onSubmit,
  onCancel,
  autoFocus = true,
}: {
  initialValue: string;
  onSubmit: (value: string) => void;
  onCancel: () => void;
  autoFocus?: boolean;
}) {
  const [value, setValue] = useState(initialValue);
  return (
    <input
      autoFocus={autoFocus}
      value={value}
      onChange={(e) => setValue(e.target.value)}
      onClick={(e) => e.stopPropagation()}
      onFocus={(e) => e.target.select()}
      onKeyDown={(e) => {
        e.stopPropagation();
        if (e.key === 'Enter') {
          const trimmed = value.trim();
          if (trimmed) onSubmit(trimmed);
          else onCancel();
        } else if (e.key === 'Escape') {
          onCancel();
        }
      }}
      onBlur={() => onCancel()}
      className="min-w-0 flex-1 rounded border border-brand bg-terminal-bg px-1 py-0.5 text-xs text-terminal-text outline-none"
    />
  );
}

interface TreeNodeProps {
  terminalId: string;
  entry: FileTreeEntry;
  depth: number;
  activePath: string | null;
  onOpenFile: (path: string) => void;
  onDeletedPath: (path: string) => void;
  onRenamedPath: (from: string, to: string) => void;
  /** Refreshes the listing that renders this node (parent dir's children, or the root). */
  onMutated: () => void;
}

function TreeNode({
  terminalId,
  entry,
  depth,
  activePath,
  onOpenFile,
  onDeletedPath,
  onRenamedPath,
  onMutated,
}: TreeNodeProps) {
  const [expanded, setExpanded] = useState(false);
  const [children, setChildren] = useState<FileTreeEntry[] | null>(null);
  const [loading, setLoading] = useState(false);
  const [renaming, setRenaming] = useState(false);
  const [creating, setCreating] = useState<'file' | 'dir' | null>(null);
  const [error, setError] = useState<string | null>(null);

  const loadChildren = useCallback(() => {
    setLoading(true);
    setError(null);
    listFileTree(terminalId, entry.path)
      .then((res) => setChildren(res.entries))
      .catch((err) => setError(err instanceof Error ? err.message : 'failed to load'))
      .finally(() => setLoading(false));
  }, [terminalId, entry.path]);

  const toggle = () => {
    if (!entry.isDir) {
      onOpenFile(entry.path);
      return;
    }
    if (!expanded && children === null) loadChildren();
    setExpanded((v) => !v);
  };

  const handleRename = (newName: string) => {
    setRenaming(false);
    if (newName === entry.name) return;
    const parent = entry.path.includes('/') ? entry.path.slice(0, entry.path.lastIndexOf('/')) : '';
    const to = parent ? `${parent}/${newName}` : newName;
    renameFileEntry(terminalId, entry.path, to)
      .then(() => {
        onRenamedPath(entry.path, to);
        onMutated();
      })
      .catch((err) => setError(err instanceof Error ? err.message : 'rename failed'));
  };

  const handleDelete = (e: React.MouseEvent) => {
    e.stopPropagation();
    if (!window.confirm(`Delete "${entry.name}"?${entry.isDir ? ' This removes everything inside it.' : ''}`)) return;
    deleteFileEntry(terminalId, entry.path)
      .then(() => {
        onDeletedPath(entry.path);
        onMutated();
      })
      .catch((err) => setError(err instanceof Error ? err.message : 'delete failed'));
  };

  const handleCreate = (name: string) => {
    const type = creating;
    setCreating(null);
    if (!type) return;
    const path = `${entry.path}/${name}`;
    createFileEntry(terminalId, path, type)
      .then(() => loadChildren())
      .catch((err) => setError(err instanceof Error ? err.message : 'create failed'));
  };

  return (
    <div>
      <div
        role="button"
        tabIndex={0}
        onClick={toggle}
        style={{ paddingLeft: `${depth * 14 + 6}px` }}
        className={cn(
          'group flex h-6 items-center gap-1 pr-1 text-xs hover:bg-terminal-border/60',
          activePath === entry.path ? 'bg-terminal-border text-terminal-text' : 'text-terminal-output',
        )}
      >
        {entry.isDir ? (
          loading ? (
            <IconLoader2 size={13} className="shrink-0 animate-spin" />
          ) : expanded ? (
            <IconChevronDown size={13} className="shrink-0" />
          ) : (
            <IconChevronRight size={13} className="shrink-0" />
          )
        ) : (
          <span className="w-[13px] shrink-0" />
        )}
        {entry.isDir ? (
          <IconFolder size={13} className="shrink-0 text-terminal-path" />
        ) : (
          <IconFile size={13} className="shrink-0 text-terminal-output/70" />
        )}
        {renaming ? (
          <InlineInput initialValue={entry.name} onSubmit={handleRename} onCancel={() => setRenaming(false)} />
        ) : (
          <span className="min-w-0 flex-1 truncate">{entry.name}</span>
        )}
        {!renaming && (
          <span className="hidden shrink-0 items-center gap-0.5 group-hover:flex">
            {entry.isDir && (
              <>
                <button
                  aria-label="New file"
                  title="New file"
                  onClick={(e) => {
                    e.stopPropagation();
                    setExpanded(true);
                    if (children === null) loadChildren();
                    setCreating('file');
                  }}
                  className="rounded p-0.5 hover:bg-terminal-border hover:text-terminal-text"
                >
                  <IconFilePlus size={12} />
                </button>
                <button
                  aria-label="New folder"
                  title="New folder"
                  onClick={(e) => {
                    e.stopPropagation();
                    setExpanded(true);
                    if (children === null) loadChildren();
                    setCreating('dir');
                  }}
                  className="rounded p-0.5 hover:bg-terminal-border hover:text-terminal-text"
                >
                  <IconFolderPlus size={12} />
                </button>
              </>
            )}
            <button
              aria-label="Rename"
              title="Rename"
              onClick={(e) => {
                e.stopPropagation();
                setRenaming(true);
              }}
              className="rounded p-0.5 hover:bg-terminal-border hover:text-terminal-text"
            >
              <IconPencil size={12} />
            </button>
            <button
              aria-label="Delete"
              title="Delete"
              onClick={handleDelete}
              className="rounded p-0.5 hover:bg-terminal-error hover:text-white"
            >
              <IconTrash size={12} />
            </button>
          </span>
        )}
      </div>

      {error && (
        <div style={{ paddingLeft: `${depth * 14 + 22}px` }} className="text-[11px] text-terminal-error">
          {error}
        </div>
      )}

      {entry.isDir && expanded && (
        <div>
          {creating && (
            <div
              style={{ paddingLeft: `${(depth + 1) * 14 + 6}px` }}
              className="flex h-6 items-center gap-1 pr-1 text-xs"
            >
              <span className="w-[13px] shrink-0" />
              {creating === 'dir' ? (
                <IconFolder size={13} className="shrink-0 text-terminal-path" />
              ) : (
                <IconFile size={13} className="shrink-0 text-terminal-output/70" />
              )}
              <InlineInput initialValue="" onSubmit={handleCreate} onCancel={() => setCreating(null)} />
            </div>
          )}
          {children?.map((child) => (
            <TreeNode
              key={child.path}
              terminalId={terminalId}
              entry={child}
              depth={depth + 1}
              activePath={activePath}
              onOpenFile={onOpenFile}
              onDeletedPath={onDeletedPath}
              onRenamedPath={onRenamedPath}
              onMutated={loadChildren}
            />
          ))}
        </div>
      )}
    </div>
  );
}

export function FileTree({ terminalId, activePath, onOpenFile, onDeletedPath, onRenamedPath }: FileTreeProps) {
  const [entries, setEntries] = useState<FileTreeEntry[] | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [creating, setCreating] = useState<'file' | 'dir' | null>(null);

  const refresh = useCallback(() => {
    listFileTree(terminalId, '')
      .then((res) => setEntries(res.entries))
      .catch((err) => setError(err instanceof Error ? err.message : 'failed to load'));
  }, [terminalId]);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const handleCreate = (name: string) => {
    const type = creating;
    setCreating(null);
    if (!type) return;
    createFileEntry(terminalId, name, type)
      .then(() => refresh())
      .catch((err) => setError(err instanceof Error ? err.message : 'create failed'));
  };

  return (
    <div className="flex h-full w-60 shrink-0 flex-col overflow-y-auto border-r border-terminal-border bg-terminal-bg-toolbar">
      <div className="flex shrink-0 items-center justify-between px-2 py-1.5">
        <span className="text-[11px] font-semibold uppercase tracking-wide text-terminal-output/60">
          Explorer
        </span>
        <div className="flex gap-0.5">
          <button
            aria-label="New file"
            title="New file"
            onClick={() => setCreating('file')}
            className="rounded p-0.5 text-terminal-output hover:bg-terminal-border hover:text-terminal-text"
          >
            <IconFilePlus size={13} />
          </button>
          <button
            aria-label="New folder"
            title="New folder"
            onClick={() => setCreating('dir')}
            className="rounded p-0.5 text-terminal-output hover:bg-terminal-border hover:text-terminal-text"
          >
            <IconFolderPlus size={13} />
          </button>
        </div>
      </div>

      {error && <div className="px-2 pb-1 text-[11px] text-terminal-error">{error}</div>}

      {creating && (
        <div className="flex h-6 items-center gap-1 px-2 pl-3.5">
          {creating === 'dir' ? (
            <IconFolder size={13} className="shrink-0 text-terminal-path" />
          ) : (
            <IconFile size={13} className="shrink-0 text-terminal-output/70" />
          )}
          <InlineInput initialValue="" onSubmit={handleCreate} onCancel={() => setCreating(null)} />
        </div>
      )}

      {entries === null && !error && (
        <div className="flex items-center gap-2 px-2 py-2 text-xs text-terminal-output/60">
          <IconLoader2 size={13} className="animate-spin" /> Loading…
        </div>
      )}

      {entries?.map((entry) => (
        <TreeNode
          key={entry.path}
          terminalId={terminalId}
          entry={entry}
          depth={0}
          activePath={activePath}
          onOpenFile={onOpenFile}
          onDeletedPath={onDeletedPath}
          onRenamedPath={onRenamedPath}
          onMutated={refresh}
        />
      ))}
    </div>
  );
}
