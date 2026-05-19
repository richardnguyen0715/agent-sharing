import { useCallback } from 'react';
import { useAuthStore } from '@/stores/authStore';
import * as api from '@/services/api';
import type { UserRole } from '@/types';

export function useAuth() {
  const {
    user,
    isAuthenticated,
    isLoading,
    isInitialized,
    error,
    login: setAuth,
    logout: clearAuth,
    setError,
    setLoading,
  } = useAuthStore();

  const login = useCallback(async (username: string, password: string) => {
    setError(null);
    setLoading(true);
    try {
      const result = await api.login(username, password);
      setAuth(result.user, result.tokens);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Login failed';
      setError(message);
      throw err;
    } finally {
      setLoading(false);
    }
  }, [setAuth, setError, setLoading]);

  const register = useCallback(async (username: string, password: string, role: UserRole = 'operator') => {
    setError(null);
    setLoading(true);
    try {
      const result = await api.register(username, password, role);
      setAuth(result.user, result.tokens);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Registration failed';
      setError(message);
      throw err;
    } finally {
      setLoading(false);
    }
  }, [setAuth, setError, setLoading]);

  const logout = useCallback(() => {
    clearAuth();
  }, [clearAuth]);

  const clearError = useCallback(() => {
    setError(null);
  }, [setError]);

  return {
    user,
    isAuthenticated,
    isLoading,
    isInitialized,
    error,
    login,
    register,
    logout,
    clearError,
    isAdmin: user?.role === 'admin',
  };
}
