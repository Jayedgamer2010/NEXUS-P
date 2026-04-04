import React from 'react';
import { createBrowserRouter, Navigate } from 'react-router-dom';
import { useAuthStore } from '../store/authStore';
import Layout from '../components/layout/Layout';
import Login from '../pages/auth/Login';
import Dashboard from '../pages/admin/Dashboard';
import Servers from '../pages/admin/Servers';
import ServerDetail from '../pages/admin/ServerDetail';
import Nodes from '../pages/admin/Nodes';
import NodeDetail from '../pages/admin/NodeDetail';
import Users from '../pages/admin/Users';
import Eggs from '../pages/admin/Eggs';
import NotFound from '../pages/errors/404';
import Forbidden from '../pages/errors/403';

// Protected route wrapper
function ProtectedRoute({ children, requireAdmin = false }: { children: React.ReactNode; requireAdmin?: boolean }) {
  const { isAuthenticated, user } = useAuthStore();

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  if (requireAdmin && user?.role !== 'admin') {
    return <Navigate to="/403" replace />;
  }

  return <Layout>{children}</Layout>;
}

export const router = createBrowserRouter([
  {
    path: '/login',
    element: <Login />,
  },
  {
    path: '/403',
    element: <Forbidden />,
  },
  {
    path: '/404',
    element: <NotFound />,
  },
  {
    path: '/admin',
    element: <ProtectedRoute requireAdmin={true}><Navigate to="/admin/dashboard" replace /></ProtectedRoute>,
  },
  {
    path: '/admin/dashboard',
    element: <ProtectedRoute requireAdmin={true}><Dashboard /></ProtectedRoute>,
  },
  {
    path: '/admin/servers',
    element: <ProtectedRoute requireAdmin={true}><Servers /></ProtectedRoute>,
  },
  {
    path: '/admin/servers/:id',
    element: <ProtectedRoute requireAdmin={true}><ServerDetail /></ProtectedRoute>,
  },
  {
    path: '/admin/nodes',
    element: <ProtectedRoute requireAdmin={true}><Nodes /></ProtectedRoute>,
  },
  {
    path: '/admin/nodes/:id',
    element: <ProtectedRoute requireAdmin={true}><NodeDetail /></ProtectedRoute>,
  },
  {
    path: '/admin/users',
    element: <ProtectedRoute requireAdmin={true}><Users /></ProtectedRoute>,
  },
  {
    path: '/admin/eggs',
    element: <ProtectedRoute requireAdmin={true}><Eggs /></ProtectedRoute>,
  },
  {
    path: '/',
    element: <Navigate to="/admin/dashboard" replace />,
  },
  {
    path: '*',
    element: <NotFound />,
  },
]);
