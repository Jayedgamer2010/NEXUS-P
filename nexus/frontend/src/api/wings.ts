import api from './client';
import { ServerResources } from '../types';

export const wingsApi = {
  getResources: async (serverUUID: string): Promise<ServerResources> => {
    const response = await api.get<ServerResources>(`/api/client/servers/${serverUUID}/resources`);
    return response.data;
  },
};
