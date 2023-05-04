import { render } from '@testing-library/react';
import jest from 'jest-mock';
import React from 'react';
import { vi } from 'vitest';

//import { customRender } from '../../../../utils/test-utils';
import { Project } from '../../types/Project';
import ProjectsTable from '../ProjectsTable';

const renderComponent = (props) => {
  return render(<ProjectsTable {...props} />);
};
vi.mock('react-router-dom', async (importOrginial: any) => ({
  ...importOrginial,
  useNavigate: () => jest.fn(),
}));

vi.mock('@lingui/macro', async (importOrginial: any) => ({
  ...importOrginial,
  t: () => jest.fn(),
  Trans: () => jest.fn(),
}));
describe('ProjectsTable', () => {
  // mock useNavigate hook

  //   vi.mock('react-table', async (importOrginial: any) => ({
  //     ...importOrginial,
  //     useTable: () => jest.fn(),
  //     useFilters: () => jest.fn(),
  //     useSortBy: () => jest.fn(),
  //     usePagination: () => jest.fn(),
  //     Row: () => jest.fn(),

  //     Column: () => jest.fn(),
  //   }));

  const MockProjectProps = {
    projects: [
      {
        id: '1',
        name: 'Project 1',
        tenant_id: '1',
        modified: '2021-01-01',
        description: 'Project 1 description',
        status: 'Active',
        created: '2021-01-01',
      },
      {
        id: '2',
        name: 'Project 2',
        description: 'Project 2 description',
        status: 'Inactive',
        created: '2021-01-02',
      },
    ] as Project[],
  };

  it('should render table with correct columns', () => {
    const { container } = renderComponent(MockProjectProps);

    expect(container).toMatchSnapshot();
    // const table = screen.getByRole('table');
    // expect(table).toBeInTheDocument();
    // expect(screen.getByText('Project ID')).toBeInTheDocument();
    // expect(screen.getByText('Project Name')).toBeInTheDocument();
    // expect(screen.getByText('Description')).toBeInTheDocument();
    // expect(screen.getByText('Status')).toBeInTheDocument();
    // expect(screen.getByText('Date Created')).toBeInTheDocument();
  });
});
