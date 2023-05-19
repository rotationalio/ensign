import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen, fireEvent, act } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

import type { Project } from '@/features/projects/types/Project';

import EditProjectForm from '../EditProjectForm';
// import selectEvent from 'react-select-event';

const renderComponent = (props) => {
  const queryClient = new QueryClient();
  const wrapper = ({ children }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
  return render(<EditProjectForm {...props} />, { wrapper });
};

vi.mock('@lingui/macro', () => ({
  t: (str) => str,
  Trans: ({ children }) => children,
}));
// vi.mock('@/features/members/hooks/useMembers', () => ({
//   useMembers: () => ({
//     members: [
//       {
//         id: '1',
//         name: 'test',
//         status: 'active',
//       },
//       {
//         id: '2',
//         name: 'test2',
//         status: 'active',
//       },
//     ],
//   }),
// }));
const propsMock = {
  initialValues: {
    id: '1',
    name: 'name test',
    description: 'description test',
    status: 'completed',
    created: '11-11-2021',
    modified: '11-11-2021',
    owner: {
      id: '1',
      name: 'test',
      status: 'active',
    },
  } as Project,
  handleModalClose: vi.fn(),
};
describe('EditProjectForm', () => {
  it('should render the form', () => {
    renderComponent(propsMock);

    expect(screen.getByTestId('update-project-form')).toBeInTheDocument();
  });

  it('should disable the current name field', () => {
    renderComponent(propsMock);

    expect(screen.getByTestId('prj-current-name')).toBeDisabled();
  });

  it('should render the new name field', () => {
    renderComponent(propsMock);

    expect(screen.getByTestId('prj-new-name')).toBeInTheDocument();
  });

  it('should render the description field', () => {
    renderComponent(propsMock);

    expect(screen.getByTestId('prj-description')).toBeInTheDocument();
  });

  it('should disable the submit button if the name is empty and the description is the same', () => {
    const props = {
      initialValues: {
        id: '1',
        name: '',
        description: 'description test',
      } as Project,
      handleModalClose: vi.fn(),
    };
    renderComponent(props);

    expect(screen.getByTestId('update-project-submit')).toBeDisabled();
  });
});
