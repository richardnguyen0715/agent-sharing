import type { ReactNode } from 'react';

type BadgeVariant = 'default' | 'success' | 'warning' | 'error';

interface BadgeProps {
  variant?: BadgeVariant;
  children: ReactNode;
}

const VARIANT_CLASSES: Record<BadgeVariant, string> = {
  default: 'bg-terminal-border text-terminal-text',
  success: 'bg-status-green/20 text-status-green',
  warning: 'bg-status-amber/20 text-status-amber',
  error: 'bg-status-red/20 text-status-red',
};

export function Badge({ variant = 'default', children }: BadgeProps) {
  return (
    <span
      className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium ${VARIANT_CLASSES[variant]}`}
    >
      {children}
    </span>
  );
}
