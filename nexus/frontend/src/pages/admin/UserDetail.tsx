import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import DataTable from '../../components/ui/DataTable';
import Modal from '../../components/ui/Modal';
import ConfirmDialog from '../../components/ui/ConfirmDialog';
import { adminApi } from '../../api/admin';
import type { User, Server } from '../../types';
import { formatDate } from '../../utils/format';
import './UserDetail.css';

export default function UserDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const numericId = parseInt(id || '0', 10);

  const [user, setUser] = useState<User | null>(null);
  const [servers, setServers] = useState<Server[]>([]);
  const [loading, setLoading] = useState(true);

  const [showEditModal, setShowEditModal] = useState(false);
  const [editCoins, setEditCoins] = useState(0);
  const [editRole, setEditRole] = useState<'admin' | 'client'>('client');
  const [editPassword, setEditPassword] = useState('');

  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
  const [deleting, setDeleting] = useState(false);

  useEffect(() => {
    loadData();
  }, [id]);

  const loadData = async () => {
    if (!numericId) return;
    setLoading(true);
    try {
      const [usersRes, serversRes] = await Promise.all([
        adminApi.getUsers(1, 1000),
        adminApi.getServers(1, 1000),
      ]);

      const foundUser = usersRes.data.find((u) => u.id === numericId);
      if (!foundUser) {
        navigate('/admin/users');
        return;
      }

      setUser(foundUser);
      setEditCoins(foundUser.coins);
      setEditRole(foundUser.role);
      setServers(serversRes.data.filter((s) => s.user_id === numericId));
    } catch (error) {
      console.error('Failed to load user:', error);
      navigate('/admin/users');
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    if (!user) return;
    try {
      await adminApi.updateUser(user.id, {
        role: editRole,
        coins: editCoins,
        ...(editPassword ? { password: editPassword } : {}),
      });
      setShowEditModal(false);
      setEditPassword('');
      loadData();
    } catch (error) {
      console.error('Failed to update user:', error);
      alert('Failed to update user');
    }
  };

  const handleDelete = async () => {
    if (!user) return;
    setDeleting(true);
    try {
      await adminApi.deleteUser(user.id);
      navigate('/admin/users');
    } catch (error) {
      console.error('Failed to delete user:', error);
      alert('Failed to delete user');
    } finally {
      setDeleting(false);
      setShowDeleteConfirm(false);
    }
  };

  if (loading) return <div className="loading">Loading user...</div>;
  if (!user) return null;

  return (
    <div className="user-detail">
      {/* Header */}
      <div className="user-header">
        <div className="user-info">
          <h1>{user.username}</h1>
          <div className="user-meta">
            <span className="meta-item">{user.email}</span>
            <span className={`role-badge role-${user.role}`}>{user.role.toUpperCase()}</span>
            <span className="meta-item">{user.coins} coins</span>
          </div>
        </div>
        <div className="user-actions">
          <button className="ghost-btn" onClick={() => {
            setEditCoins(user.coins);
            setEditRole(user.role);
            setEditPassword('');
            setShowEditModal(true);
          }}>
            Edit User
          </button>
          <button className="danger-btn" onClick={() => setShowDeleteConfirm(true)}>
            Delete User
          </button>
        </div>
      </div>

      {/* Profile info */}
      <div className="profile-card">
        <h2>Profile</h2>
        <div className="profile-grid">
          <div className="profile-field">
            <span className="profile-label">UUID</span>
            <span className="profile-value"><code>{user.uuid}</code></span>
          </div>
          <div className="profile-field">
            <span className="profile-label">Created</span>
            <span className="profile-value">{formatDate(user.created_at)}</span>
          </div>
          <div className="profile-field">
            <span className="profile-label">Language</span>
            <span className="profile-value">{user.language || 'en'}</span>
          </div>
          <div className="profile-field">
            <span className="profile-label">Servers</span>
            <span className="profile-value">{servers.length}</span>
          </div>
        </div>
      </div>

      {/* Servers table */}
      <div className="servers-section">
        <h2>Servers ({servers.length})</h2>
        <DataTable
          columns={[
            {
              key: 'uuid_short',
              header: 'UUID',
              render: (val: string) => val ? val.substring(0, 8) : '-',
            },
            { key: 'name', header: 'Name' },
            {
              key: 'status',
              header: 'Status',
              render: (val: string) => (
                <span className={`status-badge status-${val}`}>{val}</span>
              ),
            },
            { key: 'memory', header: 'RAM', render: (val: number) => `${val} MB` },
            { key: 'cpu', header: 'CPU', render: (val: number) => `${val}%` },
          ]}
          data={servers}
          loading={false}
          emptyMessage="This user has no servers"
        />
      </div>

      {/* Edit Modal */}
      <Modal
        isOpen={showEditModal}
        onClose={() => { setShowEditModal(false); setEditPassword(''); }}
        title="Edit User"
        footer={
          <>
            <button onClick={() => { setShowEditModal(false); setEditPassword(''); }}>Cancel</button>
            <button className="primary" onClick={handleSave}>Save</button>
          </>
        }
      >
        <div className="form-row">
          <label>Username</label>
          <input value={user.username} disabled />
        </div>
        <div className="form-row">
          <label>Email</label>
          <input value={user.email} disabled />
        </div>
        <div className="form-row">
          <label>Role</label>
          <select value={editRole} onChange={(e) => setEditRole(e.target.value as 'admin' | 'client')}>
            <option value="client">Client</option>
            <option value="admin">Admin</option>
          </select>
        </div>
        <div className="form-row">
          <label>Coins</label>
          <input
            type="number"
            value={editCoins}
            onChange={(e) => setEditCoins(parseInt(e.target.value) || 0)}
          />
          <p className="helper-text">Current: {user.coins}</p>
        </div>
        <div className="form-row">
          <label>New Password (leave blank to keep)</label>
          <input
            type="password"
            value={editPassword}
            onChange={(e) => setEditPassword(e.target.value)}
            placeholder="Enter new password"
          />
        </div>
      </Modal>

      {/* Delete Confirmation */}
      <ConfirmDialog
        isOpen={showDeleteConfirm}
        onClose={() => setShowDeleteConfirm(false)}
        onConfirm={handleDelete}
        title="Delete User"
        message={`Delete "${user.username}"? This will also delete all ${servers.length} server(s) belonging to this user. This action cannot be undone.`}
        confirmLabel="Delete"
        loading={deleting}
      />
    </div>
  );
}
