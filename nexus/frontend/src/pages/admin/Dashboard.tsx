import { useEffect, useState } from 'react';
import StatCard from '../../components/ui/StatCard';
import DataTable from '../../components/ui/DataTable';
import { adminApi } from '../../api/admin';
import { Server } from '../../types';
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
      const [serversRes] = await Promise.all([
        adminApi.getServers(1, 100),
      ]);
      setServers(serversRes.data);

      // Get counts from first page total
      setNodesCount(serversRes.total); // Placeholder - would need separate node count API
      setUsersCount(serversRes.total); // Placeholder - would need separate user count API
    } catch (error) {
      console.error('Failed to load dashboard:', error);
    } finally {
      setLoading(false);
    }
  };

  const runningServers = servers.filter(s => s.status === 'running').length;

  const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString();
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
          subtitle={`${servers.length > 0 ? Math.round((runningServers / servers.length) * 100) : 0}% uptime`}
        />
      </div>

      <div className="recent-servers">
        <h2 className="section-title">Recent Servers</h2>
        <DataTable
          columns={[
            { key: 'name', header: 'Name' },
            { key: 'node', header: 'Node', render: (_, row) => row.node?.name || 'N/A' },
            { key: 'status', header: 'Status', render: (status) => (
              <span className={`status-badge ${status}`}>{status}</span>
            )},
            { key: 'user', header: 'Owner', render: (_, row) => row.user_id },
            { key: 'created_at', header: 'Created', render: (value) => formatDate(value) },
          ]}
          data={servers.slice(0, 10)}
          loading={loading}
        />
      </div>
    </div>
  );
}
