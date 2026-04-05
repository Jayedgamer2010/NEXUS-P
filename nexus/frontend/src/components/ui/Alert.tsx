import { useState } from 'react'

interface AlertProps {
  type: 'success' | 'error' | 'warning' | 'info'
  message: string
  dismissible?: boolean
}

const icons: Record<string, string> = {
  success: '+',
  error: '!',
  warning: '!',
  info: 'i',
}

export default function Alert({ type, message, dismissible = true }: AlertProps) {
  const [visible, setVisible] = useState(true)
  if (!visible) return null

  return (
    <div className={`nx-alert nx-alert--${type}`}>
      <span style={{ fontWeight: 600, fontSize: 16, lineHeight: 1 }}>{icons[type]}</span>
      <div style={{ flex: 1 }}>{message}</div>
      {dismissible && (
        <button className="nx-alert-dismiss" onClick={() => setVisible(false)}>x</button>
      )}
    </div>
  )
}
