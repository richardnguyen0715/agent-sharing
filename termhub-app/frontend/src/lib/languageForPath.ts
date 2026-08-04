const EXT_LANGUAGE: Record<string, string> = {
  js: 'javascript',
  jsx: 'javascript',
  mjs: 'javascript',
  cjs: 'javascript',
  ts: 'typescript',
  tsx: 'typescript',
  py: 'python',
  go: 'go',
  rs: 'rust',
  rb: 'ruby',
  php: 'php',
  java: 'java',
  c: 'c',
  h: 'c',
  cpp: 'cpp',
  cc: 'cpp',
  hpp: 'cpp',
  cs: 'csharp',
  json: 'json',
  yaml: 'yaml',
  yml: 'yaml',
  toml: 'ini',
  md: 'markdown',
  markdown: 'markdown',
  html: 'html',
  htm: 'html',
  css: 'css',
  scss: 'scss',
  less: 'less',
  sh: 'shell',
  bash: 'shell',
  zsh: 'shell',
  sql: 'sql',
  xml: 'xml',
  graphql: 'graphql',
  lua: 'lua',
  swift: 'swift',
  kt: 'kotlin',
  kts: 'kotlin',
  proto: 'proto',
  env: 'ini',
  ini: 'ini',
  conf: 'ini',
};

/** Maps a file path to a Monaco language id for syntax highlighting. */
export function languageForPath(path: string): string {
  const base = path.split('/').pop() ?? path;
  const lower = base.toLowerCase();
  if (lower === 'dockerfile') return 'dockerfile';
  if (lower === 'makefile') return 'makefile';
  const dot = base.lastIndexOf('.');
  if (dot <= 0) return 'plaintext';
  return EXT_LANGUAGE[lower.slice(dot + 1)] ?? 'plaintext';
}
