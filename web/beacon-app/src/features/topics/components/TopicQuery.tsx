import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';
import React, { useState } from 'react';

import { ChevronDown } from '@/components/icons/chevron-down';

const TopicQuery = () => {
  const [open, setOpen] = useState<boolean>(true);
  const toggleHandler = () => setOpen(!open);

  return (
    <div data-testid="topic-query-title" className="mt-10">
      <button className="mb-4 flex h-5 place-items-center gap-2" onClick={toggleHandler}>
        <Heading as="h1" className=" text-lg font-semibold">
          <Trans>Topic Query</Trans>
        </Heading>
        <ChevronDown />
      </button>

      {open && (
        <div className="flex space-x-1">
          <Trans>
            <p>
              Coming soon! Query the topic for insights with{' '}
              <span>
                <a
                  href="#"
                  target="_blank"
                  rel="noreferrer"
                  className="font-bold text-[#1F4CED] underline hover:!underline"
                >
                  EnSQL
                </a>
              </span>{' '}
              e.g. the latest event or last 10 events.
            </p>
          </Trans>
        </div>
      )}
    </div>
  );
};

export default TopicQuery;
