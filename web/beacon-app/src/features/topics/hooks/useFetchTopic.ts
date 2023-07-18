import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants/queryKeys';

import { topicRequest } from '../api/topicDetailApiService';
import { TopicQuery } from '../types/topicService';

export function useFetchTopic(topicID: string): TopicQuery {
  const query = useQuery([RQK.TOPIC, topicID], () => topicRequest(axiosInstance)(topicID), {
    enabled: !!topicID,
  });

  return {
    getTopic: query.refetch,
    hasTopicFailed: query.isError,
    isFetchingTopic: query.isLoading,
    topic: query.data,
    wasTopicFetched: query.isSuccess,
    error: query.error,
  };
}
