import { useEffect, useRef, useCallback } from 'react';
import {
  createTerminalSocket,
  sendInput,
  sendResize,
  type WSEventHandler,
} from '@/services/websocket';

export function useWebSocket(
  terminalId: string | null,
  onMessage: WSEventHandler,
  onClose?: () => void,
) {
  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    if (!terminalId) return;

    const ws = createTerminalSocket(terminalId, onMessage, onClose);
    wsRef.current = ws;

    return () => {
      ws.close();
      wsRef.current = null;
    };
  }, [terminalId, onMessage, onClose]);

  const send = useCallback((data: string) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      sendInput(wsRef.current, data);
    }
  }, []);

  const resize = useCallback((cols: number, rows: number) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      sendResize(wsRef.current, cols, rows);
    }
  }, []);

  return { send, resize };
}
