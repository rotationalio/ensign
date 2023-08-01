import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { Request } from '@/application/api/ApiService';
import { getValidApiResponse } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

import { TopicsResponse } from '../types/topicService';
export function topicsRequest(request: Request): ApiAdapters['getTopics'] {
  return async (projectID: string) => {
    const response = (await request(`${APP_ROUTE.PROJECTS}/${projectID}/topics`, {
      method: 'GET',
    })) as any;

    return getValidApiResponse<TopicsResponse>(response);
  };
}
