import api from '../client'
import { Server, PaginatedData, PowerAction } from '../../types'

export interface CreateServerData {
  name: string
  user_id: number
  node_id: number
  egg_id: number
  allocation_id: number
  memory: number
  disk: number
  cpu: number
  swap?: number
  io?: number
  startup?: string
  image?: string
  database_limit?: number
  allocation_limit?: number
  backup_limit?: number
}

export interface UpdateServerData extends Partial<Server> {
  password?: string
}

export interface ListServersParams {
  page?: number
  per_page?: number
  search?: string
}

export interface PowerResponse {
  result: string
}

export interface CreateServerResponse {
  id: number
  uuid: string
  status: string
}

// ─── Admin Servers ─────────────────────────────────────────────────────────────

export const getAll = async (params?: ListServersParams): Promise<PaginatedData<Server>> => {
  return api.get('/admin/servers', { params })
}

export const create = async (data: CreateServerData): Promise<CreateServerResponse> => {
  return api.post('/admin/servers', data)
}

export const getById = async (id: number): Promise<Server> => {
  return api.get(`/admin/servers/${id}`)
}

export const update = async (id: number, data: UpdateServerData): Promise<Server> => {
  return api.patch(`/admin/servers/${id}`, data)
}

export const del = async (uuid: string): Promise<void> => {
  return api.delete(`/admin/servers/${uuid}`)
}

export const power = async (uuid: string, action: PowerAction): Promise<PowerResponse> => {
  return api.post(`/admin/servers/${uuid}/power`, { action })
}

export const suspend = async (uuid: string): Promise<void> => {
  return api.post(`/admin/servers/${uuid}/suspend`)
}

export const unsuspend = async (uuid: string): Promise<void> => {
  return api.post(`/admin/servers/${uuid}/unsuspend`)
}

export const reinstall = async (uuid: string): Promise<void> => {
  return api.post(`/admin/servers/${uuid}/reinstall`)
}
