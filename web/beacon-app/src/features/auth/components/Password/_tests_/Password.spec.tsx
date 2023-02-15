import { fireEvent, render, screen } from '@testing-library/react';

import Password from '../Password';

describe('Password render', () => {
  it('should show icon with eye closed, hidden password text, and accessible button name "Show Password" on render', () => {
    render(<Password />);
    const button = screen.getByTestId('button');

    expect(screen.getByTestId('hidePassword')).toBeVisible;
    expect(screen.getByTestId('password')).toHaveAttribute('type', 'password');
    expect(button).toHaveAccessibleName('Show Password');
  });
});

describe('Show password', () => {
  it('should show icon with open eye, password text, and accessible button name "Hide Password" when clicks on icon', () => {
    render(<Password />);
    const button = screen.getByTestId('button');

    fireEvent.click(button);
    expect(screen.getByTestId('showPassword')).toBeVisible;
    expect(screen.getByTestId('password')).toHaveAttribute('type', 'text');
    expect(button).toHaveAccessibleName('Hide Password');
  });
});

describe('Toggle icon', () => {
  it('should show icon with eye closed, hidden password text, and accessible button name "Show Password" when user clicks on eye icon a 2nd time', () => {
    render(<Password />);
    const button = screen.getByTestId('button');

    fireEvent.doubleClick(button);
    expect(screen.getByTestId('hidePassword')).toBeVisible;
    expect(screen.getByTestId('password')).toHaveAttribute('type', 'password');
    expect(button).toHaveAccessibleName('Show Password');
  });
});
