import { Trans } from '@lingui/macro';
import { useEffect } from 'react';

import useUserLoader from '@/features/members/loaders/userLoader';
import { useUpdateMember } from '@/features/onboarding/hooks/useUpdateMember';
import { useOrgStore } from '@/store';

import StepCounter from '../StepCounter';
import WorkspaceForm from './form';
const WorkspaceStep = () => {
  const increaseStep = useOrgStore((state: any) => state.increaseStep) as any;
  const { member } = useUserLoader();
  const { updateMember, wasMemberUpdated, isUpdatingMember, reset } = useUpdateMember();

  const submitFormHandler = (values: any) => {
    const { organization, name, profession_segment } = member;
    const requestPayload = {
      memberID: member?.id,
      payload: {
        organization,
        name,
        profession_segment,
        workspace: values.workspace,
      },
    };
    console.log(requestPayload);
    updateMember(requestPayload);
  };

  useEffect(() => {
    if (wasMemberUpdated) {
      increaseStep();
      reset();
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
          initialValues={{
            workspace: member?.workspace, // we may need to remove rotational.app from the name
          }}
        />
      </div>
    </>
  );
};

export default WorkspaceStep;
