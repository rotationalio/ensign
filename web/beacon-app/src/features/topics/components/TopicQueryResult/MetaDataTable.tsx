import { t } from '@lingui/macro';
import { Table } from '@rotational/beacon-core';
import { ErrorBoundary } from '@sentry/react';

const MetaDataTable = () => {
  const initialColumns: any = [
    { Header: t`Key`, accessor: 'key' },
    {
      Header: t`Value`,
      accessor: 'value',
    },
  ];
  return (
    <div className="mx-4">
      <ErrorBoundary
        fallback={
          <div className="item-center my-auto flex w-full text-center font-bold text-danger-500">
            <p>
              Sorry we are having trouble listing the meta data for your event, please refresh the
              page and try again.
            </p>
          </div>
        }
      >
        {/* TODO: Add getMetaData util */}
        <Table columns={initialColumns} data={[]} />
      </ErrorBoundary>
    </div>
  );
};

export default MetaDataTable;
