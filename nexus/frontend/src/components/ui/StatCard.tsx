import './StatCard.css';

interface StatCardProps {
  title: string;
  value: string | number;
  icon?: React.ReactNode;
  accentColor?: string;
  subtitle?: string;
}

export default function StatCard({ title, value, icon, accentColor = '#7c3aed', subtitle }: StatCardProps) {
  const displayValue = value ?? 0;
  return (
    <div className="stat-card">
      <div className="stat-card-header">
        <div className="stat-icon" style={{ color: accentColor }}>
          {icon}
        </div>
        <div className="stat-title">{title}</div>
      </div>
      <div className="stat-value" style={{ color: accentColor }}>
        {displayValue}
      </div>
      {subtitle && <div className="stat-subtitle">{subtitle}</div>}
    </div>
  );
}
