// Test the UpdateAlert component
// Must mock the useFetchStatus hook
// The appConfig should be mocked as well
// If the version from the useFetchStatus hook is different from the version from the appConfig, the UpdateAlert should be rendered

import { render, screen } from '@testing-library/react';
import React from 'react';
import { describe, expect, it, vi } from 'vitest';

import UpdateAlert from '../Sidebar/UpdateAlert';

vi.mock('@/features/home/hooks/useFetchStatus', () => ({
  useFetchStatus: () => ({
    status: {
      version: '0.11.0',
    },
  }),
}));

vi.mock('@/appConfig', () => ({
  appConfig: {
    version: '0.10.0',
  },
}));

vi.mock('@lingui/macro', async (importOrginial: any) => ({
  ...importOrginial,
  t: () => vi.fn(),
  Trans: () => vi.fn(),
}));

const renderComponent = () => {
  return render(<UpdateAlert />);
};

describe('UpdateAlert', () => {
  it('should render if the version number returned from the useFetchStatus hook is different from the version number from the appConfig', () => {
    const { container } = renderComponent();
    expect(container).toMatchSnapshot();
    expect(screen.getByTestId('update-alert-btn')).toBeInTheDocument();
  });

  /* it('should refresh the page when the update button is clicked', () => {
        renderComponent();
        const updateBtn = screen.getByTestId('update-alert-btn');
        fireEvent.click(updateBtn);

        // Test reload functionality.
        expect(window.location.reload).toHaveBeenCalled();
    }); */
});
