import { useMutation } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { queryClient } from '@/application/config/react-query';
import { RQK } from '@/constants';

import { createProjectAPI } from '../api/createProjectAPI';
import { ProjectMutation } from '../types/createProjectService';

export function useCreateProject(): ProjectMutation {
  const mutation = useMutation(createProjectAPI(axiosInstance), {
    retry: 0,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [RQK.PROJECTS] });
      queryClient.invalidateQueries({ queryKey: [RQK.QUICK_VIEW] });
      queryClient.invalidateQueries({ queryKey: [RQK.PROJECT_QUICK_VIEW] });
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
