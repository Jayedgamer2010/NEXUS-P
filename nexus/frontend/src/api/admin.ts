import api from './client';
import type { User, Server, Node, Egg, PowerAction } from '../types';

// ── Helper types ────────────────────────────────────────────────────────────────
// After the interceptor, axios returns `res.data` directly, which is the backend
// envelope: { success, data: T, meta?: {...} }
// Axios still types it as AxiosResponse, but at runtime it IS the data.
// We cast the AxiosResponse generic to the unwrapped BackendEnvelope.

interface BackendEnvelope<T> {
  success: boolean;
  data: T;
  meta?: { total: number; per_page: number; current_page: number; last_page: number };
}

type PageArgs = Parameters<typeof api.get>[1][];

function paged<T>(url: string, page: number, limit: number): Promise<BackendEnvelope<T[]>> {
  return api.get<any, BackendEnvelope<T[]>>(url, { params: { page, limit } } as PageArgs[0] as any);
}

function post<T>(url: string, body?: Record<string, unknown>): Promise<BackendEnvelope<T>> {
  return api.post<any, BackendEnvelope<T>>(url, body);
}

function patch<T>(url: string, body: Record<string, unknown>): Promise<BackendEnvelope<T>> {
  return api.patch<any, BackendEnvelope<T>>(url, body);
}

function get<T>(url: string): Promise<BackendEnvelope<T>> {
  return api.get<any, BackendEnvelope<T>>(url);
}

function delete_1(url: string): Promise<void> {
  return api.delete(url) as Promise<void>;
}

// ── Paginated list wrapper ──────────────────────────────────────────────────────
export interface PaginatedResult<T> {
  data: T[];
  total: number;
  current_page: number;
  per_page: number;
  last_page: number;
}

function unwrapPaginated<T>(env: BackendEnvelope<T[]>): PaginatedResult<T> {
  return {
    data: env.data ?? [],
    total: env.meta?.total ?? 0,
    current_page: env.meta?.current_page ?? 1,
    per_page: env.meta?.per_page ?? 20,
    last_page: env.meta?.last_page ?? 1,
  };
}

// ── API object ──────────────────────────────────────────────────────────────────
export const adminApi = {
  // Users
  getUsers: (page = 1, limit = 20) => paged<User>('/api/admin/users', page, limit).then(unwrapPaginated),

  createUser: (data: Record<string, unknown>) => post<User>('/api/admin/users', data).then(e => e.data),

  updateUser: (id: number, data: Record<string, unknown>) => patch<User>(`/api/admin/users/${id}`, data).then(e => e.data),

  deleteUser: (id: number) => delete_1(`/api/admin/users/${id}`),

  // Nodes
  getNodes: (page = 1, limit = 20) => paged<Node>('/api/admin/nodes', page, limit).then(unwrapPaginated),

  createNode: (data: Record<string, unknown>) => post<Node>('/api/admin/nodes', data).then(e => e.data),

  updateNode: (id: number, data: Record<string, unknown>) => patch<Node>(`/api/admin/nodes/${id}`, data).then(e => e.data),

  deleteNode: (id: number) => delete_1(`/api/admin/nodes/${id}`),

  // Node allocations
  getNodeAllocations: (nodeId: number) => get<any[]>(`/api/admin/nodes/${nodeId}/allocations`).then(e => e.data),

  createAllocation: (nodeId: number, data: { ip: string; port: number }) =>
    post<any>(`/api/admin/nodes/${nodeId}/allocations`, data).then(e => e.data),

  deleteAllocation: (allocId: number) => delete_1(`/api/admin/allocations/${allocId}`),

  // Servers
  getServers: (page = 1, limit = 20) => paged<Server>('/api/admin/servers', page, limit).then(unwrapPaginated),

  createServer: (data: Record<string, unknown>) => post<Server>('/api/admin/servers', data).then(e => e.data),

  updateServer: (id: number, data: Record<string, unknown>) => patch<Server>(`/api/admin/servers/${id}`, data).then(e => e.data),

  deleteServer: (uuid: string) => delete_1(`/api/admin/servers/${uuid}`),

  powerServer: (uuid: string, action: PowerAction) =>
    post<null>(`/api/admin/servers/${uuid}/power`, { action }).then(e => e.data),

  suspendServer: (uuid: string) => post<null>(`/api/admin/servers/${uuid}/suspend`).then(e => e.data),

  unsuspendServer: (uuid: string) => post<null>(`/api/admin/servers/${uuid}/unsuspend`).then(e => e.data),

  reinstallServer: (uuid: string) => post<null>(`/api/admin/servers/${uuid}/reinstall`).then(e => e.data),

  // Eggs
  getEggs: (page = 1, limit = 20) => paged<Egg>('/api/admin/eggs', page, limit).then(unwrapPaginated),

  createEgg: (data: Record<string, unknown>) => post<Egg>('/api/admin/eggs', data).then(e => e.data),

  updateEgg: (id: number, data: Record<string, unknown>) => patch<Egg>(`/api/admin/eggs/${id}`, data).then(e => e.data),

  deleteEgg: (id: number) => delete_1(`/api/admin/eggs/${id}`),

  // Stats
  getStats: () => get<any>('/api/admin/stats').then(e => e.data),

  // Server detail by ID (single server fetch)
  getServerById: (id: number) => get<Server>(`/api/admin/servers/${id}`).then(e => e.data),
};
