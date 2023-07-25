import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';
import React, { useState } from 'react';
import { SlArrowDown, SlArrowRight } from 'react-icons/sl';

const AdvancedTopicPolicy = () => {
  const [open, setOpen] = useState<boolean>(true);
  const toggleHandler = () => setOpen(!open);

  return (
    <div data-testid="topic-query-title" className="mt-10">
      <button className="mb-4 flex h-5 place-items-center gap-3" onClick={toggleHandler}>
        <Heading as="h1" className=" text-lg font-semibold">
          <Trans>Advanced Topic Policy Management</Trans>
        </Heading>
        {open ? <SlArrowDown /> : <SlArrowRight />}
      </button>

      {open && (
        <div className="flex space-x-1">
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