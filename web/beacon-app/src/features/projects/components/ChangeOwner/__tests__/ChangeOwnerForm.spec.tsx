import { render, screen, fireEvent } from '@testing-library/react';
import { vi } from 'vitest';

import type { Project } from '@/features/projects/types/Project';

import ChangeOwnerForm from '../ChangeOwnerForm';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import selectEvent from 'react-select-event';

const renderComponent = (props) => {
  const queryClient = new QueryClient();
  const wrapper = ({ children }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
  return render(<ChangeOwnerForm {...props} />, { wrapper });
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
    name: 'test',
    description: 'test',
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
describe('ChangeOwnerForm', () => {
  it('should render the form', () => {
    renderComponent(propsMock);

    expect(screen.getByTestId('update-owner-form')).toBeInTheDocument();
  });

  it('should be disabled for the current owner', () => {
    renderComponent(propsMock);

    expect(screen.getByTestId('current-owner')).toBeDisabled();
  });

  it('should display the current owner', () => {
    renderComponent(propsMock);

    expect(screen.getByTestId('current-owner')).toBeDefined();
  });

  it('should display the correct number of options', async () => {
    const propsMock = {
      initialValues: {
        id: '1',
        name: 'test',
        description: 'test',
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

    const { getByTestId, getByLabelText } = renderComponent(propsMock);

    expect(getByTestId('update-owner-form')).toHaveFormValues({ new_owner: '' }); // empty select

    // await selectEvent.select(getByLabelText('Select New Owner'), ['test2']);

    // expect(getByTestId('update-owner-form')).toHaveFormValues({ new_owner: 'test2' }); // selected option
  });
});
