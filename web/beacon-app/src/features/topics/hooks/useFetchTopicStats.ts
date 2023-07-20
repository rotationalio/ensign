import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants/queryKeys';

import topicStatsApiRequest from '../api/topicStatsApiService';
import { TopicStatsQuery } from '../types/topicService';

function useFetchTopicStats(topicID: string): TopicStatsQuery {
  const query = useQuery(
    [RQK.TOPIC_STATS, topicID],
    () => topicStatsApiRequest(axiosInstance)(topicID),
    {
      enabled: !!topicID,
    }
  );

  return {
    getTopicStats: query.refetch,
    hasTopicStatsFailed: query.isError,
    isFetchingTopicStats: query.isLoading,
    topicStats: query.data,
    wasTopicStatsFetched: query.isSuccess,
    error: query.error,
  };
}

export default useFetchTopicStats;
