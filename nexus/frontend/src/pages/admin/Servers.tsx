import { useState } from 'react'
import { useApi } from '../../hooks/useApi'
import { serversApi } from '../../api/admin/servers'
import { usersApi } from '../../api/admin/users'
import { nodesApi } from '../../api/admin/nodes'
import { eggsApi } from '../../api/admin/eggs'
import type { Server, User, Node, Egg } from '../../types'
import DataTable from '../../components/ui/DataTable'
import Button from '../../components/ui/Button'
import Modal from '../../components/ui/Modal'
import Input from '../../components/ui/Input'
import StatusBadge from '../../components/ui/StatusBadge'
import ConfirmDialog from '../../components/ui/ConfirmDialog'
import { formatMB } from '../../utils/format'
import { useNavigate } from 'react-router-dom'

export default function Servers() {
  const navigate = useNavigate()
  const [page, setPage] = useState(1)
  const [showModal, setShowModal] = useState(false)
  const [deleteTarget, setDeleteTarget] = useState<Server | null>(null)
  const [form, setForm] = useState({
    name: '', user_id: '', node_id: '', egg_id: '',
    memory: 1024, disk: 1024, cpu: 0,
    startup: '', image: '',
  })
  const [formErrors, setFormErrors] = useState<Record<string, string>>({})
  const [submitting, setSubmitting] = useState(false)

  const { data, loading, refetch } = useApi<{ data: Server[]; meta: any }>(
    () => serversApi.getAll(page),
    [page]
  )

  const { data: usersData } = useApi<User[]>(() => usersApi.getAll(1), [])
  const { data: nodesData } = useApi<Node[]>(() => nodesApi.getAll(), [])
  const { data: eggsData } = useApi<Egg[]>(() => eggsApi.getAll(), [])

  const handleCreate = async () => {
    setFormErrors({})
    if (!form.name || !form.user_id || !form.node_id || !form.egg_id) {
      setFormErrors({ general: 'All fields are required' })
      return
    }
    setSubmitting(true)
    try {
      await serversApi.create({
        name: form.name,
        user_id: Number(form.user_id),
        node_id: Number(form.node_id),
        egg_id: Number(form.egg_id),
        memory: Number(form.memory),
        disk: Number(form.disk),
        cpu: Number(form.cpu),
        startup: form.startup || undefined,
        image: form.image || undefined,
      })
      setShowModal(false)
      setForm({ name: '', user_id: '', node_id: '', egg_id: '', memory: 1024, disk: 1024, cpu: 0, startup: '', image: '' })
      refetch()
    } catch (err: any) {
      setFormErrors({ general: err.response?.data?.message || 'Failed to create server' })
    } finally {
      setSubmitting(false)
    }
  }

  const handleDelete = async () => {
    if (!deleteTarget) return
    try {
      await serversApi.delete(deleteTarget.id)
      refetch()
    } catch { /* ignore */ }
  }

  const handleSuspend = async (server: Server) => {
    try {
      if (server.suspended) {
        await serversApi.unsuspend(server.id)
      } else {
        await serversApi.suspend(server.id)
      }
      refetch()
    } catch { /* ignore */ }
  }

  const handleEggSelect = (eggId: string) => {
    const egg = eggsData?.find((e) => String(e.id) === eggId)
    if (egg) {
      setForm((f) => ({ ...f, ...f, startup: egg.startup, image: egg.docker_image }))
    }
  }

  const columns = [
    { key: 'uuid_short', header: 'UUID', render: (row: Server) => (
      <span className="nx-mono" style={{ color: '#4b5563' }}>{row.uuid_short}</span>
    )},
    { key: 'name', header: 'Name', render: (row: Server) => (
      <span style={{ fontWeight: 500 }}>{row.name}</span>
    )},
    { key: 'node', header: 'Node', render: (row: Server) => (
      <span>{row.node?.name}</span>
    )},
    { key: 'user', header: 'Owner', render: (row: Server) => (
      <span>{row.user?.username}</span>
    )},
    { key: 'memory', header: 'RAM', render: (row: Server) => (
      <span>{formatMB(row.memory)}</span>
    )},
    { key: 'cpu', header: 'CPU', render: (row: Server) => (
      <span>{row.cpu}%</span>
    )},
    { key: 'status', header: 'Status', render: (row: Server) => (
      <StatusBadge status={row.suspended ? 'suspended' : row.status} />
    )},
    { key: 'actions', header: 'Actions', render: (row: Server) => (
      <div style={{ display: 'flex', gap: 6 }}>
        <Button variant="ghost" size="sm" onClick={() => navigate('/admin/servers/' + row.id)}>
          View
        </Button>
        <Button variant="ghost" size="sm" onClick={() => handleSuspend(row)}>
          {row.suspended ? 'Unsuspend' : 'Suspend'}
        </Button>
        <Button variant="ghost" size="sm" onClick={() => setDeleteTarget(row)}>
          <span style={{ color: '#ef4444' }}>Delete</span>
        </Button>
      </div>
    ), width: '250px' },
  ]

  return (
    <div>
      <div className="nx-page-header">
        <h2>Servers</h2>
        <Button onClick={() => setShowModal(true)}>Create Server</Button>
      </div>

      <div className="card">
        <DataTable
          columns={columns}
          data={data?.data ?? []}
          loading={loading}
          pagination={data?.meta ? {
            current: data.meta.current_page,
            total: data.meta.total,
            perPage: data.meta.per_page,
            onPageChange: setPage,
          } : undefined}
        />
      </div>

      <Modal isOpen={showModal} onClose={() => setShowModal(false)} title="Create Server" size="lg">
        {formErrors.general && (
          <div style={{ marginBottom: 16, color: '#ef4444', fontSize: 13 }}>{formErrors.general}</div>
        )}
        <Input
          label="Server Name"
          value={form.name}
          onChange={(e) => setForm({ ...form, name: e.target.value })}
          required
        />
        <div className="nx-form-row">
          <div>
            <label className="nx-input-label">Owner</label>
            <select
              className="nx-input"
              value={form.user_id}
              onChange={(e) => setForm({ ...form, user_id: e.target.value })}
            >
              <option value="">Select owner...</option>
              {usersData?.map((u) => (
                <option key={u.id} value={u.id}>{u.username}</option>
              ))}
            </select>
          </div>
          <div>
            <label className="nx-input-label">Node</label>
            <select
              className="nx-input"
              value={form.node_id}
              onChange={(e) => setForm({ ...form, node_id: e.target.value })}
            >
              <option value="">Select node...</option>
              {nodesData?.map((n) => (
                <option key={n.id} value={n.id}>{n.name}</option>
              ))}
            </select>
          </div>
        </div>
        <div>
          <label className="nx-input-label">Egg</label>
          <select
            className="nx-input"
            value={form.egg_id}
            onChange={(e) => { handleEggSelect(e.target.value); setForm({ ...form, egg_id: e.target.value }) }}
          >
            <option value="">Select egg...</option>
            {eggsData?.map((egg) => (
              <option key={egg.id} value={egg.id}>{egg.name}</option>
            ))}
          </select>
        </div>
        <div className="nx-form-row">
          <Input
            label="Memory (MB)"
            type="number"
            value={form.memory}
            onChange={(e) => setForm({ ...form, memory: Number(e.target.value) })}
            min={128}
          />
          <Input
            label="Disk (MB)"
            type="number"
            value={form.disk}
            onChange={(e) => setForm({ ...form, disk: Number(e.target.value) })}
            min={256}
          />
        </div>
        <Input
          label="CPU %"
          type="number"
          value={form.cpu}
          onChange={(e) => setForm({ ...form, cpu: Number(e.target.value) })}
          min={0}
          max={10000}
          helper="Set to 0 for unlimited"
        />
        <Input
          label="Docker Image"
          value={form.image}
          onChange={(e) => setForm({ ...form, image: e.target.value })}
        />
        <Input
          label="Startup Command"
          value={form.startup}
          onChange={(e) => setForm({ ...form, startup: e.target.value })}
        />
        <div style={{ display: 'flex', gap: 10, marginTop: 12 }}>
          <Button variant="ghost" onClick={() => setShowModal(false)}>Cancel</Button>
          <Button loading={submitting} onClick={handleCreate}>Create Server</Button>
        </div>
      </Modal>

      <ConfirmDialog
        isOpen={!!deleteTarget}
        onClose={() => setDeleteTarget(null)}
        onConfirm={handleDelete}
        title="Delete Server"
        message="Are you sure you want to delete this server? This action cannot be undone."
        confirmText="Delete"
      />
    </div>
  )
}
