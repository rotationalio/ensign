import { Table } from "@rotational/beacon-core";
import { useCreateTenant } from "../hooks/useCreateTenant";

export default function TenantTable() {  
  const { tenant } = useCreateTenant();
  
    return (
        <>
        <Table
  columns={[
    {
      Header: 'Tenant Name',
      accessor: 'tenantName'
    },
    {
      Header: 'Environment Label',
      accessor: 'environmentLabel'
    },
    {
      Header: 'Status',
      accessor: 'status'
    },
/*     {
        Header: 'Cloud Provider',
        accessor: 'cloudProvider',
      },
      {
        Header: 'Region',
        accessor: 'region',
      }, */
    {
      Header: 'Date Created',
      accessor: 'dateCreated',
    },
    /* {
        Header: 'Actions',
        accessor: 'actions',
      }, */
  ]}
  data={[
    {
      /* actions: [
        {
          label: 'Edit',
          onClick: () => {}
        },
        {
          label: 'Delete',
          onClick: function noRefCheck() {}
        }
      ], */
      tenantName: `${tenant.name}`,
      environmentLabel: `${tenant.environment_type}`,
      status: 'Active',
      dateCreated: `${tenant.date_created}`
    },
  ]}
/>
        </>
    )
}
