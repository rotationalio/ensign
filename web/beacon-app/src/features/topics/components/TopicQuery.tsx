import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';
// import { useAnimate, useInView } from 'framer-motion';
import React, { useState } from 'react';
import { SlArrowDown, SlArrowRight } from 'react-icons/sl';

import TopicQueryResult from './TopicQueryResult';
type TopicNameProps = {
  name: string;
};

import QueryInput from './QueryInput';
const TopicQuery = ({ name }: TopicNameProps) => {
  const [open, setOpen] = useState<boolean>(true);

  const toggleHandler = () => setOpen(!open);

  return (
    <div data-testid="topic-query-title" className="mt-10">
      <button className="mb-4 flex h-5 place-items-center gap-3" onClick={toggleHandler}>
        <Heading as="h1" className=" text-lg font-semibold">
          <Trans>Topic Query</Trans>
        </Heading>
        {open ? <SlArrowDown /> : <SlArrowRight />}
      </button>

      {open && (
        <>
          <div className="flex space-x-1">
            <p>
              <Trans>
                Coming soon! Query the topic for insights with EnSQL e.g. the latest event or last
                10 events.
              </Trans>
            </p>
          </div>
          <QueryInput name={name} />
          <TopicQueryResult result={[]} />
        </>
      )}
    </div>
  );
};

export default TopicQuery;
