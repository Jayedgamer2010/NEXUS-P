import { useEffect, useRef } from 'react'
import type { Terminal } from '@xterm/xterm'

export function useConsole(terminal: Terminal | null, serverUUID: string) {
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const attemptsRef = useRef(0)
  const maxAttempts = 5

  const connect = () => {
    if (!terminal || !serverUUID) return

    const token = (() => {
      try {
        const stored = localStorage.getItem('nexus-auth')
        return stored ? JSON.parse(stored)?.state?.token : null
      } catch { return null }
    })()

    const apiUrl = import.meta.env.VITE_API_URL || window.location.origin
    const wsUrl = apiUrl.replace('https://', 'wss://').replace('http://', 'ws://')
    const url = wsUrl + '/ws/console?server=' + serverUUID + '&token=' + token

    terminal.writeln('\r\n\x1b[33mConnecting to server console...\x1b[0m')

    const ws = new WebSocket(url)
    wsRef.current = ws

    ws.onopen = () => {
      attemptsRef.current = 0
      terminal.writeln('\x1b[32mConnected.\x1b[0m\r\n')
    }

    ws.onmessage = (event) => {
      terminal.write(event.data)
    }

    ws.onclose = () => {
      if (attemptsRef.current < maxAttempts) {
        attemptsRef.current++
        terminal.writeln('\r\n\x1b[31mDisconnected. Reconnecting in 3s...\x1b[0m')
        reconnectRef.current = setTimeout(connect, 3000)
      } else {
        terminal.writeln('\r\n\x1b[31mFailed to connect after ' + maxAttempts + ' attempts.\x1b[0m')
      }
    }

    ws.onerror = () => {
      terminal.writeln('\r\n\x1b[31mConnection error.\x1b[0m')
    }
  }

  const sendCommand = (command: string) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({ event: 'send command', args: [command] }))
    }
  }

  useEffect(() => {
    connect()
    return () => {
      if (reconnectRef.current) clearTimeout(reconnectRef.current)
      wsRef.current?.close()
    }
  }, [terminal, serverUUID])

  return { sendCommand }
}
