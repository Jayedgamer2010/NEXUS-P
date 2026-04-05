import { STATUS_COLORS, STATUS_TEXT_COLORS } from '../../utils/constants'
import Spinner from './Spinner'

interface StatusBadgeProps {
  status: string
}

export default function StatusBadge({ status }: StatusBadgeProps) {
  const bgColor = STATUS_COLORS[status] || 'bg-gray-500'
  const textColor = STATUS_TEXT_COLORS[status] || 'text-gray-400'
  const isInstalling = status === 'installing'

  return (
    <span className={`nx-badge ${bgColor.replace('bg-', 'nx-')}`} style={{
      background: '#0d0d17',
      border: `1px solid ${textColor.includes('red') ? '#991b1b' : textColor.includes('green') ? '#166534' : textColor.includes('yellow') ? '#854d0e' : textColor.includes('orange') ? '#9a3412' : '#374151'}`,
    }}>
      <span className="nx-badge-dot" style={{
        background: textColor.includes('green') ? '#22c55e' : textColor.includes('yellow') ? '#f59e0b' : textColor.includes('red') ? '#ef4444' : textColor.includes('gray') ? '#6b7280' : '#f97316',
        animation: isInstalling ? 'pulse 1s infinite' : 'none',
      }} />
      <span style={{ color: textColor.includes('green') ? '#22c55e' : textColor.includes('yellow') ? '#f59e0b' : textColor.includes('red') ? '#ef4444' : textColor.includes('gray') ? '#9ca3af' : '#f97316' }}>
        {status.replace('_', ' ')}
      </span>
    </span>
  )
}
