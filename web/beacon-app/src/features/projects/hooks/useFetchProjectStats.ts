import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants/queryKeys';

import { projectStatsApiRequest } from '../api/projectStatsApiService';
import { ProjectDetailQuery } from '../types/projectService';
export function useFetchProjectStats(projectID: string): ProjectDetailQuery {
  const query = useQuery(
    [RQK.PROJECT, projectID],
    () => projectStatsApiRequest(axiosInstance)(projectID),
    {
      enabled: !!projectID,
    }
  );

  return {
    hasProjectFailed: query.isError,
    isFetchingProject: query.isLoading,
    project: query.data,
    wasProjectFetched: query.isSuccess,
    error: query.error,
  };
}
