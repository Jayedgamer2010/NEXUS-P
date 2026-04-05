import { useState } from 'react'
import { useApi } from '../../hooks/useApi'
import { usersApi } from '../../api/admin/users'
import type { User } from '../../types'
import DataTable from '../../components/ui/DataTable'
import Button from '../../components/ui/Button'
import Modal from '../../components/ui/Modal'
import Input from '../../components/ui/Input'
import ConfirmDialog from '../../components/ui/ConfirmDialog'
import Alert from '../../components/ui/Alert'
import { formatDate, formatRelativeTime } from '../../utils/format'
import Spinner from '../../components/ui/Spinner'

export default function Users() {
  const { data, loading, refetch } = useApi<{ data: User[]; meta?: any }>(
    () => usersApi.getAll(1),
    []
  )
  const [showModal, setShowModal] = useState(false)
  const [editTarget, setEditTarget] = useState<User | null>(null)
  const [deleteTarget, setDeleteTarget] = useState<User | null>(null)
  const [form, setForm] = useState({ username: '', email: '', password: '', role: 'client', coins: '' })
  const [formError, setFormError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const openCreate = () => {
    setForm({ username: '', email: '', password: '', role: 'client', coins: '' })
    setEditTarget(null)
    setShowModal(true)
  }

  const openEdit = (user: User) => {
    setForm({
      username: user.username,
      email: user.email,
      password: '',
      role: user.root_admin ? 'admin' : user.role,
      coins: '',
    })
    setEditTarget(user)
    setShowModal(true)
  }

  const handleSubmit = async () => {
    setFormError('')
    setSubmitting(true)
    try {
      if (editTarget) {
        const data: any = {
          username: form.username,
          email: form.email,
        }
        if (form.password) data.password = form.password
        if (form.coins) data.coins = Number(form.coins)
        await usersApi.update(editTarget.id, data)
      } else {
        if (!form.password) {
          setFormError('Password is required for new users')
          setSubmitting(false)
          return
        }
        await usersApi.create(form)
      }
      setShowModal(false)
      setEditTarget(null)
      refetch()
    } catch (err: any) {
      setFormError(err.response?.data?.message || 'Operation failed')
    } finally {
      setSubmitting(false)
    }
  }

  const handleDelete = async () => {
    try {
      if (deleteTarget) await usersApi.delete(deleteTarget.id)
      refetch()
    } catch {
      // may fail if user has servers
    } finally {
      setDeleteTarget(null)
    }
  }

  const columns = [
    { key: 'username', header: 'Username', render: (row: User) => (
      <span style={{ fontWeight: 500 }}>{row.username}</span>
    )},
    { key: 'email', header: 'Email', render: (row: User) => (
      <span style={{ color: '#9ca3af' }}>{row.email}</span>
    )},
    { key: 'role', header: 'Role', render: (row: User) => (
      <span className={`role-badge role-badge--${row.root_admin ? 'admin' : 'client'}`}>
        {row.root_admin ? 'Admin' : row.role}
      </span>
    )},
    { key: 'coins', header: 'Coins', render: (row: User) => (
      <span>{row.coins ?? 0}</span>
    )},
    { key: 'created_at', header: 'Created', render: (row: User) => (
      <span>{formatRelativeTime(row.created_at)}</span>
    )},
    { key: 'actions', header: 'Actions', render: (row: User) => (
      <div style={{ display: 'flex', gap: 6 }}>
        <Button variant="ghost" size="sm" onClick={() => openEdit(row)}>
          Edit
        </Button>
        <Button variant="ghost" size="sm" onClick={() => setDeleteTarget(row)}>
          <span style={{ color: '#ef4444' }}>Delete</span>
        </Button>
      </div>
    ), width: '130px' },
  ]

  const title = editTarget ? 'Edit User' : 'Create User'

  return (
    <div>
      <div className="nx-page-header" style={{ marginBottom: 24 }}>
        <h2>Users</h2>
        <Button onClick={openCreate}>Create User</Button>
      </div>

      <div className="card">
        <DataTable columns={columns} data={data?.data ?? []} loading={loading} />
      </div>

      <Modal isOpen={showModal} onClose={() => setShowModal(false)} title={title}>
        {formError && <Alert type="error" message={formError} dismissible={false} />}
        <Input
          label="Username"
          value={form.username}
          onChange={(e) => setForm({ ...form, username: e.target.value })}
          required
        />
        <Input
          label="Email"
          type="email"
          value={form.email}
          onChange={(e) => setForm({ ...form, email: e.target.value })}
          required
        />
        <Input
          label={editTarget ? 'New Password (leave empty to keep)' : 'Password'}
          type="password"
          value={form.password}
          onChange={(e) => setForm({ ...form, password: e.target.value })}
          required={!editTarget}
        />
        <Input
          label="Adjust Coins"
          type="number"
          value={form.coins}
          onChange={(e) => setForm({ ...form, coins: e.target.value })}
          helper={editTarget ? "Enter negative to deduct" : "Leave empty for 0"}
        />
        <div className="nx-form-group">
          <label className="nx-input-label">Role</label>
          <select className="nx-input" value={form.role} onChange={(e) => setForm({ ...form, role: e.target.value })}>
            <option value="client">Client</option>
            <option value="admin">Admin</option>
          </select>
        </div>
        <div style={{ display: 'flex', gap: 10, marginTop: 12 }}>
          <Button variant="ghost" onClick={() => setShowModal(false)}>Cancel</Button>
          <Button loading={submitting} onClick={handleSubmit}>
            {editTarget ? 'Save Changes' : 'Create User'}
          </Button>
        </div>
      </Modal>

      <ConfirmDialog
        isOpen={!!deleteTarget}
        onClose={() => setDeleteTarget(null)}
        onConfirm={handleDelete}
        title="Delete User"
        message={`Are you sure you want to delete ${deleteTarget?.username ?? 'this user'}?`}
        confirmText="Delete"
      />
    </div>
  )
}
