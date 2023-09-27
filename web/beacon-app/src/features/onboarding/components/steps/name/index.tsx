import { Trans } from '@lingui/macro';
import { useEffect } from 'react';

import { useFetchProfile } from '@/features/members/hooks/useFetchProfile';
import { useUpdateProfile } from '@/features/members/hooks/useUpdateProfile';
import { getOnboardingStepsData } from '@/features/onboarding/shared/utils';
import { useOrgStore } from '@/store';

import StepCounter from '../StepCounter';
import NameForm from './form';
const NameStep = () => {
  const { profile: userProfile } = useFetchProfile();
  const { updateProfile, wasProfileUpdated, isUpdatingProfile, reset } = useUpdateProfile();
  const increaseStep = useOrgStore((state: any) => state.increaseStep) as any;
  const submitFormHandler = (values: any) => {
    const payload = {
      memberID: userProfile?.id,
      payload: {
        ...getOnboardingStepsData(userProfile),
        name: values.name,
      },
    };
    updateProfile(payload);
  };

  useEffect(() => {
    if (wasProfileUpdated) {
      reset();
      increaseStep();
    }
  }, [wasProfileUpdated, increaseStep, reset]);
  return (
    <>
      <StepCounter />
      <p className="mt-4 font-bold">
        <Trans>What's your name?</Trans>
      </p>
      <p className="my-4">
        <Trans>
          Adding your name will make it easier for your teammates to collaborate with you.
        </Trans>
      </p>
      <div className="w-full">
        <NameForm
          onSubmit={submitFormHandler}
          isSubmitting={isUpdatingProfile}
          initialValues={{
            name: userProfile?.name,
          }}
        />
      </div>
    </>
  );
};

export default NameStep;
