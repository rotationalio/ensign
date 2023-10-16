import { fireEvent, render, screen } from '@testing-library/react';

import TextField from './TextField';

describe('<TextField />', () => {
  it('should display correclty the label', () => {
    render(<TextField label="Name" />);
    const input = screen.getByLabelText('Name');
    expect(input).toBeInTheDocument();
  });

  it('should apply size props', () => {
    render(<TextField size="medium" />);
    const input = screen.getByTestId('input');
    expect(input).toHaveClass('beacon-text-base');
  });

  it('should display the right icon', () => {
    const rightIcon = <span>X</span>;
    render(<TextField rightIcon={rightIcon} />);
    const icon = screen.getByTestId('right-icon');
    expect(icon).toBeInTheDocument();
  });

  it('should display the description', () => {
    render(<TextField description="Enter your name" />);
    const description = screen.getByText('Enter your name');
    expect(description).toBeInTheDocument();
  });

  it('should display error message', () => {
    render(<TextField errorMessage="This field is required" />);
    const error = screen.getByText('This field is required');
    expect(error).toBeInTheDocument();
  });

  it('should update value on change event', () => {
    render(<TextField label="Name" />);
    const input: any = screen.getByLabelText('Name');
    fireEvent.change(input, { target: { value: 'John' } });
    expect(input.value).toBe('John');
  });
});
