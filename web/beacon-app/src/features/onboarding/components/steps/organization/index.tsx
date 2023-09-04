import { Trans } from '@lingui/macro';

import StepCounter from '../StepCounter';
import OrganizationForm from './form';

const OrganizationStep = () => {
  const handleSubmitOrganizationForm = (values: any) => {
    console.log(values);
  };
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
      <OrganizationForm onSubmit={handleSubmitOrganizationForm} />
    </>
  );
};

export default OrganizationStep;
