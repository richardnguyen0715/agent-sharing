import { create } from 'zustand';
import type { User, AuthTokens } from '@/types';
import * as api from '@/services/api';

const STORAGE_KEY_ACCESS = 'termhub_access_token';
const STORAGE_KEY_REFRESH = 'termhub_refresh_token';

function persistTokens(tokens: AuthTokens | null): void {
  if (tokens) {
    localStorage.setItem(STORAGE_KEY_ACCESS, tokens.accessToken);
    localStorage.setItem(STORAGE_KEY_REFRESH, tokens.refreshToken);
  } else {
    localStorage.removeItem(STORAGE_KEY_ACCESS);
    localStorage.removeItem(STORAGE_KEY_REFRESH);
  }
}

function loadTokens(): AuthTokens | null {
  const accessToken = localStorage.getItem(STORAGE_KEY_ACCESS);
  const refreshToken = localStorage.getItem(STORAGE_KEY_REFRESH);
  if (accessToken && refreshToken) {
    return { accessToken, refreshToken };
  }
  return null;
}

interface AuthState {
  user: User | null;
  tokens: AuthTokens | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  isInitialized: boolean;
  error: string | null;

  login: (user: User, tokens: AuthTokens) => void;
  logout: () => void;
  refresh: (tokens: AuthTokens) => void;
  setUser: (user: User) => void;
  setError: (error: string | null) => void;
  setLoading: (loading: boolean) => void;
  initialize: () => Promise<void>;
}

export const useAuthStore = create<AuthState>((set, get) => ({
  user: null,
  tokens: loadTokens(),
  isAuthenticated: false,
  isLoading: false,
  isInitialized: false,
  error: null,

  login: (user, tokens) => {
    persistTokens(tokens);
    set({
      user,
      tokens,
      isAuthenticated: true,
      error: null,
    });
  },

  logout: () => {
    persistTokens(null);
    set({
      user: null,
      tokens: null,
      isAuthenticated: false,
      error: null,
    });
  },

  refresh: (tokens) => {
    persistTokens(tokens);
    set({ tokens });
  },

  setUser: (user) =>
    set({ user }),

  setError: (error) =>
    set({ error }),

  setLoading: (isLoading) =>
    set({ isLoading }),

  initialize: async () => {
    const { tokens, isInitialized } = get();

    // Don't re-initialize
    if (isInitialized) return;

    if (!tokens?.accessToken) {
      set({ isInitialized: true, isLoading: false });
      return;
    }

    set({ isLoading: true });

    try {
      const user = await api.getMe();
      set({
        user,
        isAuthenticated: true,
        isLoading: false,
        isInitialized: true,
      });
    } catch {
      // Token invalid — try refresh
      if (tokens.refreshToken) {
        try {
          const newTokens = await api.refreshToken(tokens.refreshToken);
          persistTokens(newTokens);
          set({ tokens: newTokens });

          const user = await api.getMe();
          set({
            user,
            isAuthenticated: true,
            isLoading: false,
            isInitialized: true,
          });
          return;
        } catch {
          // Refresh also failed
        }
      }

      // All attempts failed — clear auth
      persistTokens(null);
      set({
        user: null,
        tokens: null,
        isAuthenticated: false,
        isLoading: false,
        isInitialized: true,
      });
    }
  },
}));
