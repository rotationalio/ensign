import { useMutation } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { queryClient } from '@/application/config/react-query';
import { RQK } from '@/constants';

import { createProjectAPIKey } from '../api/createApiKey';
import { APIKeyMutation } from '../types/createApiKeyService';
export function useCreateProjectAPIKey(): APIKeyMutation {
  const mutation = useMutation(createProjectAPIKey(axiosInstance), {
    retry: 0,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [RQK.PROJECT_API_KEYS] });
      queryClient.invalidateQueries({ queryKey: [RQK.QUICK_VIEW] });
      queryClient.invalidateQueries({ queryKey: [RQK.API_KEYS] });
      queryClient.invalidateQueries({ queryKey: [RQK.PROJECT_QUICK_VIEW] });
      queryClient.invalidateQueries({ queryKey: [RQK.PROJECTS] });
      queryClient.invalidateQueries({ queryKey: [RQK.PROJECT] });
    },
  });
  return {
    createProjectNewKey: mutation.mutate,
    reset: mutation.reset,
    key: mutation.data as APIKeyMutation['key'],
    hasKeyFailed: mutation.isError,
    wasKeyCreated: mutation.isSuccess,
    isCreatingKey: mutation.isLoading,
    error: mutation.error,
  };
}
