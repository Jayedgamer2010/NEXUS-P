import api from './client';
import { User, Server, Node, Egg, PaginatedResponse, ApiResponse, PowerAction } from '../types';

export const adminApi = {
  // Users
  getUsers: async (page = 1, limit = 20): Promise<PaginatedResponse<User>> => {
    return api.get<PaginatedResponse<User>>('/api/admin/users', { page, limit });
  },

  createUser: async (data: {
    username: string;
    email: string;
    password: string;
    role: 'admin' | 'client';
    coins?: number;
  }): Promise<ApiResponse<User>> => {
    return api.post<User>('/api/admin/users', data);
  },

  updateUser: async (id: number, data: Partial<User>): Promise<ApiResponse<User>> => {
    return api.patch<User>(`/api/admin/users/${id}`, data);
  },

  deleteUser: async (id: number): Promise<ApiResponse<null>> => {
    return api.delete<null>(`/api/admin/users/${id}`);
  },

  // Nodes
  getNodes: async (page = 1, limit = 20): Promise<PaginatedResponse<Node>> => {
    return api.get<PaginatedResponse<Node>>('/api/admin/nodes', { page, limit });
  },

  createNode: async (data: {
    name: string;
    fqdn: string;
    scheme?: 'http' | 'https';
    wings_port?: number;
    memory: number;
    memory_overalloc: number;
    disk: number;
    disk_overalloc: number;
    token_id: string;
    token: string;
  }): Promise<ApiResponse<Node>> => {
    return api.post<Node>('/api/admin/nodes', data);
  },

  updateNode: async (id: number, data: Partial<Node>): Promise<ApiResponse<Node>> => {
    return api.patch<Node>(`/api/admin/nodes/${id}`, data);
  },

  deleteNode: async (id: number): Promise<ApiResponse<null>> => {
    return api.delete<null>(`/api/admin/nodes/${id}`);
  },

  getNodeStats: async (id: number): Promise<ApiResponse<{
    total_memory: number;
    total_disk: number;
    memory_limit: number;
    disk_limit: number;
    server_count: number;
  }>> => {
    return api.get(`/api/admin/nodes/${id}/stats`);
  },

  // Servers
  getServers: async (page = 1, limit = 20): Promise<PaginatedResponse<Server>> => {
    return api.get<PaginatedResponse<Server>>('/api/admin/servers', { page, limit });
  },

  createServer: async (data: {
    name: string;
    user_id: number;
    node_id: number;
    egg_id: number;
    allocation_id: number;
    memory: number;
    disk: number;
    cpu: number;
    startup?: string;
    image?: string;
  }): Promise<ApiResponse<{ id: number; uuid: string; status: string }>> => {
    return api.post('/api/admin/servers', data);
  },

  updateServer: async (id: number, data: Partial<Server>): Promise<ApiResponse<Server>> => {
    return api.patch<Server>(`/api/admin/servers/${id}`, data);
  },

  deleteServer: async (id: number): Promise<ApiResponse<null>> => {
    return api.delete<null>(`/api/admin/servers/${id}`);
  },

  powerServer: async (id: number, action: PowerAction): Promise<ApiResponse<null>> => {
    return api.post<null>(`/api/admin/servers/${id}/power`, { action });
  },

  // Eggs
  getEggs: async (page = 1, limit = 20): Promise<PaginatedResponse<Egg>> => {
    return api.get<PaginatedResponse<Egg>>('/api/admin/eggs', { page, limit });
  },

  createEgg: async (data: {
    name: string;
    description?: string;
    docker_image: string;
    startup_command: string;
  }): Promise<ApiResponse<Egg>> => {
    return api.post<Egg>('/api/admin/eggs', data);
  },

  updateEgg: async (id: number, data: Partial<Egg>): Promise<ApiResponse<Egg>> => {
    return api.patch<Egg>(`/api/admin/eggs/${id}`, data);
  },

  deleteEgg: async (id: number): Promise<ApiResponse<null>> => {
    return api.delete<null>(`/api/admin/eggs/${id}`);
  },
};
