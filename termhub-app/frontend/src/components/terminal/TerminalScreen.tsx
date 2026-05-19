interface TerminalScreenProps {
  terminalRef: (el: HTMLDivElement | null) => void;
}

export function TerminalScreen({ terminalRef }: TerminalScreenProps) {
  return (
    <div
      ref={terminalRef}
      className="flex-1 bg-terminal-bg p-1"
      data-terminal-container
    />
  );
}
