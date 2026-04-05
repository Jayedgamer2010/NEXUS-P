import { ReactNode } from 'react'
import Spinner from './Spinner'

interface StatCardProps {
  label: string
  value: string | number
  icon?: ReactNode
  accent?: string
  loading?: boolean
  trend?: { value: number; positive: boolean }
}

export default function StatCard({ label, value, icon, accent = '#7c3aed', loading, trend }: StatCardProps) {
  return (
    <div className="nx-stat-card" style={{ borderTopColor: accent }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
        <div>
          {loading ? (
            <>
              <div className="nx-skeleton" style={{ height: 12, width: 60, marginBottom: 8 }} />
              <div className="nx-skeleton" style={{ height: 24, width: 40 }} />
            </>
          ) : (
            <>
              <div style={{ fontSize: 12, color: '#6b7280', marginBottom: 6, textTransform: 'uppercase', letterSpacing: 0.5 }}>
                {label}
              </div>
              <div style={{ fontSize: 28, fontWeight: 700 }}>{value}</div>
              {trend && (
                <div style={{ fontSize: 12, color: trend.positive ? '#22c55e' : '#ef4444', marginTop: 4 }}>
                  {trend.positive ? '+' : ''}{trend.value}% from last week
                </div>
              )}
            </>
          )}
        </div>
        {!loading && icon && <div style={{ color: accent, opacity: 0.6 }}>{icon}</div>}
      </div>
    </div>
  )
}
