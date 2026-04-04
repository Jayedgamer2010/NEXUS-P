import React from 'react';
import './DataTable.css';

interface Column<T> {
  key: keyof T | string;
  header: string;
  render?: (value: unknown, row: T) => React.ReactNode;
  width?: string;
}

interface DataTableProps<T> {
  columns: Column<T>[];
  data: T[];
  loading?: boolean;
  emptyMessage?: string;
  pagination?: {
    page: number;
    limit: number;
    total: number;
    onPageChange: (page: number) => void;
  };
}

export default function DataTable<T extends { id: number }>({
  columns,
  data,
  loading = false,
  emptyMessage = 'No records found',
  pagination,
}: DataTableProps<T>) {
  const totalPages = pagination ? Math.ceil(pagination.total / pagination.limit) : 1;
  const startIdx = pagination ? (pagination.page - 1) * pagination.limit + 1 : data.length > 0 ? 1 : 0;
  const endIdx = pagination ? Math.min(pagination.page * pagination.limit, pagination.total) : data.length;
  const colCount = columns.length;

  if (loading) {
    return (
      <div className="data-table-container">
        <table className="data-table">
          <thead>
            <tr>
              {columns.map((col, idx) => (
                <th key={idx} style={{ width: col.width }}>{col.header}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {Array.from({ length: 5 }).map((_, rowIdx) => (
              <tr key={rowIdx} className="skeleton-row">
                {columns.map((_, colIdx) => (
                  <td key={colIdx}>
                    <div className="skeleton-block" style={{ width: '80%', height: '16px' }} />
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    );
  }

  return (
    <div className="data-table-container">
      <table className="data-table">
        <thead>
          <tr>
            {columns.map((col, idx) => (
              <th key={idx} style={{ width: col.width }}>{col.header}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {data.length === 0 ? (
            <tr>
              <td colSpan={colCount} className="empty-cell">
                <div className="empty-state">
                  <div className="empty-icon">
                    <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5">
                      <circle cx="12" cy="12" r="10" />
                      <line x1="8" y1="12" x2="16" y2="12" />
                    </svg>
                  </div>
                  <p>{emptyMessage}</p>
                </div>
              </td>
            </tr>
          ) : (
            data.map((row) => (
              <tr key={row.id}>
                {columns.map((col, colIdx) => {
                  const value = typeof col.key === 'string' ? (row as any)[col.key] : undefined;
                  return (
                    <td key={colIdx}>
                      {col.render ? col.render(value, row) : value}
                    </td>
                  );
                })}
              </tr>
            ))
          )}
        </tbody>
      </table>

      {pagination && totalPages > 0 && (
        <div className="pagination">
          <div className="pagination-info">
            Showing {startIdx}&ndash;{endIdx} of {pagination.total}
          </div>
          <button
            className="pagination-btn"
            disabled={pagination.page <= 1}
            onClick={() => pagination.onPageChange(pagination.page - 1)}
          >
            &larr; Prev
          </button>
          <span className="pagination-page-indicator">
            {pagination.page} / {totalPages}
          </span>
          <button
            className="pagination-btn"
            disabled={pagination.page >= totalPages}
            onClick={() => pagination.onPageChange(pagination.page + 1)}
          >
            Next &rarr;
          </button>
        </div>
      )}
    </div>
  );
}
