import { dynamicActivate } from '@/I18n';
import { act, render, screen } from '@/utils/test-utils';

import LandingHeader from './LandingHeader';

describe('LandingHeader', () => {
  it('renders the logo and links', async () => {
    act(() => {
      dynamicActivate('en');
    });
    render(<LandingHeader />);

    const logo = screen.getByTestId('logo');
    expect(logo).toBeInTheDocument();

    const starterPlanLink = screen.getByText('Starter Plan');
    expect(starterPlanLink).toBeInTheDocument();
  });
});
