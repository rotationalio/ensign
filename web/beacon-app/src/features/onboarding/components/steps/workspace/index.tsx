import { Trans } from '@lingui/macro';
import { useEffect } from 'react';

import { useUpdateMember } from '@/features/members/hooks/useUpdateMember';
import useUserLoader from '@/features/members/loaders/userLoader';
import { useOrgStore } from '@/store';

import { getOnboardingStepsData } from '../../../shared/utils';
import StepCounter from '../StepCounter';
import WorkspaceForm from './form';
const WorkspaceStep = () => {
  const increaseStep = useOrgStore((state: any) => state.increaseStep) as any;
  const { member } = useUserLoader();
  const { updateMember, wasMemberUpdated, isUpdatingMember, reset, error } = useUpdateMember();

  const hasError = error && error.response.status === 400; // this means the workspace is already taken by another user

  const submitFormHandler = (values: any) => {
    const requestPayload = {
      memberID: member?.id,
      payload: {
        ...getOnboardingStepsData(member),
        workspace: values.workspace,
      },
    };
    console.log(requestPayload);
    updateMember(requestPayload);
    increaseStep();
  };

  // move to next step if member was updated
  useEffect(() => {
    if (wasMemberUpdated) {
      reset();
      increaseStep();
    }
  }, [wasMemberUpdated, increaseStep, reset]);

  return (
    <>
      <StepCounter />
      <div className="flex flex-col justify-center ">
        <p className="text-base font-bold">
          <Trans>Now let’s create your workspace URL.</Trans>
        </p>
        <p className="pt-3 text-base">
          <Trans>
            Your workspace URL should be unique, short, and recognizable. We suggest using the slug
            we created for you, but you can change it now because you can’t change it later. It must
            be letters, numbers or dashes only.
          </Trans>
        </p>

        <WorkspaceForm
          onSubmit={submitFormHandler}
          isSubmitting={isUpdatingMember}
          hasError={hasError}
          initialValues={{
            workspace: member?.workspace,
          }}
        />
      </div>
    </>
  );
};

export default WorkspaceStep;
