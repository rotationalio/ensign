/* eslint-disable testing-library/await-async-utils */
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen, waitFor } from '@testing-library/react';
import { vi } from 'vitest';

import * as useFetchTenants from '@/features/tenants/hooks/useFetchTenants';

import DashLayout from '../DashLayout';

const queryClient = new QueryClient();

const renderDashLayout = () => {
  return render(
    <QueryClientProvider client={queryClient}>
      <DashLayout />
    </QueryClientProvider>
  );
};
vi.mock('invariant');
// mock t tag from lingui
vi.mock('@lingui/macro', () => ({
  t: (str) => str,
  Trans: ({ children }) => children,
}));

// mock useFetchTenants hook
vi.mock('@/features/tenants/hooks/useFetchTenants', () => ({
  __esModule: true,
  useFetchTenants: () => ({
    tenants: {
      tenants: [
        {
          id: 1,
          name: 'test',
        },
      ],
    },
    wasTenantsFetched: true,
  }),
}));

// mock useDropdownMenu
vi.mock('@/components/MenuDropdown/useDropdownMenu', () => ({
  __esModule: true,
  useDropdownMenu: () => ({
    menuItems: [
      {
        label: 'test',
        value: 'test',
      },
    ],
  }),
}));

// mock router hooks
vi.mock('react-router-dom', () => ({
  useNavigate: () => vi.fn(),
  Link: ({ children }) => children,
  useLocation: () => ({
    pathname: '/test',
  }),
  NavLink: ({ children }) => children,
}));

// mock useOrgStore
vi.mock('@/store', () => ({
  __esModule: true,
  useOrgStore: {
    getState: () => ({
      org: {
        name: 'test',
      },
      setTenantID: vi.fn(),
      setOrgName: vi.fn(),
    }),
  },
}));

describe('DashLayout', () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });

  it('should render the component', () => {
    const { container } = renderDashLayout();
    expect(container).toMatchSnapshot();
  });

  it('should render the error component when the tenant is not loaded correctly', () => {
    vi.spyOn(useFetchTenants, 'useFetchTenants')
      .mockImplementation()
      .mockReturnValue({
        tenants: {
          tenants: null,
        },
        wasTenantsFetched: true,
      });

    renderDashLayout();

    waitFor(() => {
      // check if testid is rendered
      expect(
        screen.getByText(
          'Something went wrong. Please contact us at support@rotational.io for assistance.'
        )
      );
    });
  });
});
