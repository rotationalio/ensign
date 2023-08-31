import { Trans } from '@lingui/macro';
import { Button, Heading, TextField } from '@rotational/beacon-core';

const nameStep = () => {
  return (
    <>
      <Heading as="h1" className="text-lg font-bold">
        <Trans>Step 3 of 4</Trans>
      </Heading>
      <p className="mt-4 font-bold">
        <Trans>What's your name?</Trans>
      </p>
      <p className="my-4">
        <Trans>
          Adding your name will make it easier for your teammates to collaborate with you.
        </Trans>
      </p>
      <div className="max-w-6xl">
        <TextField
          fullWidth
          placeholder="Ex. Haley Smith"
          label="Name"
          labelClassName="sr-only"
          className="mb-4"
        />
      </div>
      <Button variant="secondary">
        <Trans>Next</Trans>
      </Button>
    </>
  );
};

export default nameStep;
