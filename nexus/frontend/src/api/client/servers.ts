import api from '../client';
import type { Server, ServerResources } from '../../types';

interface BackendEnvelope<T> {
  success: boolean;
  data: T;
  message?: string;
}

function get<T>(url: string): Promise<BackendEnvelope<T>> {
  return api.get<any, BackendEnvelope<T>>(url);
}

function post(url: string, body?: Record<string, unknown>): Promise<BackendEnvelope<null>> {
  return api.post<any, BackendEnvelope<null>>(url, body);
}

/**
 * Get all servers for the authenticated user.
 */
export async function getMyServers(): Promise<Server[]> {
  const res = await get<Server[]>('/api/client/servers');
  return res.data;
}

/**
 * Get a single server by UUID.
 */
export async function getMyServer(uuid: string): Promise<Server> {
  const res = await get<Server>(`/api/client/servers/${uuid}`);
  return res.data;
}

/**
 * Get live resource usage for a server.
 */
export async function getServerResources(uuid: string): Promise<ServerResources> {
  const res = await get<ServerResources>(`/api/client/servers/${uuid}/resources`);
  return res.data;
}

/**
 * Send a power action to a server.
 */
export async function sendPowerAction(uuid: string, action: string): Promise<void> {
  await post(`/api/client/servers/${uuid}/power`, { action });
}
