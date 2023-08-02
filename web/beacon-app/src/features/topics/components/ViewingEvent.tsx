import { Heading } from '@rotational/beacon-core';

import { CardListItem } from '@/components/common/CardListItem';

import MetaDataTable from './MetaDataTable';

const ViewingEvent = () => {
  return (
    <div className="mt-4">
      <Heading as="h2" className=" text-lg font-semibold">
        Query Results
      </Heading>
      <CardListItem>
        <p>Viewing Event: 1 of 10</p>
        <p className="pt-2 font-semibold">Meta Data</p>
        <MetaDataTable />
      </CardListItem>
    </div>
  );
};

export default ViewingEvent;
