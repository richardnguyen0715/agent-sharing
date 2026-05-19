import { useEffect, useRef, useCallback, useState } from 'react';
import { Terminal } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import {
  createTerminalSocket,
  sendInput,
  sendResize,
} from '@/services/websocket';
import type { WSMessage } from '@/types';

interface UseTerminalOptions {
  terminalId: string | null;
}

export function useTerminal({ terminalId }: UseTerminalOptions) {
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

    // Connect WebSocket
    const onMessage = (msg: WSMessage) => {
      if ((msg.type === 'output' || msg.type === 'stdout') && msg.data) {
        terminal.write(msg.data);
      }
    };

    const onClose = () => {
      setIsConnected(false);
    };

    const ws = createTerminalSocket(terminalId, onMessage, onClose);
    wsRef.current = ws;

    ws.onopen = () => {
      setIsConnected(true);
      // Send initial resize
      fitAddon.fit();
      sendResize(ws, terminal.cols, terminal.rows);
    };

    // Terminal data (keyboard input) → WebSocket
    const dataDisposable = terminal.onData((data) => {
      if (ws.readyState === WebSocket.OPEN) {
        sendInput(ws, data);
      }
    });

    // Terminal resize → WebSocket
    const resizeDisposable = terminal.onResize(({ cols, rows }) => {
      setDimensions({ cols, rows });
      if (ws.readyState === WebSocket.OPEN) {
        sendResize(ws, cols, rows);
      }
    });

    // Container resize → fit terminal
    const resizeObserver = new ResizeObserver(() => {
      fitAddon.fit();
    });
    resizeObserver.observe(container);

    return () => {
      resizeObserver.disconnect();
      dataDisposable.dispose();
      resizeDisposable.dispose();
      ws.close();
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

  return {
    terminalRef: attach,
    containerRef,
    isConnected,
    dimensions,
    fit,
    terminal: termRef,
  };
}
