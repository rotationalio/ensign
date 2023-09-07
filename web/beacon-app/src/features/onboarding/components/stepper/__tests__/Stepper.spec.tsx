import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render } from '@testing-library/react';
import { vi } from 'vitest';

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

// mock userloader

vi.mock('userLoader', () => ({
  __esModule: true,
  default: () => ({
    member: {
      invited: true,
      organization: 'organization',
      workspace: 'workspace',
    },
  }),
}));

describe('Stepper', () => {
  it('should render the component with default value', () => {
    const { container } = renderComponent();
    expect(container).toMatchSnapshot();
  });
});
