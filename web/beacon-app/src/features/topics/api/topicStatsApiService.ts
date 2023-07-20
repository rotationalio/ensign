import invariant from 'invariant';

import { ApiAdapters } from '@/application/api/ApiAdapters';
import { getValidApiResponse, Request } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';

import { TopicQuickViewResponse } from '../types/topicService';

function topicStatsApiRequest(request: Request): ApiAdapters['getTopicStats'] {
  return async (topicID: string) => {
    invariant(topicID, 'topic id is required');
    const response = (await request(`${APP_ROUTE.TOPICS}/${topicID}/stats`, {
      method: 'GET',
    })) as any;
    return getValidApiResponse<TopicQuickViewResponse>(response);
  };
}

export default topicStatsApiRequest;
