import apiClient from './client'
import type { ApiResponse, User } from '../types'

export const authApi = {
  login: (email: string, password: string) =>
    apiClient.post<ApiResponse<{ user: User; token: string }>>('/api/auth/login', { email, password }),

  register: (username: string, email: string, password: string) =>
    apiClient.post<ApiResponse<{ user: User; token: string }>>('/api/auth/register', { username, email, password }),

  me: () =>
    apiClient.get<ApiResponse<{ user: User }>>('/api/auth/me'),
}
