import { useState, useEffect, useCallback } from 'react';
import { wingsApi } from '../api/wings';
import { ServerResources } from '../types';

interface UseServerStatsOptions {
  serverUUID: string;
  interval?: number;
}

export function useServerStats({ serverUUID, interval = 3000 }: UseServerStatsOptions) {
  const [stats, setStats] = useState<ServerResources | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchStats = useCallback(async () => {
    try {
      const data = await wingsApi.getResources(serverUUID);
      setStats(data);
      setError(null);
    } catch (err: any) {
      console.error('Failed to fetch server stats:', err);
      setError(err.message || 'Failed to fetch stats');
    } finally {
      setLoading(false);
    }
  }, [serverUUID]);

  useEffect(() => {
    fetchStats();

    const timer = setInterval(() => {
      fetchStats();
    }, interval);

    return () => clearInterval(timer);
  }, [fetchStats, interval]);

  return {
    stats,
    loading,
    error,
    refetch: fetchStats,
  };
}
