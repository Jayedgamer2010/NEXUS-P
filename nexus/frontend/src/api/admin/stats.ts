import apiClient from '../client'
import type { ApiResponse, AdminStats } from '../../types'

export const statsApi = {
  get: () => apiClient.get<ApiResponse<AdminStats>>('/api/admin/stats'),
}
