import { Trans } from '@lingui/macro';
import { useEffect } from 'react';

import { useFetchProfile } from '@/features/members/hooks/useFetchProfile';
import { useUpdateProfile } from '@/features/members/hooks/useUpdateProfile';
import { getOnboardingStepsData, isInvitedUser } from '@/features/onboarding/shared/utils';
import { useOrgStore } from '@/store';

import StepCounter from '../StepCounter';
import OrganizationForm from './form';

const OrganizationStep = () => {
  const { profile } = useFetchProfile();
  const { updateProfile, wasProfileUpdated, isUpdatingProfile, error, reset } = useUpdateProfile();

  const state = useOrgStore((state: any) => state) as any;

  // Display error if organization name is already taken.
  const hasError = error && error.response.status === 409;
  const isInvited = isInvitedUser(profile);
  const submitFormHandler = (values: any) => {
    if (isInvited) {
      state.increaseStep();
      return;
    }
    const payload = {
      memberID: profile?.id,
      payload: {
        ...getOnboardingStepsData(profile),
        organization: values.organization,
      },
    };
    updateProfile(payload);
  };

  useEffect(() => {
    if (wasProfileUpdated) {
      state.setTempData({ tempData: { organization: profile?.organization } });
      reset();
      state.increaseStep();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [wasProfileUpdated, state.increaseStep]);

  return (
    <>
      <StepCounter />
      <p className="mt-4 font-bold">
        <Trans>What's the name of your team?</Trans>
      </p>
      <p className="my-4">
        <Trans>
          This will be the name of your workspace where you create projects and collaborate, so
          choose something you and your teammates will recognize.
        </Trans>
      </p>

      <OrganizationForm
        onSubmit={submitFormHandler}
        shouldDisableInput={isInvited}
        isSubmitting={isUpdatingProfile}
        initialValues={{ organization: profile?.organization }}
        hasError={hasError}
      />
    </>
  );
};

export default OrganizationStep;
