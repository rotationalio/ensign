import * as Sentry from '@sentry/react';
import { useQuery } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { RQK } from '@/constants/queryKeys';

import { projectRequest } from '../api/projectDetailApiService';
import { ProjectDetailQuery } from '../types/projectService';
export function useFetchProject(projectID: string): ProjectDetailQuery {
  const query = useQuery([RQK.PROJECT, projectID], () => projectRequest(axiosInstance)(projectID), {
    enabled: !!projectID,
    onError: (error: any) => {
      // stop logging 401 & 403 errors to sentry
      if (error.response.status !== 401 && error.response.status !== 403) {
        Sentry.captureException(error);
      }
    },
  });

  return {
    hasProjectFailed: query.isError,
    isFetchingProject: query.isLoading,
    project: query.data,
    wasProjectFetched: query.isSuccess,
    error: query.error,
  };
}
