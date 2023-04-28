import { useMutation } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { queryClient } from '@/application/config/react-query';
import { RQK } from '@/constants';

import { createProjectAPI } from '../api/createProjectAPI';
import { ProjectMutation } from '../types/createProjectService';
export function useCreateProjectAPIKey(): ProjectMutation {
  const mutation = useMutation(createProjectAPI(axiosInstance), {
    retry: 0,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [RQK.PROJECTS] });
    },
  });
  return {
    createNewProject: mutation.mutate,
    reset: mutation.reset,
    project: mutation.data as ProjectMutation['project'],
    hasProjectFailed: mutation.isError,
    wasProjectCreated: mutation.isSuccess,
    isCreatingProject: mutation.isLoading,
    error: mutation.error,
  };
}
