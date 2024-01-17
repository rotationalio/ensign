import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { render, screen } from '@testing-library/react';
import React from 'react';
import { describe, expect, it, vi } from 'vitest';

import OnboardingLayout from '../OnboardingLayout';

const queryClient = new QueryClient();

const renderOnboardingLayout = () => {
  return render(
    <QueryClientProvider client={queryClient}>
      <OnboardingLayout />
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

describe('Onboarding layout', () => {
  it('should render', () => {
    const { container } = renderOnboardingLayout();
    expect(container).toMatchSnapshot();
    expect(screen.getByTestId('onboarding-layout')).toBeInTheDocument();
  });

  it('should display the onboarding sidebar', () => {
    renderOnboardingLayout();
    expect(screen.getByTestId('onboarding-sidebar')).toBeInTheDocument();
  });
});
