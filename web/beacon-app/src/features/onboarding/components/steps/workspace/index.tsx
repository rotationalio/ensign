import { Trans } from '@lingui/macro';
import { useEffect } from 'react';

import { useFetchProfile } from '@/features/members/hooks/useFetchProfile';
import { useUpdateProfile } from '@/features/members/hooks/useUpdateProfile';
import { useOrgStore } from '@/store';

import { getOnboardingStepsData, isInvitedUser } from '../../../shared/utils';
import StepCounter from '../StepCounter';
import WorkspaceForm from './form';
const WorkspaceStep = () => {
  const increaseStep = useOrgStore((state: any) => state.increaseStep) as any;
  const { profile } = useFetchProfile();
  const isInvited = isInvitedUser(profile);
  const { updateProfile, wasProfileUpdated, isUpdatingProfile, reset, error } = useUpdateProfile();

  const hasError = error && error.response.status === 400; // this means the workspace is already taken by another user

  const submitFormHandler = (values: any) => {
    if (isInvited) {
      increaseStep();
      return;
    }
    const requestPayload = {
      memberID: profile?.id,
      payload: {
        ...getOnboardingStepsData(profile),
        workspace: values.workspace,
      },
    };
    console.log(requestPayload);
    updateProfile(requestPayload);
  };

  // move to next step if member was updated
  useEffect(() => {
    if (wasProfileUpdated) {
      reset();
      increaseStep();
    }
  }, [wasProfileUpdated, increaseStep, reset]);

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
          isSubmitting={isUpdatingProfile}
          shouldDisableInput={isInvited}
          hasError={hasError}
          initialValues={{
            workspace: profile?.workspace,
          }}
        />
      </div>
    </>
  );
};

export default WorkspaceStep;
