import api from './client'
import { User } from '../types'

export interface LoginCredentials {
  email: string
  password: string
}

export interface RegisterCredentials {
  username: string
  email: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

// ─── Authentication ────────────────────────────────────────────────────────────

export const login = async (
  email: string,
  password: string
): Promise<LoginResponse> => {
  return api.post('/auth/login', { email, password })
}

export const register = async (
  username: string,
  email: string,
  password: string
): Promise<LoginResponse> => {
  return api.post('/auth/register', { username, email, password })
}

export const getMe = async (): Promise<User> => {
  return api.get('/auth/me')
}

// Named export object for components that import { authApi }
export const authApi = {
  login: async (email: string, password: string) => {
    const res = await api.post<{ token: string; user: User }>('/auth/login', { email, password });
    return { success: true, data: res.data, message: 'Login successful' };
  },
  register: async (username: string, email: string, password: string) => {
    const res = await api.post<{ token: string; user: User }>('/auth/register', { username, email, password });
    return { success: true, data: res.data, message: 'Account created' };
  },
  getMe: async () => api.get<User>('/auth/me'),
};
