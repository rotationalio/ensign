import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen } from '@testing-library/react';
import React from 'react';
import { describe, expect, it, vi } from 'vitest';

import SandboxLayout from '../SandboxLayout';

const queryClient = new QueryClient();

const renderSandboxLayout = () => {
  return render(
    <QueryClientProvider client={queryClient}>
      <SandboxLayout />
    </QueryClientProvider>
  );
};

// Mock t and Trans tags from lingui.
vi.mock('invariant');
vi.mock('@lingui/macro', () => ({
  t: (str) => str,
  Trans: ({ children }) => children,
}));

// Mock react router hooks.
vi.mock('react-router-dom', () => ({
  useNavigate: () => vi.fn(),
  Link: ({ children }) => children,
  useLocation: () => ({
    pathname: '/test',
  }),
  NavLink: ({ children }) => children,
}));

describe('Sandbox layout', () => {
  it('should render', () => {
    const { container } = renderSandboxLayout();
    expect(container).toMatchSnapshot();
    expect(screen.getByTestId('sandbox-layout')).toBeInTheDocument();
  });

  it('should display the sandbox sidebar', () => {
    renderSandboxLayout();
    expect(screen.getByTestId('sandbox-sidebar')).toBeInTheDocument();
  });
});
