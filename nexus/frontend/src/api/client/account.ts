import api from '../client';
import type { User } from '../../types';

interface BackendEnvelope<T> {
  success: boolean;
  data: T;
  message?: string;
}

function get<T>(url: string): Promise<BackendEnvelope<T>> {
  return api.get<any, BackendEnvelope<T>>(url);
}

function patch<T>(url: string, body: Record<string, unknown>): Promise<BackendEnvelope<T>> {
  return api.patch<any, BackendEnvelope<T>>(url, body);
}

/**
 * Get the authenticated user's account info.
 */
export async function getAccount(): Promise<User> {
  const res = await get<User>('/api/client/account');
  return res.data;
}

/**
 * Update the authenticated user's account.
 */
export async function updateAccount(data: {
  first_name?: string;
  last_name?: string;
  language?: string;
  current_password?: string;
  new_password?: string;
}): Promise<User> {
  const res = await patch<User>('/api/client/account', data);
  return res.data;
}
