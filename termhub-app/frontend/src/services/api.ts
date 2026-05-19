import { useAuthStore } from '@/stores/authStore';
import type { AdminMetrics, AuthTokens, AuditLog, Terminal, User, UserRole } from '@/types';

const BASE_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080/api';

class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

/** Flag to prevent infinite refresh loops */
let isRefreshing = false;

async function request<T>(
  path: string,
  options: RequestInit = {},
  skipAuth = false,
): Promise<T> {
  const state = useAuthStore.getState();
  const tokens = state.tokens;

  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string> | undefined),
  };

  if (!skipAuth && tokens?.accessToken) {
    headers['Authorization'] = `Bearer ${tokens.accessToken}`;
  }

  const response = await fetch(`${BASE_URL}${path}`, {
    ...options,
    headers,
  });

  // Handle 401: attempt token refresh once, then retry
  if (response.status === 401 && !skipAuth && !isRefreshing && tokens?.refreshToken) {
    isRefreshing = true;
    try {
      const refreshed = await request<AuthTokens>(
        '/auth/refresh',
        {
          method: 'POST',
          body: JSON.stringify({ refreshToken: tokens.refreshToken }),
        },
        true,
      );
      useAuthStore.getState().refresh(refreshed);

      // Retry original request with new token
      headers['Authorization'] = `Bearer ${refreshed.accessToken}`;
      const retryResponse = await fetch(`${BASE_URL}${path}`, {
        ...options,
        headers,
      });

      if (!retryResponse.ok) {
        const body = await retryResponse.text();
        throw new ApiError(retryResponse.status, body || retryResponse.statusText);
      }

      if (retryResponse.status === 204) {
        return undefined as T;
      }

      return retryResponse.json() as Promise<T>;
    } catch {
      // Refresh failed — force logout
      useAuthStore.getState().logout();
      throw new ApiError(401, 'Session expired. Please log in again.');
    } finally {
      isRefreshing = false;
    }
  }

  if (!response.ok) {
    let message = response.statusText;
    try {
      const body: unknown = await response.json();
      if (body && typeof body === 'object' && 'error' in body) {
        message = String((body as { error: string }).error);
      } else if (body && typeof body === 'object' && 'message' in body) {
        message = String((body as { message: string }).message);
      }
    } catch {
      // body wasn't JSON — fall back to statusText
    }
    throw new ApiError(response.status, message);
  }

  // Handle 204 No Content
  if (response.status === 204) {
    return undefined as T;
  }

  return response.json() as Promise<T>;
}

// ── Auth ─────────────────────────────────────────────────────────────

export function login(username: string, password: string) {
  return request<{ user: User; tokens: AuthTokens }>('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ username, password }),
  });
}

export function register(username: string, password: string, role: UserRole = 'operator') {
  return request<{ user: User; tokens: AuthTokens }>('/auth/register', {
    method: 'POST',
    body: JSON.stringify({ username, password, role }),
  });
}

export function refreshToken(token: string) {
  return request<AuthTokens>(
    '/auth/refresh',
    {
      method: 'POST',
      body: JSON.stringify({ refreshToken: token }),
    },
    true,
  );
}

export function getMe() {
  return request<User>('/auth/me');
}

// ── Terminals ────────────────────────────────────────────────────────

export interface CreateTerminalRequest {
  title: string;
  cwd?: string;
  shell?: string;
  colorTag?: string;
}

export function createTerminal(data: CreateTerminalRequest) {
  return request<Terminal>('/terminals', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export function listTerminals() {
  return request<Terminal[]>('/terminals');
}

export function deleteTerminal(id: string) {
  return request<void>(`/terminals/${id}`, {
    method: 'DELETE',
  });
}

// ── Admin ────────────────────────────────────────────────────────────

export function adminListUsers() {
  return request<User[]>('/admin/users');
}

export function adminListTerminals() {
  return request<Terminal[]>('/admin/terminals');
}

export function adminGetAuditLogs(limit = 100) {
  return request<AuditLog[]>(`/admin/audit?limit=${limit}`);
}

export function adminGetMetrics() {
  return request<AdminMetrics>('/admin/metrics');
}

export { ApiError };
