export const STATUS_COLORS: Record<string, string> = {
  running: 'bg-green-500',
  offline: 'bg-gray-500',
  installing: 'bg-yellow-500',
  install_failed: 'bg-red-500',
  suspended: 'bg-orange-500',
}

export const STATUS_TEXT_COLORS: Record<string, string> = {
  running: 'text-green-400',
  offline: 'text-gray-400',
  installing: 'text-yellow-400',
  install_failed: 'text-red-400',
  suspended: 'text-orange-400',
}

export const POWER_ACTIONS = [
  { action: 'start', label: 'Start', color: 'bg-green-600 hover:bg-green-700' },
  { action: 'restart', label: 'Restart', color: 'bg-yellow-600 hover:bg-yellow-700' },
  { action: 'stop', label: 'Stop', color: 'bg-red-600 hover:bg-red-700' },
  { action: 'kill', label: 'Kill', color: 'bg-red-900 hover:bg-red-950' },
]
