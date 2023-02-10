import invariant from 'invariant';

import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { Request } from '@/application/api/ApiService';
import { getValidApiError, getValidApiResponse } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';
import type { Topics } from '@/features/topics/types/topicService';
export function topicsRequest(request: Request): ApiAdapters['getTopics'] {
  return async (projectID: string) => {
    invariant(projectID, 'projectID is required');

    try {
      const response = (await request(`${APP_ROUTE.PROJECTS}/${projectID}/topics`, {
        method: 'GET',
      })) as any;

      return getValidApiResponse<Topics>(response);
    } catch (error: any) {
      getValidApiError(error);
    }
  };
}
