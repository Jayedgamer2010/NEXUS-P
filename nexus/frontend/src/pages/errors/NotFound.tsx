import { Link } from 'react-router-dom'

export default function NotFound() {
  return (
    <div style={{
      minHeight: '100vh',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      background: '#0a0a0f',
      flexDirection: 'column',
      gap: 16,
    }}>
      <h1 style={{
        fontSize: 72,
        fontWeight: 700,
        background: 'linear-gradient(135deg, #7c3aed, #3b82f6)',
        WebkitBackgroundClip: 'text',
        WebkitTextFillColor: 'transparent',
        fontFamily: 'Inter, sans-serif',
        marginBottom: 8,
      }}>
        404
      </h1>
      <p style={{ color: '#6b7280', fontSize: 16, marginBottom: 24 }}>
        The page you are looking for does not exist.
      </p>
      <Link to="/admin/dashboard" className="nx-btn nx-btn--primary nx-btn--lg">
        Back to Dashboard
      </Link>
    </div>
  )
}
