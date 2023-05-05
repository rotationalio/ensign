import { AxiosResponse } from 'axios';

import { ApiAdapters } from '@/application/api/ApiAdapters';
import { getValidApiResponse, Request } from '@/application/api/ApiService';
import { APP_ROUTE } from '@/constants';
import { Topic } from '@/features/topics/types/topicService';

import { NewTopicDTO } from '../types/createTopicService';

export function createProjectTopic(request: Request): ApiAdapters['createProjectTopic'] {
  return async ({ projectID, topic_name }: NewTopicDTO) => {
    const response = (await request(`${APP_ROUTE.PROJECTS}/${projectID}/topics`, {
      method: 'POST',
      data: JSON.stringify({
        topic_name,
      }),
    })) as unknown as AxiosResponse;

    return getValidApiResponse<Topic>(response);
  };
}
