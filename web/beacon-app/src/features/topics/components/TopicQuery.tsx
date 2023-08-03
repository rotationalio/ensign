import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';
// import { useAnimate, useInView } from 'framer-motion';
import React, { useState } from 'react';
import { SlArrowDown, SlArrowRight } from 'react-icons/sl';

import TopicQueryResult from './TopicQueryResult';
type TopicNameProps = {
  data: any;
};

import { Link } from 'react-router-dom';

import { EXTRENAL_LINKS } from '@/application';
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

  const queryHandler = (values: any) => {
    getProjectQuery({
      projectID: ProjectID,
      query: values.query,
    } as any);
    setQuery(values.query);
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
          <div className="flex space-x-1">
            <p>
              <Trans>
                Query the topic for insights with{' '}
                <Link
                  to={EXTRENAL_LINKS.ENSQL}
                  className="font-semibold text-[#1D65A6] underline"
                  target="_blank"
                >
                  EnSQL
                </Link>{' '}
                e.g. the latest event or last 5 events. The maximum number of query results is 10.
                Use our{' '}
                <Link
                  to={EXTRENAL_LINKS.SDKs}
                  className="font-semibold text-[#1D65A6] underline"
                  target="_blank"
                >
                  SDKs
                </Link>{' '}
                for more results.
              </Trans>
            </p>
          </div>
          <QueryInput defaultEnSQL={query} queryHandler={queryHandler} />
          <TopicQueryResult result={[]} />
        </>
      )}
    </div>
  );
};

export default TopicQuery;
