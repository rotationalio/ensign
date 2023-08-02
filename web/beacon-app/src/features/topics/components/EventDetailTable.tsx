import { Table } from '@rotational/beacon-core';
import { ErrorBoundary } from '@sentry/react';
import { useMemo } from 'react';
import { useParams } from 'react-router-dom';

import { useFetchTopicEvents } from '../hooks/useFetchTopicEvents';
import { getEventDetailColumns } from '../utils';

const EventDetailTable = () => {
  const param = useParams();
  const { id: topicID } = param as { id: string };
  const { topicEvents, isFetchingTopicEvents } = useFetchTopicEvents(topicID);
  const initialColumns = useMemo(() => getEventDetailColumns(), []) as any;

  return (
    <div className="mx-4">
      <ErrorBoundary
        fallback={
          <div className="item-center my-auto flex w-full text-center font-bold text-danger-500">
            <p>
              Sorry we are having trouble listing event details for your topic, please refresh the
              page and try again.
            </p>
          </div>
        }
      >
        <Table columns={initialColumns} data={topicEvents} isLoading={isFetchingTopicEvents} />
      </ErrorBoundary>
    </div>
  );
};

export default EventDetailTable;
