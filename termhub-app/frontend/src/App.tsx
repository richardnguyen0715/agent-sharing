import { useEffect } from 'react';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { IconLoader2 } from '@tabler/icons-react';
import { LoginPage } from '@/pages/LoginPage';
import { MainPage } from '@/pages/MainPage';
import { AdminPage } from '@/pages/AdminPage';
import { RecordingsPage } from '@/pages/RecordingsPage';
import { NewTerminalModal } from '@/components/modals/NewTerminalModal';
import { useAuthStore } from '@/stores/authStore';
import { useUIStore } from '@/stores/uiStore';
import type { ReactNode } from 'react';

function LoadingScreen() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-terminal-bg">
      <div className="flex flex-col items-center gap-3">
        <IconLoader2 size={32} className="animate-spin text-brand" />
        <p className="text-sm text-terminal-output">Loading…</p>
      </div>
    </div>
  );
}

function ProtectedRoute({ children }: { children: ReactNode }) {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
  const isLoading = useAuthStore((s) => s.isLoading);
  const isInitialized = useAuthStore((s) => s.isInitialized);

  if (!isInitialized || isLoading) {
    return <LoadingScreen />;
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return <>{children}</>;
}

function AdminRoute({ children }: { children: ReactNode }) {
  const user = useAuthStore((s) => s.user);
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
  const isLoading = useAuthStore((s) => s.isLoading);
  const isInitialized = useAuthStore((s) => s.isInitialized);

  if (!isInitialized || isLoading) {
    return <LoadingScreen />;
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (user?.role !== 'admin') {
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
}

function ModalLayer() {
  const modal = useUIStore((s) => s.modal);

  if (modal === 'newTerminal') {
    return <NewTerminalModal />;
  }

  return null;
}

function AppInitializer({ children }: { children: ReactNode }) {
  const initialize = useAuthStore((s) => s.initialize);
  const isInitialized = useAuthStore((s) => s.isInitialized);

  useEffect(() => {
    void initialize();
  }, [initialize]);

  if (!isInitialized) {
    return <LoadingScreen />;
  }

  return <>{children}</>;
}

export function App() {
  return (
    <BrowserRouter>
      <AppInitializer>
        <Routes>
          <Route path="/login" element={<LoginPage />} />
          <Route
            path="/"
            element={
              <ProtectedRoute>
                <MainPage />
              </ProtectedRoute>
            }
          />
          <Route
            path="/admin"
            element={
              <AdminRoute>
                <AdminPage />
              </AdminRoute>
            }
          />
          <Route
            path="/recordings"
            element={
              <ProtectedRoute>
                <RecordingsPage />
              </ProtectedRoute>
            }
          />
        </Routes>
        <ModalLayer />
      </AppInitializer>
    </BrowserRouter>
  );
}
