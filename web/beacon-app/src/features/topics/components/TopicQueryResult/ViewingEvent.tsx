import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';

import { CardListItem } from '@/components/common/CardListItem';

import MetaDataTable from './MetaDataTable';

interface MetaDataProps {
  totalResults: number;
  totalEvents: string;
  counter: number;
}

const ViewingEvent = ({ totalResults, totalEvents, counter }: MetaDataProps) => {
  const totalEventsString = totalEvents ? `1 of ${totalEvents}` : '0 of 0'; // TODO:  implement event pagination counter
  const totalResultsPaginateString =
    totalResults > 0 ? ` ${counter} of ${totalResults}` : '0 results of 0 total';

  return (
    <div className="mt-8">
      <Heading as="h2" className="mb-2 text-lg font-semibold">
        <Trans>
          Query Results
          <span className="ml-1 font-normal"> - {totalResultsPaginateString}</span>
        </Trans>
      </Heading>
      <CardListItem className="!rounded-none">
        <p>
          <Trans>
            Viewing Event:
            {totalEventsString}
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
