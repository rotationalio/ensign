import { render, screen } from '@testing-library/react';
import React from 'react';
import { describe, expect, it, vi } from 'vitest';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import SandboxLayout from '../SandboxLayout';
import DashLayout from '../DashLayout';
import OnboardingLayout from '../OnboardingLayout';

const queryClient = new QueryClient();

const wrapper = ({ children }: any) => (
  <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
);

const renderOnboardingLayout = () => {
  return render(<OnboardingLayout />, { wrapper });
}

const renderSandboxLayout = () => {
  return render(<SandboxLayout />, { wrapper });
};

const renderDashLayout = () => {
  return render(<DashLayout />, { wrapper });
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

describe('Onboarding layout', () => {
  it('should render', () => {
    renderOnboardingLayout();
    expect(screen.getByTestId('onboarding-layout')).toBeInTheDocument();
    expect(screen.getByTestId('onboarding-sidebar')).toBeInTheDocument();
  });
});

describe('Sandbox layout', () => {
    it('should render', () => {
      renderSandboxLayout();
      expect(screen.getByTestId('sandbox-layout')).toBeInTheDocument();
      expect(screen.getByTestId('sandbox-sidebar')).toBeInTheDocument();
    });
});