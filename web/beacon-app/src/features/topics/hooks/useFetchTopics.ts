import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants/queryKeys';

import { topicsRequest } from '../api/topicsApiService';
import { TopicsQuery } from '../types/topicService';

export function useFetchTopics(projectID: string): TopicsQuery {
  const query = useQuery(
    [RQK.TOPICS, projectID] as const,
    () => topicsRequest(axiosInstance)(projectID),
    {
      enabled: !!projectID,
      refetchOnWindowFocus: false,
      refetchOnMount: true,
      // set stale time to 15 minutes
    }
  );

  return {
    getTopics: query.refetch,
    hasTopicsFailed: query.isError,
    isFetchingTopics: query.isLoading,
    topics: query.data as TopicsQuery['topics'],
    wasTopicsFetched: query.isSuccess,
    error: query.error,
  };
}
