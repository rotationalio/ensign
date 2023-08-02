import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';

import { CardListItem } from '@/components/common/CardListItem';

import MetaDataTable from './MetaDataTable';

interface MetaDataProps {
  totalResults: number;
  totalEvents: string;
}

const ViewingEvent = ({ totalResults, totalEvents }: MetaDataProps) => {
  return (
    <div className="mt-8">
      <Heading as="h2" className="mb-2 text-lg font-semibold">
        <Trans>
          Query Results
          <span className="ml-1 font-normal">
            {' '}
            - {String(totalResults ?? 'N/A')} of {totalEvents ?? 'N/A'}
          </span>
        </Trans>
      </Heading>
      <CardListItem className="!rounded-none">
        <p>
          <Trans>Viewing Event: 1 of 10</Trans>
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
