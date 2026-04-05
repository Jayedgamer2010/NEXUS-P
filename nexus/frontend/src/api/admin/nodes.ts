import apiClient from '../client'

export const nodesApi = {
  getAll: () => apiClient.get('/api/admin/nodes'),
  getOne: (id: number) => apiClient.get('/api/admin/nodes/' + id),
  create: (data: any) => apiClient.post('/api/admin/nodes', data),
  update: (id: number, data: any) => apiClient.patch('/api/admin/nodes/' + id, data),
  delete: (id: number) => apiClient.delete('/api/admin/nodes/' + id),
  getAllocations: (id: number) => apiClient.get('/api/admin/nodes/' + id + '/allocations'),
  addAllocation: (id: number, data: any) => apiClient.post('/api/admin/nodes/' + id + '/allocations', data),
  deleteAllocation: (allocId: number) => apiClient.delete('/api/admin/allocations/' + allocId),
}
