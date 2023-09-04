import * as Sentry from '@sentry/react';
import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants/queryKeys';

import { topicEventsRequest } from '../api/topicEventsApiService';
import { type TopicEvents, TopicEventsQuery } from '../types/topicEventsService';

export function useFetchTopicEvents(topicID: string): TopicEventsQuery {
  const eventID = `events-${topicID}`; // we already have a query key for topic, so we need to make a new one for events
  const query = useQuery([RQK.TOPIC, eventID], () => topicEventsRequest(axiosInstance)(topicID), {
    enabled: !!eventID,
    onError: (error) => {
      Sentry.captureException(error);
    },
  });

  return {
    getTopicEvents: query.refetch,
    hasTopicEventsFailed: query.isError,
    isFetchingTopicEvents: query.isLoading,
    topicEvents: query.data as TopicEvents[],
    wasTopicEventsFetched: query.isSuccess,
    error: query.error,
  };
}
