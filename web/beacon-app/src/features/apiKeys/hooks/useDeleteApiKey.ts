import { useMutation } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { queryClient } from '@/application/config/react-query';
import { RQK } from '@/constants';

import { deleteAPIKeyRequest } from '../api/deleteApiKeyApi';
import { DeleteAPIKeyMutation } from '../types/deleteApiKeyService';
export function useDeleteAPIKey(apiKey: string): DeleteAPIKeyMutation {
  const mutation = useMutation(() => deleteAPIKeyRequest(axiosInstance)(apiKey), {
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
    deleteApiKey: mutation.mutate,
    reset: mutation.reset,
    key: mutation.data as DeleteAPIKeyMutation['key'],
    hasKeyDeletedFailed: mutation.isError,
    wasKeyDeleted: mutation.isSuccess,
    isDeletingKey: mutation.isLoading,
    error: mutation.error,
  };
}
