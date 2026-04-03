import { useEffect, useState } from 'react';
import DataTable from '../../components/ui/DataTable';
import Modal from '../../components/ui/Modal';
import { adminApi } from '../../api/admin';
import { Egg } from '../../types';
import './Eggs.css';

export default function Eggs() {
  const [eggs, setEggs] = useState<Egg[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [showModal, setShowModal] = useState(false);
  const [editingEgg, setEditingEgg] = useState<Egg | null>(null);
  const [formData, setFormData] = useState({
    name: '',
    description: '',
    docker_image: '',
    startup_command: '',
  });

  useEffect(() => {
    loadEggs();
  }, [page]);

  const loadEggs = async () => {
    setLoading(true);
    try {
      const response = await adminApi.getEggs(page, 20);
      setEggs(response.data);
      setTotal(response.total);
    } catch (error) {
      console.error('Failed to load eggs:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (editingEgg) {
        await adminApi.updateEgg(editingEgg.id, formData);
      } else {
        await adminApi.createEgg(formData);
      }
      setShowModal(false);
      setEditingEgg(null);
      setFormData({ name: '', description: '', docker_image: '', startup_command: '' });
      loadEggs();
    } catch (error) {
      console.error('Failed to save egg:', error);
      alert('Failed to save egg');
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Delete this egg?')) return;
    try {
      await adminApi.deleteEgg(id);
      loadEggs();
    } catch (error) {
      console.error('Failed to delete egg:', error);
      alert('Failed to delete egg');
    }
  };

  const openEdit = (egg: Egg) => {
    setEditingEgg(egg);
    setFormData({
      name: egg.name,
      description: egg.description || '',
      docker_image: egg.docker_image,
      startup_command: egg.startup_command,
    });
    setShowModal(true);
  };

  const columns = [
    { key: 'name', header: 'Name' },
    { key: 'description', header: 'Description', render: (desc: string) => desc?.substring(0, 50) + (desc?.length > 50 ? '...' : '') || '-' },
    { key: 'docker_image', header: 'Docker Image' },
    { key: 'actions', header: 'Actions', render: (_, row: Egg) => (
      <div className="table-actions">
        <button onClick={() => openEdit(row)}>Edit</button>
        <button className="danger" onClick={() => handleDelete(row.id)}>Delete</button>
      </div>
    )},
  ];

  return (
    <div className="eggs-page">
      <div className="page-header">
        <h1>Eggs</h1>
        <button className="primary-btn" onClick={() => setShowModal(true)}>
          + Create Egg
        </button>
      </div>

      <DataTable
        columns={columns}
        data={eggs}
        loading={loading}
        pagination={{ page, limit: 20, total, onPageChange: setPage }}
      />

      <Modal
        isOpen={showModal}
        onClose={() => { setShowModal(false); setEditingEgg(null); }}
        title={editingEgg ? 'Edit Egg' : 'Create Egg'}
        footer={
          <>
            <button onClick={() => { setShowModal(false); setEditingEgg(null); }}>Cancel</button>
            <button className="primary" onClick={handleSubmit}>{editingEgg ? 'Save' : 'Create'}</button>
          </>
        }
      >
        <form onSubmit={handleSubmit}>
          <div className="form-row">
            <label>Name</label>
            <input name="name" value={formData.name} onChange={(e) => setFormData({...formData, name: e.target.value})} required />
          </div>
          <div className="form-row">
            <label>Description</label>
            <textarea name="description" rows={3} value={formData.description} onChange={(e) => setFormData({...formData, description: e.target.value})} />
          </div>
          <div className="form-row">
            <label>Docker Image</label>
            <input name="docker_image" value={formData.docker_image} onChange={(e) => setFormData({...formData, docker_image: e.target.value})} required />
          </div>
          <div className="form-row">
            <label>Startup Command</label>
            <textarea name="startup_command" rows={4} value={formData.startup_command} onChange={(e) => setFormData({...formData, startup_command: e.target.value})} required />
          </div>
        </form>
      </Modal>
    </div>
  );
}
