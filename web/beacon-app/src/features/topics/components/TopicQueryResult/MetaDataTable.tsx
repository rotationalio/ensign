import { Trans } from '@lingui/macro';
import { Table } from '@rotational/beacon-core';
import { ErrorBoundary } from '@sentry/react';

import { getProjectQueryMetaData } from '../../utils';

type MetaDataTableProps = {
  results: any;
};

const MetaDataTable = ({ results }: MetaDataTableProps) => {
  const initialColumns: any = [
    {
      Header: '',
      accessor: 'key',
    },
    {
      Header: '',
      accessor: 'value',
    },
  ];
  return (
    <div className="mx-4">
      <ErrorBoundary
        fallback={
          <div className="item-center my-auto flex w-full text-center font-bold text-danger-500">
            <p>
              <Trans>
                Sorry we are having trouble listing the meta data for your event, please refresh the
                page and try again.
              </Trans>
            </p>
          </div>
        }
      >
        <Table columns={initialColumns} data={getProjectQueryMetaData(results)} />
      </ErrorBoundary>
    </div>
  );
};

export default MetaDataTable;
