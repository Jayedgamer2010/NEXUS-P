import apiClient from '../client'

export const usersApi = {
  getAll: (page = 1, search = '') =>
    apiClient.get('/api/admin/users', { params: { page, search } }),
  getOne: (id: number) =>
    apiClient.get('/api/admin/users/' + id),
  create: (data: any) =>
    apiClient.post('/api/admin/users', data),
  update: (id: number, data: any) =>
    apiClient.patch('/api/admin/users/' + id, data),
  delete: (id: number) =>
    apiClient.delete('/api/admin/users/' + id),
}
