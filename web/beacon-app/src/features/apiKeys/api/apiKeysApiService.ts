import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { Request } from '@/application/api/ApiService';
import { getValidApiResponse } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

export function apiKeysRequest(request: Request): ApiAdapters['getApiKeys'] {
  return async (projectID: string) => {
    const response = (await request(`${APP_ROUTE.PROJECTS}/${projectID}/apikeys`, {
      method: 'GET',
    })) as any;

    return getValidApiResponse<any>(response);
  };
}
