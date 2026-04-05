import apiClient from '../client'

export const clientAccountApi = {
  get: () => apiClient.get('/api/client/account'),
  update: (data: any) => apiClient.patch('/api/client/account', data),
}
