import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants';

import { projectsRequest } from '../api/projectListAPI';
import { ProjectsQuery } from '../types/projectService';

export function useFetchProjects(): ProjectsQuery {
  const query = useQuery([RQK.PROJECT_LIST], projectsRequest(axiosInstance), {
    refetchOnWindowFocus: false,
    refetchOnMount: true,
    // set stale time to 15 minutes
    // TODO: Change stale time
    staleTime: 1000 * 60 * 15,
  });

  return {
    getProjects: query.refetch,
    hasProjectsFailed: query.isError,
    isFetchingProjects: query.isLoading,
    projects: query.data,
    wasProjectsFetched: query.isSuccess,
    error: query.error,
  };
}
