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

import QueryForm from './QueryForm';
const TopicQuery = ({ data }: TopicNameProps) => {
  const { topic_name: name, project_id: ProjectID } = data ?? {};
  const DEFAULT_QUERY = `SELECT * FROM ${name} LIMIT 1`;
  const [open, setOpen] = useState<boolean>(true);
  const [query, setQuery] = useState<string>(DEFAULT_QUERY);

  const { getProjectQuery } = useProjectQuery();
  // TODO: refactor this logic with sc-19702
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
          <QueryForm defaultEnSQL={query} queryHandler={queryHandler} />
          <TopicQueryResult result={[]} />
        </>
      )}
    </div>
  );
};

export default TopicQuery;
