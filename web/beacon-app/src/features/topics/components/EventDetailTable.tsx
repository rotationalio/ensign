import { Trans } from '@lingui/macro';
import { Heading, Table } from '@rotational/beacon-core';
import { ErrorBoundary } from '@sentry/react';
import { useMemo } from 'react';
import { useParams } from 'react-router-dom';

import { getTopicEventsMockData } from '../__mocks__';
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
              <Trans>
                Sorry we are having trouble listing event details for your topic, please refresh the
                page and try again.
              </Trans>
            </p>
          </div>
        }
      >
        <Heading as="h2" className="text-xl font-semibold">
          <Trans>Topic Usage</Trans>
        </Heading>
        <Table
          columns={initialColumns}
          data-testId="event-detail-table"
          data={topicEvents || getTopicEventsMockData()} // for now we are using mock data until we have the API ready
          isLoading={isFetchingTopicEvents}
          data-cy="event-detail-table"
        />
      </ErrorBoundary>
    </div>
  );
};

export default EventDetailTable;
