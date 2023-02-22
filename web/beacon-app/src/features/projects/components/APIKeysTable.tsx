import { t, Trans } from '@lingui/macro';
import { Table, Toast } from '@rotational/beacon-core';

import { TableHeading } from '@/components/common/TableHeader';
import { useFetchApiKeys } from '@/features/apiKeys/hooks/useFetchApiKeys';

export const APIKeysTable = () => {
  const { apiKeys, isFetchingApiKeys, hasApiKeysFailed, error } = useFetchApiKeys();

  if (isFetchingApiKeys) {
    // TODO: add loading state
    return (
      <div>
        <Trans>Loading...</Trans>
      </div>
    );
  }

  if (error) {
    return (
      <Toast
        isOpen={hasApiKeysFailed}
        variant="danger"
        title={t`Sorry we are having trouble fetching your topics, please try again later.`}
        description={(error as any)?.response?.data?.error}
      />
    );
  }

  const getApiKeys = (apikeys: any) => {
    if (!apikeys) return [];
    return Object.keys(apiKeys).map((key) => {
      const { name, owner, permissions, modifiers, created } = apiKeys[key];
      return { name, owner, permissions, modifiers, created };
    }) as any;
  };

  return (
    <div>
      <TableHeading>
        <Trans>API Keys</Trans>
      </TableHeading>
      <Table
        className="w-full"
        columns={[
          { Header: t`Name`, accessor: 'name' },
          { Header: t`Permissions`, accessor: 'permissions' },
          { Header: t`Owner`, accessor: 'owner' },
          { Header: t`Last Used`, accessor: 'modifiers' },
          { Header: t`Date Created`, accessor: 'created' },
        ]}
        data={getApiKeys(apiKeys)}
      />
    </div>
  );
};

export default APIKeysTable;
