import { Table } from '@rotational/beacon-core';
import { ErrorBoundary } from '@sentry/react';
import { useMemo } from 'react';

import { getEventDetailColumns } from '../utils';

const EventDetailTable = () => {
  const initialColumns = useMemo(() => getEventDetailColumns(), []) as any;
  return (
    <>
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
        <Table columns={initialColumns} data={[]} />
      </ErrorBoundary>
    </>
  );
};

export default EventDetailTable;
