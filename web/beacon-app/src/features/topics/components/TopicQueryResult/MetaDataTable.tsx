import { Trans } from '@lingui/macro';
import { Table } from '@rotational/beacon-core';
import { ErrorBoundary } from '@sentry/react';

import { getProjectQueryMetaData } from '../../utils';

type MetaDataTableProps = {
  metadataResult: any;
};

const MetaDataTable = ({ metadataResult }: MetaDataTableProps) => {
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
        <div className="max-h-[150px] overflow-y-auto overflow-x-hidden">
          <Table
            columns={initialColumns}
            data={getProjectQueryMetaData(metadataResult)}
            theadClassName="hidden"
            tdClassName="first:font-bold"
            data-cy="query-meta-data-table"
          />
        </div>
      </ErrorBoundary>
    </div>
  );
};

export default MetaDataTable;
