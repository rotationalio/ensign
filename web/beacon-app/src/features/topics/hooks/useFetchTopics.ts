import * as Sentry from '@sentry/react';
import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants/queryKeys';

import { topicsRequest } from '../api/topicsApiService';
import { TopicsQuery } from '../types/topicService';

export function useFetchTopics(projectID: string): TopicsQuery {
  const query = useQuery([RQK.TOPICS, projectID], () => topicsRequest(axiosInstance)(projectID), {
    enabled: !!projectID,
    onError: (error: any) => {
      // stop logging 401 & 403 errors to sentry
      if (error?.response?.status !== 401 && error?.response?.status !== 403) {
        Sentry.captureException(error);
      }
    },
  });

  return {
    getTopics: query.refetch,
    hasTopicsFailed: query.isError,
    isFetchingTopics: query.isLoading,
    topics: query.data as any,
    wasTopicsFetched: query.isSuccess,
    error: query.error,
  };
}
