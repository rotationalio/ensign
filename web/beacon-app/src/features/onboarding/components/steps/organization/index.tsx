import { Trans } from '@lingui/macro';
import { Button, Heading, TextField } from '@rotational/beacon-core';

const OrganizationStep = () => {
  return (
    <>
      <Heading as="h1" className="text-lg font-bold">
        Step 1 of 4
      </Heading>
      <p className="mt-4 font-bold">
        <Trans>What's the name of your team or organization?</Trans>
      </p>
      <p className="my-4">
        <Trans>
          This will be the name of your workspace where you create projects and collaborate, so
          choose something you and your teammates will recognize.
        </Trans>
      </p>
      <div className="max-w-6xl">
        <TextField
          fullWidth
          placeholder="Ex. Rotational Labs"
          label="Team or Organization Name"
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

export default OrganizationStep;
