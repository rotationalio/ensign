import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants';

import { projectsRequest } from '../api/projectListAPI';
import { ProjectsQuery } from '../types/projectService';

export function useFetchProjects(): ProjectsQuery {
  const query = useQuery([RQK.PROJECTS], projectsRequest(axiosInstance));

  return {
    getProjects: query.refetch,
    hasProjectsFailed: query.isError,
    isFetchingProjects: query.isLoading,
    projects: query.data,
    wasProjectsFetched: query.isSuccess,
    error: query.error,
  };
}
