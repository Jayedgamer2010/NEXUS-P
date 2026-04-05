import { useState } from 'react'
import { useParams } from 'react-router-dom'
import { useApi } from '../../hooks/useApi'
import { nodesApi } from '../../api/admin/nodes'
import type { Node, Allocation } from '../../types'
import DataTable from '../../components/ui/DataTable'
import Button from '../../components/ui/Button'
import Modal from '../../components/ui/Modal'
import Input from '../../components/ui/Input'
import Spinner from '../../components/ui/Spinner'
import ConfirmDialog from '../../components/ui/ConfirmDialog'
import { formatMB, formatDate } from '../../utils/format'
import StatusBadge from '../../components/ui/StatusBadge'

interface AllocationExt extends Allocation {
  server_name?: string
}

export default function NodeDetail() {
  const { id } = useParams<{ id: string }>()
  const nodeId = Number(id)
  const [showAllocationModal, setShowAllocationModal] = useState(false)
  const [allocForm, setAllocForm] = useState({ ip: '', port: '' })
  const [deleteAllocId, setDeleteAllocId] = useState<number | null>(null)
  const [submitting, setSubmitting] = useState(false)

  const { data: node, loading: nodeLoading } = useApi<Node>(
    () => nodesApi.getOne(nodeId),
    [nodeId]
  )

  const { data: allocations, loading: allocLoading, refetch: refetchAllocs } = useApi<AllocationExt[]>(
    () => nodesApi.getAllocations(nodeId),
    [nodeId]
  )

  const handleAddAllocation = async () => {
    if (!allocForm.port) return
    setSubmitting(true)
    try {
      await nodesApi.addAllocation(nodeId, {
        ip: allocForm.ip || '0.0.0.0',
        port: Number(allocForm.port),
      })
      setShowAllocationModal(false)
      setAllocForm({ ip: '', port: '' })
      refetchAllocs()
    } catch {
    } finally {
      setSubmitting(false)
    }
  }

  const handleDeleteAllocation = async () => {
    if (!deleteAllocId) return
    try {
      await nodesApi.deleteAllocation(deleteAllocId)
      refetchAllocs()
    } catch {
    } finally {
      setDeleteAllocId(null)
    }
  }

  const allocColumns = [
    { key: 'ip', header: 'IP', render: (row: AllocationExt) => (
      <span className="nx-mono">{row.ip || '0.0.0.0'}</span>
    )},
    { key: 'port', header: 'Port', render: (row: AllocationExt) => (
      <span className="nx-mono">{row.port}</span>
    )},
    { key: 'alias', header: 'Alias', render: (row: AllocationExt) => (
      <span>{row.ip_alias || '-'}</span>
    )},
    { key: 'assigned', header: 'Assigned To', render: (row: AllocationExt) => (
      <span style={{ color: row.server_id ? '#ffffff' : '#6b7280' }}>
        {row.server_id ? (row.server_name || `Server #${row.server_id}`) : 'Unassigned'}
      </span>
    )},
    { key: 'actions', header: 'Actions', render: (row: AllocationExt) => (
      row.server_id ? null : (
        <Button variant="ghost" size="sm" onClick={() => setDeleteAllocId(row.id)}>
          <span style={{ color: '#ef4444' }}>Delete</span>
        </Button>
      )
    ), width: '90px' },
  ]

  if (nodeLoading || !node) {
    return <div style={{ padding: 40, textAlign: 'center' }}><Spinner size="lg" /></div>
  }

  return (
    <div>
      {/* Node Info */}
      <div style={{ marginBottom: 24 }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginBottom: 16 }}>
          <h2 style={{ fontSize: 22, fontWeight: 700 }}>{node.name}</h2>
          <StatusBadge status={node.maintenance_mode ? 'suspended' : 'running'} />
        </div>
        <div className="nx-grid-4">
          <div className="nx-stat-card">
            <div style={{ fontSize: 12, color: '#6b7280', marginBottom: 4 }}>FQDN</div>
            <div style={{ fontSize: 14, fontWeight: 500 }} className="nx-mono">{node.fqdn}</div>
          </div>
          <div className="nx-stat-card">
            <div style={{ fontSize: 12, color: '#6b7280', marginBottom: 4 }}>Memory</div>
            <div style={{ fontSize: 20, fontWeight: 600 }}>{formatMB(node.memory)}</div>
          </div>
          <div className="nx-stat-card">
            <div style={{ fontSize: 12, color: '#6b7280', marginBottom: 4 }}>Disk</div>
            <div style={{ fontSize: 20, fontWeight: 600 }}>{formatMB(node.disk)}</div>
          </div>
          <div className="nx-stat-card">
            <div style={{ fontSize: 12, color: '#6b7280', marginBottom: 4 }}>Ports</div>
            <div style={{ fontSize: 14, fontWeight: 500 }}>
              Daemon: {node.daemon_listen} / SFTP: {node.daemon_sftp}
            </div>
          </div>
        </div>
      </div>

      {/* Allocations */}
      <div className="card">
        <div className="nx-section-title">
          <span>Allocations</span>
          <Button variant="ghost" size="sm" onClick={() => setShowAllocationModal(true)}>
            Add Allocation
          </Button>
        </div>
        <DataTable columns={allocColumns} data={allocations ?? []} loading={allocLoading} />
      </div>

      {/* Add Allocation Modal */}
      <Modal isOpen={showAllocationModal} onClose={() => setShowAllocationModal(false)} title="Add Allocation" size="sm">
        <Input
          label="IP Address"
          value={allocForm.ip}
          onChange={(e) => setAllocForm({ ...allocForm, ip: e.target.value })}
          placeholder="Leave empty for 0.0.0.0"
        />
        <Input
          label="Port"
          type="number"
          value={allocForm.port}
          onChange={(e) => setAllocForm({ ...allocForm, port: e.target.value })}
          placeholder="25565"
        />
        <div style={{ display: 'flex', gap: 10, marginTop: 12 }}>
          <Button variant="ghost" onClick={() => setShowAllocationModal(false)}>Cancel</Button>
          <Button loading={submitting} onClick={handleAddAllocation}>Add</Button>
        </div>
      </Modal>

      <ConfirmDialog
        isOpen={!!deleteAllocId}
        onClose={() => setDeleteAllocId(null)}
        onConfirm={handleDeleteAllocation}
        title="Delete Allocation"
        message="Are you sure you want to delete this allocation?"
        confirmText="Delete"
      />
    </div>
  )
}
