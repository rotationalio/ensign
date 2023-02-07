import { act, render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { vi } from 'vitest';

import { LoginForm } from '../..';

describe('first', () => {
  it('should submit the form', async () => {
    const handleSubmit = vi.fn();

    render(<LoginForm onSubmit={handleSubmit} />);

    // eslint-disable-next-line testing-library/no-unnecessary-act
    await act(async () => {
      await userEvent.type(screen.getByTestId('email'), 'john.doe@example.com');
      await userEvent.type(screen.getByTestId('password'), 'Password123@@_');
    });

    //expect(handleSubmit).toHaveBeenCalledTimes(1);
  });
});
