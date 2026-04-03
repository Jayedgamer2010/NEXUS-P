import api from './client';
import { User, ApiResponse } from '../types';

export const authApi = {
  login: async (email: string, password: string): Promise<ApiResponse<{ token: string; user: User }>> => {
    return api.post<{ token: string; user: User }>('/api/auth/login', {
      email,
      password,
    });
  },

  register: async (
    username: string,
    email: string,
    password: string,
    role?: string
  ): Promise<ApiResponse<{ token: string; user: User }>> => {
    return api.post<{ token: string; user: User }>('/api/auth/register', {
      username,
      email,
      password,
      ...(role && { role }),
    });
  },

  me: async (): Promise<ApiResponse<User>> => {
    return api.get<User>('/api/auth/me');
  },

  logout: async (): Promise<ApiResponse<null>> => {
    return api.post<null>('/api/auth/logout');
  },
};
