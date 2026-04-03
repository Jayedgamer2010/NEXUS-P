import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import DataTable from '../../components/ui/DataTable';
import Modal from '../../components/ui/Modal';
import { adminApi } from '../../api/admin';
import { Server } from '../../types';
import { PowerAction } from '../../types';
import './Servers.css';

export default function Servers() {
  const [servers, setServers] = useState<Server[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [selectedServer, setSelectedServer] = useState<Server | null>(null);
  const [nodes, setNodes] = useState<{id: number; name: string}[]>([]);
  const [eggs, setEggs] = useState<{id: number; name: string}[]>([]);
  const [users, setUsers] = useState<{id: number; username: string}[]>([]);
  const limit = 20;

  const navigate = useNavigate();

  useEffect(() => {
    loadData();
  }, [page]);

  const loadData = async () => {
    setLoading(true);
    try {
      const [serversRes, nodesRes, eggsRes, usersRes] = await Promise.all([
        adminApi.getServers(page, limit),
        adminApi.getNodes(1, 100),
        adminApi.getEggs(1, 100),
        adminApi.getUsers(1, 100),
      ]);
      setServers(serversRes.data);
      setTotal(serversRes.total);
      setNodes(nodesRes.data.map(n => ({ id: n.id, name: n.name })));
      setEggs(eggsRes.data.map(e => ({ id: e.id, name: e.name })));
      setUsers(usersRes.data.map(u => ({ id: u.id, username: u.username })));
    } catch (error) {
      console.error('Failed to load servers:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateServer = async (formData: FormData) => {
    try {
      await adminApi.createServer({
        name: formData.get('name') as string,
        user_id: parseInt(formData.get('user_id') as string),
        node_id: parseInt(formData.get('node_id') as string),
        egg_id: parseInt(formData.get('egg_id') as string),
        allocation_id: 1, // TODO: fetch available allocations
        memory: parseInt(formData.get('memory') as string),
        disk: parseInt(formData.get('disk') as string),
        cpu: parseInt(formData.get('cpu') as string),
      });
      setShowCreateModal(false);
      loadData();
    } catch (error) {
      console.error('Failed to create server:', error);
      alert('Failed to create server');
    }
  };

  const handleDeleteServer = async () => {
    if (!selectedServer) return;
    try {
      await adminApi.deleteServer(selectedServer.id);
      setShowDeleteModal(false);
      setSelectedServer(null);
      loadData();
    } catch (error) {
      console.error('Failed to delete server:', error);
      alert('Failed to delete server');
    }
  };

  const handlePowerAction = async (serverId: number, action: PowerAction) => {
    try {
      await adminApi.powerServer(serverId, action);
    } catch (error) {
      console.error('Failed to send power action:', error);
      alert('Failed to send power action');
    }
  };

  const columns = [
    { key: 'uuid', header: 'UUID (short)', render: (uuid: string) => uuid.substring(0, 8) + '...' },
    { key: 'name', header: 'Name' },
    { key: 'node', header: 'Node', render: (_, row) => row.node?.name || 'N/A' },
    { key: 'user', header: 'Owner', render: (_, row) => row.user_id },
    { key: 'memory', header: 'RAM', render: (mem: number) => `${(mem / 1024).toFixed(1)} GB` },
    { key: 'cpu', header: 'CPU', render: (cpu: number) => `${cpu}%` },
    { key: 'status', header: 'Status', render: (status: string) => (
      <span className={`status-badge status-${status}`}>{status}</span>
    )},
    { key: 'actions', header: 'Actions', render: (_, row) => (
      <div className="table-actions">
        <button onClick={() => navigate(`/admin/servers/${row.uuid}`)}>View</button>
        <button onClick={() => handlePowerAction(row.id, row.status === 'running' ? 'stop' : 'start')}>
          {row.status === 'running' ? 'Stop' : 'Start'}
        </button>
        <button className="danger" onClick={() => {
          setSelectedServer(row);
          setShowDeleteModal(true);
        }}>Delete</button>
      </div>
    )},
  ];

  return (
    <div className="servers-page">
      <div className="page-header">
        <h1>Servers</h1>
        <button className="primary-btn" onClick={() => setShowCreateModal(true)}>
          + Create Server
        </button>
      </div>

      <DataTable
        columns={columns}
        data={servers}
        loading={loading}
        pagination={{
          page,
          limit,
          total,
          onPageChange: setPage,
        }}
      />

      {/* Create Server Modal */}
      <Modal
        isOpen={showCreateModal}
        onClose={() => setShowCreateModal(false)}
        title="Create Server"
        footer={
          <>
            <button onClick={() => setShowCreateModal(false)}>Cancel</button>
            <button className="primary" onClick={() => {
              const form = document.getElementById('create-server-form') as HTMLFormElement;
              if (form) form.requestSubmit();
            }}>Create</button>
          </>
        }
      >
        <form id="create-server-form" onSubmit={(e) => {
          e.preventDefault();
          handleCreateServer(new FormData(e.target as HTMLFormElement));
        }}>
          <div className="form-row">
            <label>Name</label>
            <input name="name" required />
          </div>
          <div className="form-row">
            <label>Owner</label>
            <select name="user_id" required>
              <option value="">Select user...</option>
              {users.map(u => (
                <option key={u.id} value={u.id}>{u.username}</option>
              ))}
            </select>
          </div>
          <div className="form-row">
            <label>Node</label>
            <select name="node_id" required>
              <option value="">Select node...</option>
              {nodes.map(n => (
                <option key={n.id} value={n.id}>{n.name}</option>
              ))}
            </select>
          </div>
          <div className="form-row">
            <label>Egg</label>
            <select name="egg_id" required>
              <option value="">Select egg...</option>
              {eggs.map(e => (
                <option key={e.id} value={e.id}>{e.name}</option>
              ))}
            </select>
          </div>
          <div className="form-row">
            <label>Memory (MB)</label>
            <input name="memory" type="number" required />
          </div>
          <div className="form-row">
            <label>Disk (MB)</label>
            <input name="disk" type="number" required />
          </div>
          <div className="form-row">
            <label>CPU (%)</label>
            <input name="cpu" type="number" required />
          </div>
        </form>
      </Modal>

      {/* Delete Confirmation Modal */}
      <Modal
        isOpen={showDeleteModal}
        onClose={() => setShowDeleteModal(false)}
        title="Delete Server"
        footer={
          <>
            <button onClick={() => setShowDeleteModal(false)}>Cancel</button>
            <button className="danger" onClick={handleDeleteServer}>Delete</button>
          </>
        }
      >
        <p>Are you sure you want to delete server "{selectedServer?.name}"? This action cannot be undone.</p>
      </Modal>
    </div>
  );
}
