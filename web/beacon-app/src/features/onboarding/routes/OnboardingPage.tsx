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

  useEffect(() => {
    // if the user is not onboarded, we need to set the onboarding step
    if (member && !currentStep) {
      const step = getCurrentStepFromMember(member);
      orgDataState.setOnboardingStep(step);
    }
  }, [member, currentStep, orgDataState]);

  return (
    <AppLayout>
      <Step />
    </AppLayout>
  );
};

export default OnboardingPage;
