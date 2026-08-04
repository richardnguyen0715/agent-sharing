type ClassValue = string | number | boolean | null | undefined;

/** Concatenate CSS class names, filtering out falsy values. */
export function cn(...inputs: ClassValue[]): string {
  return inputs.filter(Boolean).join(' ');
}

/** Format a date string to a human-readable locale string. */
export function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleString();
}

/** Generate a short random ID. */
export function randomId(): string {
  return Math.random().toString(36).substring(2, 10);
}
