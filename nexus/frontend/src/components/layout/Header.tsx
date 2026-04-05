import { useAuthStore } from '../../store/authStore'

interface HeaderProps {
  title: string
}

export default function Header({ title }: HeaderProps) {
  const { user } = useAuthStore()

  const initials = user ? user.username.substring(0, 2).toUpperCase() : '?'

  return (
    <div className="nx-header">
      <div className="nx-header-title">{title}</div>
      {user && (
        <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
          <div
            style={{
              width: 32,
              height: 32,
              borderRadius: '50%',
              background: '#7c3aed',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              fontSize: 13,
              fontWeight: 600,
            }}
          >
            {initials}
          </div>
          <div>
            <div style={{ fontSize: 13, fontWeight: 500 }}>{user.username}</div>
            <div style={{ fontSize: 10, color: '#6b7280' }}>{user.root_admin ? 'Root Admin' : user.role}</div>
          </div>
        </div>
      )}
    </div>
  )
}
