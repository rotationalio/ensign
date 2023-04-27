/* eslint-disable unused-imports/no-unused-imports */
import { fireEvent, screen, waitFor } from '@testing-library/react';
// import userEvent from '@testing-library/user-event';
import React from 'react';
import { vi } from 'vitest';

import { dynamicActivate } from '../../../../../I18n';
import { customRender } from '../../../../../utils/test-utils';
import { RegistrationForm } from '../..';

vi.mock('react-router-dom');

describe('RegistrationForm', () => {
  beforeEach(() => {
    dynamicActivate();
  });
  beforeAll(() => {
    Object.defineProperty(window, 'matchMedia', {
      writable: true,
      value: vi.fn().mockImplementation((query) => ({
        matches: false,
        media: query,
        onchange: null,
        addListener: vi.fn(), // Deprecated
        removeListener: vi.fn(), // Deprecated
        addEventListener: vi.fn(),
        removeEventListener: vi.fn(),
        dispatchEvent: vi.fn(),
      })),
    });
  });

  it.todo('should render form initial value');

  // it('should render form initial value', () => {
  //   const handleSubmit = vi.fn();
  //   customRender(<RegistrationForm onSubmit={handleSubmit} />);

  //   expect(screen.getByTestId('name')).toHaveAttribute('value', '');
  //   expect(screen.getByTestId('email')).toHaveAttribute('value', '');
  //   expect(screen.getByTestId('password')).toHaveAttribute('value', '');
  //   expect(screen.getByTestId('pwcheck')).toHaveAttribute('value', '');
  //   expect(screen.getByTestId('organization')).toHaveAttribute('value', '');
  //   expect(screen.getByTestId('terms_agreement')).not.toBeChecked();
  // });

  // describe('Name', () => {
  //   it('should display error message when name field is empty', async () => {
  //     const handleSubmit = vi.fn();
  //     customRender(<RegistrationForm onSubmit={handleSubmit} />);

  //     expect(screen.getByTestId('name')).toHaveAttribute('value', '');
  //     userEvent.click(screen.getByRole('button', { name: /create starter account/i }));

  //     await waitFor(() => {
  //       expect(screen.getByText(/The name is required./i)).toBeInTheDocument();
  //     });
  //   });
  // });

  // describe('Password eye icon renders', () => {
  //   it('should show icon with a closed eye, hidden password text, and accessible button name "Show Password" on render', () => {
  //     const handleSubmit = vi.fn();
  //     customRender(<RegistrationForm onSubmit={handleSubmit} />);
  //     const button = screen.getByTestId('button');

  //     userEvent.type(screen.getByTestId('pwcheck'), '');
  //     userEvent.click(screen.getByRole('button', { name: /create starter account/i }));

  //     await waitFor(() => {
  //       expect(screen.getByText(/Please re-enter your password to confirm./i)).toBeInTheDocument();
  //     });
  //   });
  // });

  // describe('Show password', () => {
  //   it('should show icon with open eye, password text, and accessible button name "Hide Password" when clicks on icon', () => {
  //     const handleSubmit = vi.fn();
  //     customRender(<RegistrationForm onSubmit={handleSubmit} />);
  //     const button = screen.getByTestId('button');

  //     fireEvent.click(button);
  //     expect(screen.getByTestId('showPassword')).toBeVisible;
  //     expect(screen.getByTestId('password')).toHaveAttribute('type', 'text');
  //     expect(button).toHaveAccessibleName('Hide Password');
  //   });
  // });

  // describe('Hide password', () => {
  //   it('should show icon with eye closed, hidden password text, and accessible button name "Show Password" when user clicks on eye icon a 2nd time', () => {
  //     const handleSubmit = vi.fn();
  //     customRender(<RegistrationForm onSubmit={handleSubmit} />);
  //     const button = screen.getByTestId('button');

  //     fireEvent.doubleClick(button);
  //     expect(screen.getByTestId('hidePassword')).toBeVisible;
  //     expect(screen.getByTestId('password')).toHaveAttribute('type', 'password');
  //     expect(button).toHaveAccessibleName('Show Password');
  //   });
  // });

  // afterEach(() => {
  //   vi.unmock('react-router-dom');
  // });
});
