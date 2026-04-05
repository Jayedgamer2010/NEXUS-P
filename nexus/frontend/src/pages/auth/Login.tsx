import { useState, useEffect, FormEvent } from 'react'
import { useNavigate } from 'react-router-dom'
import { authApi } from '../../api/auth'
import { useAuthStore } from '../../store/authStore'
import Button from '../../components/ui/Button'
import Input from '../../components/ui/Input'
import Alert from '../../components/ui/Alert'

export default function Login() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const { isAuthenticated } = useAuthStore()
  const navigate = useNavigate()

  useEffect(() => {
    if (isAuthenticated) navigate('/admin/dashboard', { replace: true })
  }, [isAuthenticated, navigate])

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError(null)

    try {
      const res = await authApi.login(email, password)
      const data = res.data.data
      useAuthStore.getState().login(data.token, data.user)
      navigate('/admin/dashboard')
    } catch (err: any) {
      const msg = err.response?.data?.message || err.message || 'Login failed. Please try again.'
      setError(msg)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{
      minHeight: '100vh',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      background: '#0a0a0f',
      padding: 20,
    }}>
      <div style={{
        width: 420,
        maxWidth: '100%',
        background: '#0d0d17',
        border: '1px solid #1e1e30',
        borderRadius: 12,
        padding: 48,
      }}>
        <div style={{
          textAlign: 'center',
          marginBottom: 32,
        }}>
          <h1 style={{
            fontSize: 36,
            fontWeight: 700,
            background: 'linear-gradient(135deg, #7c3aed, #3b82f6)',
            WebkitBackgroundClip: 'text',
            WebkitTextFillColor: 'transparent',
            fontFamily: 'Inter, sans-serif',
            marginBottom: 8,
          }}>
            NEXUS
          </h1>
          <p style={{ fontSize: 13, color: '#6b7280', margin: 0 }}>
            Game Server Management Panel
          </p>
        </div>

        <div style={{ height: 1, background: '#1e1e30', marginBottom: 24 }} />

        <form onSubmit={handleSubmit}>
          <Input
            label="Email Address"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            autoComplete="email"
            placeholder="admin@example.com"
          />

          <Input
            label="Password"
            type={showPassword ? 'text' : 'password'}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            autoComplete="current-password"
            placeholder="Enter your password"
          />

          <div style={{ display: 'flex', gap: 6, marginTop: -8, marginBottom: 16 }}>
            <label style={{ fontSize: 12, color: '#6b7280', cursor: 'pointer' }}>
              <input
                type="checkbox"
                checked={showPassword}
                onChange={(e) => setShowPassword(e.target.checked)}
                style={{ marginRight: 4 }}
              />
              Show password
            </label>
          </div>

          {error && (
            <div style={{ marginBottom: 16 }}>
              <Alert type="error" message={error} dismissible={false} />
            </div>
          )}

          <Button
            type="submit"
            variant="primary"
            size="lg"
            loading={loading}
            style={{ width: '100%', height: 48 }}
          >
            Sign In
          </Button>
        </form>
      </div>
    </div>
  )
}
