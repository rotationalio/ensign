import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';

import { CardListItem } from '@/components/common/CardListItem';

import MetaDataTable from './MetaDataTable';

interface TopicQueryResultProps {
  result: any;
}

const ViewingEvent = ({ result }: TopicQueryResultProps) => {
  const { data } = result;

  const totalResults = data?.results?.length ?? 'N/A';

  return (
    <div className="mt-8">
      <Heading as="h2" className="mb-2 text-lg font-semibold">
        <Trans>
          Query Results
          <span className="ml-1 font-normal">
            {' '}
            - {totalResults} of {data?.total_events ?? 'N/A'}
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
