import { useState, useEffect, useRef } from 'react'
import { clientServersApi } from '../api/client/servers'
import type { ServerResources } from '../types'

export function useServerStats(uuid: string, interval = 3000) {
  const [stats, setStats] = useState<ServerResources | null>(null)
  const [error, setError] = useState<string | null>(null)
  const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null)

  useEffect(() => {
    if (!uuid) return

    const fetchStats = async () => {
      try {
        const res = await clientServersApi.getResources(uuid)
        setStats(res.data?.data ?? res.data)
        setError(null)
      } catch {
        setError('Failed to fetch stats')
      }
    }

    fetchStats()
    intervalRef.current = setInterval(fetchStats, interval)

    return () => {
      if (intervalRef.current) clearInterval(intervalRef.current)
    }
  }, [uuid, interval])

  return { stats, error }
}
