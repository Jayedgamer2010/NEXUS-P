import apiClient from '../client'

export const clientServersApi = {
  getAll: () => apiClient.get('/api/client/servers'),
  getOne: (uuid: string) => apiClient.get('/api/client/servers/' + uuid),
  getResources: (uuid: string) => apiClient.get('/api/client/servers/' + uuid + '/resources'),
  power: (uuid: string, action: string) =>
    apiClient.post('/api/client/servers/' + uuid + '/power', { action }),
}
