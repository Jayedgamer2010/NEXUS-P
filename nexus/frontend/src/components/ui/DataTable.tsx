import React from 'react';
import './DataTable.css';

interface Column<T> {
  key: keyof T | string;
  header: string;
  render?: (value: any, row: T) => React.ReactNode;
  width?: string;
}

interface DataTableProps<T> {
  columns: Column<T>[];
  data: T[];
  loading?: boolean;
  pagination?: {
    page: number;
    limit: number;
    total: number;
    onPageChange: (page: number) => void;
  };
  onRowAction?: (action: string, row: T) => void;
  actionColumn?: {
    label: string;
    actions: Array<{ key: string; label: string; icon?: string }>;
  };
}

export default function DataTable<T extends { id: number }>({
  columns,
  data,
  loading = false,
  pagination,
  onRowAction,
  actionColumn,
}: DataTableProps<T>) {
  const getCellValue = (row: T, column: Column<T>) => {
    const value = typeof column.key === 'string' ? (row as any)[column.key] : undefined;
    if (column.render) {
      return column.render(value, row);
    }
    return value;
  };

  const totalPages = pagination ? Math.ceil(pagination.total / pagination.limit) : 1;

  return (
    <div className="data-table-container">
      <table className="data-table">
        <thead>
          <tr>
            {columns.map((col, idx) => (
              <th key={idx} style={{ width: col.width }}>
                {col.header}
              </th>
            ))}
            {actionColumn && <th>Actions</th>}
          </tr>
        </thead>
        <tbody>
          {loading ? (
            <tr>
              <td colSpan={columns.length + (actionColumn ? 1 : 0)} className="loading-cell">
                Loading...
              </td>
            </tr>
          ) : data.length === 0 ? (
            <tr>
              <td colSpan={columns.length + (actionColumn ? 1 : 0)} className="empty-cell">
                No data available
              </td>
            </tr>
          ) : (
            data.map((row) => (
              <tr key={row.id}>
                {columns.map((col, idx) => (
                  <td key={idx}>{getCellValue(row, col)}</td>
                ))}
                {actionColumn && (
                  <td className="actions-cell">
                    {actionColumn.actions.map((action) => (
                      <button
                        key={action.key}
                        className="action-btn"
                        onClick={() => onRowAction?.(action.key, row)}
                        title={action.label}
                      >
                        {action.icon || action.label}
                      </button>
                    ))}
                  </td>
                )}
              </tr>
            ))
          )}
        </tbody>
      </table>

      {pagination && totalPages > 1 && (
        <div className="pagination">
          <button
            className="pagination-btn"
            disabled={pagination.page <= 1}
            onClick={() => pagination.onPageChange(pagination.page - 1)}
          >
            Previous
          </button>
          <span className="pagination-info">
            Page {pagination.page} of {totalPages}
          </span>
          <button
            className="pagination-btn"
            disabled={pagination.page >= totalPages}
            onClick={() => pagination.onPageChange(pagination.page + 1)}
          >
            Next
          </button>
        </div>
      )}
    </div>
  );
}
