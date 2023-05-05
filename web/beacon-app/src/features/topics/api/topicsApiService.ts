import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { Request } from '@/application/api/ApiService';
import { getValidApiError, getValidApiResponse } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';
export function topicsRequest(request: Request): ApiAdapters['getTopics'] {
  return async (projectID: string) => {
    try {
      const response = (await request(`${APP_ROUTE.PROJECTS}/${projectID}/topics`, {
        method: 'GET',
      })) as any;

      return getValidApiResponse<any>(response);
    } catch (error: any) {
      getValidApiError(error);
    }
  };
}
