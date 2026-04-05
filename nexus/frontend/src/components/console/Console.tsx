import { useEffect, useRef } from 'react'
import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'
import { WebLinksAddon } from '@xterm/addon-web-links'
import '@xterm/xterm/css/xterm.css'

interface ConsoleProps {
  serverUUID: string
  onTerminalReady?: (terminal: Terminal) => void
}

export default function Console({ serverUUID, onTerminalReady }: ConsoleProps) {
  const containerRef = useRef<HTMLDivElement>(null)
  const terminalRef = useRef<Terminal | null>(null)

  useEffect(() => {
    if (!containerRef.current) return

    const terminal = new Terminal({
      convertEol: true,
      cursorBlink: true,
      fontFamily: 'JetBrains Mono, Menlo, monospace',
      fontSize: 13,
      scrollback: 1000,
      theme: {
        background: '#0a0a0f',
        foreground: '#e2e8f0',
        cursor: '#7c3aed',
        selectionBackground: '#7c3aed40',
      },
    })

    const fitAddon = new FitAddon()
    const webLinksAddon = new WebLinksAddon()

    terminal.loadAddon(fitAddon)
    terminal.loadAddon(webLinksAddon)
    terminal.open(containerRef.current)
    fitAddon.fit()

    const ro = new ResizeObserver(() => fitAddon.fit())
    ro.observe(containerRef.current)

    terminalRef.current = terminal
    if (onTerminalReady) onTerminalReady(terminal)

    return () => {
      ro.disconnect()
      terminal.dispose()
    }
  }, [serverUUID])

  return (
    <div
      ref={containerRef}
      style={{
        width: '100%',
        minHeight: 400,
        flex: 1,
        borderRadius: 8,
        overflow: 'hidden',
      }}
    />
  )
}
