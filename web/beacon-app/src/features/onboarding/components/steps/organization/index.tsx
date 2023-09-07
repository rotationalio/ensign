import { Trans } from '@lingui/macro';
import { useEffect } from 'react';

import { useUpdateMember } from '@/features/members/hooks/useUpdateMember';
import useUserLoader from '@/features/members/loaders/userLoader';
import { getOnboardingStepsData } from '@/features/onboarding/shared/utils';
import { useOrgStore } from '@/store';

import StepCounter from '../StepCounter';
import OrganizationForm from './form';

const OrganizationStep = () => {
  const { member } = useUserLoader();
  const { updateMember, wasMemberUpdated, isUpdatingMember, error, reset } = useUpdateMember();
  const increaseStep = useOrgStore((state: any) => state.increaseStep) as any;

  // Display error if organization name is already taken.
  const hasError = error && error.response.status === 409;

  const submitFormHandler = (values: any) => {
    const payload = {
      memberID: member?.id,
      payload: {
        ...getOnboardingStepsData(member),
        organization: values.organization,
      },
    };
    updateMember(payload);
    increaseStep();
  };

  useEffect(() => {
    if (wasMemberUpdated) {
      reset();
      increaseStep();
    }
  }, [wasMemberUpdated, increaseStep, reset]);

  return (
    <>
      <StepCounter />
      <p className="mt-4 font-bold">
        <Trans>What's the name of your team or organization?</Trans>
      </p>
      <p className="my-4">
        <Trans>
          This will be the name of your workspace where you create projects and collaborate, so
          choose something you and your teammates will recognize.
        </Trans>
      </p>
      <OrganizationForm
        onSubmit={submitFormHandler}
        isSubmitting={isUpdatingMember}
        initialValues={{ organization: member?.organization }}
        hasError={hasError}
      />
    </>
  );
};

export default OrganizationStep;
