import type { WSMessage } from '@/types';
import { useAuthStore } from '@/stores/authStore';

const WS_URL = import.meta.env.VITE_WS_URL ?? 'ws://localhost:8080/api/ws';

export type WSEventHandler = (message: WSMessage) => void;

export function createTerminalSocket(
  terminalId: string,
  onMessage: WSEventHandler,
  onClose?: () => void,
): WebSocket {
  const tokens = useAuthStore.getState().tokens;
  const url = `${WS_URL}/terminals/${terminalId}?token=${tokens?.accessToken ?? ''}`;

  const ws = new WebSocket(url);

  ws.onmessage = (event: MessageEvent) => {
    try {
      const msg = JSON.parse(event.data as string) as WSMessage;
      onMessage(msg);
    } catch {
      // If not JSON, treat as raw stdout
      onMessage({ type: 'stdout', data: event.data as string });
    }
  };

  ws.onclose = () => {
    onClose?.();
  };

  return ws;
}

export function sendInput(ws: WebSocket, data: string): void {
  const msg: WSMessage = { type: 'stdin', data };
  ws.send(JSON.stringify(msg));
}

export function sendResize(ws: WebSocket, cols: number, rows: number): void {
  const msg: WSMessage = { type: 'resize', cols, rows };
  ws.send(JSON.stringify(msg));
}
