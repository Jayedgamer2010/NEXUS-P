import apiClient from '../client'

export const serversApi = {
  getAll: (page = 1) =>
    apiClient.get('/api/admin/servers', { params: { page } }),
  getOne: (id: number) =>
    apiClient.get('/api/admin/servers/' + id),
  create: (data: any) =>
    apiClient.post('/api/admin/servers', data),
  update: (id: number, data: any) =>
    apiClient.patch('/api/admin/servers/' + id, data),
  delete: (id: number) =>
    apiClient.delete('/api/admin/servers/' + id),
  power: (id: number, action: string) =>
    apiClient.post('/api/admin/servers/' + id + '/power', { action }),
  suspend: (id: number) =>
    apiClient.post('/api/admin/servers/' + id + '/suspend'),
  unsuspend: (id: number) =>
    apiClient.post('/api/admin/servers/' + id + '/unsuspend'),
}
