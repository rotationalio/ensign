import { useMutation } from '@tanstack/react-query';

import axiosInstance from '@/application/api/ApiService';
import { queryClient } from '@/application/config/react-query';
import { RQK } from '@/constants';

import { onboardingStepAPI } from '../api/onboardingStepApi';
import { OnboardingMemberUpdateMutation } from '../types/onboardingServices';

export function useOnboarding(): OnboardingMemberUpdateMutation {
  const mutation = useMutation(onboardingStepAPI(axiosInstance), {
    retry: 0,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: [RQK.MEMBER_LIST] });
    },
  });
  return {
    updateMember: mutation.mutate,
    reset: mutation.reset,
    member: mutation.data as OnboardingMemberUpdateMutation['member'],
    hasMemberFailed: mutation.isError,
    wasMemberUpdated: mutation.isSuccess,
    isUpdatingMember: mutation.isLoading,
    error: mutation.error,
  };
}
