import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';

import { CardListItem } from '@/components/common/CardListItem';
import { ProjectQuery } from '@/features/projects/types/projectQueryService';

import MetaDataTable from './MetaDataTable';

interface MetaDataProps {
  data: ProjectQuery;
}

const ViewingEvent = ({ data }: MetaDataProps) => {
  const queryResults = String(data?.results?.length);
  return (
    <div className="mt-8">
      <Heading as="h2" className="mb-2 text-lg font-semibold">
        <Trans>
          Query Results
          <span className="ml-1 font-normal">
            {' '}
            - {queryResults ?? 'N/A'} of {data?.total_events ?? 'N/A'}
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
        <MetaDataTable results={data?.results} />
      </CardListItem>
    </div>
  );
};

export default ViewingEvent;
