import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen } from '@testing-library/react';
import { vi } from 'vitest';

import type { Project } from '@/features/projects/types/Project';

import RenameProjectModal from '../RenameProjectModal';
const renderComponent = (props) => {
  const queryClient = new QueryClient();
  const wrapper = ({ children }) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
  return render(<RenameProjectModal {...props} />, { wrapper });
};

vi.mock('@lingui/macro', () => ({
  t: (str) => str,
}));

const projectMock = {
  id: '1',
  name: 'test',
  description: 'test',
  status: 'completed',
  created: '11-11-2021',
  modified: '11-11-2021',
} as Project;

describe('RenameProjectModal', () => {
  it('the modal should display ', () => {
    const propsMock = {
      open: true,
      handleModalClose: vi.fn(),
      project: projectMock,
    };

    renderComponent(propsMock);

    expect(screen.getByTestId('rename-project-modal')).toBeInTheDocument();
  });

  it('the modal should not display ', () => {
    const propsMock = {
      open: false,
      handleModalClose: vi.fn(),
      project: projectMock,
    };

    renderComponent(propsMock);

    expect(screen.queryByTestId('rename-project-modal')).not.toBeInTheDocument();
  });
});
