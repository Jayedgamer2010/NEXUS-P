import { createBrowserRouter, Navigate, Outlet } from 'react-router-dom'
import Login from '../pages/auth/Login'
import Layout from '../components/layout/Layout'
import Dashboard from '../pages/admin/Dashboard'
import Servers from '../pages/admin/Servers'
import ServerDetail from '../pages/admin/ServerDetail'
import Nodes from '../pages/admin/Nodes'
import NodeDetail from '../pages/admin/NodeDetail'
import Users from '../pages/admin/Users'
import Eggs from '../pages/admin/Eggs'
import NotFound from '../pages/errors/NotFound'
import Forbidden from '../pages/errors/Forbidden'
import { useAuthStore } from '../store/authStore'

// Protected route wrapper
function AdminRoute() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)
  const user = useAuthStore((s) => s.user)

  if (!isAuthenticated) return <Navigate to="/login" replace />
  if (!user?.root_admin && user?.role !== 'admin') return <Navigate to="/403" replace />
  return <Outlet />
}

function GuestRoute() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated)
  if (isAuthenticated) return <Navigate to="/admin/dashboard" replace />
  return <Outlet />
}

export const router = createBrowserRouter([
  {
    path: '/login',
    element: <GuestRoute />,
    children: [{ index: true, element: <Login /> }],
  },
  {
    path: '/admin',
    element: <AdminRoute />,
    children: [
      {
        path: '/admin',
        element: <Layout />,
        children: [
          { path: '', element: <Navigate to="/admin/dashboard" replace /> },
          { path: 'dashboard', element: <Dashboard /> },
          { path: 'servers', element: <Servers /> },
          { path: 'servers/:id', element: <ServerDetail /> },
          { path: 'nodes', element: <Nodes /> },
          { path: 'nodes/:id', element: <NodeDetail /> },
          { path: 'users', element: <Users /> },
          { path: 'eggs', element: <Eggs /> },
        ],
      },
    ],
  },
  { path: '/403', element: <Forbidden /> },
  { path: '/404', element: <NotFound /> },
  { path: '/', element: <Navigate to="/admin/dashboard" replace /> },
  { path: '*', element: <NotFound /> },
])
