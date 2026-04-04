import './Spinner.css';

interface SpinnerProps {
  size?: 'sm' | 'md' | 'lg';
  color?: string;
}

export default function Spinner({ size = 'md', color }: SpinnerProps) {
  const sizeMap = { sm: '20px', md: '36px', lg: '56px' };

  return (
    <div className="spinner-container" style={{ width: sizeMap[size], height: sizeMap[size] }}>
      <div
        className="spinner"
        style={{
          borderColor: color ? `${color}33` : undefined,
          borderTopColor: color || '#7c3aed',
        }}
      />
    </div>
  );
}
