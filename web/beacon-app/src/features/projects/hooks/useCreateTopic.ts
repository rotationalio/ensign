import { useMutation } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { queryClient } from '@/application/config/react-query';
import { RQK } from '@/constants';
import { createProjectTopic } from '@/features/projects/api/createTopicApiService';
import { TopicMutation } from '@/features/projects/types/createTopicService';

export function useCreateTopic(): TopicMutation {
  const mutation = useMutation(createProjectTopic(axiosInstance), {
    retry: 0,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [RQK.TOPICS] });
      queryClient.invalidateQueries({ queryKey: [RQK.QUICK_VIEW] });
      queryClient.invalidateQueries({ queryKey: [RQK.PROJECT_QUICK_VIEW] });
    },
  });
  return {
    createTopic: mutation.mutate,
    reset: mutation.reset,
    topic: mutation.data as TopicMutation['topic'],
    hasTopicFailed: mutation.isError,
    wasTopicCreated: mutation.isSuccess,
    isCreatingTopic: mutation.isLoading,
    error: mutation.error,
  };
}
