import api from '../client'
import { User, PaginatedData } from '../../types'

export interface CreateUserData {
  username: string
  email: string
  password: string
  role: 'admin' | 'client'
  root_admin?: boolean
  coins?: number
  name_first?: string
  name_last?: string
}

export interface UpdateUserData extends Partial<User> {
  password?: string
}

export interface ListUsersParams {
  page?: number
  per_page?: number
  search?: string
}

// ─── Admin Users ───────────────────────────────────────────────────────────────

export const getAll = async (params?: ListUsersParams): Promise<PaginatedData<User>> => {
  return api.get('/admin/users', { params })
}

export const create = async (data: CreateUserData): Promise<User> => {
  return api.post('/admin/users', data)
}

export const getById = async (id: number): Promise<User> => {
  return api.get(`/admin/users/${id}`)
}

export const update = async (id: number, data: UpdateUserData): Promise<User> => {
  return api.patch(`/admin/users/${id}`, data)
}

export const del = async (id: number): Promise<void> => {
  return api.delete(`/admin/users/${id}`)
}
