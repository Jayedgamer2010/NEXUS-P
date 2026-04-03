import { Outlet, useLocation } from 'react-router-dom';
import Sidebar from './Sidebar';
import Header from './Header';
import './Layout.css';

const pageTitleMap: Record<string, string> = {
  '/admin/dashboard': 'Dashboard',
  '/admin/servers': 'Servers',
  '/admin/nodes': 'Nodes',
  '/admin/users': 'Users',
  '/admin/eggs': 'Eggs',
};

export default function Layout() {
  const location = useLocation();
  const title = pageTitleMap[location.pathname] || 'NEXUS';

  return (
    <div className="layout">
      <Sidebar />
      <main className="main-content">
        <Header title={title} />
        <div className="page-content">
          <Outlet />
        </div>
      </main>
    </div>
  );
}
