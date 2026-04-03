import { useEffect, useState } from 'react';
import DataTable from '../../components/ui/DataTable';
import Modal from '../../components/ui/Modal';
import { adminApi } from '../../api/admin';
import { Node } from '../../types';
import './Nodes.css';

export default function Nodes() {
  const [nodes, setNodes] = useState<Node[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [showModal, setShowModal] = useState(false);
  const [formData, setFormData] = useState({
    name: '',
    fqdn: '',
    scheme: 'https',
    wings_port: 8080,
    memory: 0,
    memory_overalloc: 0,
    disk: 0,
    disk_overalloc: 0,
    token_id: '',
    token: '',
  });

  useEffect(() => {
    loadNodes();
  }, [page]);

  const loadNodes = async () => {
    setLoading(true);
    try {
      const response = await adminApi.getNodes(page, 20);
      setNodes(response.data);
      setTotal(response.total);
    } catch (error) {
      console.error('Failed to load nodes:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await adminApi.createNode(formData);
      setShowModal(false);
      setFormData({
        name: '',
        fqdn: '',
        scheme: 'https',
        wings_port: 8080,
        memory: 0,
        memory_overalloc: 0,
        disk: 0,
        disk_overalloc: 0,
        token_id: '',
        token: '',
      });
      loadNodes();
    } catch (error) {
      console.error('Failed to create node:', error);
      alert('Failed to create node');
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Delete this node?')) return;
    try {
      await adminApi.deleteNode(id);
      loadNodes();
    } catch (error) {
      console.error('Failed to delete node:', error);
      alert('Failed to delete node');
    }
  };

  const columns = [
    { key: 'name', header: 'Name' },
    { key: 'fqdn', header: 'FQDN' },
    { key: 'memory', header: 'Memory', render: (mem: number) => `${(mem / 1024).toFixed(1)} GB` },
    { key: 'disk', header: 'Disk', render: (disk: number) => `${(disk / 1024).toFixed(1)} GB` },
    { key: 'servers_count', header: 'Servers' }, // would come from stats API
    { key: 'status', header: 'Status', render: () => <span className="status-online">Online</span> },
    { key: 'actions', header: 'Actions', render: (_, row: Node) => (
      <button className="danger" onClick={() => handleDelete(row.id)}>Delete</button>
    )},
  ];

  return (
    <div className="nodes-page">
      <div className="page-header">
        <h1>Nodes</h1>
        <button className="primary-btn" onClick={() => setShowModal(true)}>
          + Create Node
        </button>
      </div>

      <DataTable
        columns={columns}
        data={nodes}
        loading={loading}
        pagination={{ page, limit: 20, total, onPageChange: setPage }}
      />

      <Modal
        isOpen={showModal}
        onClose={() => setShowModal(false)}
        title="Create Node"
        footer={
          <>
            <button onClick={() => setShowModal(false)}>Cancel</button>
            <button className="primary" onClick={handleSubmit}>Create</button>
          </>
        }
      >
        <form onSubmit={handleSubmit}>
          <div className="form-row">
            <label>Name</label>
            <input name="name" value={formData.name} onChange={(e) => setFormData({...formData, name: e.target.value})} required />
          </div>
          <div className="form-row">
            <label>FQDN</label>
            <input name="fqdn" value={formData.fqdn} onChange={(e) => setFormData({...formData, fqdn: e.target.value})} required />
          </div>
          <div className="form-row">
            <label>Scheme</label>
            <select name="scheme" value={formData.scheme} onChange={(e) => setFormData({...formData, scheme: e.target.value as 'http' | 'https'})}>
              <option value="https">HTTPS</option>
              <option value="http">HTTP</option>
            </select>
          </div>
          <div className="form-row">
            <label>Wings Port</label>
            <input name="wings_port" type="number" value={formData.wings_port} onChange={(e) => setFormData({...formData, wings_port: parseInt(e.target.value)})} required />
          </div>
          <div className="form-row">
            <label>Memory (MB)</label>
            <input name="memory" type="number" value={formData.memory} onChange={(e) => setFormData({...formData, memory: parseInt(e.target.value)})} required />
          </div>
          <div className="form-row">
            <label>Disk (MB)</label>
            <input name="disk" type="number" value={formData.disk} onChange={(e) => setFormData({...formData, disk: parseInt(e.target.value)})} required />
          </div>
          <div className="form-row">
            <label>Token ID</label>
            <input name="token_id" value={formData.token_id} onChange={(e) => setFormData({...formData, token_id: e.target.value})} required />
          </div>
          <div className="form-row">
            <label>Token</label>
            <input name="token" type="password" value={formData.token} onChange={(e) => setFormData({...formData, token: e.target.value})} required />
          </div>
        </form>
      </Modal>
    </div>
  );
}
