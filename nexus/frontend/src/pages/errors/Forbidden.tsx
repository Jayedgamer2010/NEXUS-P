import { Link, useNavigate } from 'react-router-dom'
import Button from '../../components/ui/Button'

export default function Forbidden() {
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
        background: 'linear-gradient(135deg, #ef4444, #f59e0b)',
        WebkitBackgroundClip: 'text',
        WebkitTextFillColor: 'transparent',
        fontFamily: 'Inter, sans-serif',
        marginBottom: 8,
      }}>
        403
      </h1>
      <p style={{ color: '#6b7280', fontSize: 16, marginBottom: 24 }}>
        You do not have permission to access this page.
      </p>
      <Link to="/admin/dashboard">
        <Button>Back to Dashboard</Button>
      </Link>
    </div>
  )
}
