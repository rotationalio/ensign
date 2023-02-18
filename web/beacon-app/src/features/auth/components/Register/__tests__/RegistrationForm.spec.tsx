import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';

import { RegistrationForm } from '../..';

vi.mock('react-router-dom');

describe('first', () => {
  it('should render form initial value', () => {
    const handleSubmit = vi.fn();
    render(<RegistrationForm onSubmit={handleSubmit} />);

    expect(screen.getByTestId('name')).toHaveAttribute('value', '');
    expect(screen.getByTestId('email')).toHaveAttribute('value', '');
    expect(screen.getByTestId('password')).toHaveAttribute('value', '');
    expect(screen.getByTestId('pwcheck')).toHaveAttribute('value', '');
    expect(screen.getByTestId('organization')).toHaveAttribute('value', '');
    expect(screen.getByTestId('terms_agreement')).not.toBeChecked();
  });

  describe('Name', () => {
    it('should display error message when name field is empty', async () => {
      const handleSubmit = vi.fn();
      render(<RegistrationForm onSubmit={handleSubmit} />);

      expect(screen.getByTestId('name')).toHaveAttribute('value', '');
      userEvent.click(screen.getByRole('button', { name: /create starter account/i }));

      await waitFor(() => {
        expect(screen.getByText(/The name is required./i)).toBeInTheDocument();
      });
    });
  });

  describe('Password eye icon renders', () => {
    it('should show icon with a closed eye, hidden password text, and accessible button name "Show Password" on render', () => {
      const handleSubmit = vi.fn();
      render(<RegistrationForm onSubmit={handleSubmit} />);
      const button = screen.getByTestId('button');

      expect(screen.getByTestId('hidePassword')).toBeVisible;
      expect(screen.getByTestId('password')).toHaveAttribute('type', 'password');
      expect(button).toHaveAccessibleName('Show Password');
    });
  });

  describe('Show password', () => {
    it('should show icon with open eye, password text, and accessible button name "Hide Password" when clicks on icon', () => {
      const handleSubmit = vi.fn();
      render(<RegistrationForm onSubmit={handleSubmit} />);
      const button = screen.getByTestId('button');

      fireEvent.click(button);
      expect(screen.getByTestId('showPassword')).toBeVisible;
      expect(screen.getByTestId('password')).toHaveAttribute('type', 'text');
      expect(button).toHaveAccessibleName('Hide Password');
    });
  });

  describe('Hide password', () => {
    it('should show icon with eye closed, hidden password text, and accessible button name "Show Password" when user clicks on eye icon a 2nd time', () => {
      const handleSubmit = vi.fn();
      render(<RegistrationForm onSubmit={handleSubmit} />);
      const button = screen.getByTestId('button');

      fireEvent.doubleClick(button);
      expect(screen.getByTestId('hidePassword')).toBeVisible;
      expect(screen.getByTestId('password')).toHaveAttribute('type', 'password');
      expect(button).toHaveAccessibleName('Show Password');
    });
  });

  // describe('Password', () => {
  //   it('should display error message when password field is empty', async () => {
  //     const handleSubmit = vi.fn();
  //     render(<RegistrationForm onSubmit={handleSubmit} />);

  //     userEvent.type(screen.getByTestId('password'), '');
  //     userEvent.click(screen.getByRole('button', { name: /create starter account/i }));

  //     await waitFor(() => {
  //       expect(screen.getByText(/The password is required./i)).toBeInTheDocument();
  //     });
  //   });

  //   it('should display error message when confirm password field is empty', async () => {
  //     const handleSubmit = vi.fn();
  //     render(<RegistrationForm onSubmit={handleSubmit} />);

  //     userEvent.type(screen.getByTestId('pwcheck'), '');
  //     userEvent.click(screen.getByRole('button', { name: /create starter account/i }));

  //     await waitFor(() => {
  //       expect(screen.getByText(/The confirm password is required./i)).toBeInTheDocument();
  //     });
  //   });

  //   it('should display error message when organization field is empty', async () => {
  //     const handleSubmit = vi.fn();
  //     render(<RegistrationForm onSubmit={handleSubmit} />);

  //     userEvent.type(screen.getByTestId('organization'), '');
  //     userEvent.click(screen.getByRole('button', { name: /create starter account/i }));

  //     await waitFor(() => {
  //       expect(screen.getByText(/The organization is required./i)).toBeInTheDocument();
  //     });
  //   });

  //   it('should display error message when domain field is empty', async () => {
  //     const handleSubmit = vi.fn();
  //     render(<RegistrationForm onSubmit={handleSubmit} />);

  //     userEvent.type(screen.getByTestId('domain'), '');
  //     userEvent.click(screen.getByRole('button', { name: /create starter account/i }));

  //     await waitFor(() => {
  //       expect(screen.getByText(/The domain is required./i)).toBeInTheDocument();
  //     });
  //   });
  // });

  it('should submit the form', async () => {
    const handleSubmit = vi.fn();

    render(<RegistrationForm onSubmit={handleSubmit} />);

    // eslint-disable-next-line testing-library/no-unnecessary-act
    await act(async () => {
      await userEvent.type(screen.getByTestId('name'), 'John Doe');
      await userEvent.type(screen.getByTestId('email'), 'john.doe@example.com');
      await userEvent.type(screen.getByTestId('password'), 'Password123@@_');
      await userEvent.type(screen.getByTestId('pwcheck'), 'Password123@@_');
      await userEvent.type(screen.getByTestId('organization'), 'Acme Inc.');
      await userEvent.type(screen.getByTestId('domain'), 'acme.io');
      fireEvent.click(screen.getByTestId('terms_agreement'));

      fireEvent.click(screen.getByRole('button', { name: /create starter account/i }));
    });

    expect(handleSubmit).toHaveBeenCalledTimes(1);

    expect(handleSubmit.mock.calls[0][0]).toStrictEqual({
      name: 'John Doe',
      email: 'john.doe@example.com',
      password: 'Password123@@_',
      pwcheck: 'Password123@@_',
      organization: 'Acme Inc.',
      domain: 'acme.io',
      terms_agreement: true,
      privacy_agreement: true,
    });
  });

  afterEach(() => {
    vi.unmock('react-router-dom');
  });
});
