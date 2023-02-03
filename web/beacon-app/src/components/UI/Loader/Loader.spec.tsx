import { render, screen } from '@testing-library/react';

import Loader from './Loader';

describe('Loader', () => {
  it('renders without crashing', () => {
    render(<Loader />);
  });

  it('displays a loading spinner', () => {
    render(<Loader />);
    const spinner = screen.getByTestId('loading-spinner');

    expect(spinner).toBeInTheDocument();
  });

  it('displays a label if provided', () => {
    const label = 'Loading...';
    render(<Loader label={label} />);
    const labelElement = screen.getByText(label);

    expect(labelElement).toBeInTheDocument();
  });

  it('applies additional labelProps if provided', () => {
    const label = 'Loading...';
    const labelProps = { className: 'text-red' };
    render(<Loader label={label} labelProps={labelProps} />);
    const labelElement = screen.getByText(label);

    expect(labelElement).toHaveClass('text-red');
  });
});
