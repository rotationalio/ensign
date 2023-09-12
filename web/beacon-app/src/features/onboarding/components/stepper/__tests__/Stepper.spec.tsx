import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen } from '@testing-library/react';
import { vi } from 'vitest';

import * as userLoader from '@/features/members/loaders/userLoader';

import { WORKSPACE_DOMAIN_BASE } from '../../../shared/constants';
import Stepper from '../Stepper';
const renderComponent = () => {
  const queryClient = new QueryClient();
  const wrapper = ({ children }: any) => (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
  return render(<Stepper />, { wrapper });
};

// mock t tag from lingui
vi.mock('@lingui/macro', () => ({
  t: (str) => str,
}));

vi.mock('@/features/members/loaders/userLoader', () => ({
  __esModule: true,
  default: () => ({
    member: {
      invited: false,
      organization: 'test',
      workspace: 'test',
    },
  }),
}));

describe('Stepper', () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });

  it('should render the component with default value', () => {
    vi.spyOn(userLoader, 'default')
      .mockImplementation()
      .mockReturnValue({
        member: {
          invited: false,
          organization: 'test',
          workspace: 'test',
        },
      });

    const { container } = renderComponent();
    expect(container).toMatchSnapshot();
  });

  it('should render the component for invited user', () => {
    const MockMember = {
      invited: true,
      organization: 'invited organization',
      workspace: 'invited-workspace',
    };

    vi.spyOn(userLoader, 'default')
      .mockImplementation()
      .mockReturnValue({
        member: {
          ...MockMember,
        },
      });

    renderComponent();

    expect(screen.getByText('Organization')).toBeInTheDocument();
    expect(screen.getByText('invited organization')).toBeInTheDocument();

    expect(screen.getByText('Workspace URL')).toBeInTheDocument();
    expect(screen.getByText(`${WORKSPACE_DOMAIN_BASE}${MockMember.workspace}`)).toBeInTheDocument();
  });
});
