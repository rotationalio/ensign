import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';
// import { useAnimate, useInView } from 'framer-motion';
import React, { useState } from 'react';
import { SlArrowDown, SlArrowUp } from 'react-icons/sl';

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

  const { getProjectQuery, isCreatingProjectQuery, projectQuery, error } = useProjectQuery();

  console.log('[] projectQuery', projectQuery);

  const handleSubmitProjectQuery = (values: any) => {
    const payload = {
      ...values,
      projectID: ProjectID,
    };
    getProjectQuery(payload);
  };

  const toggleHandler = () => setOpen(!open);

  return (
    <div data-testid="topic-query-title" className="mt-10">
      <button className="mb-4 flex h-5 place-items-center gap-3" onClick={toggleHandler}>
        <Heading as="h1" className=" text-lg font-semibold">
          <Trans>Topic Query</Trans>
        </Heading>
        {open ? <SlArrowDown /> : <SlArrowUp />}
      </button>

      {open && (
        <>
          <TopicQueryInfo />
          <QueryForm
            defaultEnSQL={DEFAULT_QUERY}
            onSubmit={handleSubmitProjectQuery}
            isSubmitting={isCreatingProjectQuery}
          />
          <TopicQueryResult data={projectQuery} error={error} />
        </>
      )}
    </div>
  );
};

export default TopicQuery;
