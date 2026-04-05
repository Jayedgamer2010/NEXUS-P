import { useEffect } from 'react'
import { Server, Network, Users, Activity } from 'lucide-react'
import { useApi } from '../../hooks/useApi'
import { statsApi } from '../../api/admin/stats'
import type { AdminStats } from '../../types'
import StatCard from '../../components/ui/StatCard'
import { formatDate, formatRelativeTime } from '../../utils/format'
import StatusBadge from '../../components/ui/StatusBadge'

export default function Dashboard() {
  const { data, loading, refetch } = useApi<AdminStats>(() => statsApi.get(), [])

  useEffect(() => {
    const interval = setInterval(refetch, 30000)
    return () => clearInterval(interval)
  }, [])

  return (
    <div>
      <div className="nx-grid-4" style={{ marginBottom: 24 }}>
        <StatCard label="Total Servers" value={data?.servers ?? 0} icon={<Server />} loading={loading} />
        <StatCard label="Total Nodes" value={data?.nodes ?? 0} icon={<Network />} loading={loading} />
        <StatCard label="Total Users" value={data?.users ?? 0} icon={<Users />} loading={loading} />
        <StatCard
          label="Running Servers"
          value={data?.running_servers ?? 0}
          icon={<Activity />}
          accent="#22c55e"
          loading={loading}
        />
      </div>

      <div className="nx-grid-2">
        <div className="card">
          <div className="nx-section-title">Recent Servers</div>
          {loading ? (
            <div style={{ padding: 20, textAlign: 'center', color: '#6b7280' }}>
              Loading...
            </div>
          ) : data?.recent_servers && data.recent_servers.length > 0 ? (
            <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
              {data.recent_servers.map((server) => (
                <div key={server.id} style={{
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'space-between',
                  padding: '8px 10px',
                  background: '#0a0a0f',
                  borderRadius: 6,
                  fontSize: 13,
                }}>
                  <div>
                    <div style={{ fontWeight: 500 }}>{server.name}</div>
                    <div style={{ fontSize: 11, color: '#6b7280' }}>{server.node?.name}</div>
                  </div>
                  <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
                    <StatusBadge status={server.status} />
                    <span style={{ fontSize: 11, color: '#4b5563' }}>{formatRelativeTime(server.created_at)}</span>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div style={{ padding: 20, textAlign: 'center', color: '#6b7280', fontSize: 13 }}>
              No servers yet
            </div>
          )}
        </div>

        <div className="card">
          <div className="nx-section-title">Recent Users</div>
          {loading ? (
            <div style={{ padding: 20, textAlign: 'center', color: '#6b7280' }}>
              Loading...
            </div>
          ) : data?.recent_users && data.recent_users.length > 0 ? (
            <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
              {data.recent_users.map((user) => (
                <div key={user.id} style={{
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'space-between',
                  padding: '8px 10px',
                  background: '#0a0a0f',
                  borderRadius: 6,
                  fontSize: 13,
                }}>
                  <div>
                    <div style={{ fontWeight: 500 }}>{user.username}</div>
                    <div style={{ fontSize: 11, color: '#6b7280' }}>{user.email}</div>
                  </div>
                  <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
                    <span className={`role-badge role-badge--${user.root_admin ? 'admin' : 'client'}`}>
                      {user.root_admin ? 'Admin' : user.role}
                    </span>
                    <span style={{ fontSize: 11, color: '#4b5563' }}>{formatRelativeTime(user.created_at)}</span>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div style={{ padding: 20, textAlign: 'center', color: '#6b7280', fontSize: 13 }}>
              No users yet
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
