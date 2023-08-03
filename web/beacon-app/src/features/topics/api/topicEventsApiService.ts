import type { ApiAdapters } from '@/application/api/ApiAdapters';
import type { Request } from '@/application/api/ApiService';
import { getValidApiResponse } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

import { Topic } from '../types/topicService';
export function topicEventsRequest(request: Request): ApiAdapters['getTopicEvents'] {
  return async (topicID: string) => {
    const response = (await request(`${APP_ROUTE.TOPICS}/${topicID}/events`, {
      method: 'GET',
    })) as any;

    return getValidApiResponse<Topic>(response);
  };
}
