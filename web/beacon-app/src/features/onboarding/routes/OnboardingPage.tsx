import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application';
import AppLayout from '@/components/layout/AppLayout';
import { useFetchProfile } from '@/features/members/hooks/useFetchProfile';
import { useOrgStore } from '@/store';

import Step from '../components/Step';
import { ONBOARDING_STATUS, ONBOARDING_STEPS } from '../shared/constants';
import { getCurrentStepFromMember, isInvitedUser } from '../shared/utils';

const OnboardingPage = () => {
  const { profile: userProfile } = useFetchProfile();
  const isInvited = isInvitedUser(userProfile);
  const orgDataState = useOrgStore.getState() as any;
  const { currentStep } = orgDataState?.onboarding || null;
  const navigate = useNavigate();

  useEffect(() => {
    // if the user is not onboarded, we need to set the onboarding step
    if (userProfile && !currentStep) {
      if (isInvited) {
        // set to 3 if the user is invited since step 1 and 2 are already done
        orgDataState?.setOnboardingStep(ONBOARDING_STEPS.NAME);
      } else {
        const step = getCurrentStepFromMember(userProfile);
        orgDataState.setOnboardingStep(step);
      }
    }
  }, [userProfile, currentStep, orgDataState, isInvited]);

  // if onboarding status change then redirect to home page

  useEffect(() => {
    if (userProfile?.onboarding_status === ONBOARDING_STATUS.ACTIVE) {
      navigate(PATH_DASHBOARD.HOME);
    }
  }, [userProfile, navigate]);

  return (
    <AppLayout>
      <Step />
    </AppLayout>
  );
};

export default OnboardingPage;
