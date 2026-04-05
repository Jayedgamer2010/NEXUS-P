import { useState } from 'react'
import { useApi } from '../../hooks/useApi'
import { nodesApi } from '../../api/admin/nodes'
import type { Node } from '../../types'
import DataTable from '../../components/ui/DataTable'
import Button from '../../components/ui/Button'
import Modal from '../../components/ui/Modal'
import Input from '../../components/ui/Input'
import ConfirmDialog from '../../components/ui/ConfirmDialog'
import Spinner from '../../components/ui/Spinner'
import Alert from '../../components/ui/Alert'
import { formatMB } from '../../utils/format'
import { useNavigate } from 'react-router-dom'

export default function Nodes() {
  const navigate = useNavigate()
  const [showModal, setShowModal] = useState(false)
  const [editTarget, setEditTarget] = useState<Node | null>(null)
  const [deleteTarget, setDeleteTarget] = useState<Node | null>(null)
  const [form, setForm] = useState({ name: '', description: '', fqdn: '', scheme: 'https', memory: 4096, disk: 51200, daemon_listen: 8080, daemon_sftp: 2022, memory_overallocate: 10, disk_overallocate: 10 })
  const [submitting, setSubmitting] = useState(false)
  const [formError, setFormError] = useState('')

  const { data, loading, refetch } = useApi<Node[]>(() => nodesApi.getAll(), [])

  const openCreate = () => {
    setForm({ name: '', description: '', fqdn: '', scheme: 'https', memory: 4096, disk: 51200, daemon_listen: 8080, daemon_sftp: 2022, memory_overallocate: 10, disk_overallocate: 10 })
    setEditTarget(null)
    setShowModal(true)
  }

  const openEdit = (node: Node) => {
    setForm({
      name: node.name,
      description: node.description,
      fqdn: node.fqdn,
      scheme: node.scheme,
      memory: node.memory,
      disk: node.disk,
      daemon_listen: node.daemon_listen,
      daemon_sftp: node.daemon_sftp,
      memory_overallocate: node.memory_overallocate,
      disk_overallocate: node.disk_overallocate,
    })
    setEditTarget(node)
    setShowModal(true)
  }

  const handleSubmit = async () => {
    setFormError('')
    setSubmitting(true)
    try {
      if (editTarget) {
        await nodesApi.update(editTarget.id, form)
      } else {
        await nodesApi.create(form)
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
    if (!deleteTarget) return
    try {
      await nodesApi.delete(deleteTarget.id)
      refetch()
    } catch {
      // ignore
    } finally {
      setDeleteTarget(null)
    }
  }

  const columns = [
    { key: 'name', header: 'Name', render: (row: Node) => (
      <span style={{ fontWeight: 500 }}>{row.name}</span>
    )},
    { key: 'fqdn', header: 'FQDN', render: (row: Node) => (
      <span className="nx-mono" style={{ color: '#9ca3af' }}>{row.scheme}://{row.fqdn}:{row.daemon_listen}</span>
    )},
    { key: 'memory', header: 'Memory', render: (row: Node) => (
      <span>{formatMB(row.memory)}</span>
    )},
    { key: 'disk', header: 'Disk', render: (row: Node) => (
      <span>{formatMB(row.disk)}</span>
    )},
    { key: 'status', header: 'Status', render: (_row: Node) => (
      <span className="nx-badge">
        <span className="nx-badge-dot" style={{ background: '#22c55e' }} />
        <span style={{ color: '#22c55e' }}>Online</span>
      </span>
    )},
    { key: 'actions', header: 'Actions', render: (row: Node) => (
      <div style={{ display: 'flex', gap: 6 }}>
        <Button variant="ghost" size="sm" onClick={() => navigate('/admin/nodes/' + row.id)}>
          View
        </Button>
        <Button variant="ghost" size="sm" onClick={() => openEdit(row)}>
          Edit
        </Button>
        <Button variant="ghost" size="sm" onClick={() => setDeleteTarget(row)}>
          <span style={{ color: '#ef4444' }}>Delete</span>
        </Button>
      </div>
    ), width: '180px' },
  ]

  const title = editTarget ? 'Edit Node' : 'Create Node'

  return (
    <div>
      <div className="nx-page-header" style={{ marginBottom: 24 }}>
        <h2>Nodes</h2>
        <Button onClick={openCreate}>Create Node</Button>
      </div>

      <div className="card">
        <DataTable columns={columns} data={data ?? []} loading={loading} />
      </div>

      <Modal isOpen={showModal} onClose={() => setShowModal(false)} title={title} size="lg">
        {formError && <Alert type="error" message={formError} dismissible={false} />}
        <div className="nx-form-row">
          <Input label="Name" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
          <Input label="FQDN" value={form.fqdn} onChange={(e) => setForm({ ...form, fqdn: e.target.value })} required placeholder="node.example.com" />
        </div>
        <Input label="Description" value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} />
        <div className="nx-form-row">
          <Input label="Memory (MB)" type="number" value={form.memory} onChange={(e) => setForm({ ...form, memory: Number(e.target.value) })} min={0} />
          <Input label="Memory Overallocation (%)" type="number" value={form.memory_overallocate} onChange={(e) => setForm({ ...form, memory_overallocate: Number(e.target.value) })} />
        </div>
        <div className="nx-form-row">
          <Input label="Disk (MB)" type="number" value={form.disk} onChange={(e) => setForm({ ...form, disk: Number(e.target.value) })} min={0} />
          <Input label="Disk Overallocation (%)" type="number" value={form.disk_overallocate} onChange={(e) => setForm({ ...form, disk_overallocate: Number(e.target.value) })} />
        </div>
        <div className="nx-form-row">
          <Input label="Daemon Port" type="number" value={form.daemon_listen} onChange={(e) => setForm({ ...form, daemon_listen: Number(e.target.value) })} />
          <Input label="SFTP Port" type="number" value={form.daemon_sftp} onChange={(e) => setForm({ ...form, daemon_sftp: Number(e.target.value) })} />
        </div>
        <div>
          <label className="nx-input-label">Scheme</label>
          <select className="nx-input" value={form.scheme} onChange={(e) => setForm({ ...form, scheme: e.target.value })}>
            <option value="https">HTTPS</option>
            <option value="http">HTTP</option>
          </select>
        </div>
        <div style={{ display: 'flex', gap: 10, marginTop: 12 }}>
          <Button variant="ghost" onClick={() => setShowModal(false)}>Cancel</Button>
          <Button loading={submitting} onClick={handleSubmit}>
            {editTarget ? 'Save Changes' : 'Create Node'}
          </Button>
        </div>
      </Modal>

      <ConfirmDialog
        isOpen={!!deleteTarget}
        onClose={() => setDeleteTarget(null)}
        onConfirm={handleDelete}
        title="Delete Node"
        message="Are you sure you want to delete this node? Associated servers may be affected."
        confirmText="Delete"
      />
    </div>
  )
}
