import { useEffect, useRef, useCallback, useState } from 'react';
import { Terminal } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import {
  createTerminalSocket,
  sendInput,
  sendResize,
} from '@/services/websocket';
import { uploadTerminalFile } from '@/services/api';
import type { WSMessage } from '@/types';

/** Shell-quotes a path so spaces/special characters paste safely into the terminal. */
function shellQuote(path: string): string {
  return `'${path.replace(/'/g, "'\\''")}'`;
}

interface UseTerminalOptions {
  terminalId: string | null;
  active?: boolean;
}

export function useTerminal({ terminalId, active = true }: UseTerminalOptions) {
  const termRef = useRef<Terminal | null>(null);
  const fitRef = useRef<FitAddon | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const containerRef = useRef<HTMLDivElement | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const [dimensions, setDimensions] = useState({ cols: 80, rows: 24 });
  const [containerReady, setContainerReady] = useState(false);

  const terminalIdRef = useRef(terminalId);
  terminalIdRef.current = terminalId;

  const attach = useCallback((el: HTMLDivElement | null) => {
    containerRef.current = el;
    setContainerReady(el !== null);
  }, []);

  // Initialize xterm + websocket when terminalId and container are available
  useEffect(() => {
    const container = containerRef.current;
    if (!container || !terminalId || !containerReady) return;

    const terminal = new Terminal({
      fontFamily: '"JetBrains Mono", monospace',
      fontSize: 14,
      lineHeight: 1.2,
      theme: {
        background: '#0D1117',
        foreground: '#E6EDF3',
        cursor: '#58A6FF',
        cursorAccent: '#0D1117',
        selectionBackground: '#264F78',
        black: '#484F58',
        red: '#F85149',
        green: '#3FB950',
        yellow: '#D29922',
        blue: '#58A6FF',
        magenta: '#BC8CFF',
        cyan: '#39C5CF',
        white: '#B1BAC4',
        brightBlack: '#6E7681',
        brightRed: '#FF7B72',
        brightGreen: '#56D364',
        brightYellow: '#E3B341',
        brightBlue: '#79C0FF',
        brightMagenta: '#D2A8FF',
        brightCyan: '#56D4DD',
        brightWhite: '#F0F6FC',
      },
      cursorBlink: true,
      allowProposedApi: true,
      scrollback: 5000,
    });

    const fitAddon = new FitAddon();
    terminal.loadAddon(fitAddon);
    terminal.open(container);

    // Small delay to let the DOM settle before fitting
    requestAnimationFrame(() => {
      fitAddon.fit();
      setDimensions({ cols: terminal.cols, rows: terminal.rows });
    });

    termRef.current = terminal;
    fitRef.current = fitAddon;

    let disposed = false;
    let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
    let reconnectAttempt = 0;

    const onMessage = (msg: WSMessage) => {
      if ((msg.type === 'output' || msg.type === 'stdout') && msg.data) {
        terminal.write(msg.data);
      }
    };

    // Connect (and reconnect) the WebSocket. Handlers below always read
    // wsRef.current rather than closing over a specific socket instance, so
    // input keeps flowing to whichever connection is currently live instead
    // of silently going nowhere after a reconnect.
    const connect = () => {
      if (disposed) return;

      const ws = createTerminalSocket(terminalId, onMessage, () => {
        setIsConnected(false);
        if (disposed) return;
        // Reconnect with capped exponential backoff — recovers automatically
        // from idle timeouts / network blips instead of leaving the terminal
        // silently unable to accept input until the page is reloaded.
        const delay = Math.min(1000 * 2 ** reconnectAttempt, 10000);
        reconnectAttempt += 1;
        reconnectTimer = setTimeout(connect, delay);
      });
      wsRef.current = ws;

      ws.onopen = () => {
        reconnectAttempt = 0;
        setIsConnected(true);
        fitAddon.fit();
        sendResize(ws, terminal.cols, terminal.rows);
      };
    };
    connect();

    // Terminal data (keyboard input) → WebSocket
    const dataDisposable = terminal.onData((data) => {
      const ws = wsRef.current;
      if (ws && ws.readyState === WebSocket.OPEN) {
        sendInput(ws, data);
      }
    });

    // Terminal resize → WebSocket
    const resizeDisposable = terminal.onResize(({ cols, rows }) => {
      setDimensions({ cols, rows });
      const ws = wsRef.current;
      if (ws && ws.readyState === WebSocket.OPEN) {
        sendResize(ws, cols, rows);
      }
    });

    // Container resize → fit terminal
    const resizeObserver = new ResizeObserver(() => {
      fitAddon.fit();
    });
    resizeObserver.observe(container);

    // Pasting/dropping files (e.g. an image copied from the clipboard) can't
    // go into a terminal as raw bytes — there's no such thing as "pasting an
    // image" into a text stream. Instead, upload the file and paste its path
    // so the user can reference it in a shell command. Listened for in the
    // capture phase so it runs before xterm's own paste handler on its
    // hidden textarea (which would otherwise just paste nothing/garbage for
    // a file-only clipboard item).
    const uploadAndPastePaths = async (files: File[]) => {
      for (const file of files) {
        try {
          const result = await uploadTerminalFile(terminalId, file);
          terminal.paste(`${shellQuote(result.path)} `);
        } catch (err) {
          const message = err instanceof Error ? err.message : 'upload failed';
          terminal.write(`\r\n\x1b[31m[upload failed: ${file.name} — ${message}]\x1b[0m\r\n`);
        }
      }
    };

    const handlePaste = (e: ClipboardEvent) => {
      const files = e.clipboardData?.files;
      if (!files || files.length === 0) return; // plain text — let xterm handle it
      e.preventDefault();
      e.stopPropagation();
      void uploadAndPastePaths(Array.from(files));
    };
    container.addEventListener('paste', handlePaste, true);

    const handleDrop = (e: DragEvent) => {
      const files = e.dataTransfer?.files;
      if (!files || files.length === 0) return;
      e.preventDefault();
      void uploadAndPastePaths(Array.from(files));
    };
    const handleDragOver = (e: DragEvent) => {
      if (e.dataTransfer?.types.includes('Files')) e.preventDefault();
    };
    container.addEventListener('drop', handleDrop);
    container.addEventListener('dragover', handleDragOver);

    return () => {
      disposed = true;
      if (reconnectTimer) clearTimeout(reconnectTimer);
      resizeObserver.disconnect();
      container.removeEventListener('paste', handlePaste, true);
      container.removeEventListener('drop', handleDrop);
      container.removeEventListener('dragover', handleDragOver);
      dataDisposable.dispose();
      resizeDisposable.dispose();
      wsRef.current?.close();
      wsRef.current = null;
      terminal.dispose();
      termRef.current = null;
      fitRef.current = null;
      setIsConnected(false);
    };
  }, [terminalId, containerReady]);

  const fit = useCallback(() => {
    fitRef.current?.fit();
  }, []);

  // When this pane becomes visible, re-fit (it had 0 size while hidden) and
  // re-sync size to the PTY. The buffer is preserved because the terminal is
  // kept mounted, so the existing prompt/output stays on screen.
  useEffect(() => {
    if (!active) return;
    const terminal = termRef.current;
    const fitAddon = fitRef.current;
    const ws = wsRef.current;
    if (!terminal || !fitAddon) return;

    const raf = requestAnimationFrame(() => {
      fitAddon.fit();
      setDimensions({ cols: terminal.cols, rows: terminal.rows });
      if (ws && ws.readyState === WebSocket.OPEN) {
        sendResize(ws, terminal.cols, terminal.rows);
      }
      terminal.scrollToBottom();
      terminal.focus();
    });
    return () => cancelAnimationFrame(raf);
  }, [active, containerReady]);

  return {
    terminalRef: attach,
    containerRef,
    isConnected,
    dimensions,
    fit,
    terminal: termRef,
  };
}
