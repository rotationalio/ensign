import { Trans } from '@lingui/macro';
import { Heading } from '@rotational/beacon-core';

import { CardListItem } from '@/components/common/CardListItem';

import MetaDataTable from './MetaDataTable';

interface TopicQueryResultProps {
    result: any;
  }

const ViewingEvent = ({result}: TopicQueryResultProps) => {
    const { data } = result;
  return (
    <div className="mt-4">
      <Heading as="h2" className=" mb-2 text-lg font-semibold">
        <Trans>Query Results</Trans>
        <span className="font-normal"> - 8 of {data?.total_events}</span>
      </Heading>
      <CardListItem className="!rounded-none">
        <p>Viewing Event: 1 of 10</p>
        <p className="pt-2 font-semibold">Meta Data</p>
        <MetaDataTable />
      </CardListItem>
    </div>
  );
};

export default ViewingEvent;
