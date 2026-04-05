import apiClient from '../client'

export const eggsApi = {
  getAll: () => apiClient.get('/api/admin/eggs'),
  getOne: (id: number) => apiClient.get('/api/admin/eggs/' + id),
  create: (data: any) => apiClient.post('/api/admin/eggs', data),
  update: (id: number, data: any) => apiClient.patch('/api/admin/eggs/' + id, data),
  delete: (id: number) => apiClient.delete('/api/admin/eggs/' + id),
}
