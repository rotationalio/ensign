import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';

import { CardListItem } from '@/components/common/CardListItem';

import { getEventsPaginationCounter, getQueryPaginationCounter } from '../../utils';
import MetaDataTable from './MetaDataTable';
interface MetaDataProps {
  totalResults: number;
  totalEvents: string;
  counter: number;
  metadataResult: any;
}

const ViewingEvent = ({ totalResults, totalEvents, counter, metadataResult }: MetaDataProps) => {
  return (
    <div className="mt-8">
      <Heading as="h2" className="mb-2 text-lg font-semibold" data-cy="topic-query-results">
        <Trans>
          Query Results -
          <span className="ml-1 font-normal" data-testid="query-result-count">
            {getQueryPaginationCounter(+totalResults, +totalEvents)}
          </span>
        </Trans>
      </Heading>
      <CardListItem className="!rounded-none">
        <p data-testid="view-event" data-cy="viewing-event-results">
          <Trans>Viewing Event: {getEventsPaginationCounter(counter, +totalResults)}</Trans>
        </p>
        <p className="pt-2 font-semibold">
          <Trans>Meta Data</Trans>
        </p>
        <MetaDataTable metadataResult={metadataResult} />
      </CardListItem>
    </div>
  );
};

export default ViewingEvent;
