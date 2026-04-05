import Spinner from './Spinner';
import './StatCard.css';

interface StatCardProps {
  title: string;
  value: string | number;
  icon?: React.ReactNode;
  accentColor?: string;
  subtitle?: string;
  loading?: boolean;
}

export default function StatCard({ title, value, icon, accentColor = '#7c3aed', subtitle, loading = false }: StatCardProps) {
  const displayValue = value ?? 0;
  return (
    <div className="stat-card">
      <div className="stat-card-header">
        {icon && (
          <div className="stat-icon" style={{ color: accentColor }}>
            {icon}
          </div>
        )}
        <div className="stat-title">{title}</div>
      </div>
      {loading ? (
        <div className="stat-value-loading">
          <Spinner size="sm" color={accentColor} />
        </div>
      ) : (
        <div className="stat-value" style={{ color: accentColor }}>
          {displayValue}
        </div>
      )}
      {subtitle && <div className="stat-subtitle">{subtitle}</div>}
    </div>
  );
}
