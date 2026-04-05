import { useState } from 'react'
import { useApi } from '../../hooks/useApi'
import { eggsApi } from '../../api/admin/eggs'
import type { Egg } from '../../types'
import DataTable from '../../components/ui/DataTable'
import Button from '../../components/ui/Button'
import Modal from '../../components/ui/Modal'
import Input from '../../components/ui/Input'
import ConfirmDialog from '../../components/ui/ConfirmDialog'
import Alert from '../../components/ui/Alert'
import Spinner from '../../components/ui/Spinner'

export default function Eggs() {
  const { data, loading, refetch } = useApi<Egg[]>(() => eggsApi.getAll(), [])
  const [showModal, setShowModal] = useState(false)
  const [editTarget, setEditTarget] = useState<Egg | null>(null)
  const [deleteTarget, setDeleteTarget] = useState<Egg | null>(null)
  const [form, setForm] = useState({ name: '', author: 'NEXUS', description: '', docker_image: '', startup: '', config_stop: 'stop' })
  const [formError, setFormError] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const openCreate = () => {
    setForm({ name: '', author: 'NEXUS', description: '', docker_image: '', startup: '', config_stop: 'stop' })
    setEditTarget(null)
    setShowModal(true)
  }

  const openEdit = (egg: Egg) => {
    setForm({
      name: egg.name,
      author: egg.author,
      description: egg.description,
      docker_image: egg.docker_image,
      startup: egg.startup,
      config_stop: egg.config_stop || 'stop',
    })
    setEditTarget(egg)
    setShowModal(true)
  }

  const handleSubmit = async () => {
    setFormError('')
    if (!form.name || !form.docker_image || !form.startup) {
      setFormError('Name, Docker Image, and Startup Command are required')
      return
    }
    setSubmitting(true)
    try {
      if (editTarget) {
        await eggsApi.update(editTarget.id, form)
      } else {
        await eggsApi.create(form)
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
      if (deleteTarget) await eggsApi.delete(deleteTarget.id)
      refetch()
    } catch {
      // may fail if servers are using the egg
    } finally {
      setDeleteTarget(null)
    }
  }

  const columns = [
    { key: 'name', header: 'Name', render: (row: Egg) => (
      <span style={{ fontWeight: 500 }}>{row.name}</span>
    )},
    { key: 'author', header: 'Author', render: (row: Egg) => (
      <span>{row.author}</span>
    )},
    { key: 'docker_image', header: 'Docker Image', render: (row: Egg) => (
      <span className="nx-mono" style={{ color: '#9ca3af', maxWidth: 250, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap', display: 'block' }}>
        {row.docker_image}
      </span>
    )},
    { key: 'actions', header: 'Actions', render: (row: Egg) => (
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

  return (
    <div>
      <div className="nx-page-header" style={{ marginBottom: 24 }}>
        <h2>Eggs</h2>
        <Button onClick={openCreate}>Create Egg</Button>
      </div>

      <div className="card">
        <DataTable columns={columns} data={data ?? []} loading={loading} />
      </div>

      <Modal isOpen={showModal} onClose={() => setShowModal(false)} title={editTarget ? 'Edit Egg' : 'Create Egg'} size="lg">
        {formError && <Alert type="error" message={formError} dismissible={false} />}
        <div className="nx-form-row">
          <Input label="Name" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
          <Input label="Author" value={form.author} onChange={(e) => setForm({ ...form, author: e.target.value })} />
        </div>
        <Input
          label="Description"
          value={form.description}
          onChange={(e) => setForm({ ...form, description: e.target.value })}
        />
        <Input
          label="Docker Image"
          value={form.docker_image}
          onChange={(e) => setForm({ ...form, docker_image: e.target.value })}
          required
          helper="e.g. ghcr.io/pterodactyl/yolks:java_17"
        />
        <Input
          label="Startup Command"
          value={form.startup}
          onChange={(e) => setForm({ ...form, startup: e.target.value })}
          required
          helper="e.g. java -Xms128M -XX:MaxRAMPercentage=95.0 -jar {{SERVER_JARFILE}}"
        />
        <Input
          label="Stop Command"
          value={form.config_stop}
          onChange={(e) => setForm({ ...form, config_stop: e.target.value })}
          helper="Command used to gracefully stop the server"
        />
        <div style={{ display: 'flex', gap: 10, marginTop: 12 }}>
          <Button variant="ghost" onClick={() => setShowModal(false)}>Cancel</Button>
          <Button loading={submitting} onClick={handleSubmit}>
            {editTarget ? 'Save Changes' : 'Create Egg'}
          </Button>
        </div>
      </Modal>

      <ConfirmDialog
        isOpen={!!deleteTarget}
        onClose={() => setDeleteTarget(null)}
        onConfirm={handleDelete}
        title="Delete Egg"
        message={`Are you sure you want to delete "${deleteTarget?.name ?? 'this egg'}"? Servers using this egg may be affected.`}
        confirmText="Delete"
      />
    </div>
  )
}
