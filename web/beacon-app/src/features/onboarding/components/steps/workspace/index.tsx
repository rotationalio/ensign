import { Trans } from '@lingui/macro';

import { useOrgStore } from '@/store';

import StepCounter from '../StepCounter';
import WorkspaceForm from './form';
const WorkspaceStep = () => {
  const increaseStep = useOrgStore((state: any) => state.increaseStep) as any;
  const submitFormHandler = (values: any) => {
    console.log(values);
    increaseStep();
  };
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

        <WorkspaceForm onSubmit={submitFormHandler} />
      </div>
    </>
  );
};

export default WorkspaceStep;
