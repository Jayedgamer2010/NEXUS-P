import { Outlet, useLocation } from 'react-router-dom'
import Sidebar from './Sidebar'
import Header from './Header'

const titles: Record<string, string> = {
  '/admin/dashboard': 'Dashboard',
  '/admin/servers': 'Servers',
  '/admin/nodes': 'Nodes',
  '/admin/users': 'Users',
  '/admin/eggs': 'Eggs',
}

export default function Layout() {
  const location = useLocation()

  // Find the best matching title
  let title = 'Admin'
  for (const [path, label] of Object.entries(titles)) {
    if (location.pathname === path || location.pathname.startsWith(path + '/')) {
      title = label
      // For detail pages, append info
      if (location.pathname.split('/').length > 3 && path !== '/admin/nodes' && path !== '/admin/eggs') {
        title = label.replace(/s$/, '') + ' Detail'
      }
      break
    }
  }

  return (
    <div className="nx-layout">
      <Sidebar />
      <div className="nx-content">
        <Header title={title} />
        <div className="nx-page">
          <Outlet />
        </div>
      </div>
    </div>
  )
}
