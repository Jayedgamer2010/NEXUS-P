import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import DataTable from '../../components/ui/DataTable';
import Modal from '../../components/ui/Modal';
import { adminApi } from '../../api/admin';
import type { User } from '../../types';
import './Users.css';

export default function Users() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [showModal, setShowModal] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
    role: 'client' as 'admin' | 'client',
    coins: 0,
  });

  const navigate = useNavigate();

  useEffect(() => {
    loadUsers();
  }, [page]);

  const loadUsers = async () => {
    setLoading(true);
    try {
      const response = await adminApi.getUsers(page, 20);
      setUsers(response.data);
      setTotal(response.total);
    } catch (error) {
      console.error('Failed to load users:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (editingUser) {
        await adminApi.updateUser(editingUser.id, formData);
      } else {
        await adminApi.createUser(formData);
      }
      setShowModal(false);
      setEditingUser(null);
      setFormData({ username: '', email: '', password: '', role: 'client', coins: 0 });
      loadUsers();
    } catch (error) {
      console.error('Failed to save user:', error);
      alert('Failed to save user');
    }
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Delete this user?')) return;
    try {
      await adminApi.deleteUser(id);
      loadUsers();
    } catch (error) {
      console.error('Failed to delete user:', error);
      alert('Failed to delete user');
    }
  };

  const openEdit = (user: User) => {
    setEditingUser(user);
    setFormData({
      username: user.username,
      email: user.email,
      password: '',
      role: user.role,
      coins: user.coins,
    });
    setShowModal(true);
  };

  const columns = [
    { key: 'username', header: 'Username' },
    { key: 'email', header: 'Email' },
    {
      key: 'role',
      header: 'Role',
      render: (role: string) => (
        <span className={`role-badge role-${role}`}>{role.toUpperCase()}</span>
      ),
    },
    { key: 'coins', header: 'Coins' },
    { key: 'servers_count', header: 'Servers', render: () => '0' },
    {
      key: 'created_at',
      header: 'Created',
      render: (val: string) => new Date(val).toLocaleDateString(),
    },
    { key: 'actions', header: 'Actions', render: (_val: string, row: User) => (
      <div className="table-actions">
        <button onClick={() => navigate(`/admin/users/${row.id}`)}>View</button>
        <button onClick={() => openEdit(row)}>Edit</button>
        <button className="danger" onClick={() => handleDelete(row.id)}>Delete</button>
      </div>
    )},
  ];

  return (
    <div className="users-page">
      <div className="page-header">
        <h1>Users</h1>
        <button className="primary-btn" onClick={() => setShowModal(true)}>
          + Create User
        </button>
      </div>

      <DataTable
        columns={columns}
        data={users}
        loading={loading}
        pagination={{ page, limit: 20, total, onPageChange: setPage }}
      />

      <Modal
        isOpen={showModal}
        onClose={() => { setShowModal(false); setEditingUser(null); }}
        title={editingUser ? 'Edit User' : 'Create User'}
        footer={
          <>
            <button onClick={() => { setShowModal(false); setEditingUser(null); }}>Cancel</button>
            <button className="primary" onClick={handleSubmit}>{editingUser ? 'Save' : 'Create'}</button>
          </>
        }
      >
        <form onSubmit={handleSubmit}>
          <div className="form-row">
            <label>Username</label>
            <input name="username" value={formData.username} onChange={(e) => setFormData({...formData, username: e.target.value})} required />
          </div>
          <div className="form-row">
            <label>Email</label>
            <input name="email" type="email" value={formData.email} onChange={(e) => setFormData({...formData, email: e.target.value})} required />
          </div>
          <div className="form-row">
            <label>Password {editingUser && '(leave blank to keep)'}</label>
            <input name="password" type="password" value={formData.password} onChange={(e) => setFormData({...formData, password: e.target.value})} {...(!editingUser ? {required: true} : {})} />
          </div>
          <div className="form-row">
            <label>Role</label>
            <select name="role" value={formData.role} onChange={(e) => setFormData({...formData, role: e.target.value as 'admin' | 'client'})}>
              <option value="client">Client</option>
              <option value="admin">Admin</option>
            </select>
          </div>
          <div className="form-row">
            <label>Coins</label>
            <input name="coins" type="number" value={formData.coins} onChange={(e) => setFormData({...formData, coins: parseInt(e.target.value)})} />
          </div>
        </form>
      </Modal>
    </div>
  );
}
