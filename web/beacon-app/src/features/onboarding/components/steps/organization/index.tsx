import { Trans } from '@lingui/macro';
import { TextField } from '@rotational/beacon-core';

import StepButtons from '../../StepButtons';
import StepCounter from '../StepCounter';

const OrganizationStep = () => {
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
      <TextField
        fullWidth
        placeholder="Ex. Rotational Labs"
        label="Team or Organization Name"
        labelClassName="sr-only"
        className="mb-4"
      />
      <StepButtons />
    </>
  );
};

export default OrganizationStep;
