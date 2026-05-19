import { useState, type FormEvent } from 'react';
import { Navigate, useNavigate } from 'react-router-dom';
import { IconTerminal2, IconLoader2, IconEye, IconEyeOff } from '@tabler/icons-react';
import { useAuth } from '@/hooks/useAuth';
import type { UserRole } from '@/types';

type Mode = 'login' | 'register';

const ROLE_OPTIONS: { value: UserRole; label: string }[] = [
  { value: 'operator', label: 'Operator' },
  { value: 'viewer', label: 'Viewer' },
];

export function LoginPage() {
  const [mode, setMode] = useState<Mode>('login');
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [role, setRole] = useState<UserRole>('operator');
  const [showPassword, setShowPassword] = useState(false);
  const [localError, setLocalError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const { login, register, isAuthenticated, isLoading } = useAuth();
  const navigate = useNavigate();

  // Already authenticated — redirect to home
  if (isAuthenticated && !isLoading) {
    return <Navigate to="/" replace />;
  }

  const switchMode = () => {
    setMode((m) => (m === 'login' ? 'register' : 'login'));
    setLocalError(null);
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setLocalError(null);

    if (!username.trim() || !password.trim()) {
      setLocalError('Username and password are required.');
      return;
    }

    if (mode === 'register' && password.length < 6) {
      setLocalError('Password must be at least 6 characters.');
      return;
    }

    setSubmitting(true);
    try {
      if (mode === 'login') {
        await login(username, password);
      } else {
        await register(username, password, role);
      }
      navigate('/', { replace: true });
    } catch (err) {
      setLocalError(err instanceof Error ? err.message : 'Something went wrong');
    } finally {
      setSubmitting(false);
    }
  };

  const isDisabled = submitting || isLoading;

  return (
    <div className="flex min-h-screen items-center justify-center bg-terminal-bg px-4">
      {/* Subtle grid background effect */}
      <div className="pointer-events-none fixed inset-0 opacity-[0.03]"
        style={{
          backgroundImage: 'linear-gradient(rgba(230,237,243,0.1) 1px, transparent 1px), linear-gradient(90deg, rgba(230,237,243,0.1) 1px, transparent 1px)',
          backgroundSize: '40px 40px',
        }}
      />

      <div className="relative w-full max-w-[420px]">
        {/* Brand header */}
        <div className="mb-8 flex flex-col items-center gap-3">
          <div className="flex h-14 w-14 items-center justify-center rounded-xl bg-brand/15 ring-1 ring-brand/30">
            <IconTerminal2 size={28} className="text-brand" />
          </div>
          <div className="text-center">
            <h1 className="text-2xl font-bold tracking-tight text-terminal-text">
              TermHub
            </h1>
            <p className="mt-1 text-sm text-terminal-output">
              {mode === 'login'
                ? 'Sign in to your terminal workspace'
                : 'Create a new account'}
            </p>
          </div>
        </div>

        {/* Card */}
        <div className="rounded-xl border border-terminal-border bg-terminal-bg-toolbar p-8 shadow-2xl shadow-black/40">
          <form onSubmit={handleSubmit} className="space-y-5">
            {/* Username */}
            <div>
              <label
                htmlFor="username"
                className="mb-1.5 block text-sm font-medium text-terminal-output"
              >
                Username
              </label>
              <input
                id="username"
                type="text"
                autoComplete="username"
                autoFocus
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                disabled={isDisabled}
                className="w-full rounded-lg border border-terminal-border bg-terminal-bg px-3.5 py-2.5 text-sm text-terminal-text placeholder-terminal-output/50 outline-none transition-colors focus:border-brand focus:ring-1 focus:ring-brand/40 disabled:opacity-50"
                placeholder="Enter your username"
                required
              />
            </div>

            {/* Password */}
            <div>
              <label
                htmlFor="password"
                className="mb-1.5 block text-sm font-medium text-terminal-output"
              >
                Password
              </label>
              <div className="relative">
                <input
                  id="password"
                  type={showPassword ? 'text' : 'password'}
                  autoComplete={mode === 'login' ? 'current-password' : 'new-password'}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  disabled={isDisabled}
                  className="w-full rounded-lg border border-terminal-border bg-terminal-bg px-3.5 py-2.5 pr-10 text-sm text-terminal-text placeholder-terminal-output/50 outline-none transition-colors focus:border-brand focus:ring-1 focus:ring-brand/40 disabled:opacity-50"
                  placeholder="Enter your password"
                  required
                />
                <button
                  type="button"
                  tabIndex={-1}
                  onClick={() => setShowPassword((v) => !v)}
                  className="absolute right-2.5 top-1/2 -translate-y-1/2 rounded p-0.5 text-terminal-output/60 hover:text-terminal-text"
                  aria-label={showPassword ? 'Hide password' : 'Show password'}
                >
                  {showPassword ? <IconEyeOff size={16} /> : <IconEye size={16} />}
                </button>
              </div>
            </div>

            {/* Role selector (register only) */}
            {mode === 'register' && (
              <div>
                <label
                  htmlFor="role"
                  className="mb-1.5 block text-sm font-medium text-terminal-output"
                >
                  Role
                </label>
                <select
                  id="role"
                  value={role}
                  onChange={(e) => setRole(e.target.value as UserRole)}
                  disabled={isDisabled}
                  className="w-full rounded-lg border border-terminal-border bg-terminal-bg px-3.5 py-2.5 text-sm text-terminal-text outline-none transition-colors focus:border-brand focus:ring-1 focus:ring-brand/40 disabled:opacity-50"
                >
                  {ROLE_OPTIONS.map((opt) => (
                    <option key={opt.value} value={opt.value}>
                      {opt.label}
                    </option>
                  ))}
                </select>
              </div>
            )}

            {/* Error */}
            {localError && (
              <div className="rounded-lg border border-terminal-error/30 bg-terminal-error/10 px-3.5 py-2.5 text-sm text-terminal-error">
                {localError}
              </div>
            )}

            {/* Submit */}
            <button
              type="submit"
              disabled={isDisabled}
              className="flex w-full items-center justify-center gap-2 rounded-lg bg-brand py-2.5 text-sm font-semibold text-white transition-colors hover:bg-brand-hover focus:outline-none focus:ring-2 focus:ring-brand/50 focus:ring-offset-2 focus:ring-offset-terminal-bg-toolbar disabled:cursor-not-allowed disabled:opacity-50"
            >
              {submitting ? (
                <>
                  <IconLoader2 size={16} className="animate-spin" />
                  {mode === 'login' ? 'Signing in…' : 'Creating account…'}
                </>
              ) : (
                mode === 'login' ? 'Sign In' : 'Create Account'
              )}
            </button>
          </form>

          {/* Mode toggle */}
          <div className="mt-6 border-t border-terminal-border pt-5 text-center">
            <p className="text-sm text-terminal-output">
              {mode === 'login' ? "Don't have an account?" : 'Already have an account?'}{' '}
              <button
                type="button"
                onClick={switchMode}
                className="font-medium text-brand transition-colors hover:text-brand-hover"
              >
                {mode === 'login' ? 'Register' : 'Sign In'}
              </button>
            </p>
          </div>
        </div>

        {/* Footer */}
        <p className="mt-6 text-center text-xs text-terminal-output/50">
          Secure terminal management platform
        </p>
      </div>
    </div>
  );
}
