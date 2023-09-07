import { useEffect } from 'react';

import AppLayout from '@/components/layout/AppLayout';
import useUserLoader from '@/features/members/loaders/userLoader';
import { useOrgStore } from '@/store';

import Step from '../components/Step';
import { ONBOARDING_STEPS } from '../shared/constants';
import { getCurrentStepFromMember, isInvitedUser } from '../shared/utils';

const OnboardingPage = () => {
  const { member } = useUserLoader();
  const isInvited = isInvitedUser(member);
  const orgDataState = useOrgStore.getState() as any;
  const { currentStep } = orgDataState?.onboarding || null;

  useEffect(() => {
    // if the user is not onboarded, we need to set the onboarding step
    if (member && !currentStep) {
      if (isInvited) {
        // set to 3 if the user is invited since step 1 and 2 are already done
        orgDataState?.setOnboardingStep(ONBOARDING_STEPS.NAME);
      } else {
        const step = getCurrentStepFromMember(member);
        orgDataState.setOnboardingStep(step);
      }
    }
  }, [member, currentStep, orgDataState, isInvited]);

  return (
    <AppLayout>
      <Step />
    </AppLayout>
  );
};

export default OnboardingPage;
