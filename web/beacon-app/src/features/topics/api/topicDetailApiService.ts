import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { Request } from '@/application/api/ApiService';
import { getValidApiResponse } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

import { Topic } from '../types/topicService';
export function topicRequest(request: Request): ApiAdapters['getTopic'] {
  return async (topicID: string) => {
    const response = (await request(`${APP_ROUTE.TOPICS}/${topicID}`, {
      method: 'GET',
    })) as any;

    return getValidApiResponse<Topic>(response);
  };
}
