import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';
import React, { useState } from 'react';
import { SlArrowDown, SlArrowUp } from 'react-icons/sl';

const AdvancedTopicPolicy = () => {
  const [open, setOpen] = useState<boolean>(true);
  const toggleHandler = () => setOpen(!open);

  return (
    <div data-testid="topic-query-title" className="mt-10" data-cy="topic-mgmt">
      <button
        className="mb-4 flex h-5 place-items-center gap-3"
        onClick={toggleHandler}
        data-cy="topic-mgmt-heading"
      >
        <Heading as="h1" className=" text-lg font-semibold">
          <Trans>Advanced Topic Policy Management</Trans>
        </Heading>
        {open ? (
          <SlArrowDown data-cy="topic-mgmt-carat-down" />
        ) : (
          <SlArrowUp data-cy="topic-mgmt-carat-up" />
        )}
      </button>

      {open && (
        <div className="flex space-x-1" data-cy="topic-mgmt-content">
          <p>
            <Trans>
              Coming soon! Customize and manage topic policies. The topic must be in the “Ready”
              state to be edited.
            </Trans>
          </p>
        </div>
      )}
    </div>
  );
};

export default AdvancedTopicPolicy;
