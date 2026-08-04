export type UserRole = 'admin' | 'operator' | 'viewer' | 'restricted';

export interface User {
  id: string;
  username: string;
  role: UserRole;
  createdAt: string;
}

export interface Terminal {
  id: string;
  ownerId: string;
  title: string;
  cwd: string;
  shell: string;
  status: 'active' | 'idle' | 'terminated';
  colorTag: string;
  createdAt: string;
}

export interface TerminalTab {
  terminalId: string;
  title: string;
  status: 'active' | 'idle';
  colorTag: string;
}

export interface AuditLog {
  id: string;
  userId: string;
  terminalId?: string;
  action: string;
  metadata: Record<string, unknown>;
  createdAt: string;
}

export interface WSMessage {
  type: 'stdin' | 'stdout' | 'resize' | 'output' | 'error';
  data?: string;
  cols?: number;
  rows?: number;
}

export interface AuthTokens {
  accessToken: string;
  refreshToken: string;
}

export interface AdminMetrics {
  totalUsers: number;
  totalTerminals: number;
  activeTerminals: number;
  activeConnections: number;
}

export interface FileTreeEntry {
  name: string;
  path: string;
  isDir: boolean;
  size: number;
  modTime: string;
}
