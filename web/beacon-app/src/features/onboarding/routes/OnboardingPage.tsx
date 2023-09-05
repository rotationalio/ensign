import { useEffect } from 'react';

import AppLayout from '@/components/layout/AppLayout';
import useUserLoader from '@/features/members/loaders/userLoader';
import { useOrgStore } from '@/store';

import Step from '../components/Step';
import { getCurrentStepFromMember } from '../shared/utils';

const OnboardingPage = () => {
  const { member } = useUserLoader();
  const orgDataState = useOrgStore.getState() as any;
  const { currentStep } = orgDataState?.onboarding || null;

  // set the current step

  useEffect(() => {
    if (member && !currentStep) {
      const step = getCurrentStepFromMember(member);
      orgDataState.setOnboardingStep(step);
    }
  }, [member, currentStep, orgDataState]);

  // if current is null set to 1

  useEffect(() => {
    if (!currentStep) {
      orgDataState.setOnboardingStep(1);
    }
  }, [currentStep, orgDataState]);

  return (
    <AppLayout>
      <Step />
    </AppLayout>
  );
};

export default OnboardingPage;
