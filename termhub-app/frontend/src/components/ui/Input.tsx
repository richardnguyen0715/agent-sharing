import { forwardRef, type InputHTMLAttributes } from 'react';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, id, className = '', ...props }, ref) => {
    return (
      <div className="space-y-1">
        {label && (
          <label htmlFor={id} className="block text-sm text-terminal-output">
            {label}
          </label>
        )}
        <input
          ref={ref}
          id={id}
          className={`w-full rounded border border-terminal-border bg-terminal-bg px-3 py-2 text-terminal-text outline-none transition-colors focus:border-brand ${
            error ? 'border-terminal-error' : ''
          } ${className}`}
          {...props}
        />
        {error && <p className="text-xs text-terminal-error">{error}</p>}
      </div>
    );
  },
);

Input.displayName = 'Input';
