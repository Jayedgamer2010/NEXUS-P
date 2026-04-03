import { STATUS_COLORS } from '../../types';
import './StatusBadge.css';

interface StatusBadgeProps {
  status: string;
}

export default function StatusBadge({ status }: StatusBadgeProps) {
  const color = STATUS_COLORS[status] || '#6b7280';
  const label = status.charAt(0).toUpperCase() + status.slice(1);

  return (
    <span
      className="status-badge"
      style={{ backgroundColor: color }}
    >
      {label}
    </span>
  );
}
