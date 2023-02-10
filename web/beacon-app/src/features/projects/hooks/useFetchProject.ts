import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants/queryKeys';

import { projectRequest } from '../api/projectDetailApiService';
import { ProjectDetailQuery } from '../types/projectService';
export function useFetchProject(id: string): ProjectDetailQuery {
  const query = useQuery([RQK.PROJECT, id] as const, () => projectRequest(axiosInstance)(id), {
    enabled: !!id,
    refetchOnWindowFocus: false,
    refetchOnMount: true,
    // set stale time to 15 minutes
    // TODO: Change stale time
    staleTime: 1000 * 60 * 15,
  });

  return {
    hasProjectFailed: query.isError,
    isFetchingProject: query.isLoading,
    project: query.data as ProjectDetailQuery['project'],
    wasProjectFetched: query.isSuccess,
    error: query.error,
  };
}
