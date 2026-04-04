import { NavLink } from 'react-router-dom';
import { useLocation } from 'react-router-dom';
import { useAuthStore } from '../../store/authStore';
import './Sidebar.css';

export default function Sidebar() {
  const { user, logout } = useAuthStore();
  const location = useLocation();

  const navItems = [
    { to: '/admin/dashboard', label: 'Dashboard', icon: '📊' },
    { to: '/admin/servers', label: 'Servers', icon: '🖥️' },
    { to: '/admin/nodes', label: 'Nodes', icon: '🌐' },
    { to: '/admin/users', label: 'Users', icon: '👥' },
    { to: '/admin/eggs', label: 'Eggs', icon: '🥚' },
  ];

  return (
    <aside className="sidebar">
      <div className="sidebar-header">
        <div className="logo">NEXUS</div>
      </div>

      <nav className="sidebar-nav">
        {navItems.map((item) => {
          const isActive = location.pathname === item.to || location.pathname.startsWith(item.to + '/');
          return (
            <NavLink
              key={item.to}
              to={item.to}
              className={`nav-item${isActive ? ' active' : ''}`}
            >
              <span className="nav-icon">{item.icon}</span>
              <span className="nav-label">{item.label}</span>
            </NavLink>
          );
        })}
      </nav>

      <div className="sidebar-footer">
        <div className="user-menu">
          <div className="user-avatar">
            {user?.username?.charAt(0).toUpperCase() || 'U'}
          </div>
          <div className="user-info">
            <div className="user-name">{user?.username}</div>
            <div className="user-role">{user?.role}</div>
          </div>
          <button className="logout-btn" onClick={logout} title="Logout">
            →
          </button>
        </div>
      </div>
    </aside>
  );
}
