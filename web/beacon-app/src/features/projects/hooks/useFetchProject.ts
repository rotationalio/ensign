import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants/queryKeys';

import { projectRequest } from '../api/projectDetailApiService';
import { ProjectDetailQuery } from '../types/projectService';
export function useFetchProject({ id }: any): ProjectDetailQuery {
  const query = useQuery([RQK.PROJECT, id] as const, () => projectRequest(axiosInstance)(id), {
    enabled: !!id,
  });

  return {
    hasProjectFailed: query.isError,
    isFetchingProject: query.isLoading,
    project: query.data as ProjectDetailQuery['project'],
    wasProjectFetched: query.isSuccess,
    error: query.error,
  };
}
