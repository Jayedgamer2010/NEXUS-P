import { useState, useCallback, useEffect } from 'react'
import { useParams } from 'react-router-dom'
import { Terminal } from '@xterm/xterm'
import { useApi } from '../../hooks/useApi'
import { useServerStats } from '../../hooks/useServerStats'
import { useConsole } from '../../hooks/useConsole'
import { serversApi } from '../../api/admin/servers'
import type { Server } from '../../types'
import Button from '../../components/ui/Button'
import StatusBadge from '../../components/ui/StatusBadge'
import Console from '../../components/console/Console'
import ConsoleInput from '../../components/console/ConsoleInput'
import ConfirmDialog from '../../components/ui/ConfirmDialog'
import { formatMB, formatCPU, formatBytes, formatUptime } from '../../utils/format'
import { POWER_ACTIONS } from '../../utils/constants'
import Spinner from '../../components/ui/Spinner'

export default function ServerDetail() {
  const { id } = useParams<{ id: string }>()
  const serverId = Number(id)
  const [loadingPower, setLoadingPower] = useState<string | null>(null)
  const [confirmAction, setConfirmAction] = useState<{ action: string; label: string } | null>(null)
  const [terminal, setTerminal] = useState<Terminal | null>(null)
  const [key, setKey] = useState(0)

  const { data: server, loading, refetch } = useApi<Server>(
    () => serversApi.getOne(serverId),
    [serverId]
  )

  const { stats, error: statsError } = useServerStats(server?.uuid ?? '', 3000)
  const { sendCommand } = useConsole(terminal, server?.uuid ?? '')

  const handleTerminalReady = useCallback((t: Terminal) => {
    setTerminal(t)
  }, [])

  const handlePower = async (action: string) => {
    if (!server) return
    setLoadingPower(action)
    try {
      await serversApi.power(server.id, action)
      setTimeout(refetch, 2000)
    } catch { /* ignore */ } finally {
      setLoadingPower(null)
    }
  }

  const handleConfirmPower = (action: string, label: string) => {
    if (action === 'stop' || action === 'kill') {
      setConfirmAction({ action, label })
    } else {
      handlePower(action)
    }
  }

  const handlePowerConfirm = () => {
    if (confirmAction) handlePower(confirmAction.action)
  }

  const handleRefreshConsole = () => {
    setKey((k) => k + 1)
  }

  if (loading || !server) {
    return <div style={{ padding: 40, textAlign: 'center' }}><Spinner size="lg" /></div>
  }

  return (
    <div>
      {/* Header */}
      <div style={{ display: 'flex', alignItems: 'flex-start', justifyContent: 'space-between', marginBottom: 24 }}>
        <div>
          <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 8 }}>
            <h2 style={{ fontSize: 22, fontWeight: 700 }}>{server.name}</h2>
            <StatusBadge status={server.suspended ? 'suspended' : server.status} />
          </div>
          <div style={{ display: 'flex', gap: 8, alignItems: 'center' }}>
            <span className="nx-mono" style={{ color: '#4b5563' }}>{server.uuid_short}</span>
            {server.node && <span className="nx-tag">{server.node.name}</span>}
            {server.egg && <span className="nx-tag">{server.egg.name}</span>}
            {server.user && <span className="nx-tag">Owner: {server.user.username}</span>}
          </div>
        </div>
      </div>

      {/* Power Actions */}
      <div style={{ marginBottom: 24 }}>
        <div className="nx-section-title">Power</div>
        <div className="power-btns">
          {POWER_ACTIONS.map(({ action, label, color }) => (
            <Button
              key={action}
              variant={action === 'start' ? 'primary' : action === 'restart' ? 'warning' : 'danger'}
              size="sm"
              loading={loadingPower === action}
              onClick={() => handleConfirmPower(action, label)}
            >
              {label}
            </Button>
          ))}
          <Button variant="ghost" size="sm" onClick={handleRefreshConsole}>
            Refresh Console
          </Button>
        </div>
      </div>

      {/* Stats */}
      <div className="nx-grid-4" style={{ marginBottom: 24 }}>
        <div className="nx-stat-card">
          <div style={{ fontSize: 12, color: '#6b7280', marginBottom: 4 }}>CPU</div>
          <div style={{ fontSize: 20, fontWeight: 600 }}>
            {stats ? formatCPU(stats.cpu_absolute) : '--'}
          </div>
        </div>
        <div className="nx-stat-card">
          <div style={{ fontSize: 12, color: '#6b7280', marginBottom: 4 }}>RAM</div>
          <div style={{ fontSize: 20, fontWeight: 600 }}>
            {stats ? formatBytes(stats.memory_bytes) : '--'}/ {formatMB(server.memory)}
          </div>
        </div>
        <div className="nx-stat-card">
          <div style={{ fontSize: 12, color: '#6b7280', marginBottom: 4 }}>Disk</div>
          <div style={{ fontSize: 20, fontWeight: 600 }}>
            {stats ? formatBytes(stats.disk_bytes) : '--'}/ {formatMB(server.disk)}
          </div>
        </div>
        <div className="nx-stat-card">
          <div style={{ fontSize: 12, color: '#6b7280', marginBottom: 4 }}>Uptime</div>
          <div style={{ fontSize: 20, fontWeight: 600 }}>
            {stats && stats.uptime > 0 ? formatUptime(stats.uptime) : 'Offline'}
          </div>
        </div>
      </div>

      {/* Server Info */}
      <div className="nx-grid-2" style={{ marginBottom: 24 }}>
        <div className="card">
          <div className="nx-section-title">Configuration</div>
          <div style={{ fontSize: 13, lineHeight: 2 }}>
            <div style={{ display: 'flex', justifyContent: 'space-between' }}>
              <span style={{ color: '#6b7280' }}>Memory</span>
              <span>{formatMB(server.memory)}</span>
            </div>
            <div style={{ display: 'flex', justifyContent: 'space-between' }}>
              <span style={{ color: '#6b7280' }}>Disk</span>
              <span>{formatMB(server.disk)}</span>
            </div>
            <div style={{ display: 'flex', justifyContent: 'space-between' }}>
              <span style={{ color: '#6b7280' }}>CPU</span>
              <span>{server.cpu}%</span>
            </div>
            <div style={{ display: 'flex', justifyContent: 'space-between' }}>
              <span style={{ color: '#6b7280' }}>Image</span>
              <span className="nx-mono" style={{ color: '#9ca3af' }}>{server.image || server.egg?.docker_image}</span>
            </div>
          </div>
        </div>
        <div className="card">
          <div className="nx-section-title">Connection</div>
          {server.allocation ? (
            <div style={{ fontSize: 13, lineHeight: 2 }}>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: '#6b7280' }}>Address</span>
                <span className="nx-mono">{server.allocation.ip}</span>
              </div>
              <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                <span style={{ color: '#6b7280' }}>Port</span>
                <span className="nx-mono">{server.allocation.port}</span>
              </div>
              {server.allocation.ip_alias && (
                <div style={{ display: 'flex', justifyContent: 'space-between' }}>
                  <span style={{ color: '#6b7280' }}>Alias</span>
                  <span className="nx-mono">{server.allocation.ip_alias}</span>
                </div>
              )}
            </div>
          ) : (
            <span style={{ color: '#6b7280', fontSize: 13 }}>No allocation assigned</span>
          )}
        </div>
      </div>

      {/* Console */}
      <div style={{ marginBottom: 24 }}>
        <div className="nx-section-title">Console</div>
        <div key={key}>
          <Console serverUUID={server.uuid} onTerminalReady={handleTerminalReady} />
          <ConsoleInput onSend={sendCommand} />
        </div>
      </div>

      <ConfirmDialog
        isOpen={!!confirmAction}
        onClose={() => setConfirmAction(null)}
        onConfirm={handlePowerConfirm}
        title="Confirm Power Action"
        message={`Are you sure you want to ${confirmAction?.label?.toLowerCase()} this server?`}
        confirmText={confirmAction?.label ?? 'Confirm'}
      />
    </div>
  )
}
