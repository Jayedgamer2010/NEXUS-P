import { useEffect, useState } from 'react';
import StatCard from '../../components/ui/StatCard';
import DataTable from '../../components/ui/DataTable';
import { adminApi } from '../../api/admin';
import type { Server } from '../../types';
import './Dashboard.css';

export default function Dashboard() {
  const [servers, setServers] = useState<Server[]>([]);
  const [nodesCount, setNodesCount] = useState(0);
  const [usersCount, setUsersCount] = useState(0);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    setLoading(true);
    try {
      const [serversRes, nodesRes, usersRes] = await Promise.all([
        adminApi.getServers(1, 100),
        adminApi.getNodes(1, 1),
        adminApi.getUsers(1, 1),
      ]);
      setServers(serversRes.data);
      setNodesCount(nodesRes.total);
      setUsersCount(usersRes.total);
    } catch (error) {
      console.error('Failed to load dashboard:', error);
    } finally {
      setLoading(false);
    }
  };

  const runningServers = servers.filter(s => s.status === 'running').length;
  const uptimePercent = servers.length > 0 ? Math.round((runningServers / servers.length) * 100) : 0;

  const formatDate = (dateStr: string) => {
    try {
      return new Date(dateStr).toLocaleDateString();
    } catch {
      return dateStr;
    }
  };

  return (
    <div className="dashboard">
      <div className="stats-grid">
        <StatCard
          title="Total Servers"
          value={servers.length}
          icon="🖥️"
          accentColor="#7c3aed"
        />
        <StatCard
          title="Total Nodes"
          value={nodesCount}
          icon="🌐"
          accentColor="#3b82f6"
        />
        <StatCard
          title="Total Users"
          value={usersCount}
          icon="👥"
          accentColor="#10b981"
        />
        <StatCard
          title="Running Servers"
          value={runningServers}
          icon="▶️"
          accentColor="#10b981"
          subtitle={`${uptimePercent}% of ${servers.length} total`}
        />
      </div>

      <div className="recent-servers">
        <h2 className="section-title">Recent Servers</h2>
        <DataTable
          columns={[
            { key: 'name', header: 'Name' },
            {
              key: 'node',
              header: 'Node',
              render: (_val: string, row: Server) => row.node?.name || 'N/A',
            },
            {
              key: 'status',
              header: 'Status',
              render: (status: string) => (
                <span className={`status-badge ${status}`}>{status}</span>
              ),
            },
            {
              key: 'user',
              header: 'Owner',
              render: (_val: string, row: Server) => row.user?.username || `ID: ${row.user_id}`,
            },
            {
              key: 'created_at',
              header: 'Created',
              render: (val: string) => formatDate(val),
            },
          ]}
          data={servers.slice(0, 10)}
          loading={loading}
        />
      </div>
    </div>
  );
}
