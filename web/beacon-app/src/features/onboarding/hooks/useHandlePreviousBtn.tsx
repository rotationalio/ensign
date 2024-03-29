import { useEffect } from 'react';

import { useFetchProfile } from '@/features/members/hooks/useFetchProfile';
import { useUpdateProfile } from '@/features/members/hooks/useUpdateProfile';
import { useOrgStore } from '@/store';

import { ONBOARDING_STEPS } from '../shared/constants';
import { getOnboardingStepsData } from '../shared/utils';
type Props = {
  values: any;
  currentStep: number;
  isInvited: boolean;
};

const useHandlePreviousBtn = () => {
  const state = useOrgStore((state: any) => state) as any;
  const { profile } = useFetchProfile();
  const { updateProfile, isUpdatingProfile, wasProfileUpdated, reset } = useUpdateProfile();

  const handlePrevious = ({ values, currentStep, isInvited }: Props) => {
    if (!currentStep || currentStep === ONBOARDING_STEPS.ORGANIZATION) return;
    const requestPayload = {
      payload: {
        ...getOnboardingStepsData(profile),
      },
    };
    if (currentStep === ONBOARDING_STEPS.WORKSPACE) {
      requestPayload.payload = {
        ...requestPayload.payload,
        workspace: values.workspace,
      };
    }

    if (currentStep === ONBOARDING_STEPS.PREFERENCE) {
      requestPayload.payload = {
        ...requestPayload.payload,
        developer_segment: values?.developer_segment?.map((item: any) => item.value),
        profession_segment: values?.profession_segment,
      };
    }

    if (currentStep === ONBOARDING_STEPS.NAME) {
      requestPayload.payload = {
        ...requestPayload.payload,
        name: values.name,
      };
    }

    if (currentStep === ONBOARDING_STEPS.PREFERENCE) {
      state.decrementStep();
      return;
    }

    if (isInvited && currentStep === ONBOARDING_STEPS.WORKSPACE) {
      state.decrementStep();
      return;
    }

    updateProfile(requestPayload);
  };

  useEffect(() => {
    if (wasProfileUpdated) {
      // console.log('[] wasProfileUpdated', wasProfileUpdated);
      state.resetTempData();
      state.decrementStep();
    }
  }, [wasProfileUpdated, state, reset]);

  return { handlePrevious, isLoading: isUpdatingProfile };
};

export default useHandlePreviousBtn;
