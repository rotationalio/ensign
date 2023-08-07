import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';

import { CardListItem } from '@/components/common/CardListItem';

import { getEventsPaginationCounter, getQueryPaginationCounter } from '../../utils';
import MetaDataTable from './MetaDataTable';
interface MetaDataProps {
  totalResults: number;
  totalEvents: string;
  counter: number;
}

const ViewingEvent = ({ totalResults, totalEvents, counter }: MetaDataProps) => {
  return (
    <div className="mt-8">
      <Heading as="h2" className="mb-2 text-lg font-semibold">
        <Trans>
          Query Results
          <span className="ml-1 font-normal" data-testid="query-result-count">
            {' '}
            - {getQueryPaginationCounter(counter, +totalResults)}
          </span>
        </Trans>
      </Heading>
      <CardListItem className="!rounded-none">
        <p data-testid="view-event">
          <Trans>
            Viewing Event:
            {getEventsPaginationCounter(1, +totalEvents)}
          </Trans>
        </p>
        <p className="pt-2 font-semibold">
          <Trans>Meta Data</Trans>
        </p>
        <MetaDataTable />
      </CardListItem>
    </div>
  );
};

export default ViewingEvent;
