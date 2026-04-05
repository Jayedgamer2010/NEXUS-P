import { Link, useLocation, useNavigate } from 'react-router-dom'
import { LayoutDashboard, Server, Network, Users, Package, LogOut } from 'lucide-react'
import { useAuthStore } from '../../store/authStore'

const navItems = [
  { path: '/admin/dashboard', label: 'Dashboard', icon: LayoutDashboard },
  { path: '/admin/servers', label: 'Servers', icon: Server },
  { path: '/admin/nodes', label: 'Nodes', icon: Network },
  { path: '/admin/users', label: 'Users', icon: Users },
  { path: '/admin/eggs', label: 'Eggs', icon: Package },
]

export default function Sidebar() {
  const location = useLocation()
  const navigate = useNavigate()
  const { user, logout } = useAuthStore()

  const isActive = (path: string) => location.pathname.startsWith(path)

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  return (
    <div className="nx-sidebar">
      <div className="nx-sidebar-logo">NEXUS</div>

      <nav className="nx-sidebar-nav">
        {navItems.map((item) => {
          const Icon = item.icon
          const active = isActive(item.path)
          return (
            <Link
              key={item.path}
              to={item.path}
              className={`nx-sidebar-item ${active ? 'nx-sidebar-item--active' : ''}`}
            >
              <Icon />
              {item.label}
            </Link>
          )
        })}
      </nav>

      <div className="nx-sidebar-footer">
        {user && (
          <>
            <div className="nx-sidebar-user">
              <div className="nx-sidebar-avatar">
                {user.username.charAt(0).toUpperCase()}
              </div>
              <div className="nx-sidebar-user-info">
                <div className="nx-sidebar-user-name">{user.username}</div>
                <div className="nx-sidebar-user-role">{user.root_admin ? 'Root Admin' : user.role}</div>
              </div>
            </div>
            <button className="nx-logout-btn" onClick={handleLogout}>
              <LogOut />
              Logout
            </button>
          </>
        )}
      </div>
    </div>
  )
}
