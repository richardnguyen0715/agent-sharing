import '@/lib/monacoSetup';
import { useCallback, useEffect, useState } from 'react';
import Editor, { type OnMount } from '@monaco-editor/react';
import { IconFile, IconX, IconLoader2, IconDeviceFloppy, IconAlertTriangle, IconCode } from '@tabler/icons-react';
import { cn } from '@/lib/utils';
import { languageForPath } from '@/lib/languageForPath';
import { getFileContent, saveFileContent } from '@/services/api';
import { FileTree } from './FileTree';

interface OpenFile {
  path: string;
  content: string;
  originalContent: string;
  loading: boolean;
  saving: boolean;
  error?: string;
}

interface EditorViewProps {
  terminalId: string;
}

export function EditorView({ terminalId }: EditorViewProps) {
  const [openFiles, setOpenFiles] = useState<OpenFile[]>([]);
  const [activePath, setActivePath] = useState<string | null>(null);

  // Different terminal = different filesystem root — start clean.
  useEffect(() => {
    setOpenFiles([]);
    setActivePath(null);
  }, [terminalId]);

  const loadFile = useCallback(
    (path: string) => {
      getFileContent(terminalId, path)
        .then((res) => {
          setOpenFiles((prev) =>
            prev.map((f) => (f.path === path ? { ...f, content: res.content, originalContent: res.content, loading: false } : f)),
          );
        })
        .catch((err) => {
          setOpenFiles((prev) =>
            prev.map((f) =>
              f.path === path ? { ...f, loading: false, error: err instanceof Error ? err.message : 'failed to load' } : f,
            ),
          );
        });
    },
    [terminalId],
  );

  const openFile = useCallback(
    (path: string) => {
      setActivePath(path);
      setOpenFiles((prev) => {
        if (prev.some((f) => f.path === path)) return prev;
        return [...prev, { path, content: '', originalContent: '', loading: true, saving: false }];
      });
      loadFile(path);
    },
    [loadFile],
  );

  const closeFile = useCallback(
    (path: string, e?: React.MouseEvent) => {
      e?.stopPropagation();
      const file = openFiles.find((f) => f.path === path);
      if (file && file.content !== file.originalContent) {
        if (!window.confirm(`Discard unsaved changes to "${path.split('/').pop()}"?`)) return;
      }
      setOpenFiles((prev) => {
        const next = prev.filter((f) => f.path !== path);
        if (activePath === path) {
          setActivePath(next.at(-1)?.path ?? null);
        }
        return next;
      });
    },
    [openFiles, activePath],
  );

  const setContent = useCallback((path: string, content: string) => {
    setOpenFiles((prev) => prev.map((f) => (f.path === path ? { ...f, content } : f)));
  }, []);

  const saveFile = useCallback(
    (path: string) => {
      const file = openFiles.find((f) => f.path === path);
      if (!file || file.content === file.originalContent) return;
      setOpenFiles((prev) => prev.map((f) => (f.path === path ? { ...f, saving: true, error: undefined } : f)));
      saveFileContent(terminalId, path, file.content)
        .then(() => {
          setOpenFiles((prev) =>
            prev.map((f) => (f.path === path ? { ...f, originalContent: f.content, saving: false } : f)),
          );
        })
        .catch((err) => {
          setOpenFiles((prev) =>
            prev.map((f) =>
              f.path === path ? { ...f, saving: false, error: err instanceof Error ? err.message : 'save failed' } : f,
            ),
          );
        });
    },
    [openFiles, terminalId],
  );

  const handleDeletedPath = useCallback((path: string) => {
    setOpenFiles((prev) => {
      const next = prev.filter((f) => f.path !== path && !f.path.startsWith(`${path}/`));
      setActivePath((current) => (current && !next.some((f) => f.path === current) ? (next[0]?.path ?? null) : current));
      return next;
    });
  }, []);

  const handleRenamedPath = useCallback((from: string, to: string) => {
    setOpenFiles((prev) => prev.map((f) => (f.path === from ? { ...f, path: to } : f)));
    setActivePath((current) => (current === from ? to : current));
  }, []);

  const handleEditorMount: OnMount = (editor, monaco) => {
    editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyS, () => {
      const path = editor.getModel()?.uri.path.replace(/^\//, '');
      if (path) saveFile(path);
    });
  };

  const activeFile = openFiles.find((f) => f.path === activePath) ?? null;

  return (
    <div className="flex h-full min-h-0 w-full">
      <FileTree
        terminalId={terminalId}
        activePath={activePath}
        onOpenFile={openFile}
        onDeletedPath={handleDeletedPath}
        onRenamedPath={handleRenamedPath}
      />

      <div className="flex min-w-0 flex-1 flex-col">
        {openFiles.length > 0 && (
          <div className="flex h-8 shrink-0 items-stretch overflow-x-auto border-b border-terminal-border bg-terminal-bg-toolbar">
            {openFiles.map((f) => {
              const dirty = f.content !== f.originalContent;
              const name = f.path.split('/').pop() ?? f.path;
              return (
                <button
                  key={f.path}
                  onClick={() => setActivePath(f.path)}
                  title={f.path}
                  className={cn(
                    'group flex shrink-0 items-center gap-1.5 border-r border-terminal-border px-2.5 text-xs',
                    f.path === activePath
                      ? 'bg-terminal-bg text-terminal-text'
                      : 'text-terminal-output hover:bg-terminal-border/40',
                  )}
                >
                  <IconFile size={12} className="shrink-0" />
                  <span className="max-w-[140px] truncate">{name}</span>
                  {f.saving ? (
                    <IconLoader2 size={11} className="shrink-0 animate-spin" />
                  ) : dirty ? (
                    <span className="h-1.5 w-1.5 shrink-0 rounded-full bg-terminal-path" />
                  ) : null}
                  <span
                    role="button"
                    tabIndex={-1}
                    onClick={(e) => closeFile(f.path, e)}
                    className="rounded p-0.5 opacity-0 hover:bg-terminal-border group-hover:opacity-100"
                  >
                    <IconX size={12} />
                  </span>
                </button>
              );
            })}
          </div>
        )}

        {!activeFile && (
          <div className="flex flex-1 flex-col items-center justify-center gap-2 text-terminal-output/50">
            <IconCode size={40} className="opacity-30" />
            <p className="text-sm">Select a file to start editing</p>
          </div>
        )}

        {activeFile?.loading && (
          <div className="flex flex-1 items-center justify-center gap-2 text-sm text-terminal-output/60">
            <IconLoader2 size={16} className="animate-spin" /> Loading…
          </div>
        )}

        {activeFile?.error && !activeFile.loading && (
          <div className="flex flex-1 flex-col items-center justify-center gap-2 text-terminal-error">
            <IconAlertTriangle size={28} />
            <p className="text-sm">{activeFile.error}</p>
          </div>
        )}

        {activeFile && !activeFile.loading && !activeFile.error && (
          <div className="relative min-h-0 flex-1">
            <Editor
              key={activeFile.path}
              path={activeFile.path}
              language={languageForPath(activeFile.path)}
              value={activeFile.content}
              theme="vs-dark"
              onChange={(value) => setContent(activeFile.path, value ?? '')}
              onMount={handleEditorMount}
              options={{
                fontSize: 13,
                fontFamily: '"JetBrains Mono", monospace',
                minimap: { enabled: true },
                automaticLayout: true,
                scrollBeyondLastLine: false,
                tabSize: 2,
              }}
            />
            <button
              onClick={() => saveFile(activeFile.path)}
              disabled={activeFile.saving || activeFile.content === activeFile.originalContent}
              title="Save (Ctrl/Cmd+S)"
              className="absolute right-3 top-2 z-10 flex items-center gap-1 rounded border border-terminal-border bg-terminal-bg-toolbar px-2 py-1 text-xs text-terminal-output shadow hover:bg-terminal-border hover:text-terminal-text disabled:opacity-40"
            >
              {activeFile.saving ? <IconLoader2 size={13} className="animate-spin" /> : <IconDeviceFloppy size={13} />}
              Save
            </button>
          </div>
        )}
      </div>
    </div>
  );
}

export default EditorView;
