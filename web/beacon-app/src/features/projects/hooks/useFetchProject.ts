import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants/queryKeys';

import { projectRequest } from '../api/projectDetailApiService';
import { ProjectDetailQuery } from '../types/projectService';
export function useFetchProject(projectID: string): ProjectDetailQuery {
  const query = useQuery([RQK.PROJECT, projectID], () => projectRequest(axiosInstance)(projectID));

  return {
    hasProjectFailed: query.isError,
    isFetchingProject: query.isLoading,
    project: query.data,
    wasProjectFetched: query.isSuccess,
    error: query.error,
  };
}
