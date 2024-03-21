import * as Sentry from '@sentry/react';
import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants';

import { projectsRequest } from '../api/projectListAPI';
import { ProjectsQuery } from '../types/projectService';

export function useFetchTenantProjects(tenantID: string): ProjectsQuery {
  const query = useQuery([RQK.PROJECTS, tenantID], () => projectsRequest(axiosInstance)(tenantID), {
    enabled: !!tenantID,
    onError: (error: any) => {
      // stop logging 401 & 403 errors to sentry
      if (error?.response?.status !== 401 && error?.response?.status !== 403) {
        Sentry.captureException(error);
      }
    },
  });

  return {
    getProjects: query.refetch,
    hasProjectsFailed: query.isError,
    isFetchingProjects: query.isLoading,
    projects: query.data as ProjectsQuery['projects'],
    wasProjectsFetched: query.isSuccess,
    error: query.error,
  };
}
