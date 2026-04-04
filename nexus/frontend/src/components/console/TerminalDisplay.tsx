import { useEffect, useRef } from 'react';
import { Terminal } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import './TerminalDisplay.css';

interface TerminalDisplayProps {
  onMessage?: (message: string) => void;
  connectionStatus: 'connecting' | 'connected' | 'disconnected';
}

export default function TerminalDisplay({ onMessage, connectionStatus }: TerminalDisplayProps) {
  const terminalRef = useRef<HTMLDivElement>(null);
  const xtermRef = useRef<Terminal | null>(null);
  const fitAddonRef = useRef<FitAddon | null>(null);

  useEffect(() => {
    if (!terminalRef.current) return;

    const term = new Terminal({
      cursorBlink: true,
      fontSize: 14,
      fontFamily: 'JetBrains Mono, Fira Code, monospace',
      theme: {
        background: '#000000',
        foreground: '#10b981',
        cursor: '#10b981',
      } as any,
      scrollback: 1000,
    });

    const fitAddon = new FitAddon();
    term.loadAddon(fitAddon);
    term.open(terminalRef.current);
    requestAnimationFrame(() => fitAddon.fit());

    xtermRef.current = term;
    fitAddonRef.current = fitAddon;

    term.writeln('[NEXUS] Console initialized');

    return () => {
      term.dispose();
    };
  }, []);

  // External message handling
  useEffect(() => {
    if (xtermRef.current && onMessage) {
      // Note: the parent component handles onMessage in its useConsole hook
    }
  }, [onMessage]);

  useEffect(() => {
    const handleResize = () => {
      if (fitAddonRef.current) {
        fitAddonRef.current.fit();
      }
    };

    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  const getStatusColor = () => {
    switch (connectionStatus) {
      case 'connected':
        return '#10b981';
      case 'connecting':
        return '#f59e0b';
      default:
        return '#ef4444';
    }
  };

  const getStatusText = () => {
    switch (connectionStatus) {
      case 'connected':
        return 'Connected';
      case 'connecting':
        return 'Connecting...';
      default:
        return 'Disconnected';
    }
  };

  return (
    <div className="console-container">
      <div className="console-header">
        <div className="console-status">
          <span
            className="status-indicator"
            style={{ backgroundColor: getStatusColor() }}
          />
          {getStatusText()}
        </div>
        <div className="console-title">Server Console</div>
      </div>
      <div ref={terminalRef} className="console-terminal" />
    </div>
  );
}
