import { Trans } from '@lingui/macro';
import { TextField } from '@rotational/beacon-core';

import StepButtons from '../../StepButtons';
import StepCounter from '../StepCounter';

const nameStep = () => {
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
      <TextField
        fullWidth
        placeholder="Ex. Haley Smith"
        label="Name"
        labelClassName="sr-only"
        className="mb-4"
      />
      <StepButtons />
    </>
  );
};

export default nameStep;
