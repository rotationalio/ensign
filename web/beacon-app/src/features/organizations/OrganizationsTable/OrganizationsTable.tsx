import { Table } from "@rotational/beacon-core";

export default function OrganizationsTable() { 
     
    return (
        <>
        <div className='bg-[#F7F9FB] rounded-lg py-2'>
        <h2 className='ml-4 text-lg font-bold'>Organizations</h2>
        </div>
        <Table
        columns={[
            {
                Header: 'Organization Name',
                accessor: 'organizationName'
            },
            {
                Header: 'Organization Role',
                accessor: 'organizationRole'
            },
            {
                Header: 'Status',
                accessor: 'status'
            },
            {
                Header: 'Projects',
                accessor: 'projects',
            },
            {
                Header: '# Teammates',
                accessor: 'teammates',
            },
            {
                Header: 'Date Created',
                accessor: 'dateCreated',
            },
          /*   {
                Header: 'Actions',
                accessor: 'actions',
            }, */
        ]}
        data={[
            {
                organizationName: 'Test',
                organizationRole: "Development",
                status: 'Active',
                projects: "1",
                teammates: '0',
                dateCreated: '20230208'
            },
            ]}
        />
        </>
    )
}