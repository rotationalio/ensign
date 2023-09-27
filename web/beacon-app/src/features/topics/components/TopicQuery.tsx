import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';
// import { useAnimate, useInView } from 'framer-motion';
import React, { useEffect, useState } from 'react';
import toast from 'react-hot-toast';
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
  const [openResult, setOpenResult] = useState<boolean>(false);
  const [hasInvalidQuery, setHasInvalidQuery] = useState<boolean>(false);
  const {
    getProjectQuery,
    isCreatingProjectQuery,
    projectQuery,
    error,
    reset,
    wasProjectQueryCreated,
  } = useProjectQuery();
  const [resetQuery, setResetQuery] = useState<boolean>(false);

  const handleSubmitProjectQuery = (values: any) => {
    setHasInvalidQuery(false);
    const payload = {
      ...values,
      projectID: ProjectID,
    };

    getProjectQuery(payload);
  };

  const toggleHandler = () => setOpen(!open);

  const handleResetQuery = () => {
    setResetQuery(!resetQuery);
    reset();
  };

  useEffect(() => {
    if (isCreatingProjectQuery) {
      setResetQuery(!resetQuery);
    }
    return () => {
      setResetQuery(false);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isCreatingProjectQuery]);

  useEffect(() => {
    if (wasProjectQueryCreated) {
      setOpenResult(true);
    }
    return () => {
      setOpenResult(false);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [wasProjectQueryCreated]);

  useEffect(() => {
    if (error) {
      toast.error(`${error?.response?.data?.error}` || error.message);
    }
    return () => {
      toast.dismiss();
    };
  }, [error]);

  useEffect(() => {
    if (wasProjectQueryCreated && !!projectQuery?.error) {
      setHasInvalidQuery(true);
    }
  }, [projectQuery?.error, wasProjectQueryCreated]);

  return (
    <div data-testid="topic-query-title" className="mt-10" data-cy="topic-query-title">
      <button
        className="mb-4 flex h-5 place-items-center gap-3"
        onClick={toggleHandler}
        data-cy="topic-query-heading"
      >
        <Heading as="h1" className=" text-lg font-semibold">
          <Trans>Topic Query</Trans>
        </Heading>
        {open ? (
          <SlArrowDown data-cy="topic-query-carat-down" />
        ) : (
          <SlArrowRight data-cy="topic-query-carat-up" />
        )}
      </button>

      {open && (
        <>
          <TopicQueryInfo />
          <QueryForm
            defaultEnSQL={DEFAULT_QUERY}
            onSubmit={handleSubmitProjectQuery}
            isSubmitting={isCreatingProjectQuery}
            onReset={handleResetQuery}
          />
          {openResult && (
            <TopicQueryResult
              data={projectQuery}
              error={error}
              onReset={resetQuery}
              hasInvalidQuery={hasInvalidQuery}
            />
          )}
        </>
      )}
    </div>
  );
};

export default TopicQuery;
