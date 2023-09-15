import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';
import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

import { PATH_DASHBOARD } from '@/application';
import AppLayout from '@/components/layout/AppLayout';
import { useFetchProfile } from '@/features/members/hooks/useFetchProfile';
import { useOrgStore } from '@/store';

import Step from '../components/Step';
import { ONBOARDING_STATUS, ONBOARDING_STEPS } from '../shared/constants';
import { isInvitedUser } from '../shared/utils';

const OnboardingPage = () => {
  const { profile: userProfile } = useFetchProfile();
  const isInvited = isInvitedUser(userProfile);
  const Store = useOrgStore.getState() as any;
  const { currentStep } = Store?.onboarding || null;
  const navigate = useNavigate();

  useEffect(() => {
    // if the user is not onboarded, we need to set the onboarding step
    if (userProfile && !currentStep) {
      if (isInvited) {
        // set to 3 if the user is invited since step 1 and 2 are already done
        Store?.setOnboardingStep(ONBOARDING_STEPS.NAME);
      } else {
        // const step = getCurrentStepFromMember(userProfile);
        Store.setOnboardingStep(ONBOARDING_STEPS.ORGANIZATION);
      }
    }
  }, [userProfile, currentStep, Store, isInvited]);

  // if onboarding status change then redirect to home page

  useEffect(() => {
    if (userProfile?.onboarding_status === ONBOARDING_STATUS.ACTIVE) {
      Store.resetOnboarding();
      navigate(PATH_DASHBOARD.HOME);
    }
  }, [userProfile, navigate, Store]);

  return (
    <AppLayout>
      <Heading as="h1" className=" m-10 mt-20 px-4  text-xl font-bold xl:mt-20 xl:px-28">
        <Trans>
          Welcome to Ensign! Please complete our brief onboarding survey to get started.
        </Trans>
      </Heading>
      <Step />
    </AppLayout>
  );
};

export default OnboardingPage;
