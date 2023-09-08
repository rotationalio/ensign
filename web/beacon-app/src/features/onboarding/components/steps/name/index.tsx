import { Trans } from '@lingui/macro';
import { useEffect } from 'react';

import { useUpdateMember } from '@/features/members/hooks/useUpdateMember';
import useUserLoader from '@/features/members/loaders/userLoader';
import { getOnboardingStepsData } from '@/features/onboarding/shared/utils';
import { useOrgStore } from '@/store';

import StepCounter from '../StepCounter';
import NameForm from './form';
const NameStep = () => {
  const { member } = useUserLoader();
  const { updateMember, wasMemberUpdated, isUpdatingMember, reset } = useUpdateMember();
  const increaseStep = useOrgStore((state: any) => state.increaseStep) as any;
  const submitFormHandler = (values: any) => {
    const payload = {
      memberID: member?.id,
      payload: {
        ...getOnboardingStepsData(member),
        name: values.name,
      },
    };
    updateMember(payload);
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
        <Trans>What's your name?</Trans>
      </p>
      <p className="my-4">
        <Trans>
          Adding your name will make it easier for your teammates to collaborate with you.
        </Trans>
      </p>
      <NameForm
        onSubmit={submitFormHandler}
        isSubmitting={isUpdatingMember}
        initialValues={{
          name: member?.name,
        }}
      />
    </>
  );
};

export default NameStep;
