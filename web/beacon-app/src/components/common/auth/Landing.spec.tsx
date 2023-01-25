import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';

import LandingHeader from './LandingHeader';

describe('LandingHeader', () => {
  it('renders the logo and links', () => {
    render(<LandingHeader />, { wrapper: BrowserRouter });

    const logo = screen.getByTestId('logo');
    expect(logo).toBeInTheDocument();

    const starterPlanLink = screen.getByText('Starter Plan');
    expect(starterPlanLink).toBeInTheDocument();

    const upgradeButton = screen.getByText('Upgrade');
    expect(upgradeButton).toBeInTheDocument();
  });
});
