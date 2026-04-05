import Spinner from './Spinner'

interface Column<T> {
  key: string
  header: string
  render?: (row: T) => React.ReactNode
  width?: string
}

interface DataTableProps<T> {
  columns: Column<T>[]
  data: T[]
  loading?: boolean
  emptyMessage?: string
  pagination?: {
    current: number
    total: number
    perPage: number
    onPageChange: (page: number) => void
  }
}

export default function DataTable<T extends { id: number }>({
  columns,
  data,
  loading,
  emptyMessage = 'No records found',
  pagination,
}: DataTableProps<T>) {
  if (loading) {
    return (
      <table className="nx-table">
        <thead>
          <tr>
            {columns.map((col) => (
              <th key={col.key}>{col.header}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {Array.from({ length: 5 }).map((_, i) => (
            <tr key={i}>
              {columns.map((col) => (
                <td key={col.key}>
                  <div className="nx-skeleton" style={{ height: 14, width: '80%' }} />
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    )
  }

  if (!data.length) {
    return (
      <div style={{ textAlign: 'center', padding: '40px 0', color: '#6b7280' }}>
        <div style={{ fontSize: 32, marginBottom: 8 }}>e</div>
        {emptyMessage}
      </div>
    )
  }

  return (
    <>
      <table className="nx-table">
        <thead>
          <tr>
            {columns.map((col) => (
              <th key={col.key} style={col.width ? { width: col.width } : undefined}>
                {col.header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {data.map((row) => (
            <tr key={(row as any).id ?? (row as any).uuid}>
              {columns.map((col) => (
                <td key={col.key}>
                  {col.render ? col.render(row) : (row as any)[col.key]}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
      {pagination && (
        <div className="nx-pagination">
          <span>
            Showing {(pagination.current - 1) * pagination.perPage + 1}-
            {Math.min(pagination.current * pagination.perPage, pagination.total)} of {pagination.total} records
          </span>
          <div className="nx-pagination-btns">
            <button
              className="nx-btn nx-btn--ghost nx-btn--sm"
              disabled={pagination.current <= 1}
              onClick={() => pagination.onPageChange(pagination.current - 1)}
            >
              Previous
            </button>
            <button
              className="nx-btn nx-btn--ghost nx-btn--sm"
              disabled={pagination.current >= pagination.total}
              onClick={() => pagination.onPageChange(pagination.current + 1)}
            >
              Next
            </button>
          </div>
        </div>
      )}
    </>
  )
}
