import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import DataTable from '../../components/ui/DataTable';
import Modal from '../../components/ui/Modal';
import ConfirmDialog from '../../components/ui/ConfirmDialog';
import { adminApi } from '../../api/admin';
import type { Node, Allocation } from '../../types';
import { formatBytes } from '../../utils/format';
import './NodeDetail.css';

export default function NodeDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const numericId = parseInt(id || '0', 10);

  const [node, setNode] = useState<Node | null>(null);
  const [allocations, setAllocations] = useState<Allocation[]>([]);
  const [loading, setLoading] = useState(true);
  const [online, setOnline] = useState(false);

  // Create allocation modal state
  const [showCreateAllocation, setShowCreateAllocation] = useState(false);
  const [allocIp, setAllocIp] = useState('');
  const [allocPort, setAllocPort] = useState('');

  // Delete confirmation
  const [deleteTarget, setDeleteTarget] = useState<Allocation | null>(null);

  useEffect(() => {
    loadData();
  }, [id]);

  const loadData = async () => {
    if (!numericId) return;
    setLoading(true);
    try {
      const [nodeResult, allocResult] = await Promise.all([
        adminApi.getNodes(1, 100),
        adminApi.getNodeAllocations(numericId),
      ]);

      const foundNode = nodeResult.data.find((n) => n.id === numericId);
      if (!foundNode) {
        navigate('/admin/nodes');
        return;
      }

      setNode(foundNode);
      setOnline(foundNode.fqdn ? true : false);
      setAllocations(allocResult.map((a: any) => ({
        id: a.id,
        node_id: a.node_id,
        ip: a.ip,
        ip_alias: null,
        port: a.port,
        server_id: a.server_id,
        assigned: a.server_id !== null,
        server_name: a.server?.name,
        notes: a.notes || '',
        created_at: a.created_at,
      })));
    } catch (error) {
      console.error('Failed to load node:', error);
      navigate('/admin/nodes');
    } finally {
      setLoading(false);
    }
  };

  // Delete node
  const handleDeleteNode = async () => {
    if (!node) return;
    try {
      await adminApi.deleteNode(node.id);
      navigate('/admin/nodes');
    } catch (error) {
      console.error('Failed to delete node:', error);
      alert('Failed to delete node');
    }
  };

  // Create allocation
  const handleCreateAllocation = async () => {
    if (!node || !allocPort) return;
    try {
      await adminApi.createAllocation(node.id, {
        ip: allocIp || node.fqdn,
        port: parseInt(allocPort, 10),
      });
      setShowCreateAllocation(false);
      setAllocIp('');
      setAllocPort('');
      loadData();
    } catch (error) {
      console.error('Failed to create allocation:', error);
      alert('Failed to create allocation');
    }
  };

  // Delete allocation
  const handleDeleteAllocation = async () => {
    if (!deleteTarget) return;
    try {
      await adminApi.deleteAllocation(deleteTarget.id);
      setDeleteTarget(null);
      loadData();
    } catch (error) {
      console.error('Failed to delete allocation:', error);
      alert('Failed to delete allocation');
    }
  };

  if (loading) return <div className="loading">Loading node...</div>;
  if (!node) return null;

  const memUsed = node.used_memory || 0;
  const memTotal = node.memory || 0;
  const diskUsed = node.used_disk || 0;
  const diskTotal = node.disk || 0;
  const memPercent = memTotal > 0 ? Math.round((memUsed / memTotal) * 100) : 0;
  const diskPercent = diskTotal > 0 ? Math.round((diskUsed / diskTotal) * 100) : 0;

  return (
    <div className="node-detail">
      <div className="node-header">
        <div className="node-info">
          <h1>{node.name}</h1>
          <div className="node-meta">
            <code>{node.fqdn}</code>
            <span className={`status-badge ${online ? 'online' : 'offline'}`}>
              {online ? 'Online' : 'Offline'}
            </span>
          </div>
        </div>
        <button className="danger-btn" onClick={handleDeleteNode}>
          Delete Node
        </button>
      </div>

      {/* Resource bars */}
      <div className="resource-bars">
        <div className="resource-bar">
          <div className="resource-label">
            Memory: {formatBytes(memUsed * 1024 * 1024)} / {formatBytes(memTotal * 1024 * 1024)} ({memPercent}%)
          </div>
          <div className="bar-track">
            <div className="bar-fill bar-fill--purple" style={{ width: `${Math.min(memPercent, 100)}%` }} />
          </div>
        </div>
        <div className="resource-bar">
          <div className="resource-label">
            Disk: {formatBytes(diskUsed * 1024 * 1024)} / {formatBytes(diskTotal * 1024 * 1024)} ({diskPercent}%)
          </div>
          <div className="bar-track">
            <div className="bar-fill bar-fill--blue" style={{ width: `${Math.min(diskPercent, 100)}%` }} />
          </div>
        </div>
      </div>

      {/* Allocations */}
      <div className="allocations-section">
        <div className="section-header">
          <h2>Allocations ({allocations.length})</h2>
          <button className="primary-btn" onClick={() => setShowCreateAllocation(true)}>
            + Add Allocation
          </button>
        </div>

        <DataTable
          columns={[
            {
              key: 'ip',
              header: 'IP Address',
              render: (val: string, row) => row.assigned ? (
                <span className="alloc-assigned">{val}:{row.port}</span>
              ) : (
                <span className="alloc-unassigned">{val}:{row.port}</span>
              ),
            },
            {
              key: 'server',
              header: 'Assigned To',
              render: (val: string, row: Allocation) => row.server_name ? (
                row.server_name
              ) : (
                <span className="text-muted">Unassigned</span>
              ),
            },
            { key: 'actions', header: 'Actions', render: (val: string, row: Allocation) => (
              row.assigned ? (
                <span className="text-muted">&mdash;</span>
              ) : (
                <button className="danger-btn-sm" onClick={() => setDeleteTarget(row)}>
                  Delete
                </button>
              )
            )},
          ]}
          data={allocations}
          loading={false}
          emptyMessage="No allocations on this node"
        />
      </div>

      {/* Create Allocation Modal */}
      <Modal
        isOpen={showCreateAllocation}
        onClose={() => { setShowCreateAllocation(false); setAllocIp(''); setAllocPort(''); }}
        title="Add Allocation"
        footer={
          <>
            <button onClick={() => { setShowCreateAllocation(false); setAllocIp(''); setAllocPort(''); }}>Cancel</button>
            <button className="primary" onClick={handleCreateAllocation}>Create</button>
          </>
        }
      >
        <div className="form-row">
          <label>IP Address (optional, defaults to node FQDN)</label>
          <input value={allocIp} onChange={(e) => setAllocIp(e.target.value)} placeholder={node.fqdn} />
        </div>
        <div className="form-row">
          <label>Port</label>
          <input type="number" value={allocPort} onChange={(e) => setAllocPort(e.target.value)} required placeholder="25565" />
        </div>
      </Modal>

      {/* Delete Allocation Confirm */}
      <ConfirmDialog
        isOpen={!!deleteTarget}
        onClose={() => setDeleteTarget(null)}
        onConfirm={handleDeleteAllocation}
        title="Delete Allocation"
        message={deleteTarget ? `Delete ${deleteTarget.ip}:${deleteTarget.port}? This action cannot be undone.` : ''}
        confirmLabel="Delete"
      />
    </div>
  );
}
