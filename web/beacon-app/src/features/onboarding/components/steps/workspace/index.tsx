import { Trans } from '@lingui/macro';

import StepCounter from '../StepCounter';
import WorkspaceForm from './form';

const WorkspaceStep = () => {
  const submitFormHandler = (values: any) => {
    console.log(values);
  };
  return (
    <>
      <StepCounter />
      <div className="flex flex-col justify-center ">
        <p className="text-base font-bold">
          <Trans>Now let’s create your workspace URL</Trans>
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
