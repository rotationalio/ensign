import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { Request } from '@/application/api/ApiService';
import { getValidApiError, getValidApiResponse } from '@/application/api/ApiService';

const permissionsRequest =
  (request: Request): ApiAdapters['getPermissions'] =>
  async () => {
    const link = `/apikeys/permissions`;
    try {
      const response = (await request(`${link}`, {
        method: 'GET',
      })) as any;

      return getValidApiResponse<any>(response);
    } catch (e: any) {
      getValidApiError(e);
    }
  };

export default permissionsRequest;
