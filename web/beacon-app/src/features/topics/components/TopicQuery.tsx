import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';
// import { useAnimate, useInView } from 'framer-motion';
import React, { useState } from 'react';
import { SlArrowDown, SlArrowRight } from 'react-icons/sl';

import TopicQueryInfo from './TopicQueryInfo';
import TopicQueryResult from './TopicQueryResult';
type TopicNameProps = {
  data: any;
};

import { useProjectQuery } from '@/features/projects/hooks/useProjectQuery';

import QueryInput from './QueryInput';
const TopicQuery = ({ data }: TopicNameProps) => {
  const { topic_name: name, project_id: ProjectID } = data ?? {};
  console.log('data', data);
  const DEFAULT_QUERY = `SELECT * FROM ${name} LIMIT 10`;

  console.log('[] DEFAULT_QUERY', DEFAULT_QUERY);
  const [open, setOpen] = useState<boolean>(true);
  const [query, setQuery] = useState<string>(DEFAULT_QUERY);

  const { getProjectQuery } = useProjectQuery();

  const queryHandler = (query: string) => {
    getProjectQuery({
      ProjectID,
      query,
    } as any);
    setQuery(query);
  };

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
          <TopicQueryInfo />
          <QueryInput defaultEnSQL={query} queryHandler={queryHandler} />
          <TopicQueryResult result={[]} />
        </>
      )}
    </div>
  );
};

export default TopicQuery;
